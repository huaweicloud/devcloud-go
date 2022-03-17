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

// Package example provides an example for user how to use dms.
package example

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/huaweicloud/devcloud-go/dms"
)

func main() {
	methods := []dms.MethodInfo{
		{
			GroupId:  "groupId",
			Topics:   []string{"topic1"},
			BizGroup: "mybiz",
			Method:   handle,
		},
		{
			GroupId:  "groupId",
			Topics:   []string{"topic2"},
			BizGroup: "mySecondBiz",
			Method:   secondHandle,
		},
	}
	propsMap := map[string]*dms.Properties{
		methods[0].GetUniqueKey(): asyncProps(),
		methods[1].GetUniqueKey(): syncProps(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	consumer, err := dms.NewConsumer(ctx, methods, propsMap, &myOffsetPersist{})
	if err != nil {
		panic(err)
	}
	consumer.Consume()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("INFO: terminating: via signal")
	}
	close(sigterm)
	cancel()
	consumer.Close()
	if err = db.Close(); err != nil {
		log.Printf("close db failed, %v", err)
	}
	return
}

var addrs = []string{"127.0.0.1:9092"}

func asyncProps() *dms.Properties {
	props := dms.NewProperties()
	props.Addrs = addrs
	props.PoolSize = 20
	props.PoolTaskSize = 1000
	props.OffsetBlockSize = 32
	props.LimitPerSecond = 1000
	// every second commit
	props.AutoCommit = true
	props.CommitInterval = time.Second
	return props
}

func syncProps() *dms.Properties {
	props := dms.NewProperties()
	props.Addrs = addrs
	props.Async = false
	props.OffsetBlockSize = 32
	props.LimitPerSecond = 1000
	// every 10 msg commit
	props.AutoCommit = false
	props.CommitSize = 10
	return props
}

func handle(msg *sarama.ConsumerMessage) error {
	log.Printf("topic:%v, partition:%v, key:%v, value:%v", msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
	return nil
}

func secondHandle(msg *sarama.ConsumerMessage) error {
	log.Printf("topic:%v, partition:%v, key:%v, value:%v", msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
	return nil
}
