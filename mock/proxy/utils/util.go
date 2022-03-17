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

// Package utils Test tool
package utils

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dolthub/go-mysql-server/memory"
	mocksql "github.com/dolthub/go-mysql-server/sql"
	goredis "github.com/go-redis/redis/v8"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/mock/proxy"
	redisconfig "github.com/huaweicloud/devcloud-go/redis/config"
	mysqlconfig "github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

// mock services
var (
	EtcdMock   *mock.MockEtcd
	MysqlMocks = make([]*mock.MysqlMock, 0)
	RedisMocks = make([]*mock.RedisMock, 0)
	Proxys     = make([]*proxy.Proxy, 0)
)

func StartEtcdMock(addrs []string, dataDir string) {
	metadata := mock.NewEtcdMetadata()
	metadata.ClientAddrs = addrs
	metadata.DataDir = dataDir
	EtcdMock = &mock.MockEtcd{}
	EtcdMock.StartMockEtcd(metadata)
}

func StopEtcdMock(dataDir string) {
	var err error
	EtcdMock.StopMockEtcd()
	err = os.RemoveAll(dataDir)
	if err != nil {
		log.Println("ERROR: remove " + dataDir + " dir failed")
	}
}

func Start2MysqlMock(addrs []string) {
	var err error
	for _, addr := range addrs {
		mysqlMock := &mock.MysqlMock{
			User:         "root",
			Password:     "123456",
			Address:      addr,
			Databases:    []string{"ds0", "ds0-slave0", "ds0-slave1", "ds1", "ds1-slave0", "ds1-slave1"},
			MemDatabases: []*memory.Database{createTestDatabase("ds0", "user", addr)},
		}
		err = mysqlMock.StartMockMysql()
		if err != nil {
			log.Fatalln(err)
		}
		MysqlMocks = append(MysqlMocks, mysqlMock)
	}
}

func Stop2MysqlMock() {
	for _, mysqlMock := range MysqlMocks {
		mysqlMock.StopMockMysql()
	}
}

func createTestDatabase(dbName, tableName, address string) *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, mocksql.Schema{
		{Name: "id", Type: mocksql.Int64, Nullable: false, AutoIncrement: true, PrimaryKey: true, Source: tableName},
		{Name: "name", Type: mocksql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: mocksql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: mocksql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: mocksql.Timestamp, Nullable: false, Source: tableName},
	})

	ctx := mocksql.NewEmptyContext()
	db.AddTable(tableName, table)

	rows := []mocksql.Row{
		mocksql.NewRow(1, address+"John Doe", "jasonkay@doe.com", []string{"555-555-555"}, time.Now()),
		mocksql.NewRow(2, address+"John Doe", "johnalt@doe.com", []string{}, time.Now()),
		mocksql.NewRow(3, address+"Jane Doe", "jane@doe.com", []string{}, time.Now()),
		mocksql.NewRow(4, address+"Evil Bob", "jasonkay@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
	}

	for _, row := range rows {
		_ = table.Insert(ctx, row)
	}
	return db
}

func setSourcesNodes(mysqlAddrs []string) (map[string]*mysqlconfig.DataSourceConfiguration,
	map[string]*mysqlconfig.NodeConfiguration) {
	dataSource := make(map[string]*mysqlconfig.DataSourceConfiguration)
	nodes := make(map[string]*mysqlconfig.NodeConfiguration)
	for i, addr := range mysqlAddrs {
		stri := strconv.Itoa(i + 1)
		dataSource["ds"+stri] = &mysqlconfig.DataSourceConfiguration{
			URL:      "tcp(" + addr + ")/ds0?charset=utf8&parseTime=true",
			Username: "root",
			Password: "123456",
		}
		nodes["dc"+stri] = &mysqlconfig.NodeConfiguration{
			Master: "ds" + stri,
		}
	}
	return dataSource, nodes
}

func DCMysql(etcdAddrs, mysqlAddrs []string) *mysqlconfig.ClusterConfiguration {
	dataSource, nodes := setSourcesNodes(mysqlAddrs)
	return &mysqlconfig.ClusterConfiguration{
		Props: &mas.PropertiesConfiguration{
			AppID:        "123",
			MonitorID:    "456",
			DatabaseName: "12",
		},
		EtcdConfig: &etcd.EtcdConfiguration{
			Address:     etcdAddrs[0],
			Username:    "root",
			Password:    "root",
			HTTPSEnable: false,
		},
		RouterConfig: &mysqlconfig.RouterConfiguration{
			Nodes:  nodes,
			Active: "dc1",
		},
		DataSource: dataSource,
		Chaos: &mas.InjectionProperties{
			Active:     true,
			Duration:   50,
			Interval:   100,
			Percentage: 100,
			DelayInjection: &mas.DelayInjection{
				Active:     true,
				Percentage: 75,
				TimeMs:     1000,
				JitterMs:   500,
			},
			ErrorInjection: &mas.ErrorInjection{
				Active:     true,
				Percentage: 30,
			},
		},
	}
}

func Start2RedisMock(addrs []string) {
	var err error
	for _, addr := range addrs {
		redisMock := &mock.RedisMock{Addr: addr, Password: "123456"}
		err = redisMock.StartMockRedis()
		if err != nil {
			log.Fatalln(err)
		}
		addTestData(redisMock)
		RedisMocks = append(RedisMocks, redisMock)
	}
}

func Stop2RedisMock() {
	for _, redisMock := range RedisMocks {
		redisMock.StopMockRedis()
	}
}

func addTestData(redisMock *mock.RedisMock) {
	ctx := context.Background()
	rdb1 := goredis.NewClient(&goredis.Options{Addr: redisMock.Addr, Password: "123456"})
	rdb1.Set(ctx, "key", redisMock.Addr, 0)
	_ = rdb1.Close()
}

func setServers(redisAddrs []string) map[string]*redisconfig.ServerConfiguration {
	servers := make(map[string]*redisconfig.ServerConfiguration)
	for i, addr := range redisAddrs {
		stri := strconv.Itoa(i + 1)
		servers["ds"+stri] = &redisconfig.ServerConfiguration{
			Hosts:    addr,
			Password: "123456",
			Type:     redisconfig.ServerTypeNormal,
			Cloud:    "huawei cloud",
			Region:   "beijing",
			Azs:      "az1",
		}
	}
	return servers
}

func DCRedis(etcdAddrs, redisAddrs []string) *redisconfig.Configuration {
	servers := setServers(redisAddrs)
	configuration := &redisconfig.Configuration{
		Props: &mas.PropertiesConfiguration{
			AppID:        "123",
			MonitorID:    "456",
			DatabaseName: "789",
		},
		EtcdConfig: &etcd.EtcdConfiguration{
			Address:     etcdAddrs[0],
			Username:    "root",
			Password:    "root",
			HTTPSEnable: false,
		},
		RedisConfig: &redisconfig.RedisConfiguration{
			Servers: servers,
			ConnectionPoolConfig: &redisconfig.RedisConnectionPoolConfiguration{
				Enable: false,
			},
		},
		RouteAlgorithm: "single-read-write",
		Active:         "ds1",
		Chaos: &mas.InjectionProperties{
			Active:     true,
			Duration:   50,
			Interval:   100,
			Percentage: 100,
			DelayInjection: &mas.DelayInjection{
				Active:     true,
				Percentage: 100,
				TimeMs:     1000,
				JitterMs:   500,
			},
			ErrorInjection: &mas.ErrorInjection{
				Active:     true,
				Percentage: 30,
			},
		},
	}
	return configuration
}

func Start2Proxy(addrs []string, proxys []string, mock proxy.MockType) {
	var err error
	for i := 0; i < len(proxys); i++ {
		tProxy := proxy.NewProxy(addrs[i], proxys[i], mock)
		err = tProxy.StartProxy()
		if err != nil {
			log.Fatalln(err)
		}
		Proxys = append(Proxys, tProxy)
	}
}

func Stop2Proxy() {
	for _, tProxy := range Proxys {
		tProxy.StopProxy()
	}
}
