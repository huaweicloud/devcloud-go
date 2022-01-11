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

package redis

import (
	"context"
	"testing"

	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
	"github.com/stretchr/testify/assert"
)

func TestDevsporeClient_ActiveChanges(t *testing.T) {
	var (
		addr1 = "127.0.0.1:16379"
		addr2 = "127.0.0.1:16380"
	)
	mock.StartMockRedis(mock.RedisMetadata{Addr: addr1})
	mock.StartMockRedis(mock.RedisMetadata{Addr: addr2})
	defer mock.StopMockRedis()

	client := NewDevsporeClientWithYaml("./resources/config_for_active_change.yaml")
	assert.Equal(t, 2, len(client.configuration.RedisConfig.Servers))
	client.configuration.RedisConfig.Servers["dc1"].Options.Addr = addr1
	client.configuration.RedisConfig.Servers["dc2"].Options.Addr = addr2
	ctx := context.Background()

	// active server is dc1
	client.configuration.Active = "dc1"
	var (
		tests1Key   = "test_s1_key"
		tests1Value = "test_value"
	)
	client.Set(ctx, tests1Key, tests1Value, 0)
	s1res, _ := mock.GetMockRedisByAddr(addr1).Get(tests1Key)
	s2res, _ := mock.GetMockRedisByAddr(addr2).Get(tests1Key)
	assert.Equal(t, tests1Value, s1res)
	assert.Equal(t, "", s2res)

	// active server is dc2
	client.configuration.Active = "dc2"
	var (
		tests2Key   = "test_s2_key"
		tests2Value = "test_value"
	)
	client.Set(ctx, tests2Key, tests2Value, 0)
	s1res, _ = mock.GetMockRedisByAddr(addr1).Get(tests2Key)
	s2res, _ = mock.GetMockRedisByAddr(addr2).Get(tests2Key)
	assert.Equal(t, "", s1res)
	assert.Equal(t, tests2Value, s2res)
}

func TestDevsporeClient_ReadWriteSeparated(t *testing.T) {
	var (
		addr1 = "127.0.0.1:16379"
		addr2 = "127.0.0.1:16380"
	)
	mock.StartMockRedis(mock.RedisMetadata{Addr: addr1})
	mock.StartMockRedis(mock.RedisMetadata{Addr: addr2})
	defer mock.StopMockRedis()
	client := NewDevsporeClientWithYaml("./resources/config_for_read_write_separate.yaml")
	assert.Equal(t, 2, len(client.configuration.RedisConfig.Servers))
	client.configuration.RedisConfig.Servers["dc1"].Options.Addr = addr1
	client.configuration.RedisConfig.Servers["dc2"].Options.Addr = addr2

	assert.Equal(t, strategy.LocalReadSingleWriteMode, client.configuration.RouteAlgorithm)
	assert.Equal(t, "dc1", client.configuration.RedisConfig.Nearest)
	assert.Equal(t, "dc2", client.configuration.Active)
	ctx := context.Background()

	var (
		testKey     = "test_key"
		testValue   = "test_value"
		tests1Value = "test_s1_value"
		tests2Value = "test_s2_value"
	)
	mock.GetMockRedisByAddr(addr1).Set(testKey, tests1Value)
	mock.GetMockRedisByAddr(addr2).Set(testKey, tests2Value)
	assert.Equal(t, tests1Value, client.Get(ctx, testKey).Val())

	client.Set(ctx, testKey, testValue, 0)
	s1res, _ := mock.GetMockRedisByAddr(addr1).Get(testKey)
	s2res, _ := mock.GetMockRedisByAddr(addr2).Get(testKey)
	assert.Equal(t, tests1Value, s1res)
	assert.Equal(t, testValue, s2res)
}
