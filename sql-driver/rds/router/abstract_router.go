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

// AbstractRouter implements Router interface.
type AbstractRouter struct {
	strategy RouteStrategy
}

// Route according to router's strategy.
func (ar *AbstractRouter) Route(isSQLOnlyRead bool, runtimeCtx *RuntimeContext, exclusives map[datasource.DataSource]bool) datasource.DataSource {
	return ar.strategy.Decorate(isSQLOnlyRead, runtimeCtx, exclusives)
}
