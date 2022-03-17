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

// Package bench_test test performance comparison
package bench_test

import (
	"context"
	"strconv"
	"testing"

	goredis "github.com/go-redis/redis/v8"

	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/redis"
)

var (
	addr1          = "127.0.0.1:6379"
	addr2          = "127.0.0.1:6380"
	testKey        = "test_key"
	testValue      = "test_value"
	redisMock1     = mock.RedisMock{Addr: addr1}
	redisMock2     = mock.RedisMock{Addr: addr2}
	ctx            = context.Background()
	devsporeClient *redis.DevsporeClient
	client         *goredis.Client
)

func startService() {
	redisMock1.StartMockRedis()
	redisMock2.StartMockRedis()
	devsporeClient = redis.NewDevsporeClientWithYaml("../resources/config_for_read_write_separate.yaml")
	client = goredis.NewClient(&goredis.Options{Addr: addr1})
}

func closeService() {
	devsporeClient.Close()
	client.Close()
	redisMock1.StopMockRedis()
	redisMock2.StopMockRedis()
}

func BenchmarkDevsporeClient(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		devsporeClient.Set(ctx, testKey+strconv.Itoa(i), testValue+strconv.Itoa(i), 0)
	}
}

func BenchmarkClient(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(ctx, testKey+strconv.Itoa(i), testValue+strconv.Itoa(i), 0)
	}
}

func TestMain(m *testing.M) {
	startService()
	m.Run()
	closeService()

}
