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

// Package proxymysql Mysql TCP-based fault injection
package proxymysql

import (
	"github.com/dolthub/vitess/go/mysql"
	"github.com/huaweicloud/devcloud-go/mock/proxy"
)

type ProxyMysql struct {
	*proxy.Proxy
}

func NewProxy(server, addr string) *ProxyMysql {
	proxyMysql := proxy.NewProxy(server, addr, proxy.Mysql)
	return &ProxyMysql{
		Proxy: proxyMysql,
	}
}

// mysql chaos error
var (
	Nil            = mysql.NewSQLError(mysql.ERUnknownError, mysql.SSUnknownSQLState, "nil")
	ServerShutdown = mysql.NewSQLError(mysql.ERUnknownError, mysql.SSUnknownSQLState, "Server shutdown in progress")
	ERUnknownError = mysql.NewSQLError(mysql.ERUnknownError, mysql.SSUnknownSQLState, "unknown error")

	ERNoSuchTable    = mysql.NewSQLError(mysql.ERNoSuchTable, mysql.SSUnknownSQLState, "table not found")
	ERDbCreateExists = mysql.NewSQLError(mysql.ERDbCreateExists, mysql.SSUnknownSQLState, "can't create database; database exists")
	ERSubqueryNo1Row = mysql.NewSQLError(mysql.ERSubqueryNo1Row, mysql.SSUnknownSQLState, "the subquery returned more than 1 row")
)
