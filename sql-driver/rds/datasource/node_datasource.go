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

// NodeDataSource have read/write separation and failure retry capabilities.
type NodeDataSource struct {
	Region               string
	Name                 string
	MasterDataSource     *ActualDataSource
	SlavesDatasource     []*ActualDataSource
	LoadBalanceAlgorithm LoadBalanceAlgorithm
}

//  @param nodeName: node datasource name
//  @param loadBalanceType: random or round-robin
//  @param masterName: master datasource name
//  @param slavesName: salves datasources name
func NewNodeDataSource(nodeName, loadBalanceType, masterName string, slavesName []string,
	dataSourceConfigurationMap map[string]*config.DataSourceConfiguration) *NodeDataSource {
	var actualSlavesDatasource []*ActualDataSource
	for _, slaveName := range slavesName {
		if slaveDataSource, ok := dataSourceConfigurationMap[slaveName]; ok {
			actualSlave := NewActualDataSource(slaveName, slaveDataSource)
			actualSlavesDatasource = append(actualSlavesDatasource, actualSlave)
		}
	}
	var actualMasterDataSource *ActualDataSource
	if masterDataSource, ok := dataSourceConfigurationMap[masterName]; ok {
		actualMasterDataSource = NewActualDataSource(masterName, masterDataSource)
	}
	return &NodeDataSource{
		Name:                 nodeName,
		MasterDataSource:     actualMasterDataSource,
		SlavesDatasource:     actualSlavesDatasource,
		LoadBalanceAlgorithm: AlgorithmLoader(loadBalanceType),
	}
}

func (ns *NodeDataSource) extend() {
}
