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
 * Package router which will route an actual datasource according to route strategy and sql type.
 * the package defines two interfaces, RouteStrategy interface and Router interface.
 */

package router

import (
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
)

// Router interface
type Router interface {
	// Route return routeResult
	Route(isSQLOnlyRead bool, runtimeCtx *RuntimeContext, exclusives map[datasource.DataSource]bool) datasource.DataSource
}

type RuntimeContext struct {
	DataSource    datasource.DataSource
	InTransaction bool
	RequestId     int64
}
