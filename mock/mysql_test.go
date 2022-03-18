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
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/dolthub/go-mysql-server/memory"
	mocksql "github.com/dolthub/go-mysql-server/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMysqlMock(t *testing.T) {
	metadata := MysqlMock{
		User:         "root",
		Password:     "root",
		Address:      "127.0.0.1:3318",
		Databases:    []string{"mydb"},
		MemDatabases: []*memory.Database{createTestDatabase("mydb", "user")},
	}
	err := metadata.StartMockMysql()
	if err != nil {
		log.Fatalln(err)
	}
	defer metadata.StopMockMysql()

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3318)/mydb")
	defer func() {
		err = db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		t.Error(err)
		return
	}
	var name, email string
	err = db.QueryRow("SELECT name, email FROM user WHERE id=?", 1).Scan(&name, &email)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, name, "John Doe")
	assert.Equal(t, email, "jasonkay@doe.com")
}

func createTestDatabase(dbName, tableName string) *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, mocksql.Schema{
		{Name: "id", Type: mocksql.Int64, Nullable: false, AutoIncrement: true, PrimaryKey: true, Source: tableName},
		{Name: "name", Type: mocksql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: mocksql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: mocksql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: mocksql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := mocksql.NewEmptyContext()

	rows := []mocksql.Row{
		mocksql.NewRow(1, "John Doe", "jasonkay@doe.com", []string{"555-555-555"}, time.Now()),
		mocksql.NewRow(2, "John Doe", "johnalt@doe.com", []string{}, time.Now()),
		mocksql.NewRow(3, "Jane Doe", "jane@doe.com", []string{}, time.Now()),
		mocksql.NewRow(4, "Evil Bob", "jasonkay@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
	}

	for _, row := range rows {
		err := table.Insert(ctx, row)
		log.Println(err)
	}
	return db
}
