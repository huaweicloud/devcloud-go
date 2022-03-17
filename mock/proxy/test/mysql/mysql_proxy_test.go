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

// Package mysql proxy test cases
package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/huaweicloud/devcloud-go/mock/proxy"
	proxymysql "github.com/huaweicloud/devcloud-go/mock/proxy/proxy-mysql"
	"github.com/huaweicloud/devcloud-go/mock/proxy/utils"
	devspore "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

func TestMysqlMock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mysql")
}

var _ = Describe("Mysql", func() {
	var (
		err                  error
		etcdAddrs            = []string{"127.0.0.1:2379"}
		dataDir              = "etcd_data"
		mysqlAddrs           = []string{"127.0.0.1:13306", "127.0.0.1:13307"}
		proxyAddrs           = []string{"127.0.0.1:23306", "127.0.0.1:23307"}
		clusterConfiguration = utils.DCMysql(etcdAddrs, proxyAddrs)
		datasource           = "/mas-monitor/conf/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/datasource"
		router               = "/mas-monitor/conf/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/router"
		activekey            = "/mas-monitor/status/db/services/" + clusterConfiguration.Props.AppID + "/" + clusterConfiguration.Props.MonitorID + "/database/" + clusterConfiguration.Props.DatabaseName + "/active"
		db                   *sql.DB
		client               *clientv3.Client
		val1                 = mysqlAddrs[0] + "John Doe"
		val2                 = mysqlAddrs[1] + "John Doe"
	)
	BeforeSuite(func() {
		utils.Start2MysqlMock(mysqlAddrs)
		utils.Start2Proxy(mysqlAddrs, proxyAddrs, proxy.Mysql)
		utils.StartEtcdMock(etcdAddrs, dataDir)

		client, err = clientv3.New(clientv3.Config{Endpoints: etcdAddrs, Username: "root", Password: "root"})

		datasourceStr, err := json.Marshal(clusterConfiguration.DataSource)
		if err != nil {
			fmt.Println(err)
		}
		routerConfigStr, err := json.Marshal(clusterConfiguration.RouterConfig)
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
		_, err = client.Put(ctx, activekey, "dc1", clientv3.WithPrevKV())
		if err != nil {
			log.Println(err)
		}
		cancel()
		if clusterConfiguration.Chaos != nil {
			clusterConfiguration.Chaos.Active = false
		}
		devspore.SetClusterConfiguration(clusterConfiguration)
		db, err = sql.Open("devspore_mysql", "")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)
	})

	AfterSuite(func() {
		err = db.Close()
		if err != nil {
			log.Println(err)
		}
		err = client.Close()
		if err != nil {
			log.Println(err)
		}
		utils.StopEtcdMock(dataDir)
		utils.Stop2Proxy()
		utils.Stop2MysqlMock()
	})

	AfterEach(func() {
		utils.Proxys[0].DeleteAllRule()
	})

	It("Select", func() {
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("SelectDelay", func() {
		err = utils.Proxys[0].AddDelay("delay", 1500, 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("SelectJitter", func() {
		err = utils.Proxys[0].AddJitter("jitter", 1500, 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("SelectDrop", func() {
		err = utils.Proxys[0].AddDrop("drop", 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("SelectReturnEmpty", func() {
		err = utils.Proxys[0].AddReturnEmpty("returnEmpty", 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("SelectReturnErr", func() {
		err = utils.Proxys[0].AddReturnErr("returnErr", proxymysql.ERNoSuchTable, 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val1))
	})

	It("MockChang", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = client.Put(ctx, activekey, "dc2", clientv3.WithPrevKV())
		cancel()
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(time.Second)
	})

	It("SelectDrop-MockChang", func() {
		err = utils.Proxys[0].AddDrop("drop", 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val2))
	})

	It("SelectReturnEmpty-MockChang", func() {
		err = utils.Proxys[0].AddReturnEmpty("returnEmpty", 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val2))
	})

	It("SelectReturnErr-MockChang", func() {
		err = utils.Proxys[0].AddReturnErr("returnErr", proxymysql.ERNoSuchTable, 0, "", "select")
		if err != nil {
			log.Println(err)
		}
		var name string
		err = db.QueryRow("select name from user where id=1").Scan(&name)
		if err != nil {
			log.Println(err)
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal(val2))
	})
})
