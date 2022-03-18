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

// RedisConfiguration defines a series of redis configuration in yaml file.
type RedisConfiguration struct {
	RedisGroupName               string                            `yaml:"redisGroupName"`
	UserName                     string                            `yaml:"useName"`
	Nearest                      string                            `yaml:"nearest"`
	Servers                      map[string]*ServerConfiguration   `yaml:"servers"`
	ConnectionPoolConfig         *RedisConnectionPoolConfiguration `yaml:"connectionPool"`
	AsyncRemoteWrite             *AsyncRemoteWrite                 `yaml:"asyncRemoteWrite"`
	AsyncRemotePoolConfiguration *AsyncRemotePoolConfiguration     `yaml:"asyncRemotePool"`
}

type RedisConnectionPoolConfiguration struct {
	Enable bool `yaml:"enable"`
}

type AsyncRemoteWrite struct {
	RetryTimes int `yaml:"retryTimes"`
}

type AsyncRemotePoolConfiguration struct {
	Persist         bool   `yaml:"persist"`
	ThreadCoreSize  int    `yaml:"threadCoreSize"`
	MaximumPoolSize int    `yaml:"maximumPoolSize"`
	KeepAliveTime   int64  `yaml:"keepAliveTime"`
	TaskQueueSize   int    `yaml:"taskQueueSize"`
	PersistDir      string `yaml:"persistDir"`
}
