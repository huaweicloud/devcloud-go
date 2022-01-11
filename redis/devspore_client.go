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
	"errors"
	"fmt"
	"log"

	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

// DevsporeClient implements go-redis/UniversalClient interface which defines all redis commands, DevsporeClient includes
// configuration and a client pool, the pool will store redis client or cluster client.
type DevsporeClient struct {
	ctx           context.Context
	configuration *config.Configuration
	strategy      strategy.StrategyMode
}

// NewDevsporeClientWithYaml create a devsporeClient with yaml configuration.
func NewDevsporeClientWithYaml(yamlFilePath string) *DevsporeClient {
	configuration, err := config.LoadConfiguration(yamlFilePath)
	if err != nil {
		log.Fatalf("ERROR: create DevsporeClient failed, err [%v]", err)
		return nil
	}
	return NewDevsporeClient(configuration)
}

// NewDevsporeClient create a devsporeClient with Configuration which will assign etcd remote configuration.
func NewDevsporeClient(configuration *config.Configuration) *DevsporeClient {
	configuration.AssignRemoteConfig()
	configuration.ComputeNearestServer()
	configuration.ConvertServerConfiguration()
	if err := validateConfiguration(configuration); err != nil {
		log.Fatalf("ERROR: configration is invalid, config is [%+v], err [%v]", configuration, err)
		return nil
	}
	return &DevsporeClient{
		ctx:           context.Background(),
		strategy:      strategy.NewStrategy(configuration),
		configuration: configuration,
	}
}

// Close closes all clients in clientPool
func (c *DevsporeClient) Close() error {
	return c.strategy.Close()
}

// validateConfiguration check configuration is valid.
func validateConfiguration(configuration *config.Configuration) error {
	if configuration == nil {
		return errors.New("configuration cannot be nil")
	}
	if configuration.RedisConfig == nil {
		return errors.New("redis config cannot be nil")
	}
	if configuration.RouteAlgorithm == "" {
		return errors.New("router config cannot be null")
	}
	if configuration.EtcdConfig != nil {
		if configuration.Props == nil {
			return errors.New("props is required")
		}
		if configuration.Props.AppID == "" {
			return errors.New("appId is required")
		}
		if configuration.Props.MonitorID == "" {
			return errors.New("monitorId is required")
		}
	}
	if configuration.RedisConfig.Servers == nil || len(configuration.RedisConfig.Servers) == 0 {
		return errors.New("servers is required")
	}
	if configuration.RouteAlgorithm == strategy.DoubleWriteMode && configuration.RedisConfig.Nearest == "" {
		return fmt.Errorf("routeAlgorithm: %s required nearest setting", strategy.DoubleWriteMode)
	}
	return nil
}
