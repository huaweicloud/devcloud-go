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

package util

import "container/list"

type entry struct {
	key   int64
	value int
}

type LRUCache struct {
	cache map[int64]*list.Element
	list  *list.List
	cap   int
}

// NewLRUCache create a lruCache with input capacity
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{map[int64]*list.Element{}, list.New(), capacity}
}

// Get if key not exist, return -1
func (c *LRUCache) Get(key int64) int {
	e := c.cache[key]
	if e == nil {
		return -1
	}
	c.list.MoveToFront(e)
	return e.Value.(entry).value
}

// Put store key-val in lruCache
func (c *LRUCache) Put(key int64, value int) {
	if e := c.cache[key]; e != nil {
		e.Value = entry{key, value}
		c.list.MoveToFront(e)
		return
	}
	c.cache[key] = c.list.PushFront(entry{key, value})
	if len(c.cache) > c.cap {
		delete(c.cache, c.list.Remove(c.list.Back()).(entry).key)
	}
}
