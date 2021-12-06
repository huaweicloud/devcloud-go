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
	"log"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
)

// NodeRouteStrategy route an actual datasource
type NodeRouteStrategy struct {
}

// Decorate implements RouteStrategy
func (ns *NodeRouteStrategy) Decorate(isSQLOnlyRead bool, runtimeCtx *RuntimeContext, exclusives map[datasource.DataSource]bool) datasource.DataSource {
	if nodeDataSource, ok := runtimeCtx.DataSource.(*datasource.NodeDataSource); ok {
		actualDataSource := ns.choose(nodeDataSource, exclusives, isSQLOnlyRead, runtimeCtx.InTransaction, runtimeCtx.RequestId)
		if _, exist := exclusives[actualDataSource]; !exist {
			return actualDataSource
		}
	} else {
		log.Print("ERROR: datasource.NodeDataSource type assertion error")
	}
	return nil
}

// The write operation or transaction select master datasource, and the read operation is select from slaves datasource
// according to the load balancing algorithm. when transaction is readOnly, then select a slave datasource
func (ns *NodeRouteStrategy) choose(dataSource *datasource.NodeDataSource, exclusives map[datasource.DataSource]bool,
	isSQLOnlyRead, inTransaction bool, requestId int64) *datasource.ActualDataSource {
	if inTransaction || !isSQLOnlyRead || len(dataSource.SlavesDatasource) == 0 {
		return dataSource.MasterDataSource
	}

	for i := 0; i < len(dataSource.SlavesDatasource); i++ {
		slave := dataSource.LoadBalanceAlgorithm.GetActualDataSource(requestId, dataSource.SlavesDatasource)
		if _, ok := exclusives[slave]; !ok {
			return slave
		} else {
			// if the chosen slave is in exclusives, then continue
			continue
		}
	}
	// when all slaves exclusive, return master datasource
	return dataSource.MasterDataSource
}
