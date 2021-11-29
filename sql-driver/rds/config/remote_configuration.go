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
	"encoding/json"
	"log"
)

// RemoteDataSourceConfiguration etcd remote datasource configuration
type RemoteDataSourceConfiguration struct {
	Server   string
	Schema   string
	Username string
	Password string
	Region   string
	Cloud    string
}

// RemoteClusterConfiguration contains datasources map and router config
type RemoteClusterConfiguration struct {
	DataSources  map[string]*RemoteDataSourceConfiguration
	RouterConfig *RouterConfiguration
}

// NewRemoteClusterConfiguration with remote datasourceConfig and routerConfig
func NewRemoteClusterConfiguration(dataSourceConfigStr string, routerConfigStr string) *RemoteClusterConfiguration {
	dataSourceConfig := make(map[string]*RemoteDataSourceConfiguration)
	if dataSourceConfigStr != "" {
		if err := json.Unmarshal([]byte(dataSourceConfigStr), &dataSourceConfig); err != nil {
			log.Printf("WARNING: unmarshal dataSourceConfigStr failed, err %v", err)
		}
	} else {
		log.Printf("WARNING: dataSourceConfigStr is nil")
	}

	routerConfig := &RouterConfiguration{}
	if routerConfigStr != "" {
		if err := json.Unmarshal([]byte(routerConfigStr), routerConfig); err != nil {
			log.Printf("WARNING: unmarshal routerConfigStr failed, err %v", err)
		}
	} else {
		log.Printf("WARNING: routerConfigStr is nil")
	}

	return &RemoteClusterConfiguration{
		DataSources:  dataSourceConfig,
		RouterConfig: routerConfig,
	}
}
