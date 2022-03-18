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

package strategy

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type SingleReadWriteStrategy struct {
	abstractStrategy
}

func newSingleReadWriteStrategy(configuration *config.Configuration) *SingleReadWriteStrategy {
	return &SingleReadWriteStrategy{newAbstractStrategy(configuration)}
}

func (s *SingleReadWriteStrategy) RouteClient(opType commandType) redis.UniversalClient {
	return s.activeClient()
}

func (s *SingleReadWriteStrategy) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return s.activeClient().Watch(ctx, fn, keys...)
}
