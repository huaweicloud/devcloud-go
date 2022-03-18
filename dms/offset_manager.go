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
	"log"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/RoaringBitmap/roaring"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type OffsetManager struct {
	using         *atomic.Value
	version       byte
	startOffset   int64
	lock          sync.RWMutex
	blockCapacity int
	low           int
	high          int

	groupId   string
	topic     string
	partition int

	blocks *treemap.Map
}

func NewOffsetManager(startOffset int64, blockCapacity, partition int, groupId, topic string, version byte) *OffsetManager {
	if startOffset < 0 {
		startOffset = 0
	}
	using := &atomic.Value{}
	using.Store(false)
	return &OffsetManager{
		using:         using,
		version:       version,
		startOffset:   startOffset,
		lock:          sync.RWMutex{},
		blockCapacity: blockCapacity,
		low:           blockCapacity - 1,
		high:          -blockCapacity,
		groupId:       groupId,
		topic:         topic,
		partition:     partition,
		blocks:        treemap.NewWith(utils.Int64Comparator),
	}
}

func (m *OffsetManager) markAndCheck(absoluteOffset int64) bool {
	relativeOffset := absoluteOffset - m.startOffset
	blockKey := relativeOffset & int64(m.high)
	slotIndex := int(relativeOffset & int64(m.low))
	m.lock.Lock()
	defer m.lock.Unlock()
	var offsetNode *OffsetNode
	if val, found := m.blocks.Get(blockKey); found {
		offsetNode = val.(*OffsetNode)
	} else {
		offsetNode = NewOffsetNode(m.blockCapacity)
		m.blocks.Put(blockKey, offsetNode)
	}
	if !offsetNode.mark(slotIndex) {
		log.Printf("WARNING: group/topic/partition: %s/%s/%v, offsetNode mark %v failed, blockkey:%v",
			m.groupId, m.topic, m.partition, absoluteOffset, blockKey)
		return false
	}
	blockMinKey, blockMinVal := m.blocks.Min()
	if !reflect.DeepEqual(blockMinKey, blockKey) {
		return false
	}
	blockMinNode, ok := blockMinVal.(*OffsetNode)
	if !ok || blockMinNode == nil {
		return false
	}
	return blockMinNode.isFull()
}

func (m *OffsetManager) handleHeadBlock(offsetPersist OffsetPersist) {
	for m.checkHeadNeedPersist() {
		blockMinKey, blockMinVal := m.blocks.Min()
		newOffset := m.startOffset + blockMinKey.(int64) + int64(m.blockCapacity)
		var err error
		for i := 0; i < 3; i++ {
			if err = offsetPersist.Save(m.groupId, m.topic, m.partition, newOffset); err == nil {
				break
			}
		}
		if err != nil {
			log.Printf("WARNING: newOffset '%v' persist fail, %v", newOffset, err)
		}
		m.putNextAndPollFirst(blockMinKey, blockMinVal)
		m.using.Store(false)
	}
}

func (m *OffsetManager) handleOffsetOnCleanUp(offsetPersist OffsetPersist) int64 {
	m.lock.RLock()
	minKey, minNode := m.blocks.Min()
	m.lock.RUnlock()
	offset := minKey.(int64) + m.startOffset + int64(minNode.(*OffsetNode).maxContinuous())
	if err := offsetPersist.Save(m.groupId, m.topic, m.partition, offset); err != nil {
		log.Printf("WARNING: groupId/topic/partition %s/%s/%d persist %d fail on clean up, %v",
			m.groupId, m.topic, m.partition, offset, err)
	}
	return offset
}

func (m *OffsetManager) checkHeadNeedPersist() bool {
	var isNeed bool
	m.lock.RLock()
	if _, val := m.blocks.Min(); val != nil {
		if node, ok := val.(*OffsetNode); ok && node.isFull() {
			if reflect.DeepEqual(false, m.using.Load()) {
				isNeed = true
				m.using.Store(true)
			}
		}
	}
	m.lock.RUnlock()
	return isNeed
}

func (m *OffsetManager) putNextAndPollFirst(blockKey, blockValue interface{}) {
	m.lock.Lock()
	blockMinKey, blockMinVal := m.blocks.Min()
	if reflect.DeepEqual(blockMinKey, blockKey) && reflect.DeepEqual(blockMinVal, blockValue) {
		m.blocks.Remove(blockKey)
	}
	if _, ok := m.blocks.Get(blockKey.(int64) + int64(m.blockCapacity)); !ok {
		m.blocks.Put(blockKey.(int64)+int64(m.blockCapacity), NewOffsetNode(m.blockCapacity))
	}
	m.lock.Unlock()
}

type OffsetNode struct {
	capacity int
	lock     sync.RWMutex
	bitSet   *roaring.Bitmap
}

func NewOffsetNode(capacity int) *OffsetNode {
	return &OffsetNode{
		capacity: capacity,
		lock:     sync.RWMutex{},
		bitSet:   roaring.New(),
	}
}

func (o *OffsetNode) isFull() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return int(o.bitSet.GetCardinality()) >= o.capacity
}

func (o *OffsetNode) mark(offset int) bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	if int(o.bitSet.GetCardinality()) >= o.capacity {
		return false
	}
	if !o.bitSet.ContainsInt(offset) {
		o.bitSet.AddInt(offset)
		return true
	}
	return false
}

func (o *OffsetNode) size() int {
	return int(o.bitSet.GetCardinality())
}

func (o *OffsetNode) maxContinuous() int {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if o.isFull() {
		return o.capacity
	}
	for i := 0; i < o.capacity; i++ {
		if !o.bitSet.ContainsInt(i) {
			return i
		}
	}
	return 0
}
