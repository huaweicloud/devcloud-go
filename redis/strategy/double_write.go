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
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/file"
)

type DoubleWriteStrategy struct {
	abstractStrategy
	asyncRemoteWriteExecutor []string
	fileOperationMap         map[string]*file.Operation
	jobChan                  chan job
}

func newDoubleWriteStrategy(configuration *config.Configuration) *DoubleWriteStrategy {
	doubleWriteStrategy := &DoubleWriteStrategy{
		abstractStrategy: newAbstractStrategy(configuration),
		jobChan:          make(chan job, 0),
		fileOperationMap: make(map[string]*file.Operation),
	}
	if configuration.RedisConfig.AsyncRemotePoolConfiguration == nil {
		log.Fatalln("asyncRemotePool is required")
	}
	file.MkDirs(configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir)
	doubleWriteStrategy.createThreadPoolExecutor(configuration.RedisConfig.AsyncRemotePoolConfiguration)
	if configuration.RedisConfig.AsyncRemotePoolConfiguration.Persist {
		for name, _ := range configuration.RedisConfig.Servers {
			doubleWriteStrategy.fileOperationMap[name] = file.NewFileOperation()
		}
		go doubleWriteStrategy.asyncWrite(configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir, doubleWriteStrategy.ClientPool)
	} else {
		go doubleWriteStrategy.asyncDoubleWrite()
	}
	// add hook for double write
	doubleWriteStrategy.nearestClient().AddHook(doubleWriteStrategy)
	return doubleWriteStrategy
}

func (d *DoubleWriteStrategy) RouteClient(opType commandType) redis.UniversalClient {
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
		if d.Configuration.RedisConfig.AsyncRemotePoolConfiguration.Persist {
			d.executeAsyncPersist(ctx, cmd.Args())
		} else {
			d.executeAsyncNotPersist(ctx, cmd.Args())
		}
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

// asyncDoubleWrite Memory double-write
func (d *DoubleWriteStrategy) asyncDoubleWrite() {
	for jobs := range d.jobChan {
		for i := 0; i < d.Configuration.RedisConfig.AsyncRemoteWrite.RetryTimes; i++ {
			if c := d.remoteClient().Do(jobs.ctx, jobs.args...); c.Err() == nil {
				break
			} else {
				log.Println(jobs.args, c.Err())
			}
		}
	}
}

// asyncWrite File double-write
func (d *DoubleWriteStrategy) asyncWrite(dir string, clients map[string]redis.UniversalClient) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		filenames := file.FileListNeedReplay(dir)
		if len(filenames) == 0 {
			continue
		}
		log.Println("start replay redis file size: ", len(filenames))
		file.BatchReplay(clients, filenames)
	}
}

// createThreadPoolExecutor Memory double-write buffer creation
func (d *DoubleWriteStrategy) createThreadPoolExecutor(configuration *config.AsyncRemotePoolConfiguration) {
	if !configuration.Persist {
		d.jobChan = make(chan job, configuration.TaskQueueSize)
	}
}

// executeAsyncPersist Memory double-write command writing
func (d *DoubleWriteStrategy) executeAsyncPersist(ctx context.Context, args []interface{}) {
	var remotename string
	nearest := d.Configuration.RedisConfig.Nearest
	for name, _ := range d.Configuration.RedisConfig.Servers {
		if name != nearest {
			remotename = name
			break
		}
	}
	item := file.Item{
		args,
	}
	d.fileOperationMap[remotename].WriteFile(d.Configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir+remotename, item)
}

// executeAsyncNotPersist File double-write command writing
func (d *DoubleWriteStrategy) executeAsyncNotPersist(ctx context.Context, args []interface{}) {
	d.jobChan <- job{ctx: ctx, args: args}
}
