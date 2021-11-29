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

// RemoteRedisConfiguration is set in remote etcd, contains active server, routeAlgorithm and redis servers.
type RemoteRedisConfiguration struct {
	RouteAlgorithm string
	Active         string
	Servers        map[string]*ServerConfiguration
}

// NewRemoteRedisConfiguration create a RemoteRedisConfiguration with etcd configuration strings .
func NewRemoteRedisConfiguration(routeAlgorithm, active, serversStr string) *RemoteRedisConfiguration {
	servers := make(map[string]*ServerConfiguration)
	if serversStr != "" {
		if err := json.Unmarshal([]byte(serversStr), &servers); err != nil {
			log.Printf("WARNING: unmarshal servers [%s] failed, err: %v", serversStr, err)
		}
	}
	return &RemoteRedisConfiguration{
		RouteAlgorithm: routeAlgorithm,
		Active:         active,
		Servers:        servers,
	}
}
