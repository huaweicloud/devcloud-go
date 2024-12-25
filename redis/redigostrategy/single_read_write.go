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

package redigostrategy

import (
	"context"

	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

type SingleReadWriteRedigoStrategy struct {
	abstractRedigoStrategy
}

func newSingleReadWriteStrategy(configuration *config.Configuration) *SingleReadWriteRedigoStrategy {
	return &SingleReadWriteRedigoStrategy{newAbstractStrategy(configuration)}
}

func (s *SingleReadWriteRedigoStrategy) RouteClient(opType strategy.CommandType) *RedigoUniversalClient {
	return s.activeClient()
}

func (s *SingleReadWriteRedigoStrategy) Watch(ctx context.Context, keys ...string) error {
	return Watch(s.activeClient().Get(), keys...)
}

func (s *SingleReadWriteRedigoStrategy) Do(opType strategy.CommandType, commandName string, args ...interface{}) (reply interface{}, err error) {
	return s.RouteClient(opType).Do(commandName, args...)
}

func (s *SingleReadWriteRedigoStrategy) Pipeline(transactions bool, cmds interface{}) ([]interface{}, error) {
	return s.RouteClient(strategy.CommandTypeMulti).Pipeline(transactions, cmds)
}
