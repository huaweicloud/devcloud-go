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

package datasource

import (
	"testing"
)

func TestRoundRobinLoadBalanceAlgorithm(t *testing.T) {
	loadBalanceAlgorithm := AlgorithmLoader("ROUND_ROBIN")
	slaves := []*ActualDataSource{
		NewActualDataSource("a", nil),
		NewActualDataSource("b", nil),
		NewActualDataSource("c", nil),
		NewActualDataSource("d", nil),
	}
	var (
		requestId1 int64 = 1
		requestId2 int64 = 2
	)
	for i := 0; i < 10; i++ {
		slave := loadBalanceAlgorithm.GetActualDataSource(requestId1, slaves)
		var wantedSlave *ActualDataSource
		switch {
		case i%len(slaves) == 0:
			wantedSlave = slaves[0]
		case i%len(slaves) == 1:
			wantedSlave = slaves[1]
		case i%len(slaves) == 2:
			wantedSlave = slaves[2]
		case i%len(slaves) == 3:
			wantedSlave = slaves[3]
		}
		if slave != wantedSlave {
			t.Errorf("wanted slave is %v, actual return slave is %v", wantedSlave, slave)
		}
	}

	// test with another requestId
	slave := loadBalanceAlgorithm.GetActualDataSource(requestId2, slaves)
	if slave != slaves[1] {
		t.Errorf("wanted slave is %v, actual return slave is %v", slaves[1], slave)
	}
}

func TestRandomLoadBalanceAlgorithm(t *testing.T) {
	loadBalanceAlgorithm := AlgorithmLoader("RANDOM")
	slaves := []*ActualDataSource{
		NewActualDataSource("a", nil),
		NewActualDataSource("b", nil),
		NewActualDataSource("c", nil),
		NewActualDataSource("d", nil),
	}
	for i := 0; i < 10; i++ {
		slave := loadBalanceAlgorithm.GetActualDataSource(0, slaves)
		t.Logf("slave name:%v", slave.Name)
	}
}
