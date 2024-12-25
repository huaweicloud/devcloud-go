/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025.
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

// Package strategy defines different route redigoStrategy mode.
package redigostrategy

import (
	"context"
	"log"

	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

type RedigoStrategyMode interface {
	RouteClient(opType strategy.CommandType) *RedigoUniversalClient
	Close() error
	Watch(ctx context.Context, keys ...string) error
	Do(opType strategy.CommandType, commandName string, args ...interface{}) (reply interface{}, err error)
	Pipeline(transactions bool, cmds interface{}) ([]interface{}, error)
}

func NewStrategy(configuration *config.Configuration) RedigoStrategyMode {
	switch configuration.RouteAlgorithm {
	case strategy.SingleReadWriteMode:
		return newSingleReadWriteStrategy(configuration)
	case strategy.LocalReadSingleWriteMode:
		return newLocalReadSingleWriteStrategy(configuration)
	case strategy.SingleReadDoubleWriteMode:
		return newSingleReadWriteStrategy(configuration)
	case strategy.LocalReadDoubleWriteMode:
		return newLocalReadSingleWriteStrategy(configuration)
	default:
		log.Printf("ERROR: invalid route algorithm:%v", configuration.RouteAlgorithm)
	}
	return nil
}

func GetWriteReadCommandType(commandName string, args ...interface{}) strategy.CommandType {
	commandType := strategy.CommandTypeRead
	if strategy.IsWriteCommand(commandName, args) {
		commandType = strategy.CommandTypeWrite
	} else {
		commandType = strategy.CommandTypeRead
	}
	return commandType
}
