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
 */

/*
Package datasource defines datasource interface, and implements the interface by ClusterDataSource,
NodeDataSource, and ActualDataSource. The ClusterDataSource have auto change target datasource capabilities,
NodeDatasource have read/write separation and failure retry capabilities, ActualDataSource is an actual
datasource which contains actual dsn.
*/
package datasource

import (
	"log"
	"sync/atomic"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config/loader"
)

// ClusterDataSource which have auto change target datasource capabilities.
type ClusterDataSource struct {
	RouterConfiguration *config.RouterConfiguration
	DataSources         map[string]*NodeDataSource
	Active              string
	switchTimes         int64
	Region              string
}

// NewClusterDataSource create a clusterDataSource by yaml clusterConfiguration and remote etcd clusterConfiguration,
// and listen remote activeKey, when remote activeKey changes, change the clusterDataSource's active node.
func NewClusterDataSource(clusterConfiguration *config.ClusterConfiguration) (*ClusterDataSource, error) {
	if err := config.ValidateClusterConfiguration(clusterConfiguration); err != nil {
		return nil, err
	}
	remoteConfigurationLoader := loader.NewRemoteConfigurationLoader(clusterConfiguration.Props, clusterConfiguration.EtcdConfig)
	var (
		remoteConfiguration *config.RemoteClusterConfiguration
		region              string
	)
	if clusterConfiguration.Props != nil {
		remoteConfiguration = remoteConfigurationLoader.GetConfiguration()
		region = clusterConfiguration.Props.Region
	}
	integrationClusterConfiguration := &config.IntegrationClusterConfiguration{
		ClusterConfiguration:       clusterConfiguration,
		RemoteClusterConfiguration: remoteConfiguration,
	}

	nodeDataSourceMap := getDataSourceMap(integrationClusterConfiguration)
	routerConfig := integrationClusterConfiguration.GetRouterConfig()
	clusterDataSource := &ClusterDataSource{
		RouterConfiguration: routerConfig,
		DataSources:         nodeDataSourceMap,
		Active:              routerConfig.Active,
		Region:              region,
	}
	remoteConfigurationLoader.AddRouterListener(clusterDataSource)
	remoteConfigurationLoader.Init()
	return clusterDataSource, nil
}

// setActive if the activeKey node exists
func (cd *ClusterDataSource) setActive(activeKey string) {
	if _, ok := cd.DataSources[activeKey]; ok {
		if cd.Active == "" || activeKey != cd.Active {
			cd.Active = activeKey
			atomic.AddInt64(&cd.switchTimes, 1)
		}
	} else {
		log.Printf("WARNING: set activeKey = %s failed, because `dataSources` not exists such key\n", activeKey)
	}
}

// OnChanged implements RouterConfigurationListener interface, when remote routerConfiguration active changes,
// change the clusterDataSource's active node.
func (cd *ClusterDataSource) OnChanged(configuration *config.RouterConfiguration) {
	if configuration != nil {
		cd.setActive(configuration.Active)
	}
}

func (cd *ClusterDataSource) extend() {
}
