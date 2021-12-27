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

package example

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "user:passwd@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		log.Printf("create db failed, %v", err)
		return
	}
}

type myOffsetPersist struct {
}

func (p *myOffsetPersist) Find(groupId, topic string, partition int) (int64, error) {
	var offset int64
	err := db.QueryRow("SELECT offset FROM devspore_offset_table WHERE `group_id`=? AND `topic`=? AND `partition`=?", groupId, topic, partition).Scan(&offset)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (p *myOffsetPersist) Save(groupId, topic string, partition int, offset int64) error {
	if _, err := p.Find(groupId, topic, partition); err == nil {
		_, err := db.Exec("UPDATE devspore_offset_table SET `offset`=? WHERE `group_id`=? AND `topic`=? AND `partition`=?", offset, groupId, topic, partition)
		return err
	}
	_, err := db.Exec("INSERT INTO devspore_offset_table (`offset`, `group_id`, `topic`, `partition`) VALUES (?,?,?,?)", offset, groupId, topic, partition)
	return err
}
