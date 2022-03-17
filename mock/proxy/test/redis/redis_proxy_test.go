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

// Package redis proxy test cases
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/huaweicloud/devcloud-go/mock/proxy"
	proxyredis "github.com/huaweicloud/devcloud-go/mock/proxy/proxy-redis"
	"github.com/huaweicloud/devcloud-go/mock/proxy/utils"
	"github.com/huaweicloud/devcloud-go/redis"
)

func TestRedisMock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis")
}

var _ = Describe("Redis", func() {
	var (
		err           error
		etcdAddrs     = []string{"127.0.0.1:2379"}
		dataDir       = "etcd_data"
		redisAddrs    = []string{"127.0.0.1:16379", "127.0.0.1:16380"}
		proxyAddrs    = []string{"127.0.0.1:26379", "127.0.0.1:26380"}
		configuration = utils.DCRedis(etcdAddrs, proxyAddrs)
		servers       = "/mas-monitor/conf/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/servers"
		algorithm     = "/mas-monitor/conf/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/route-algorithm"
		activekey     = "/mas-monitor/status/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/active"
		redisClient   *redis.DevsporeClient
		client        *clientv3.Client
	)

	BeforeSuite(func() {
		utils.Start2RedisMock(redisAddrs)
		utils.Start2Proxy(redisAddrs, proxyAddrs, proxy.Redis)
		utils.StartEtcdMock(etcdAddrs, dataDir)

		client, _ = clientv3.New(clientv3.Config{Endpoints: etcdAddrs, Username: "root", Password: "root"})

		serversStr, err := json.Marshal(configuration.RedisConfig.Servers)
		if err != nil {
			fmt.Println(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, _ = client.Put(ctx, servers, string(serversStr), clientv3.WithPrevKV())
		_, _ = client.Put(ctx, algorithm, configuration.RouteAlgorithm, clientv3.WithPrevKV())
		_, _ = client.Put(ctx, activekey, "ds1", clientv3.WithPrevKV())
		cancel()
		if configuration.Chaos != nil {
			configuration.Chaos.Active = false
		}
		redisClient = redis.NewDevsporeClient(configuration)
	})

	AfterSuite(func() {
		// defer太多，放到一个defer就行
		_ = redisClient.Close()
		_ = client.Close()
		utils.StopEtcdMock(dataDir)
		utils.Stop2Proxy()
		utils.Stop2RedisMock()
	})

	AfterEach(func() {
		utils.Proxys[0].DeleteAllRule()
	})

	It("Get", func() {
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
		Expect(res1.Val()).To(Equal(utils.RedisMocks[0].Addr))
	})

	It("GetDelay", func() {
		_ = utils.Proxys[0].AddDelay("delay", 1500, 0, "", "")
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
		Expect(res1.Val()).To(Equal(utils.RedisMocks[0].Addr))
	})

	It("GetJitter", func() {
		_ = utils.Proxys[0].AddJitter("jitter", 1500, 0, "", "")
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
		Expect(res1.Val()).To(Equal(utils.RedisMocks[0].Addr))
	})

	It("GetDrop", func() {
		_ = utils.Proxys[0].AddDrop("drop", 0, "", "")
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
	})

	It("GetReturnEmpty", func() {
		_ = utils.Proxys[0].AddReturnEmpty("returnEmpty", 0, "", "")
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
	})

	It("GetReturnErr", func() {
		_ = utils.Proxys[0].AddReturnErr("returnErr", proxyredis.UnknownError, 0, "", "")
		ctx := context.Background()
		res1 := redisClient.Get(ctx, "key")
		Expect(res1.Err()).NotTo(HaveOccurred())
	})

	It("MockChang", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = client.Put(ctx, activekey, "ds2", clientv3.WithPrevKV())
		cancel()
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(time.Second)
	})

	It("GetDrop-MockChang", func() {
		_ = utils.Proxys[0].AddDrop("drop", 0, "", "")
		ctx := context.Background()
		res2 := redisClient.Get(ctx, "key")
		Expect(res2.Err()).NotTo(HaveOccurred())
		Expect(res2.Val()).To(Equal(utils.RedisMocks[1].Addr))
	})

	It("GetReturnEmpty-MockChang", func() {
		_ = utils.Proxys[0].AddReturnEmpty("returnEmpty", 0, "", "")
		ctx := context.Background()
		res2 := redisClient.Get(ctx, "key")
		Expect(res2.Err()).NotTo(HaveOccurred())
		Expect(res2.Val()).To(Equal(utils.RedisMocks[1].Addr))
	})

	It("GetReturnErr-MockChang", func() {
		_ = utils.Proxys[0].AddReturnErr("returnErr", proxyredis.UnknownError, 0, "", "")
		ctx := context.Background()
		res2 := redisClient.Get(ctx, "key")
		Expect(res2.Err()).NotTo(HaveOccurred())
		Expect(res2.Val()).To(Equal(utils.RedisMocks[1].Addr))
	})
})
