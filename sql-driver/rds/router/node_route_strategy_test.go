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

package router

import (
	"testing"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
	"github.com/stretchr/testify/assert"
)

func TestNodeRouteStrategy_QueryWithoutSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:                 "c0",
		MasterDataSource:     datasource.NewActualDataSource("master", nil),
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource}
	targetDataSource := NewNodeRouter().Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.MasterDataSource, targetDataSource)
}

func TestNodeRouteStrategy_QueryWithSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource}
	targetDataSource := NewNodeRouter().Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.NotEqual(t, nodeDataSource.MasterDataSource, targetDataSource)
}

func TestNodeRouteStrategy_QueryWithSlave_ROUND_ROBIN(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("ROUND_ROBIN"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource}
	targetDataSource := NewNodeRouter().Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.SlavesDatasource[0], targetDataSource)

	targetDataSource = NewNodeRouter().Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.SlavesDatasource[1], targetDataSource)
}

func TestNodeRouteStrategy_QueryWithSlave_WithExclusives(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource}
	exclusives := map[datasource.DataSource]bool{nodeDataSource.SlavesDatasource[0]: true}
	targetDataSource := NewNodeRouter().Route(true, runtimeCtx, exclusives)
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.SlavesDatasource[1], targetDataSource)
}

func TestNodeRouteStrategy_InsertWithSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource}
	targetDataSource := NewNodeRouter().Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.MasterDataSource, targetDataSource)
}

func TestNodeRouteStrategy_Transaction_OnlyRead_WithSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource, IsBeginTransaction: true, IsTransactionReadOnly: true}
	targetDataSource := NewNodeRouter().Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.NotEqual(t, nodeDataSource.MasterDataSource, targetDataSource)
}

func TestNodeRouteStrategy_Transaction_OnlyRead_WithoutSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:                 "c0",
		MasterDataSource:     datasource.NewActualDataSource("master", nil),
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource, IsBeginTransaction: true, IsTransactionReadOnly: true}
	targetDataSource := NewNodeRouter().Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.MasterDataSource, targetDataSource)
}

func TestNodeRouteStrategy_Transaction_WithSlave(t *testing.T) {
	nodeDataSource := &datasource.NodeDataSource{
		Name:             "c0",
		MasterDataSource: datasource.NewActualDataSource("master", nil),
		SlavesDatasource: []*datasource.ActualDataSource{
			datasource.NewActualDataSource("slave0", nil),
			datasource.NewActualDataSource("slave1", nil)},
		LoadBalanceAlgorithm: datasource.AlgorithmLoader("RANDOM"),
	}
	runtimeCtx := &RuntimeContext{DataSource: nodeDataSource, IsBeginTransaction: true}
	targetDataSource := NewNodeRouter().Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, nodeDataSource.MasterDataSource, targetDataSource)
}
