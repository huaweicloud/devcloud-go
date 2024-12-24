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

package strategy

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/file"
)

type SingelReadDoubleWriteStrategy struct {
	DoubleWriteStrategy
}

func newSingelReadDoubleWriteStrategy(configuration *config.Configuration) *SingelReadDoubleWriteStrategy {
	doubleWriteStrategy := &SingelReadDoubleWriteStrategy{
		DoubleWriteStrategy: DoubleWriteStrategy{
			abstractStrategy: newAbstractStrategy(configuration),
			jobChan:          make(chan job, 0),
			fileOperationMap: make(map[string]*file.Operation),
		},
	}
	if configuration.RedisConfig.AsyncRemotePoolConfiguration == nil {
		log.Fatalln("asyncRemotePool is required")
	}
	file.MkDirs(configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir)
	doubleWriteStrategy.createThreadPoolExecutor(configuration.RedisConfig.AsyncRemotePoolConfiguration)

	go doubleWriteStrategy.asyncDoubleWrite()
	// add hook for double write
	doubleWriteStrategy.activeClient().AddHook(doubleWriteStrategy)
	doubleWriteStrategy.noActiveClient().AddHook(doubleWriteStrategy)

	return doubleWriteStrategy
}

func (d *SingelReadDoubleWriteStrategy) RouteClient(opType CommandType) redis.UniversalClient {
	return d.activeClient()
}

func (d *SingelReadDoubleWriteStrategy) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	go func() {
		_ = d.noActiveClient().Watch(ctx, fn, keys...)
	}()
	return d.activeClient().Watch(ctx, fn, keys...)
}

// singel read double-write
func (d *SingelReadDoubleWriteStrategy) asyncDoubleWrite() {
	for jobs := range d.jobChan {
		for i := 0; i < d.Configuration.RedisConfig.AsyncRemoteWrite.RetryTimes; i++ {
			if c := d.noActiveClient().Do(jobs.ctx, jobs.args...); c.Err() == nil {
				break
			} else {
				log.Println(jobs.args, c.Err())
			}
		}
	}
}

func (d *SingelReadDoubleWriteStrategy) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (d *SingelReadDoubleWriteStrategy) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if IsWriteCommand(cmd.Name(), cmd.Args()) {
		d.executeAsyncNotPersist(ctx, cmd.Args())
	}
	return nil
}

func (d *SingelReadDoubleWriteStrategy) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (d *SingelReadDoubleWriteStrategy) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		d.noActiveClient().Do(ctx, cmd.Args()...)
	}
	return nil
}
