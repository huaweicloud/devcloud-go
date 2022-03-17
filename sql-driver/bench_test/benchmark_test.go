/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2022.
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

// Package bench_test Connection test performance comparison
package bench_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	devspore "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

var (
	url                   = "tcp(127.0.0.1:3306)/kcluster?parseTime=true"
	username              = "root"
	password              = "123456"
	sqlDb, devSporeDb     *sql.DB
	sqlOrm, devSporeOrm   orm.Ormer
	sqlGorm, devSporeGorm *gorm.DB
	err                   error
)

type Metadata struct {
	Id          int       `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name;primaryKey"`
	Defaultdata string    `gorm:"column:defaultdata"`
	List        string    `gorm:"column:list"`
	Createtime  time.Time `gorm:"column:createtime;autoCreateTime"`
	Updatetime  time.Time `gorm:"column:updatetime;autoUpdateTime"`
}

func (t *Metadata) TableName() string {
	return "metadatas"
}

func mysqlConfiguration() *config.ClusterConfiguration {
	return &config.ClusterConfiguration{
		RouterConfig: &config.RouterConfiguration{
			Nodes: map[string]*config.NodeConfiguration{
				"dc1": {
					Master: "ds1",
				},
			},
			Active: "dc1",
		},
		DataSource: map[string]*config.DataSourceConfiguration{
			"ds1": {
				URL:      url,
				Username: username,
				Password: password,
			},
		},
	}
}

func startService() {
	devspore.SetClusterConfiguration(mysqlConfiguration())
	orm.RegisterModel(new(Metadata))
	err = orm.RegisterDriver("devspore_mysql", orm.DRMySQL)
	if err != nil {
		log.Fatalln(err)
	}

	// Native Data Source
	sqlDb, err = sql.Open("mysql", username+":"+password+"@"+url)
	if err != nil {
		log.Fatalln(err)
	}
	devSporeDb, err = sql.Open("devspore_mysql", "")
	if err != nil {
		log.Fatalln(err)
	}

	// Beego-orm
	sqlOrm, err = orm.NewOrmWithDB("mysql", "sqlOrm", sqlDb)
	if err != nil {
		log.Fatalln(err)
	}
	devSporeOrm, err = orm.NewOrmWithDB("devspore_mysql", "devSporeOrm", devSporeDb)
	if err != nil {
		log.Fatalln(err)
	}

	// Gorm
	sqlGorm, err = gorm.Open(mysql.New(mysql.Config{DriverName: "mysql", DSN: username + ":" + password + "@" + url}))
	if err != nil {
		log.Fatalln(err)
	}
	devSporeGorm, err = gorm.Open(mysql.New(mysql.Config{DriverName: "devspore_mysql", DSN: ""}))
	if err != nil {
		log.Fatalln(err)
	}
}
func stopService() {
	err = sqlDb.Close()
	if err != nil {
		log.Println(err)
	}
	err = devSporeDb.Close()
	if err != nil {
		log.Println(err)
	}
}

func BenchmarkDB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		err = sqlDb.QueryRow("select defaultdata from metadatas where name=?", "bb").Scan(&res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkDevSporeDB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		err = devSporeDb.QueryRow("select defaultdata from metadatas where name=?", "bb").Scan(&res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkOrm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		if err = sqlOrm.Raw("select defaultdata from metadatas where name=?", "bb").QueryRow(&res); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkDevSporeOrm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		if err = devSporeOrm.Raw("select defaultdata from metadatas where name=?", "bb").QueryRow(&res); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkGorm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		if err = sqlGorm.Raw("select defaultdata from metadatas where name=?", "bb").Scan(&res).Error; err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkDevSporeGorm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res string
		if err = devSporeGorm.Raw("select defaultdata from metadatas where name=?", "bb").Scan(&res).Error; err != nil {
			fmt.Println(err)
		}
	}
}

func TestMain(m *testing.M) {
	startService()
	m.Run()
	stopService()
}
