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
Package config defines a series of configuration, include yaml configuration, remote configuration
and integration configuration which contains the above two configurations.
*/
package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"gopkg.in/yaml.v3"
)

// ClusterConfiguration yaml cluster configuration entity
type ClusterConfiguration struct {
	Props        *mas.PropertiesConfiguration        `yaml:"props"`
	EtcdConfig   *etcd.EtcdConfiguration             `yaml:"etcd"`
	RouterConfig *RouterConfiguration                `yaml:"router"`
	DataSource   map[string]*DataSourceConfiguration `yaml:"datasource"`
}

// RouterConfiguration yaml router configuration entity
type RouterConfiguration struct {
	Retry          *RetryConfiguration           `yaml:"retry"`
	Nodes          map[string]*NodeConfiguration `yaml:"nodes"`
	Active         string                        `yaml:"active"`
	RouteAlgorithm string                        `yaml:"routeAlgorithm"`
}

// RetryConfiguration yaml retry configuration entity
type RetryConfiguration struct {
	Times string `yaml:"times"`
	Delay string `yaml:"delay"`
}

// NodeConfiguration yaml node configuration entity
type NodeConfiguration struct {
	Weight      int      `yaml:"weight"`
	Master      string   `yaml:"master"`
	LoadBalance string   `yaml:"loadBalance"`
	Slaves      []string `yaml:"slaves"`
}

// DataSourceConfiguration contains yaml datasource configuration and remote datasource configuration
type DataSourceConfiguration struct {
	URL      string `yaml:"url"` // dsn format, `protocol(address)/dbname?param=value`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Region   string
	Server   string
	Schema   string
}

// assign combine remote datasource configuration and yaml datasource configuration
func (d *DataSourceConfiguration) assign(configuration *RemoteDataSourceConfiguration) {
	if configuration == nil {
		return
	}
	// update authentication
	if configuration.Username != "" {
		d.Username = configuration.Username
	}
	if configuration.Password != "" {
		d.Password = configuration.Password
	}
	// update region
	if configuration.Region != "" {
		d.Region = configuration.Region
	}
	// update server and schema
	if configuration.Server != "" && configuration.Schema != "" {
		d.Server = configuration.Server
		d.Schema = configuration.Schema
	}
}

// RouterConfigurationListener remote router configuration listener
type RouterConfigurationListener interface {
	// OnChanged when active data changes will call back
	OnChanged(configuration *RouterConfiguration)
}

// ValidateClusterConfiguration returns err if clusterConfiguration is invalid
func ValidateClusterConfiguration(configuration *ClusterConfiguration) error {
	if configuration == nil {
		return errors.New("clusterConfiguration cannot be nil")
	}
	if configuration.DataSource == nil {
		return errors.New("datasource config cannot be nil")
	}
	if configuration.RouterConfig == nil {
		return errors.New("router config cannot be nil")
	}
	if configuration.EtcdConfig != nil {
		if configuration.Props == nil {
			return errors.New("props cannot be nil")
		}
		if configuration.Props.AppID == "" {
			return errors.New("appId is required")
		}
		if configuration.Props.MonitorID == "" {
			return errors.New("monitorId is required")
		}
		if configuration.Props.DatabaseName == "" {
			return errors.New("databaseName is required")
		}
	}
	return nil
}

// Unmarshal yamlConfigFile to *ClusterConfiguration
func Unmarshal(yamlFilePath string) (*ClusterConfiguration, error) {
	yamlFile, err := ioutil.ReadFile(filepath.Clean(yamlFilePath))
	if err != nil {
		return nil, err
	}

	confContent := []byte(os.ExpandEnv(string(yamlFile)))
	clusterConfiguration := &ClusterConfiguration{}
	if err = yaml.Unmarshal(confContent, clusterConfiguration); err != nil {
		return nil, err
	}

	return clusterConfiguration, nil
}
