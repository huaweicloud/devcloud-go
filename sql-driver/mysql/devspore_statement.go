/*
 * Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * 2021.11.15-Changed modify the implements of statement's function.
 * 			Huawei Technologies Co., Ltd.
 */

package mysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"reflect"
)

type devsporeStmt struct {
	ctx   context.Context
	dc    *devsporeConn
	stmt  driver.Stmt
	query string
	dsn   string
}

// Close devsporeStmt
func (dsmt *devsporeStmt) Close() error {
	var err error
	if dsmt.stmt != nil {
		err = dsmt.stmt.Close()
	}
	dsmt.stmt = nil
	return err
}

// NumInput implements driver.Stmt, return query's params count.
func (dsmt *devsporeStmt) NumInput() int {
	req := &executorReq{
		ctx:        dsmt.ctx,
		query:      dsmt.query,
		methodName: StmtNumInput,
		dc:         dsmt.dc,
		dsmt:       dsmt,
	}
	resp := dsmt.dc.executor.tryExecute(req)
	if resp.err != nil {
		log.Printf("ERROR: devsporeStatement execute NumInput failed, err %v", resp.err)
	}
	return resp.numInput
}

// ColumnConverter implements driver.ColumnConverter.
func (dsmt *devsporeStmt) ColumnConverter(idx int) driver.ValueConverter {
	return converter{}
}

// CheckNamedValue implements driver.NamedValueChecker.
func (dsmt *devsporeStmt) CheckNamedValue(nv *driver.NamedValue) error {
	var err error
	nv.Value, err = converter{}.ConvertValue(nv.Value)
	return err
}

// QueryContext implements driver.StmtQueryContext.
func (dsmt *devsporeStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	req := &executorReq{
		ctx:        ctx,
		query:      dsmt.query,
		ctxArgs:    args,
		methodName: StmtQueryContext,
		dc:         dsmt.dc,
		dsmt:       dsmt,
	}
	resp := dsmt.dc.executor.tryExecute(req)
	if resp.err != nil {
		log.Printf("ERROR: devsporeStatement execute QueryContext failed, err %v", resp.err)
	}
	return resp.rows, resp.err
}

// ExecContext implements driver.StmtExecContext.
func (dsmt *devsporeStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	req := &executorReq{
		ctx:        ctx,
		query:      dsmt.query,
		ctxArgs:    args,
		methodName: StmtExecContext,
		dc:         dsmt.dc,
		dsmt:       dsmt,
	}
	resp := dsmt.dc.executor.tryExecute(req)
	if resp.err != nil {
		log.Printf("ERROR: devsporeStatement execute ExecContext failed, err %v", resp.err)
	}
	return resp.result, resp.err
}

// getStatement get an actual statement from devsporeStmt if exists or create a new statement.
func (dsmt *devsporeStmt) getStatement(ctx context.Context, dsn string) (driver.Stmt, error) {
	if dsmt.stmt != nil && dsmt.dsn == dsn {
		return dsmt.stmt, nil
	}
	dsmt.dsn = dsn
	conn, err := dsmt.dc.getConnection(ctx, dsn)
	if err != nil {
		return nil, err
	}
	conPrepareCtx, ok := conn.(driver.ConnPrepareContext)
	if !ok {
		return nil, errors.New("type assertion ConnPrepareContext failed")
	}
	stmt, err := conPrepareCtx.PrepareContext(dsmt.ctx, dsmt.query)
	if err != nil {
		return nil, err
	}
	dsmt.stmt = stmt
	return stmt, nil
}

// copy from go-sql-driver mysql
type converter struct{}

// ConvertValue mirrors the reference/default converter in database/sql/driver
// with _one_ exception.  We support uint64 with their high bit and the default
// implementation does not.  This function should be kept in sync with
// database/sql/driver defaultConverter.ConvertValue() except for that
// deliberate difference.
func (c converter) ConvertValue(v interface{}) (driver.Value, error) {
	if driver.IsValue(v) {
		return v, nil
	}

	if vr, ok := v.(driver.Valuer); ok {
		sv, err := callValuerValue(vr)
		if err != nil {
			return nil, err
		}
		if !driver.IsValue(sv) {
			return nil, fmt.Errorf("non-Value type %T returned from Value", sv)
		}
		return sv, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		} else {
			return c.ConvertValue(rv.Elem().Interface())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Slice:
		ek := rv.Type().Elem().Kind()
		if ek == reflect.Uint8 {
			return rv.Bytes(), nil
		}
		return nil, fmt.Errorf("unsupported type %T, a slice of %s", v, ek)
	case reflect.String:
		return rv.String(), nil
	}
	return nil, fmt.Errorf("unsupported type %T, a %s", v, rv.Kind())
}

var valuerReflectType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

// callValuerValue returns vr.Value(), with one exception:
// If vr.Value is an auto-generated method on a pointer type and the
// pointer is nil, it would panic at runtime in the panicwrap
// method. Treat it like nil instead.
//
// This is so people can implement driver.Value on value types and
// still use nil pointers to those types to mean nil/NULL, just like
// string/*string.
//
// This is an exact copy of the same-named unexported function from the
// database/sql package.
func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Ptr &&
		rv.IsNil() &&
		rv.Type().Elem().Implements(valuerReflectType) {
		return nil, nil
	}
	return vr.Value()
}

// Exec Deprecated
func (dsmt *devsporeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

// Query Deprecated
func (dsmt *devsporeStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, nil
}
