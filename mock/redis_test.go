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
	"context"
	"testing"

	goredis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisClusterMock(t *testing.T) {
	redisMock := RedisMock{Addr: "127.0.0.1:16379"}
	redisMock.StartMockRedis()
	cluster := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs: []string{"127.0.0.1:16379"},
	})

	ctx := context.Background()
	cluster.Set(ctx, "key", "val", 0)
	res := cluster.Get(ctx, "key")
	assert.Nil(t, res.Err())
	assert.Equal(t, res.Val(), "val")
	redisMock.StopMockRedis()
}
