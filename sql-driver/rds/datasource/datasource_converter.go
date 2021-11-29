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
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

// GetDataSourceMap by integration cluster configuration
func getDataSourceMap(integrationClusterConfiguration *config.IntegrationClusterConfiguration) map[string]*NodeDataSource {
	dataSources := integrationClusterConfiguration.GetDataSource()
	nodeConfigurationMap := integrationClusterConfiguration.GetRouterConfig().Nodes

	nodes := make(map[string]*NodeDataSource)
	for nodeName, nodeConfiguration := range nodeConfigurationMap {
		nodeDataSource := createNodeDataSource(dataSources, nodeName, nodeConfiguration)
		nodes[nodeName] = nodeDataSource
	}
	return nodes
}

func createNodeDataSource(dataSourceConfigurationMap map[string]*config.DataSourceConfiguration,
	nodeName string, config *config.NodeConfiguration) *NodeDataSource {
	var region string
	if masterDataSource, ok := dataSourceConfigurationMap[config.Master]; ok {
		region = masterDataSource.Region
	}
	nodeDataSource := NewNodeDataSource(nodeName, config.LoadBalance, config.Master, config.Slaves, dataSourceConfigurationMap)
	nodeDataSource.Region = region
	return nodeDataSource
}
