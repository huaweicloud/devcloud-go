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
Package config defines a series of devspore redis configuration, include configuration
from yaml and remote etcd.
*/
package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"gopkg.in/yaml.v3"
)

// Configuration is used to create DevsporeClient
type Configuration struct {
	Props          *mas.PropertiesConfiguration `yaml:"props"`
	EtcdConfig     *etcd.EtcdConfiguration      `yaml:"etcd"`
	RedisConfig    *RedisConfiguration          `yaml:"redis"`
	RouteAlgorithm string                       `yaml:"routeAlgorithm"`
	Active         string                       `yaml:"active"`
	Chaos          *mas.InjectionProperties     `yaml:"chaos"`
}

// OnChanged when remote etcd active key changed, change the Configuration's active server.
func (c *Configuration) OnChanged(active string) {
	c.Active = active
}

// AssignRemoteConfig will combine local configuration and remote configuration.
func (c *Configuration) AssignRemoteConfig() {
	remoteConfigurationLoader := NewRemoteConfigurationLoader(c.Props, c.EtcdConfig)
	remoteConfigurationLoader.AddRouterListener(c)
	remoteConfigurationLoader.Init()
	remoteConfiguration := remoteConfigurationLoader.GetConfiguration()
	if remoteConfiguration == nil {
		return
	}
	if remoteConfiguration.RouteAlgorithm != "" {
		c.RouteAlgorithm = remoteConfiguration.RouteAlgorithm
	}
	if remoteConfiguration.Active != "" {
		c.Active = remoteConfiguration.Active
	}
	if remoteConfiguration.Servers != nil {
		if c.RedisConfig.Servers == nil {
			c.RedisConfig.Servers = make(map[string]*ServerConfiguration)
		}
		for serverName, serverConfig := range remoteConfiguration.Servers {
			if _, ok := c.RedisConfig.Servers[serverName]; !ok {
				continue;
			}
			c.RedisConfig.Servers[serverName].Hosts = serverConfig.Hosts
			c.RedisConfig.Servers[serverName].Type = serverConfig.Type
			if len(serverConfig.Cloud) != 0 {
				c.RedisConfig.Servers[serverName].Cloud = serverConfig.Cloud
			}
			if len(serverConfig.Region) != 0 {
				c.RedisConfig.Servers[serverName].Region = serverConfig.Region
			}
			if len(serverConfig.Azs) != 0 {
				c.RedisConfig.Servers[serverName].Azs = serverConfig.Azs
			}
		}
	}
}

// ComputeNearestServer compute nearest redis server according to server's cloud, region and az.
func (c *Configuration) ComputeNearestServer() {
	if c.RedisConfig.Nearest != "" || c.Props == nil {
		return
	}
	for serverName, serverConfig := range c.RedisConfig.Servers {
		if serverConfig.Cloud == c.Props.Cloud && serverConfig.Region == c.Props.Region && serverConfig.Azs == c.Props.Azs {
			c.RedisConfig.Nearest = serverName
			return
		}
	}
	for serverName, serverConfig := range c.RedisConfig.Servers {
		if serverConfig.Cloud == c.Props.Cloud && serverConfig.Region == c.Props.Region {
			c.RedisConfig.Nearest = serverName
			return
		}
	}
	for serverName, serverConfig := range c.RedisConfig.Servers {
		if serverConfig.Cloud == c.Props.Cloud {
			c.RedisConfig.Nearest = serverName
			return
		}
	}
}

// ConvertServerConfiguration convert devspore server configuration to go-redis Options or cluster Options.
func (c *Configuration) ConvertServerConfiguration() {
	for _, serverConfig := range c.RedisConfig.Servers {
		serverConfig.convertOptions()
	}
}

// LoadConfiguration generate Configuration form yaml configuration file.
func LoadConfiguration(yamlFilePath string) (*Configuration, error) {
	realPath, err := filepath.Abs(yamlFilePath)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filepath.Clean(realPath))
	if err != nil {
		return nil, err
	}

	configuration := &Configuration{}
	if err = yaml.Unmarshal(yamlFile, configuration); err != nil {
		return nil, nil
	}
	if configuration.RedisConfig.ConnectionPoolConfig == nil {
		return configuration, nil
	}
	// clear connectionPool Config
	if !configuration.RedisConfig.ConnectionPoolConfig.Enable {
		for _, serverConfig := range configuration.RedisConfig.Servers {
			serverConfig.ConnectionPool = &ServerConnectionPoolConfiguration{}
		}
	} else {
		// if connectionPool config is nil, set default config
		for _, serverConfig := range configuration.RedisConfig.Servers {
			if serverConfig.ConnectionPool == nil {
				serverConfig.ConnectionPool = newDefaultConnectionPool()
			}
		}
	}
	return configuration, nil
}
