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

import (
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"
)

// GetWorkerIDByIp use ip low 10 bit to generate workerID
func GetWorkerIDByIp() int64 {
	ipStr, err := getIp()
	if err != nil {
		return rand.Int63n(1000)
	}
	b := net.ParseIP(ipStr).To4()
	return int64((uint64(b[2])<<8 + uint64(b[3])) % (uint64(1) << 10))
}

func getIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", err
}

const (
	nodeBits  uint8 = 10
	numBits   uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	numMax    int64 = -1 ^ (-1 << numBits)
	timeShift uint8 = nodeBits + numBits
	nodeShift uint8 = numBits
	startTime int64 = 1525705533000
)

type Node struct {
	timestamp int64
	nodeId    int64
	num       int64
	mu        sync.Mutex
}

func NewNode(nodeId int64) (*Node, error) {
	if nodeId < 0 || nodeId > nodeMax {
		return nil, errors.New("node id excess of quantity")
	}
	return &Node{
		nodeId:    nodeId,
		num:       0,
		timestamp: 0,
	}, nil
}

func (w *Node) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.num++
		if w.num > numMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.num = 0
		w.timestamp = now
	}
	return (now-startTime)<<timeShift | (w.nodeId << nodeShift) | (w.num)
}
