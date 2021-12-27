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

package dms

import (
	"bytes"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/huaweicloud/devcloud-go/common/util"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/time/rate"
)

const version = "version"

type DmsHandler struct {
	ctx           context.Context
	groupId       string
	topics        []string
	bizGroup      string
	pool          *ants.Pool
	limiter       *rate.Limiter
	properties    *Properties
	offsetPersist OffsetPersist

	// inner variables
	limiterCtx       context.Context
	version          byte
	offsetManagerMap sync.Map
	topic2methodMap  map[string]BizHandler
	consumer         sarama.ConsumerGroup
	client           sarama.Client
	msgCount         int64
	commitSize       int
	ticker           *time.Ticker
	closing          chan struct{}
}

func NewDmsHandler(
	ctx context.Context,
	methodInfo MethodInfo,
	pool *ants.Pool,
	limiter *rate.Limiter,
	properties *Properties,
	offsetPersist OffsetPersist) (*DmsHandler, error) {

	// validate properties
	if err := properties.validate(); err != nil {
		return nil, err
	}
	handler := &DmsHandler{
		ctx:              ctx,
		groupId:          methodInfo.GroupId,
		topics:           methodInfo.Topics,
		bizGroup:         methodInfo.BizGroup,
		pool:             pool,
		limiter:          limiter,
		properties:       properties,
		offsetPersist:    offsetPersist,
		version:          0,
		limiterCtx:       context.Background(),
		offsetManagerMap: sync.Map{},
		topic2methodMap:  make(map[string]BizHandler),
		commitSize:       properties.CommitSize,
	}

	properties.SaramaConfig.Consumer.Interceptors = append(properties.SaramaConfig.Consumer.Interceptors, handler)
	consumer, err := sarama.NewConsumerGroup(properties.Addrs, handler.groupId, properties.SaramaConfig)
	if err != nil {
		return nil, err
	}
	handler.consumer = consumer

	if properties.AutoCommit {
		handler.ticker = time.NewTicker(handler.properties.CommitInterval)
		handler.closing = make(chan struct{})
	}

	client, err := sarama.NewClient(properties.Addrs, properties.SaramaConfig)
	if err != nil {
		log.Printf("WARNING: sarama NewClient failed, %v", err)
		return nil, err
	}
	handler.client = client
	return handler, nil
}

func (h *DmsHandler) AddTopicToMethod(method MethodInfo) {
	for _, topic := range method.Topics {
		h.topic2methodMap[topic] = method.Method
	}
}

func (h *DmsHandler) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for {
			err := h.consumer.Consume(h.ctx, h.topics, h)
			if err != nil {
				log.Printf("ERROR: consume err, %v", err)
				time.Sleep(time.Second)
				continue
			}
			if h.ctx.Err() != nil {
				wg.Done()
				log.Printf("INFO: handler [%s]-[%s] cancel, then return!", h.groupId, h.bizGroup)
				return
			}
		}
	}()
	log.Printf("INFO: sarama [%s]-[%s] consumer up and running!...", h.groupId, h.bizGroup)
}

// Setup implements ConsumerGroupHandler interface, when set disable auto commit, Setup will obtain a valid offset of
// groupId-topic-partition from kafka broker, db and the broker's beginning offset.
func (h *DmsHandler) Setup(sess sarama.ConsumerGroupSession) error {
	if h.properties.Async && h.pool.IsClosed() {
		h.pool.Reboot()
	}
	saramaOffsetManager, err := sarama.NewOffsetManagerFromClient(h.groupId, h.client)
	if err != nil {
		log.Printf("WARNING: sarama NewOffsetManagerFromClient failed, %v", err)
		return err
	}

	partitionCount := 0
	for topic, partitions := range sess.Claims() {
		for _, partition := range partitions {
			h.initialOffsetManager(sess, saramaOffsetManager, topic, partition)
			partitionCount++
		}
	}
	if !h.properties.AutoCommit && h.properties.CommitSize <= 0 {
		h.commitSize = partitionCount * h.properties.OffsetBlockSize
	}
	if err = saramaOffsetManager.Close(); err != nil {
		log.Printf("WARNING: close sarama offset manager failed, %v", err)
	}
	if h.properties.AutoCommit {
		go h.loopCommit(sess)
	}
	return nil
}

func (h *DmsHandler) initialOffsetManager(sess sarama.ConsumerGroupSession, saramaOffsetManager sarama.OffsetManager, topic string, partition int32) {
	var offset int64
	if partitionManager, err := saramaOffsetManager.ManagePartition(topic, partition); err == nil {
		offset, _ = partitionManager.NextOffset()
		if offset < 0 {
			if initialOffset, err := h.client.GetOffset(topic, partition, h.properties.InitialOffset); err == nil {
				offset = util.MaxInt64(initialOffset, offset)
			}
		}
		partitionManager.AsyncClose()
	}
	dbOffset, err := h.offsetPersist.Find(h.groupId, topic, int(partition))
	if err == nil {
		offset = util.MaxInt64(dbOffset, offset)
		sess.MarkOffset(topic, partition, offset, "")
	}
	offsetManagerKey := topic + "-" + string(partition)
	h.offsetManagerMap.Store(offsetManagerKey, NewOffsetManager(offset, h.properties.OffsetBlockSize, int(partition), h.groupId, topic, h.version))
}

func (h *DmsHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.commitOnCleanUp(sess)
	if h.properties.AutoCommit {
		h.closing <- struct{}{}
	}
	h.version++
	log.Printf("INFO: message version increase to %v, processed %d messages", h.version, h.msgCount)
	return nil
}

func (h *DmsHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			h.limit()
			if h.ctx.Err() != nil {
				return h.ctx.Err()
			}
			h.handle(sess, msg)
		case <-sess.Context().Done():
			if h.properties.Async {
				h.pool.Release()
			}
			return nil
		}
	}
}

func (h *DmsHandler) commit(sess sarama.ConsumerGroupSession) {
	h.offsetManagerMap.Range(func(key, value interface{}) bool {
		offsetManager, ok := value.(*OffsetManager)
		if ok && offsetManager.blocks.Empty() {
			return true
		}
		offsetManager.lock.RLock()
		minKey, _ := offsetManager.blocks.Min()
		offsetManager.lock.RUnlock()
		sess.MarkOffset(offsetManager.topic, int32(offsetManager.partition), minKey.(int64)+offsetManager.startOffset, "")
		return true
	})
	sess.Commit()
}

func (h *DmsHandler) commitOnCleanUp(sess sarama.ConsumerGroupSession) {
	h.offsetManagerMap.Range(func(key, value interface{}) bool {
		offsetManager, ok := value.(*OffsetManager)
		if ok && offsetManager.blocks.Empty() {
			return true
		}
		offset := offsetManager.handleOffsetOnCleanUp(h.offsetPersist)
		sess.MarkOffset(offsetManager.topic, int32(offsetManager.partition), offset, "")
		h.offsetManagerMap.Delete(key)
		return true
	})
	sess.Commit()
}

func (h *DmsHandler) loopCommit(sess sarama.ConsumerGroupSession) {
	for {
		select {
		case <-h.ticker.C:
			h.commit(sess)
		case <-h.closing:
			return
		}
	}
}

func (h *DmsHandler) handle(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	if h.properties.Async && !h.pool.IsClosed() {
		err := h.pool.Submit(func() {
			h.innerHandler(sess, msg)
		})
		if err != nil {
			log.Printf("INFO: handle '%+v' failed, %v", msg, err)
		}
	} else {
		h.innerHandler(sess, msg)
	}
}

func (h *DmsHandler) innerHandler(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	if h.properties.Async && h.pool.IsClosed() {
		log.Printf("INFO: pool closed, abandon [%+v] !", msg)
		return
	}
	if !h.checkMsgVersionValid(msg) {
		log.Printf("INFO: invalid version, abandon [%+v] before handle msg!", msg)
		return
	}
	h.handleBiz(msg)
	atomic.AddInt64(&h.msgCount, 1)
	h.handleOffset(msg)
	if !h.properties.AutoCommit && atomic.LoadInt64(&h.msgCount)%int64(h.commitSize) == 0 {
		h.commit(sess)
	}
}

func (h *DmsHandler) handleBiz(msg *sarama.ConsumerMessage) {
	if method, ok := h.topic2methodMap[msg.Topic]; ok {
		for i := 0; i <= h.properties.BizRetryTimes; i++ {
			if err := method(msg); err == nil {
				break
			}
		}
	} else {
		log.Printf("WARNING: wrong topic [%s], abandon %+v", msg.Topic, msg)
		return
	}
}

func (h *DmsHandler) handleOffset(msg *sarama.ConsumerMessage) {
	value, ok := h.offsetManagerMap.Load(msg.Topic + "-" + string(msg.Partition))
	if !ok || value == nil {
		log.Printf("INFO: no offset manager for group/topic/partition: %s/%s/%v, abandon: %v",
			msg.Topic, msg.Topic, msg.Partition, msg)
		return
	}
	offsetManager, ok := value.(*OffsetManager)
	if !ok || (ok && offsetManager.version != h.version) {
		log.Printf("INFO: invalid version %v, abandon: %v", offsetManager.version, msg)
		return
	}
	if offsetManager.markAndCheck(msg.Offset) {
		offsetManager.handleHeadBlock(h.offsetPersist)
	}
}

func (h *DmsHandler) checkMsgVersionValid(msg *sarama.ConsumerMessage) bool {
	for _, header := range msg.Headers {
		if bytes.Equal(header.Key, []byte(version)) && !bytes.Equal(header.Value, []byte{h.version}) {
			return false
		}
	}
	return true
}

func (h *DmsHandler) limit() {
	err := h.limiter.Wait(h.limiterCtx)
	if err != nil {
		log.Printf("WARNING: limiter get token failed, %v", err)
		return
	}
}

func (h *DmsHandler) OnConsume(msg *sarama.ConsumerMessage) {
	msg.Headers = append(msg.Headers, &sarama.RecordHeader{Key: []byte(version), Value: []byte{h.version}})
}

func (h *DmsHandler) Close() error {
	var err error
	if err = h.consumer.Close(); err != nil {
		log.Printf("WARNING: close sarama consumer failed, %v", err)
	}
	if h.properties.Async && !h.pool.IsClosed() {
		h.pool.Release()
	}
	if err = h.client.Close(); err != nil {
		log.Printf("WARNING: close sarama client failed, %v", err)
	}
	if h.properties.AutoCommit {
		close(h.closing)
	}
	return err
}
