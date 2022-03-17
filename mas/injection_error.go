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
 *
 */

package mas

import (
	"io"
	"syscall"

	"github.com/dolthub/vitess/go/mysql"
)

// RuntimeError
type errorString string

func (e errorString) RuntimeError() {}

func (e errorString) Error() string {
	return "runtime error: " + string(e)
}

// RedisError redisError
type RedisError string

func (e RedisError) Error() string { return string(e) }

func (RedisError) RedisError() {}

var (
	SocketErr     = syscall.Errno(10061)
	IOErr1        = io.ErrClosedPipe
	IOErr2        = io.ErrUnexpectedEOF
	NilPointerErr = errorString("invalid memory address or nil pointer dereference")

	SQLErr        = mysql.NewSQLError(mysql.CRUnknownError, mysql.SSUnknownSQLState, "SQLErr")
	SQLTimeoutErr = mysql.NewSQLError(mysql.CRServerGone, mysql.SSUnknownSQLState, "SQLTimeoutErr")

	RedisCommandUKErr  = RedisError("ERR unknown command")
	RedisCommandArgErr = RedisError("ERR wrong number of arguments for command")
)

func MysqlErrors() []error {
	errs := []error{
		SocketErr,
		IOErr1,
		IOErr2,
		NilPointerErr,
		SQLErr,
		SQLTimeoutErr,
	}
	return errs
}

func RedisErrors() []error {
	errs := []error{
		SocketErr,
		IOErr1,
		IOErr2,
		NilPointerErr,
		RedisCommandUKErr,
		RedisCommandArgErr,
	}
	return errs
}
