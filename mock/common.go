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
 * package mock introduce three mock methods for interface, mysql and redis.
 */

package mock

// use https://github.com/stretchr/testify and https://github.com/vektra/mockery to mock interface,
// see example in common/etcd/client.go, sql-driver/rds/config/loader/remote_configuration_loader.go
// GetConfiguration() function and sql-driver/rds/config/loader/remote_configuration_loader_test.go
// TestRemoteConfigurationLoader_GetConfiguration() function.
// See more in https://bbs.huaweicloud.com/blogs/315144
