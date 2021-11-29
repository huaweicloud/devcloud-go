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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationClusterConfiguration_NoRemote(t *testing.T) {
	localConfig := createLocalClusterConfig()

	integration := &IntegrationClusterConfiguration{ClusterConfiguration: localConfig}
	assert.Equal(t, integration.GetDataSource(), localConfig.DataSource)
	assert.Equal(t, integration.GetRouterConfig(), localConfig.RouterConfig)
}

func TestIntegrationClusterConfiguration_WithRemote(t *testing.T) {
	localConfig := createLocalClusterConfig()

	remoteDs0 := &RemoteDataSourceConfiguration{
		Username: "remote_ds0_root",
		Password: "remote_root",
	}
	remoteDs1 := &RemoteDataSourceConfiguration{
		Username: "remote_ds1_root",
		Password: "remote_root",
	}
	remoteDataSources := make(map[string]*RemoteDataSourceConfiguration, 2)
	remoteDataSources["ds0"] = remoteDs0
	remoteDataSources["ds1"] = remoteDs1
	remoteConfig := &RemoteClusterConfiguration{DataSources: remoteDataSources, RouterConfig: localConfig.RouterConfig}

	integration := &IntegrationClusterConfiguration{
		ClusterConfiguration:       localConfig,
		RemoteClusterConfiguration: remoteConfig,
	}

	integrationDataSources := integration.GetDataSource()
	assert.Equal(t, integrationDataSources["ds0"].Username, "remote_ds0_root")
	assert.Equal(t, integrationDataSources["ds0"].Password, "remote_root")
	assert.Equal(t, integrationDataSources["ds1"].Username, "remote_ds1_root")
	assert.Equal(t, integrationDataSources["ds1"].Password, "remote_root")
}

func createLocalClusterConfig() *ClusterConfiguration {
	ds0 := &DataSourceConfiguration{
		URL:      "tcp(127.0.0.1:3306)/ds0",
		Username: "root",
		Password: "root",
	}
	ds1 := &DataSourceConfiguration{
		URL:      "tcp(127.0.0.1:3306)/ds1",
		Username: "root",
		Password: "root",
	}
	dataSourcesMap := make(map[string]*DataSourceConfiguration, 2)
	dataSourcesMap["ds0"] = ds0
	dataSourcesMap["ds1"] = ds1

	nodes := make(map[string]*NodeConfiguration, 2)
	nodes["az0"] = &NodeConfiguration{Master: "ds0"}
	nodes["az1"] = &NodeConfiguration{Master: "ds1"}

	router := &RouterConfiguration{Nodes: nodes}

	localConfig := &ClusterConfiguration{
		RouterConfig: router,
		DataSource:   dataSourcesMap,
	}
	return localConfig
}
