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

// Package proxyredis Redis TCP-based fault injection
package proxyredis

import (
	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/mock/proxy"
)

type ProxyRedis struct {
	*proxy.Proxy
}

func NewProxy(server, addr string) *ProxyRedis {
	proxyRedis := proxy.NewProxy(server, addr, proxy.Redis)
	return &ProxyRedis{
		Proxy: proxyRedis,
	}
}

// redis chaos error
var (
	Nil            = redis.Nil
	ServerShutdown = proxy.RedisError("Server shutdown in progress")
	UnknownError   = proxy.RedisError("UnknownError")

	CommandUKErr  = proxy.RedisError("ERR unknown command")
	CommandArgErr = proxy.RedisError("ERR wrong number of arguments for command")
)
