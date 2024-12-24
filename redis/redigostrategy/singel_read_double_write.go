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

type SingelReadDoubleWriteStrategy struct {
	DoubleWriteRedigoStrategy
}

func newSingelReadDoubleWriteStrategy(configuration *config.Configuration) *SingelReadDoubleWriteStrategy {
	doubleWriteStrategy := &SingelReadDoubleWriteStrategy{
		DoubleWriteRedigoStrategy: DoubleWriteRedigoStrategy{
			abstractRedigoStrategy: newAbstractStrategy(configuration),
			jobChan:                make(chan job, 0),
		},
	}
	if configuration.RedisConfig.AsyncRemotePoolConfiguration == nil {
		log.Fatalln("asyncRemotePool is required")
	}
	file.MkDirs(configuration.RedisConfig.AsyncRemotePoolConfiguration.PersistDir)
	doubleWriteStrategy.createThreadPoolExecutor(configuration.RedisConfig.AsyncRemotePoolConfiguration)

	go doubleWriteStrategy.asyncDoubleWrite()

	return doubleWriteStrategy
}

func (s *SingelReadDoubleWriteStrategy) RouteClient(opType strategy.CommandType) *RedigoUniversalClient {
	return s.activeClient()
}

func (d *SingelReadDoubleWriteStrategy) Watch(ctx context.Context, keys ...string) error {
	go func() {
		_ = Watch(d.noActiveClient().Get(), keys...)
	}()
	err := Watch(d.activeClient().Get(), keys...)
	return err
}

func (d *SingelReadDoubleWriteStrategy) Do(opType strategy.CommandType, commandName string, args ...interface{}) (reply interface{}, err error) {
	reply, err = d.RouteClient(opType).Do(commandName, args...)
	if strategy.IsWriteCommand(commandName, args) {
		d.executeAsyncNotPersist(context.TODO(), RedigoCommandArgs{
			CommandName: commandName,
			Args:        args,
		})
	}
	return
}

func (s *SingelReadDoubleWriteStrategy) Pipeline(transactions bool, cmds interface{}) ([]interface{}, error) {
	reply, err := s.RouteClient(strategy.CommandTypeMulti).Pipeline(transactions, cmds)
	s.executePipelineAsyncNotPersist(context.TODO(), transactions, cmds)
	return reply, err
}

// asyncDoubleWrite Memory double-write
func (s *SingelReadDoubleWriteStrategy) asyncDoubleWrite() {
	for jobs := range s.jobChan {
		switch jobs.JobType {
		case JobTypeDo:
			for i := 0; i < s.Configuration.RedisConfig.AsyncRemoteWrite.RetryTimes; i++ {
				if _, err := s.noActiveClient().DoContext(jobs.ctx, jobs.CommandName, jobs.Args...); err == nil {
					break
				} else {
					log.Printf("asyncDoubleWrite Do fail %s %v,err is %s,", jobs.CommandName, jobs.Args, err.Error())
				}
			}
		case JobTypePipeline:
			_, err := s.noActiveClient().Pipeline(jobs.transactions, jobs.cmds)
			if err != nil {
				log.Printf("asyncDoubleWrite Pipeline fail %v,err is %s,", jobs.cmds, err.Error())
			}
		default:
			log.Printf("asyncDoubleWrite not support type")
		}
	}
}
