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

func TestClusterRouteStrategy(t *testing.T) {
	node0 := &datasource.NodeDataSource{Name: "node0"}
	node1 := &datasource.NodeDataSource{Name: "node1"}
	clusterDataSource := &datasource.ClusterDataSource{
		Active: "node0",
		DataSources: map[string]*datasource.NodeDataSource{
			"node0": node0,
			"node1": node1,
		},
	}
	runtimeCtx := &RuntimeContext{DataSource: clusterDataSource}
	targetDataSource := NewClusterRouter("single-read-write").Route(
		true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.Equal(t, node0, targetDataSource)

	// change active node
	clusterDataSource.Active = "node1"
	targetDataSource = NewClusterRouter("single-read-write").Route(
		true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.Equal(t, node1, targetDataSource)
}
