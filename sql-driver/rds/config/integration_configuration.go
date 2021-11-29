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

// IntegrationClusterConfiguration combine yaml cluster datasource configuration and etcd remote cluster configuration
type IntegrationClusterConfiguration struct {
	ClusterConfiguration       *ClusterConfiguration
	RemoteClusterConfiguration *RemoteClusterConfiguration
}

// GetDataSource Get datasource configuration map through merge local cluster configuration and remote cluster configuration.
func (c *IntegrationClusterConfiguration) GetDataSource() map[string]*DataSourceConfiguration {
	local := c.ClusterConfiguration.DataSource
	for key, localConfig := range local {
		local[key] = c.combineDataSourceConfig(key, localConfig)
	}
	return local
}

// combineDataSourceConfig combine yaml datasource and remote datasource according to the input key
func (c *IntegrationClusterConfiguration) combineDataSourceConfig(key string, config *DataSourceConfiguration) *DataSourceConfiguration {
	if c.RemoteClusterConfiguration == nil {
		return config
	}
	if _, ok := c.RemoteClusterConfiguration.DataSources[key]; !ok {
		return config
	}
	config.assign(c.RemoteClusterConfiguration.DataSources[key])
	return config
}

// GetRouterConfig Get router configuration map through merge local cluster configuration and remote cluster configuration.
func (c *IntegrationClusterConfiguration) GetRouterConfig() *RouterConfiguration {
	var target = &RouterConfiguration{}
	if c.RemoteClusterConfiguration == nil || c.RemoteClusterConfiguration.RouterConfig == nil {
		target = c.ClusterConfiguration.RouterConfig
	} else {
		target = c.RemoteClusterConfiguration.RouterConfig
	}
	target.Retry = c.ClusterConfiguration.RouterConfig.Retry
	return target
}
