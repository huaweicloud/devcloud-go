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

package strategy

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type abstractStrategy struct {
	ClientPool    map[string]redis.UniversalClient
	Configuration *config.Configuration
}

func newAbstractStrategy(configuration *config.Configuration) abstractStrategy {
	strategy := abstractStrategy{
		Configuration: configuration,
		ClientPool:    make(map[string]redis.UniversalClient),
	}
	strategy.initClients()
	return strategy
}

func (a *abstractStrategy) initClients() {
	for name, serverConfig := range a.Configuration.RedisConfig.Servers {
		a.ClientPool[name] = newClient(serverConfig)
	}
}

func (a *abstractStrategy) activeClient() redis.UniversalClient {
	activeServer := a.Configuration.Active
	return a.getClientByServerName(activeServer)
}

func (a *abstractStrategy) nearestClient() redis.UniversalClient {
	nearest := a.Configuration.RedisConfig.Nearest
	return a.getClientByServerName(nearest)
}

func (a *abstractStrategy) remoteClient() redis.UniversalClient {
	nearest := a.Configuration.RedisConfig.Nearest
	for name, _ := range a.Configuration.RedisConfig.Servers {
		if name != nearest {
			return a.getClientByServerName(name)
		}
	}
	log.Println("ERROR: routeAlgorithm 'double-write' need another redis server for double write!")
	return nil
}

func (a *abstractStrategy) getClientByServerName(serverName string) redis.UniversalClient {
	if client, ok := a.ClientPool[serverName]; ok {
		return client
	}
	if serverConfig, ok := a.Configuration.RedisConfig.Servers[serverName]; ok && serverConfig != nil {
		a.ClientPool[serverName] = newClient(serverConfig)
		return a.ClientPool[serverName]
	}
	log.Printf("ERROR: server '%s' has no config!", serverName)
	return nil
}

func (a *abstractStrategy) Close() error {
	var err error
	for _, client := range a.ClientPool {
		err = client.Close()
	}
	return err
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
