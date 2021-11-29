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
 * Set load_balance_algorithm.go in datasource package to avoid import cycle.
 */

package datasource

import (
	"math/rand"
	"sync/atomic"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/util"
)

// LoadBalanceAlgorithm Read/write separation load balance algorithm interface.
type LoadBalanceAlgorithm interface {
	// GetActualDataSource Select the data source.
	GetActualDataSource(int64, []*ActualDataSource) *ActualDataSource
}

// RandomLoadBalanceAlgorithm get actualDataSource from salves randomly.
type RandomLoadBalanceAlgorithm struct {
}

// GetActualDataSource randomly
func (ra *RandomLoadBalanceAlgorithm) GetActualDataSource(requestId int64, slavesDataSource []*ActualDataSource) *ActualDataSource {
	return slavesDataSource[rand.Intn(len(slavesDataSource))]
}

// RoundRobinLoadBalanceAlgorithm get actualDataSource from salves in order.
type RoundRobinLoadBalanceAlgorithm struct {
	position      int64
	requestRecord *util.LRUCache
}

// GetActualDataSource by round-robin algorithm
func (ro *RoundRobinLoadBalanceAlgorithm) GetActualDataSource(
	requestId int64, slavesDataSource []*ActualDataSource) *ActualDataSource {
	var curIndex int
	if idx := ro.requestRecord.Get(requestId); idx != -1 {
		curIndex = (idx + 1) % len(slavesDataSource)
	} else {
		curIndex = int(ro.position) % len(slavesDataSource)
		atomic.AddInt64(&ro.position, 1)
	}
	ro.requestRecord.Put(requestId, curIndex)
	return slavesDataSource[curIndex%len(slavesDataSource)]
}

const (
	lruCacheCapacity          = 10
	LoadBalanceTypeRandom     = "RANDOM"
	LoadBalanceTypeRoundRobin = "ROUND_ROBIN"
)

// AlgorithmLoader load LoadBalanceAlgorithm by loadBalanceType.
func AlgorithmLoader(loadBalanceType string) LoadBalanceAlgorithm {
	switch loadBalanceType {
	case LoadBalanceTypeRandom:
		return &RandomLoadBalanceAlgorithm{}
	case LoadBalanceTypeRoundRobin:
		return &RoundRobinLoadBalanceAlgorithm{requestRecord: util.NewLRUCache(lruCacheCapacity)}
	default:
		return &RoundRobinLoadBalanceAlgorithm{requestRecord: util.NewLRUCache(lruCacheCapacity)}
	}
}
