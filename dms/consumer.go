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
 * Package dms implements a kafka consumer based on sarama, user can consume messages
 * asynchronous or synchronous with dms, and ensure message not lost.
 */

package dms

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/panjf2000/ants/v2"
	"golang.org/x/time/rate"
)

type Consumer struct {
	ctx                      context.Context
	methods                  []MethodInfo
	propertiesMap            map[string]*Properties
	offsetPersist            OffsetPersist
	wg                       *sync.WaitGroup
	groupIdBizGroupToHandler map[string]*DmsHandler
}

func NewConsumer(ctx context.Context, methods []MethodInfo, propertiesMap map[string]*Properties, offsetPersist OffsetPersist) (*Consumer, error) {
	dmsConsumer := &Consumer{
		ctx:                      ctx,
		methods:                  methods,
		propertiesMap:            propertiesMap,
		offsetPersist:            offsetPersist,
		wg:                       &sync.WaitGroup{},
		groupIdBizGroupToHandler: map[string]*DmsHandler{},
	}
	var (
		ok         bool
		err        error
		pool       *ants.Pool
		properties *Properties
	)
	for _, methodInfo := range methods {
		if properties, ok = dmsConsumer.propertiesMap[methodInfo.GetUniqueKey()]; !ok && properties == nil {
			log.Printf("ERROR: group[%s] bizGroup[%s] do not have properties!", methodInfo.GroupId, methodInfo.BizGroup)
			return nil, errors.New("invalid properties")
		}

		if properties.Async {
			pool, err = ants.NewPool(properties.PoolSize, ants.WithMaxBlockingTasks(properties.PoolTaskSize))
			if err != nil {
				return nil, err
			}
		}
		limiter := rate.NewLimiter(rate.Limit(properties.LimitPerSecond), properties.LimitPerSecond)
		handler, err := NewDmsHandler(ctx, methodInfo, pool, limiter, properties, offsetPersist)
		if err != nil {
			return nil, err
		}
		handler.AddTopicToMethod(methodInfo)
		dmsConsumer.groupIdBizGroupToHandler[methodInfo.GetUniqueKey()] = handler
	}
	return dmsConsumer, nil
}

func (c *Consumer) Consume() {
	for _, handler := range c.groupIdBizGroupToHandler {
		handler.Start(c.wg)
	}
}

func (c *Consumer) Close() {
	c.wg.Wait()
	for _, handler := range c.groupIdBizGroupToHandler {
		if err := handler.Close(); err != nil {
			log.Printf("INFO: close err, %v", err)
		}
	}
	log.Println("INFO: close DMS consumer.")
}
