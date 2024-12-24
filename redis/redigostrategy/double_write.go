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
	"log"

	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/file"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

// local-read-async-double-write
type DoubleWriteRedigoStrategy struct {
	abstractRedigoStrategy
	jobChan chan job
}

func newDoubleWriteStrategy(configuration *config.Configuration) *DoubleWriteRedigoStrategy {
	doubleWriteStrategy := &DoubleWriteRedigoStrategy{
		abstractRedigoStrategy: newAbstractStrategy(configuration),
		jobChan:                make(chan job, 0),
	}
	if configuration.RedisConfig.AsyncRemotePoolConfiguration == nil {
		log.Fatalln("asyncRemotePool is required")
	}
	file.MkDirs(configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir)
	doubleWriteStrategy.createThreadPoolExecutor(configuration.RedisConfig.AsyncRemotePoolConfiguration)
	go doubleWriteStrategy.asyncDoubleWrite()
	return doubleWriteStrategy
}

func (d *DoubleWriteRedigoStrategy) RouteClient(opType strategy.CommandType) *RedigoUniversalClient {
	return d.nearestClient()
}

func (d *DoubleWriteRedigoStrategy) Watch(ctx context.Context, keys ...string) error {
	go func() {
		_ = Watch(d.remoteClient().Get(), keys...)
	}()
	err := Watch(d.nearestClient().Get(), keys...)
	return err
}

func (d *DoubleWriteRedigoStrategy) Do(opType strategy.CommandType, commandName string, args ...interface{}) (reply interface{}, err error) {
	reply, err = d.RouteClient(opType).Do(commandName, args...)
	if strategy.IsWriteCommand(commandName, args) {
		d.executeAsyncNotPersist(context.TODO(), RedigoCommandArgs{
			CommandName: commandName,
			Args:        args,
		})
	}
	return
}

func (d *DoubleWriteRedigoStrategy) Pipeline(transactions bool, cmds interface{}) ([]interface{}, error) {
	reply, err := d.RouteClient(strategy.CommandTypeMulti).Pipeline(transactions, cmds)
	d.executePipelineAsyncNotPersist(context.TODO(), transactions, cmds)
	return reply, err
}

type JobType int32

const (
	JobTypeDo JobType = iota
	JobTypePipeline
)

type job struct {
	JobType
	RedigoCommandArgs
	ctx          context.Context
	cmds         interface{}
	transactions bool
}

// asyncDoubleWrite Memory double-write
func (d *DoubleWriteRedigoStrategy) asyncDoubleWrite() {
	for jobs := range d.jobChan {
		switch jobs.JobType {
		case JobTypeDo:
			for i := 0; i < d.Configuration.RedisConfig.AsyncRemoteWrite.RetryTimes; i++ {
				if _, err := d.remoteClient().DoContext(jobs.ctx, jobs.CommandName, jobs.Args...); err == nil {
					break
				} else {
					log.Printf("asyncDoubleWrite Do fail %s %v,err is %s,", jobs.CommandName, jobs.Args, err.Error())
				}
			}
		case JobTypePipeline:
			_, err := d.remoteClient().Pipeline(jobs.transactions, jobs.cmds)
			if err != nil {
				log.Printf("asyncDoubleWrite Pipeline fail %v,err is %s,", jobs.cmds, err.Error())
			}
		default:
			log.Printf("asyncDoubleWrite not support type")
		}
	}
}

// createThreadPoolExecutor Memory double-write buffer creation
func (d *DoubleWriteRedigoStrategy) createThreadPoolExecutor(configuration *config.AsyncRemotePoolConfiguration) {
	if !configuration.Persist {
		d.jobChan = make(chan job, configuration.TaskQueueSize)
	}
}

// double-write command writing
func (d *DoubleWriteRedigoStrategy) executeAsyncNotPersist(ctx context.Context, args RedigoCommandArgs) {
	d.jobChan <- job{ctx: ctx, RedigoCommandArgs: args, JobType: JobTypeDo}
}

// double-write pipeline writing
func (d *DoubleWriteRedigoStrategy) executePipelineAsyncNotPersist(ctx context.Context, transactions bool, cmds interface{}) {
	d.jobChan <- job{ctx: ctx, cmds: cmds, transactions: transactions, JobType: JobTypePipeline}
}
