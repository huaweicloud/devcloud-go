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
	"log"

	"github.com/alicebob/miniredis/v2"
)

type RedisMetadata struct {
	User     string // optional
	Password string // optional
	Addr     string
}

var redisMap = map[string]*miniredis.Miniredis{}

func StartMockRedis(metadata RedisMetadata) {
	redis := miniredis.NewMiniRedis()
	if len(metadata.User) > 0 && len(metadata.Password) > 0 {
		redis.RequireUserAuth(metadata.User, metadata.Password)
	}
	if len(metadata.User) == 0 && len(metadata.Password) > 0 {
		redis.RequireAuth(metadata.Password)
	}
	var err error
	if len(metadata.Addr) > 0 {
		err = redis.StartAddr(metadata.Addr)
	} else {
		err = redis.Start()
	}
	if err != nil {
		log.Printf("ERROR: start mock redis failed, %v", err)
		return
	}
	log.Printf("mock redis [%s] started! ", redis.Addr())
	redisMap[redis.Addr()] = redis
	return
}

func GetMockRedisByAddr(addr string) *miniredis.Miniredis {
	if redis, ok := redisMap[addr]; ok {
		return redis
	}
	log.Fatalf("ERROR: no [%s] redis", addr)
	return nil
}

func StopMockRedis() {
	if len(redisMap) == 0 {
		return
	}
	for addr, redis := range redisMap {
		redis.Close()
		log.Printf("mock redis [%s] stop! ", addr)
	}
	redisMap = map[string]*miniredis.Miniredis{}
	return
}
