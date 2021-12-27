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
	"math"
	"time"

	"github.com/Shopify/sarama"
	"github.com/huaweicloud/devcloud-go/common/util"
)

type Properties struct {
	Addrs           []string
	Async           bool
	OffsetBlockSize int
	BizRetryTimes   int
	LimitPerSecond  int

	// goroutine pool size
	PoolSize     int
	PoolTaskSize int

	SaramaConfig  *sarama.Config
	InitialOffset int64

	CommitSize     int // default partitionCount*OffsetBlockSize
	AutoCommit     bool
	CommitInterval time.Duration
}

// NewProperties return a default dms properties
func NewProperties() *Properties {
	return &Properties{
		Addrs:           []string{},
		Async:           true,
		OffsetBlockSize: 1000,
		BizRetryTimes:   0,
		LimitPerSecond:  math.MaxInt32,
		PoolSize:        20,
		PoolTaskSize:    20000,
		SaramaConfig:    sarama.NewConfig(),
		InitialOffset:   sarama.OffsetOldest,
		CommitSize:      -1,
		AutoCommit:      true,
		CommitInterval:  time.Second,
	}
}

func (p *Properties) validate() error {
	if p.Addrs == nil || len(p.Addrs) == 0 {
		return errors.New("addrs can not be empty")
	}
	if p.SaramaConfig == nil {
		return errors.New("SaramaConfig can not be nil")
	}
	if p.AutoCommit && p.CommitInterval <= 0 {
		return errors.New("auto commit must set interval")
	}
	p.normalize()
	return nil
}

func (p *Properties) normalize() {
	p.OffsetBlockSize = util.GetNearest2Power(p.OffsetBlockSize)
	p.SaramaConfig.Consumer.Offsets.Initial = p.InitialOffset
}

func (p *Properties) Clone() *Properties {
	return &Properties{
		Addrs:           p.Addrs,
		Async:           p.Async,
		OffsetBlockSize: p.OffsetBlockSize,
		BizRetryTimes:   p.BizRetryTimes,
		LimitPerSecond:  p.LimitPerSecond,
		PoolSize:        p.PoolSize,
		PoolTaskSize:    p.PoolTaskSize,
		SaramaConfig:    p.SaramaConfig,
		AutoCommit:      p.AutoCommit,
		InitialOffset:   p.InitialOffset,
		CommitInterval:  p.CommitInterval,
	}
}
