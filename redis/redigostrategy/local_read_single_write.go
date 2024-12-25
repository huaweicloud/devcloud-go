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

type LocalReadSingleWriteRedigoStrategy struct {
	abstractRedigoStrategy
}

func newLocalReadSingleWriteStrategy(configuration *config.Configuration) *LocalReadSingleWriteRedigoStrategy {
	return &LocalReadSingleWriteRedigoStrategy{newAbstractStrategy(configuration)}
}

func (l *LocalReadSingleWriteRedigoStrategy) RouteClient(opType strategy.CommandType) *RedigoUniversalClient {
	if opType == strategy.CommandTypeRead {
		return l.nearestClient()
	}
	return l.activeClient()
}

func (l *LocalReadSingleWriteRedigoStrategy) Watch(ctx context.Context, keys ...string) error {
	return Watch(l.activeClient().Get(), keys...)
}

func (l *LocalReadSingleWriteRedigoStrategy) Do(opType strategy.CommandType, commandName string, args ...interface{}) (reply interface{}, err error) {
	return l.RouteClient(opType).Do(commandName, args...)
}

func (l *LocalReadSingleWriteRedigoStrategy) Pipeline(transactions bool, cmds interface{}) ([]interface{}, error) {
	return l.RouteClient(strategy.CommandTypeWrite).Pipeline(transactions, cmds)
}
