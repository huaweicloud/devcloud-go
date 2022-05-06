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

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/huaweicloud/devcloud-go/mock/proxy/utils"
	"github.com/huaweicloud/devcloud-go/redis"
)

func TestRedisMas(t *testing.T) {
	var (
		etcdAddrs     = []string{"127.0.0.1:2379"}
		dataDir       = "etcd_data"
		redisAddrs    = []string{"127.0.0.1:16379", "127.0.0.1:16380"}
		configuration = utils.DCRedis(etcdAddrs, redisAddrs)
		servers       = "/mas-monitor/conf/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/servers"
		algorithm     = "/mas-monitor/conf/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/route-algorithm"
		activekey     = "/mas-monitor/status/dcs/services/" + configuration.Props.AppID + "/" + configuration.Props.MonitorID + "/active"
	)

	utils.Start2RedisMock(redisAddrs)
	defer utils.Stop2RedisMock()
	utils.StartEtcdMock(etcdAddrs, dataDir)
	defer utils.StopEtcdMock(dataDir)

	client, _ := clientv3.New(clientv3.Config{Endpoints: etcdAddrs, Username: "XXXX", Password: "XXXX"})
	defer func() {
		_ = client.Close()
	}()
	serversStr, err := json.Marshal(configuration.RedisConfig.Servers)
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, _ = client.Put(ctx, servers, string(serversStr), clientv3.WithPrevKV())
	_, _ = client.Put(ctx, algorithm, configuration.RouteAlgorithm, clientv3.WithPrevKV())
	cancel()

	redisClient := redis.NewDevsporeClient(configuration)
	defer func() {
		_ = redisClient.Close()
	}()
	time.Sleep(time.Second)

	for key := range configuration.RedisConfig.Servers {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		client.Put(ctx, activekey, key, clientv3.WithPrevKV())
		cancel()
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			ctx = context.Background()
			res := redisClient.Get(ctx, "key")
			fmt.Println(res)
			time.Sleep(time.Second)
		}
	}
}
