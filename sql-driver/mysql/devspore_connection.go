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
	"sync"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
)

type devsporeConn struct {
	clusterDataSource *datasource.ClusterDataSource
	cachedConn        sync.Map
	inTransaction     bool
	executor          *executor
}

// Begin Deprecated
func (dc *devsporeConn) Begin() (driver.Tx, error) {
	return nil, nil
}

// Close implements driver.Conn, will close all cached connection.
func (dc *devsporeConn) Close() error {
	var err error
	dc.cachedConn.Range(func(key, value interface{}) bool {
		err = value.(driver.Conn).Close()
		return true
	})
	return err
}

// Prepare Deprecated
func (dc *devsporeConn) Prepare(query string) (driver.Stmt, error) {
	dsmt := &devsporeStmt{
		ctx:   context.Background(),
		query: query,
		dc:    dc,
	}
	return dsmt, nil
}

// Exec Deprecated
func (dc *devsporeConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return nil, nil
}

// Query Deprecated
func (dc *devsporeConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return nil, nil
}

// Ping implements driver.Pinger interface
func (dc *devsporeConn) Ping(ctx context.Context) (err error) {
	return nil
}

// BeginTx implements driver.ConnBeginTx interface
func (dc *devsporeConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	dc.inTransaction = true
	req := &executorReq{
		ctx:        ctx,
		opts:       opts,
		methodName: BeginTx,
		dc:         dc,
	}
	resp := dc.executor.tryExecute(req)
	if resp.err != nil {
		log.Printf("ERROR: devsporeConnection execute BeginTx failed, err %v", resp.err)
		return nil, resp.err
	}

	return &devsporeTx{
		dc:       dc,
		actualTx: resp.tx,
	}, nil
}

// QueryContext implements driver.QueryerContext interface
func (dc *devsporeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	req := &executorReq{
		ctx:        ctx,
		query:      query,
		ctxArgs:    args,
		methodName: QueryContext,
		dc:         dc,
	}
	resp := dc.executor.tryExecute(req)
	if resp.err != nil && resp.err != driver.ErrSkip {
		log.Printf("ERROR: devsporeConnection execute QueryContext failed, err %v", resp.err)
	}
	return resp.rows, resp.err
}

// ExecContext implements driver.ExecerContext interface
func (dc *devsporeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	req := &executorReq{
		ctx:        ctx,
		query:      query,
		ctxArgs:    args,
		methodName: ExecContext,
		dc:         dc,
	}
	resp := dc.executor.tryExecute(req)
	if resp.err != nil && resp.err != driver.ErrSkip {
		log.Printf("ERROR: devsporeConnection execute ExecContext failed, err %v", resp.err)
	}
	return resp.result, resp.err
}

// PrepareContext implements driver.ConnPrepareContext interface
func (dc *devsporeConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	dsmt := &devsporeStmt{
		ctx:   ctx,
		query: query,
		dc:    dc,
	}
	return dsmt, nil
}

// CheckNamedValue implements driver.NamedValueChecker interface
func (dc *devsporeConn) CheckNamedValue(nv *driver.NamedValue) (err error) {
	nv.Value, err = converter{}.ConvertValue(nv.Value)
	return
}

// getConnection from cache or new connection according to dsn
func (dc *devsporeConn) getConnection(ctx context.Context, actualDSN string) (driver.Conn, error) {
	if conn, ok := dc.cachedConn.Load(actualDSN); ok {
		return conn.(driver.Conn), nil
	}
	conn, err := dc.creationConn(ctx, actualDSN)
	if err != nil {
		return nil, err
	}
	dc.cachedConn.Store(actualDSN, conn)
	return conn, nil
}

// creationConn according to dsn
func (dc *devsporeConn) creationConn(ctx context.Context, dsn string) (driver.Conn, error) {
	driCtx, ok := actualDriver.(driver.DriverContext)
	if !ok {
		return nil, errors.New("type assertion driver.DriverContext failed")
	}
	connector, err := driCtx.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}
	return connector.Connect(ctx)
}
