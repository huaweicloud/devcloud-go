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

package strategy

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type DoubleWriteStrategy struct {
	abstractStrategy
	jobChan chan job
}

func newDoubleWriteStrategy(configuration *config.Configuration) *DoubleWriteStrategy {
	doubleWriteStrategy := &DoubleWriteStrategy{
		abstractStrategy: newAbstractStrategy(configuration),
		jobChan:          make(chan job, 10000),
	}
	// add hook for double write
	doubleWriteStrategy.nearestClient().AddHook(doubleWriteStrategy)
	go doubleWriteStrategy.asyncDoubleWrite()
	return doubleWriteStrategy
}

func (d *DoubleWriteStrategy) RouteClient(optype commandType) redis.UniversalClient {
	return d.nearestClient()
}

func (d *DoubleWriteStrategy) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	go func() {
		_ = d.remoteClient().Watch(ctx, fn, keys...)
	}()
	return d.nearestClient().Watch(ctx, fn, keys...)
}

func (d *DoubleWriteStrategy) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (d *DoubleWriteStrategy) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if isWriteCommand(cmd.Name(), cmd.Args()) {
		d.jobChan <- job{ctx: ctx, args: cmd.Args()}
	}
	return nil
}

func (d *DoubleWriteStrategy) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (d *DoubleWriteStrategy) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		d.remoteClient().Do(ctx, cmd.Args()...)
	}
	return nil
}

type job struct {
	ctx  context.Context
	args []interface{}
}

func (d *DoubleWriteStrategy) asyncDoubleWrite() {
	for job := range d.jobChan {
		d.remoteClient().Do(job.ctx, job.args...)
	}
}
