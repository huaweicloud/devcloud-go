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
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

type SubscribeData struct {
	err  error
	data interface{}
}

func Subscribe(ctx context.Context, s RedigoStrategyMode, duration time.Duration, channel string) (interface{}, error) {
	conn := s.RouteClient(strategy.CommandTypeRead).Get()
	pubsubConn := redis.PubSubConn{Conn: conn}
	err := pubsubConn.Subscribe(redis.Args{}.AddFlat(channel)...)
	defer pubsubConn.Unsubscribe()
	if err != nil {
		log.Println("Subscribe fail;", err.Error())
		return nil, err
	}

	done := make(chan *SubscribeData, 1)
	go func() {
		defer pubsubConn.Close()
		for {
			switch v := pubsubConn.Receive().(type) {
			case error:
				done <- &SubscribeData{err: v}
			case redis.Message:
				done <- &SubscribeData{data: v,
					err: nil}
			case redis.Subscription:
				if v.Count == 0 {
					done <- &SubscribeData{data: nil,
						err: nil}
				}
			default:
				done <- &SubscribeData{data: v,
					err: nil}
			}
		}
	}()

	tick := time.NewTicker(duration)
	for {
		defer tick.Stop()
		select {
		case result := <-done:
			if result == nil {
				return nil, fmt.Errorf("subscribe failed")
			}
			return result.data, result.err
		case <-ctx.Done():
			if err := pubsubConn.Unsubscribe(); err != nil {
				return nil, fmt.Errorf("subscribe context done;unsubscribe fail %s", err.Error())
			}
			return nil, fmt.Errorf("subscribe context done")
		case <-tick.C:
			return nil, fmt.Errorf("subscribe timeout")
		}
	}
}

// need start and close
type SubcribeTool struct {
	callMap    map[string]SubscribeCallback
	pubSubConn *redis.PubSubConn
	stop       chan struct{}
}

type SubscribeCallback func(channel, message string)

func CreateSubcribeTool(r RedigoStrategyMode) *SubcribeTool {
	conn := r.RouteClient(strategy.CommandTypeRead).Get()
	pubsubConn := redis.PubSubConn{Conn: conn}
	tool := &SubcribeTool{
		callMap:    make(map[string]SubscribeCallback),
		pubSubConn: &pubsubConn,
	}
	tool.start()
	return tool
}

func (s *SubcribeTool) start() {
	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				switch res := s.pubSubConn.ReceiveWithTimeout(time.Duration(0)).(type) {
				case redis.Message:
					if call, ok := s.callMap[res.Channel]; ok {
						call(res.Channel, string(res.Data))
					}
				case redis.Subscription:
					log.Printf("SubcribeTool Subscription channel:%s kind:%s count:%d\n", res.Channel, res.Kind, res.Count)
				case error:
					log.Printf("ERROR: SubcribeTool receive error %s", res.Error())
					return
				}
			}

		}
	}()
}

func (s *SubcribeTool) Close() error {
	s.stop <- struct{}{}
	return s.pubSubConn.Close()
}

// add subscribe item
func (s *SubcribeTool) Subscribe(call SubscribeCallback, channel ...string) {
	err := s.pubSubConn.Subscribe(redis.Args{}.AddFlat(channel)...)
	if err != nil {
		log.Println("redis Subscribe error.")
	}

	for _, v := range channel {
		s.callMap[v] = call
	}
}

// delete subscribe item
func (s *SubcribeTool) UnSubscribe(channel ...string) {
	err := s.pubSubConn.Subscribe(redis.Args{}.AddFlat(channel)...)
	if err != nil {
		log.Println("redis Subscribe error.")
	}

	for _, v := range channel {
		delete(s.callMap, v)
	}
}
