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

import (
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
)

// Location base cluster router strategy.
type LocationBaseClusterRouteStrategy struct {
}

// Decorate return target node datasource according to sql type and clusterDatasource's region.
func (ls *LocationBaseClusterRouteStrategy) Decorate(isSQLOnlyRead bool, runtimeCtx *RuntimeContext,
	exclusives map[datasource.DataSource]bool) datasource.DataSource {
	if clusterDataSource, ok := runtimeCtx.DataSource.(*datasource.ClusterDataSource); ok {
		nodeDataSource := ls.choose(isSQLOnlyRead, clusterDataSource)
		if _, exist := exclusives[nodeDataSource]; !exist {
			return nodeDataSource
		}
	}
	return nil
}

// Determine whether to execute the local data source or remote data source based on the SQL type.
func (ls *LocationBaseClusterRouteStrategy) choose(isSQLOnlyRead bool,
	clusterDataSource *datasource.ClusterDataSource) *datasource.NodeDataSource {
	region := clusterDataSource.Region
	if isSQLOnlyRead {
		for _, nodeDataSource := range clusterDataSource.DataSources {
			if region == nodeDataSource.Region {
				return nodeDataSource
			}
		}
	} else if nodeDataSource, ok := clusterDataSource.DataSources[clusterDataSource.Active]; ok {
		return nodeDataSource
	}
	return nil
}
