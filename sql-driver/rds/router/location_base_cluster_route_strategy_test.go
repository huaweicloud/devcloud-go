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

var (
	node0             = &datasource.NodeDataSource{Name: "node0", Region: "az0"}
	node1             = &datasource.NodeDataSource{Name: "node1", Region: "az1"}
	clusterDataSource = &datasource.ClusterDataSource{
		Active: "node1",
		DataSources: map[string]*datasource.NodeDataSource{
			"node0": node0,
			"node1": node1,
		},
		Region: "az0",
	}
)

func TestLocationBaseClusterRouteStrategy_ReadSQL(t *testing.T) {
	runtimeCtx := &RuntimeContext{DataSource: clusterDataSource}
	targetDataSource := NewClusterRouter("local-read-single-write").Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node0, targetDataSource)

	// change region "az1" which contains node1
	clusterDataSource.Region = "az1"
	targetDataSource = NewClusterRouter("local-read-single-write").Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node1, targetDataSource)

	// change region "az2" which has nothing
	clusterDataSource.Region = "az2"
	targetDataSource = NewClusterRouter("local-read-single-write").Route(true, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node1, targetDataSource)
}

func TestLocationBaseClusterRouteStrategy_InsertSQL(t *testing.T) {
	clusterDataSource.Region = "az0"
	runtimeCtx := &RuntimeContext{DataSource: clusterDataSource}
	targetDataSource := NewClusterRouter("local-read-single-write").Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node1, targetDataSource)

	// change region "az1" which contains node1
	clusterDataSource.Region = "az1"
	targetDataSource = NewClusterRouter("local-read-single-write").Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node1, targetDataSource)

	// change region "az2" which has nothing
	clusterDataSource.Region = "az2"
	targetDataSource = NewClusterRouter("local-read-single-write").Route(false, runtimeCtx, make(map[datasource.DataSource]bool))
	assert.NotNil(t, targetDataSource)
	assert.Equal(t, node1, targetDataSource)
}
