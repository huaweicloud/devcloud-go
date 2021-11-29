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
 * Package redis defines DevsporeClient which implements all redis commands in
 * https://github.com/go-redis/redis/blob/master/commands.go and it provides read-write
 * separation and etcd multi-data source disaster tolerance switching capabilities,
 * user can create a DevsporeClient by yaml configuration file or by code, see details in README.md.
 */

package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

// DevsporeClient implements go-redis/UniversalClient interface which defines all redis commands, DevsporeClient includes
// configuration and a client pool, the pool will store redis client or cluster client.
type DevsporeClient struct {
	ctx           context.Context
	configuration *config.Configuration
	clientPool    map[string]redis.UniversalClient
}

// NewDevsporeClientWithYaml create a devsporeClient with yaml configuration.
func NewDevsporeClientWithYaml(yamlFilePath string) *DevsporeClient {
	configuration, err := config.LoadConfiguration(yamlFilePath)
	if err != nil {
		log.Fatalf("ERROR: create DevsporeClient failed, err [%v]", err)
		return nil
	}
	if err = config.ValidateConfiguration(configuration); err != nil {
		log.Fatalf("ERROR: configuration is invalid, config is [%+v], err [%v]", configuration, err)
		return nil
	}
	configuration.AssignRemoteConfig()
	configuration.ComputeNearestServer()
	configuration.ConvertServerConfiguration()
	return &DevsporeClient{
		ctx:           context.Background(),
		clientPool:    make(map[string]redis.UniversalClient),
		configuration: configuration,
	}
}

// NewDevsporeClient create a devsporeClient with Configuration which will assign etcd remote configuration.
func NewDevsporeClient(configuration *config.Configuration) *DevsporeClient {
	configuration.AssignRemoteConfig()
	configuration.ComputeNearestServer()
	return &DevsporeClient{
		ctx:           context.Background(),
		clientPool:    make(map[string]redis.UniversalClient),
		configuration: configuration,
	}
}

func (c *DevsporeClient) getActualClient(opType commandType) redis.UniversalClient {
	serverName := c.route(opType)
	if client, ok := c.clientPool[serverName]; ok {
		return client
	}
	if serverConfig, ok := c.configuration.RedisConfig.Servers[serverName]; ok && serverConfig != nil {
		c.clientPool[serverName] = newClient(serverConfig)
		return c.clientPool[serverName]
	}
	return nil
}

func newClient(serverConfig *config.ServerConfiguration) redis.UniversalClient {
	switch serverConfig.Type {
	case config.ServerTypeCluster:
		return redis.NewClusterClient(serverConfig.ClusterOptions)
	case config.ServerTypeNormal, config.ServerTypeMasterSlave:
		return redis.NewClient(serverConfig.Options)
	default:
		log.Printf("WARNING: invalid server type '%s'", serverConfig.Type)
		return redis.NewClient(serverConfig.Options)
	}
}

const (
	singleReadWrite      = "single-read-write"
	localReadSingleWrite = "local-read-single-write"
)

func (c *DevsporeClient) route(opType commandType) string {
	switch c.configuration.RouteAlgorithm {
	case singleReadWrite:
		return c.configuration.Active
	case localReadSingleWrite:
		if opType == commandTypeRead {
			return c.configuration.RedisConfig.Nearest
		}
		return c.configuration.Active
	default:
		log.Printf("WARNING: invalid route algorithm '%s'", c.configuration.RouteAlgorithm)
		c.configuration.RouteAlgorithm = singleReadWrite
		return c.configuration.Active
	}
}

// Close closes all clients in clientPool
func (c *DevsporeClient) Close() error {
	var err error
	for _, client := range c.clientPool {
		err = client.Close()
	}
	return err
}
