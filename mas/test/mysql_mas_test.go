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

// Package test contains mas test cases
package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/huaweicloud/devcloud-go/mock/proxy/utils"
	devspore "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

var (
	etcdAddrs            = []string{"127.0.0.1:2379"}
	dataDir              = "etcd_data"
	mysqlAddrs           = []string{"127.0.0.1:13306", "127.0.0.1:13307"}
	clusterConfiguration = utils.DCMysql(etcdAddrs, mysqlAddrs)
	datasource           = "/mas-monitor/conf/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/datasource"
	router               = "/mas-monitor/conf/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/router"
	activekey1           = "/mas-monitor/status/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/active"
)

func TestMysqlMas(t *testing.T) {
	utils.Start2MysqlMock(mysqlAddrs)
	defer utils.Stop2MysqlMock()
	utils.StartEtcdMock(etcdAddrs, dataDir)
	defer utils.StopEtcdMock(dataDir)

	client, err := clientv3.New(clientv3.Config{Endpoints: etcdAddrs, Username: "XXXX", Password: "XXXX"})
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = client.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	routerConfigStr, err := json.Marshal(clusterConfiguration.RouterConfig)
	if err != nil {
		fmt.Println(err)
	}
	datasourceStr, err := json.Marshal(clusterConfiguration.DataSource)
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, datasource, string(datasourceStr), clientv3.WithPrevKV())
	if err != nil {
		log.Println(err)
	}
	_, err = client.Put(ctx, router, string(routerConfigStr), clientv3.WithPrevKV())
	if err != nil {
		log.Println(err)
	}
	cancel()
	devspore.SetClusterConfiguration(clusterConfiguration)
	db, err := sql.Open("devspore_mysql", "")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	time.Sleep(time.Second)
	nodeTest(client, db)
}

func nodeTest(client *clientv3.Client, db *sql.DB) {
	for key := range clusterConfiguration.RouterConfig.Nodes {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.Put(ctx, activekey1, key, clientv3.WithPrevKV())
		if err != nil {
			log.Println(err)
		}
		cancel()
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			var name string
			err = db.QueryRow("select name from user where id=1").Scan(&name)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(name)
			time.Sleep(time.Second)
		}
	}
}
