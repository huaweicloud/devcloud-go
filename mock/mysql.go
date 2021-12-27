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

package mock

import (
	"fmt"
	"log"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

type MysqlMetaData struct {
	User         string
	Password     string
	Address      string
	Databases    []string
	MemDatabases []*memory.Database // see mysql_test.go Line#53-77
}

var mysqlServers []*server.Server

func StartMockMysql(metadata MysqlMetaData) {
	var databases = []sql.Database{information_schema.NewInformationSchemaDatabase()}
	for _, db := range metadata.Databases {
		databases = append(databases, memory.NewDatabase(db))
	}
	for _, db := range metadata.MemDatabases {
		databases = append(databases, db)
	}
	engine := sqle.NewDefault(sql.NewDatabaseProvider(databases...))
	config := server.Config{
		Protocol: "tcp",
		Address:  metadata.Address,
		Auth:     auth.NewNativeSingle(metadata.User, metadata.Password, auth.AllPermissions),
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		log.Printf("ERROR: create mysql server failed, %v", err)
		return
	}
	go func() {
		err = s.Start()
		if err != nil {
			log.Printf("ERROR: start mysql server failed, %v", err)
			return
		}
	}()

	fmt.Println("mysql-server started!")
	mysqlServers = append(mysqlServers, s)
	return
}

func StopMockMysql() {
	if len(mysqlServers) == 0 {
		return
	}
	for _, s := range mysqlServers {
		if err := s.Close(); err != nil {
			log.Printf("ERROR: close mysql server [%s] failed, %v", s.Listener.Addr(), err)
		}
	}
	fmt.Println("mysql-server stop!")
}
