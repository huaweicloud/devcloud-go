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

package strategy

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type abstractStrategy struct {
	ClientPool          map[string]redis.UniversalClient
	Configuration       *config.Configuration
	injectionManagement *mas.InjectionManagement
}

func newAbstractStrategy(configuration *config.Configuration) abstractStrategy {
	strategy := abstractStrategy{
		Configuration: configuration,
		ClientPool:    map[string]redis.UniversalClient{}}
	if configuration.Chaos != nil {
		strategy.injectionManagement = mas.NewInjectionManagement(configuration.Chaos)
		strategy.injectionManagement.SetError(mas.RedisErrors())
		strategy.initClients(true)
	} else {
		strategy.initClients(false)
	}
	return strategy
}

func (a *abstractStrategy) initClients(chaos bool) {
	for name, serverConfig := range a.Configuration.RedisConfig.Servers {
		client := newClient(serverConfig)
		if chaos {
			client.AddHook(a)
		}
		a.ClientPool[name] = client
	}
}

func (a *abstractStrategy) activeClient() redis.UniversalClient {
	activeServer := a.Configuration.Active
	return a.getClientByServerName(activeServer)
}

func (a *abstractStrategy) noActiveClient() redis.UniversalClient {
	activeServer := a.Configuration.Active
	for name, _ := range a.Configuration.RedisConfig.Servers {
		if name != activeServer {
			return a.getClientByServerName(name)
		}
	}
	log.Println("info: 'double-write' need another redis server for double write!")
	return nil
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
	var client redis.UniversalClient
	switch serverConfig.Type {
	case config.ServerTypeCluster:
		client = redis.NewClusterClient(serverConfig.ClusterOptions)
	case config.ServerTypeNormal:
		client = redis.NewClient(serverConfig.Options)
	case config.ServerTypeSentinel:
		client = redis.NewFailoverClient(serverConfig.FailoverOptions)
	default:
		log.Printf("WARNING: invalid server type '%s'", serverConfig.Type)
		client = redis.NewClient(serverConfig.Options)
	}
	return client
}

func (a *abstractStrategy) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	err := a.injectionManagement.Inject()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (a *abstractStrategy) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (a *abstractStrategy) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (a *abstractStrategy) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
