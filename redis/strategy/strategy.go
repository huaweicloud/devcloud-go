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

// Package strategy defines different route strategy mode.
package strategy

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type StrategyMode interface {
	RouteClient(opType CommandType) redis.UniversalClient
	Close() error
	Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error
}

func NewStrategy(configuration *config.Configuration) StrategyMode {
	switch configuration.RouteAlgorithm {
	case SingleReadWriteMode:
		return newSingleReadWriteStrategy(configuration)
	case LocalReadSingleWriteMode:
		return newLocalReadSingleWriteStrategy(configuration)
	case LocalReadDoubleWriteMode:
		return newLocalReadSingleWriteStrategy(configuration)
		return newDoubleWriteStrategy(configuration)
	case SingleReadDoubleWriteMode:
		return newSingleReadWriteStrategy(configuration)
		return newSingelReadDoubleWriteStrategy(configuration)
	default:
		log.Printf("ERROR: invalid route algorithm:%v", configuration.RouteAlgorithm)
	}
	return nil
}

type CommandType int32

const (
	CommandTypeRead CommandType = iota
	CommandTypeWrite
	CommandTypeMulti
	CommandTypeOther
)

const (
	SingleReadWriteMode       = "single-read-write"
	LocalReadSingleWriteMode  = "local-read-single-write"
	SingleReadDoubleWriteMode = "single-read-async-double-write"
	LocalReadDoubleWriteMode  = "local-read-async-double-write"
)

func IsWriteCommand(funcName string, args []interface{}) bool {
	funcName = strings.ToLower(funcName)
	if _, ok := writeCommandMap[funcName]; ok {
		return true
	}
	switch funcName {
	case "script":
		return contains(args, "flush")
	case "sort":
		return contains(args, "store")
	case "georadius":
		return contains(args, "store")
	}
	return false
}

func contains(args []interface{}, command string) bool {
	for _, arg := range args {
		if reflect.DeepEqual(arg, command) {
			return true
		}
	}
	return false
}

var writeCommandMap = map[string]bool{
	"set":              true,
	"del":              true,
	"expire":           true,
	"expireat":         true,
	"persist":          true,
	"pexpire":          true,
	"pexpireat":        true,
	"rename":           true,
	"renamenx":         true,
	"restore":          true,
	"touch":            true,
	"append":           true,
	"decr":             true,
	"decrby":           true,
	"getset":           true,
	"getex":            true,
	"getdel":           true,
	"incr":             true,
	"incrby":           true,
	"incrbyfloat":      true,
	"mset":             true,
	"msetnx":           true,
	"setex":            true,
	"setnx":            true,
	"setrange":         true,
	"setbit":           true,
	"bitop":            true,
	"hdel":             true,
	"hincrby":          true,
	"hincrbyfloat":     true,
	"hset":             true,
	"hmset":            true,
	"hsetnx":           true,
	"blpop":            true,
	"brpop":            true,
	"blpoplpush":       true,
	"linsert":          true,
	"lpop":             true,
	"lpush":            true,
	"lpushx":           true,
	"lrem":             true,
	"lset":             true,
	"ltrim":            true,
	"rpop":             true,
	"rpoplpush":        true,
	"rpush":            true,
	"rpushx":           true,
	"lmove":            true,
	"sadd":             true,
	"sdiffstore":       true,
	"sinterstore":      true,
	"smove":            true,
	"spop":             true,
	"srem":             true,
	"sunionstore":      true,
	"xadd":             true,
	"xdel":             true,
	"bzpopmax":         true,
	"bzpopmin":         true,
	"zadd":             true,
	"zincrby":          true,
	"zinterstore":      true,
	"zpopmax":          true,
	"zpopmin":          true,
	"zrem":             true,
	"zremrangebyrank":  true,
	"zremrangebyscore": true,
	"zremrangebylex":   true,
	"zunionstore":      true,
	"zdiffstore":       true,
	"pfadd":            true,
	"pfmerge":          true,
	"bgrewriteaof":     true,
	"bgsave":           true,
	"flushall":         true,
	"flushdb":          true,
	"save":             true,
	"eval":             true,
	"evalsha":          true,
	"geoadd":           true,
	"geosearchstore":   true,
}
