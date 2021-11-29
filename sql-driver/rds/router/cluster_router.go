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

package router

const (
	SingleReadWrite      = "single-read-write"
	LocalReadSingleWrite = "local-read-single-write"
)

// ClusterRouter with router strategy.
type ClusterRouter struct {
	AbstractRouter
}

// NewClusterRouter return ClusterRouter according to the route type.
func NewClusterRouter(routeType string) *ClusterRouter {
	switch routeType {
	case SingleReadWrite:
		return &ClusterRouter{AbstractRouter{strategy: &ClusterRouteStrategy{}}}
	case LocalReadSingleWrite:
		return &ClusterRouter{AbstractRouter{strategy: &LocationBaseClusterRouteStrategy{}}}
	default:
		return &ClusterRouter{AbstractRouter{strategy: &ClusterRouteStrategy{}}}
	}
}
