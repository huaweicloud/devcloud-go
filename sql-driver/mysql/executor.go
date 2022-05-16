/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2021.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package mysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/router"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/util"
)

const (
	defaultRetryTimes   = 1
	defaultRetryDelay   = 1000  // ms
	exclusiveRetryDelay = 60000 // ms
)

var (
	errNoDatasource  = errors.New("no available data source")
	errTypeAssertion = errors.New("type assertion err")
)

type methodName string

const (
	BeginTx          methodName = "BeginTx"
	QueryContext     methodName = "QueryContext"
	ExecContext      methodName = "ExecContext"
	StmtQueryContext methodName = "stmt.QueryContext"
	StmtExecContext  methodName = "stmt.ExecContext"
	StmtNumInput     methodName = "stmt.NumInput"
)

// executorReq contains all connection and statement method params
type executorReq struct {
	dc   *devsporeConn
	dsmt *devsporeStmt
	// original req params
	ctx        context.Context
	query      string
	ctxArgs    []driver.NamedValue
	opts       driver.TxOptions
	methodName methodName
}

// executorResp contains all connection and statement method return params
type executorResp struct {
	rows     driver.Rows
	result   driver.Result
	tx       driver.Tx
	numInput int

	err error
}

type executor struct {
	exclusives sync.Map
	retryTimes int
	retryDelay int // ms
}

func newExecutor(retry *config.RetryConfiguration) *executor {
	e := &executor{
		retryTimes: defaultRetryTimes,
		retryDelay: defaultRetryDelay,
		exclusives: sync.Map{},
	}
	if retry != nil && retry.Times != "" {
		if retryTimes, err := strconv.Atoi(retry.Times); err != nil {
			e.retryTimes = retryTimes
		}
	}
	if retry != nil && retry.Delay != "" {
		if retryDelay, err := strconv.Atoi(retry.Delay); err != nil {
			e.retryDelay = retryDelay
		}
	}
	return e
}

// from cluster datasource choose a node datasource
func (e *executor) tryExecute(req *executorReq) *executorResp {
	// insure parse sql only once
	isSQLOnlyRead := util.IsOnlyRead(req.query)
	// route node datasource
	clusterRuntimeCtx := &router.RuntimeContext{DataSource: req.dc.clusterDataSource}
	routeAlgorithm := req.dc.clusterDataSource.RouterConfiguration.RouteAlgorithm
	nodeTargetDataSource := router.NewClusterRouter(routeAlgorithm).Route(
		isSQLOnlyRead, clusterRuntimeCtx, map[datasource.DataSource]bool{})
	if nodeTargetDataSource == nil {
		return &executorResp{err: errNoDatasource}
	}
	return e.tryNodeExecute(req, nodeTargetDataSource, isSQLOnlyRead)
}

// from node datasource choose an actual datasource to execute connection or statement method.
func (e *executor) tryNodeExecute(req *executorReq, nodeTargetDataSource datasource.DataSource,
	isSQLOnlyRead bool) *executorResp {
	var resp = &executorResp{}
nodeRetry:
	for {
		nodeRuntimeCtx := &router.RuntimeContext{
			DataSource:    nodeTargetDataSource,
			InTransaction: req.dc.inTransaction,
			RequestId:     idGenerator.GetId(),
		}
		actualExclusives := e.filterExclusive()
		targetDataSource := router.NewNodeRouter().Route(isSQLOnlyRead, nodeRuntimeCtx, actualExclusives)
		if targetDataSource == nil {
			break
		}
		actualTargetDataSource, ok := targetDataSource.(*datasource.ActualDataSource)
		if !ok {
			e.addExclusives(req, actualTargetDataSource)
			continue
		}

		times := 0
	retry:
		for times < e.retryTimes {
			// execute
			resp = e.execute(req, actualTargetDataSource.Dsn)
			actualTargetDataSource.LastRetryTime = time.Now().UnixNano() / 1e6
			switch {
			case resp.err == nil:
				// remove actualTargetDataSource from exclusives if exists
				actualTargetDataSource.Available = true
				actualTargetDataSource.RetryTimes = 0
				e.exclusives.Delete(actualTargetDataSource)
				return resp
			case resp.err == driver.ErrSkip: // when conn.QueryContext with args, db will return driver.ErrSkip to continue
				break nodeRetry
			case !isErrorRecoverable(resp.err): // when error isn't recoverable, return directly
				break nodeRetry
			case actualTargetDataSource.Available: // retry only when datasource is available and error is recoverable
				times++
			default:
				break retry
			}
			req.dc.cachedConn.Delete(actualTargetDataSource.Dsn)
			log.Printf("WARNING: execute '%s' failed, retried times %d", req.methodName, times)
			time.Sleep(time.Millisecond * time.Duration(e.retryDelay))
		}
		e.addExclusives(req, actualTargetDataSource)
		log.Printf("WARNING: datasource '%s' is unavailable, add to exclusives", actualTargetDataSource.Name)
	}
	return resp
}

func (e *executor) addExclusives(req *executorReq, actualTargetDataSource *datasource.ActualDataSource) {
	req.dc.cachedConn.Delete(actualTargetDataSource.Dsn)
	if req.dsmt != nil {
		req.dsmt.stmt = nil
	}
	actualTargetDataSource.Available = false
	actualTargetDataSource.RetryTimes = e.retryTimes
	e.exclusives.Store(actualTargetDataSource, true)
}

// execute directly if the actual datasource is available
func (e *executor) execute(req *executorReq, dsn string) *executorResp {
	var (
		rows     driver.Rows
		result   driver.Result
		tx       driver.Tx
		conn     driver.Conn
		numInput int
		err      error
	)
	if conn, err = req.dc.getConnection(req.ctx, dsn); err != nil {
		return &executorResp{err: err}
	}
	switch req.methodName {
	// conn methods
	case BeginTx:
		if connBeginTx, ok := conn.(driver.ConnBeginTx); ok {
			tx, err = connBeginTx.BeginTx(req.ctx, req.opts)
		} else {
			err = errTypeAssertion
		}
	case QueryContext:
		if queryerCtx, ok := conn.(driver.QueryerContext); ok {
			rows, err = queryerCtx.QueryContext(req.ctx, req.query, req.ctxArgs)
		} else {
			err = errTypeAssertion
		}
	case ExecContext:
		if execerCtx, ok := conn.(driver.ExecerContext); ok {
			result, err = execerCtx.ExecContext(req.ctx, req.query, req.ctxArgs)
		} else {
			err = errTypeAssertion
		}
	// statement methods
	case StmtQueryContext:
		rows, err = stmtQueryContext(req, dsn)
	case StmtExecContext:
		result, err = stmtExecContext(req, dsn)
	case StmtNumInput:
		numInput, err = stmtNumInput(req, dsn)
	}

	return &executorResp{
		rows:     rows,
		result:   result,
		tx:       tx,
		numInput: numInput,
		err:      err,
	}
}

// filterExclusive remove data sources that have been recovered and can be retried from the blacklist
func (e *executor) filterExclusive() map[datasource.DataSource]bool {
	actualExclusives := map[datasource.DataSource]bool{}
	e.exclusives.Range(func(key, value interface{}) bool {
		actual, ok := key.(*datasource.ActualDataSource)
		if !ok || time.Now().UnixNano()/1e6-actual.LastRetryTime < exclusiveRetryDelay {
			actualExclusives[actual] = true
		}
		return true
	})
	return actualExclusives
}

// statement methods

func stmtQueryContext(req *executorReq, dsn string) (driver.Rows, error) {
	stmt, err := req.dsmt.getStatement(req.ctx, dsn)
	if err != nil {
		return nil, err
	}
	if stmtQueryCtx, ok := stmt.(driver.StmtQueryContext); ok {
		return stmtQueryCtx.QueryContext(req.ctx, req.ctxArgs)
	}
	return nil, errTypeAssertion
}

func stmtExecContext(req *executorReq, dsn string) (driver.Result, error) {
	stmt, err := req.dsmt.getStatement(req.ctx, dsn)
	if err != nil {
		return nil, err
	}
	if stmtExecCtx, ok := stmt.(driver.StmtExecContext); ok {
		return stmtExecCtx.ExecContext(req.ctx, req.ctxArgs)
	}
	return nil, errTypeAssertion
}

func stmtNumInput(req *executorReq, dsn string) (int, error) {
	stmt, err := req.dsmt.getStatement(req.ctx, dsn)
	if err != nil {
		return -1, err
	}
	return stmt.NumInput(), nil
}
