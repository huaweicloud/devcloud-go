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
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffsetNode(t *testing.T) {
	offsetNode := NewOffsetNode(10)
	for i := 0; i < 10; i++ {
		assert.Equal(t, true, offsetNode.mark(i))
	}
	for i := 0; i < 2; i++ {
		assert.Equal(t, false, offsetNode.mark(i))
	}
	for i := 10; i < 15; i++ {
		assert.Equal(t, false, offsetNode.mark(i))
	}
	assert.Equal(t, true, offsetNode.isFull())
	assert.Equal(t, 10, int(offsetNode.size()))
}

func TestOffsetManager_markAndCheck(t *testing.T) {
	offsetManager := NewOffsetManager(0, 16, 0, "groupId1", "topic1", 0)
	for i := 0; i < 15; i++ {
		assert.Equal(t, false, offsetManager.markAndCheck(int64(i)))
	}
	assert.Equal(t, true, offsetManager.markAndCheck(15))
	assert.Equal(t, false, offsetManager.markAndCheck(15))
	assert.Equal(t, true, offsetManager.markAndCheck(16))
}

func TestOffsetManager_checkHeadNeedPersist(t *testing.T) {
	offsetManager := NewOffsetManager(0, 16, 0, "groupId1", "topic1", 0)
	for i := 0; i < 15; i++ {
		offsetManager.markAndCheck(int64(i))
	}
	assert.Equal(t, false, offsetManager.checkHeadNeedPersist())
	offsetManager.markAndCheck(15)
	assert.Equal(t, true, offsetManager.checkHeadNeedPersist())
}

func TestOffsetManager_handleHeadBlock(t *testing.T) {
	offsetManager := NewOffsetManager(0, 16, 0, "groupId1", "topic1", 0)
	for i := 0; i < 20; i++ {
		offsetManager.markAndCheck(int64(i))
	}
	assert.Equal(t, 2, offsetManager.blocks.Size())

	offsetManager.handleHeadBlock(&myOffsetPersist{map[string]int64{}})
	assert.Equal(t, 1, offsetManager.blocks.Size())
	blockMinKey, blockMinVal := offsetManager.blocks.Min()
	assert.Equal(t, int64(16), blockMinKey.(int64))
	assert.Equal(t, false, blockMinVal.(*OffsetNode).isFull())
	assert.Equal(t, 4, int(blockMinVal.(*OffsetNode).size()))
}

type myOffsetPersist struct {
	cache map[string]int64
}

func (p *myOffsetPersist) Find(groupId, topic string, partition int) (int64, error) {
	key := groupId + topic + string(rune(partition))
	if offset, ok := p.cache[key]; ok {
		return offset, nil
	}
	return 0, errors.New("not found")
}

func (p *myOffsetPersist) Save(groupId, topic string, partition int, offset int64) error {
	key := groupId + topic + string(int32(partition))
	log.Printf("key:%v, offset:%d", key, offset)
	p.cache[key] = offset
	return nil
}
