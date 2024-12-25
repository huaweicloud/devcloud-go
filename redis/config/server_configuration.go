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
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/common/password"
	"github.com/huaweicloud/devcloud-go/common/util"
)

// ServerConfiguration contains yaml redis server configuration, go-redis Options and ClusterOptions(when type is cluster).
type ServerConfiguration struct {
	Hosts            string                             `yaml:"hosts"`
	Password         string                             `yaml:"password"`
	Type             string                             `yaml:"type"` // cluster, normal, sentinel
	Cloud            string                             `yaml:"cloud"`
	Region           string                             `yaml:"region"`
	Azs              string                             `yaml:"azs"`
	Db               int                                `yaml:"db"` // default 0
	MaxAttempts      int                                `yaml:"maxAttempts"`
	Timeout          int                                `yaml:"timeout"`
	SentinelPassword string                             `yaml:"sentinelPassword"` // sentinel opt
	MasterName       string                             `yaml:"masterName"`       // sentinel opt
	ConnectionPool   *ServerConnectionPoolConfiguration `yaml:"pool"`
	ClusterOptions   *redis.ClusterOptions
	FailoverOptions  *redis.FailoverOptions
	Options          *redis.Options
}

// ServerConnectionPoolConfiguration connection pool configuration
type ServerConnectionPoolConfiguration struct {
	MaxTotal                      int  `yaml:"maxTotal"`
	MaxIdle                       int  `yaml:"maxIdle"`
	MinIdle                       int  `yaml:"minIdle"`
	MaxWaitMillis                 int  `yaml:"maxWaitMillis"`
	TimeBetweenEvictionRunsMillis int  `yaml:"timeBetweenEvictionRunsMillis"`
	Lifo                          bool `yaml:"lifo"`
}

const (
	ServerTypeCluster     = "cluster"
	ServerTypeNormal      = "normal"
	ServerTypeMasterSlave = "master-slave"
	ServerTypeSentinel    = "sentinel"
)

// convertOptions convert yaml redis server configuration to go-redis Options or ClusterOptions
func (s *ServerConfiguration) convertOptions() {
	if s.Timeout == 0 {
		s.Timeout = 2000 // default 2000
	}
	timeout := time.Millisecond * time.Duration(s.Timeout)
	if s.MaxAttempts == 0 {
		s.MaxAttempts = 5 // default 5
	}
	if s.Type == ServerTypeCluster {
		clusterOpts := &redis.ClusterOptions{
			Addrs: util.ConvertAddressStrToSlice(s.Hosts, false),
		}
		if len(s.Password) > 0 {
			clusterOpts.Password = password.GetDecipher().Decode(s.Password)
		}
		clusterOpts.MaxRetries = s.MaxAttempts
		clusterOpts.DialTimeout = timeout
		clusterOpts.WriteTimeout = timeout
		clusterOpts.ReadTimeout = timeout
		if *s.ConnectionPool != (ServerConnectionPoolConfiguration{}) {
			clusterOpts.PoolSize = s.ConnectionPool.MaxTotal
			clusterOpts.MinIdleConns = s.ConnectionPool.MinIdle
			clusterOpts.IdleCheckFrequency = time.Duration(s.ConnectionPool.TimeBetweenEvictionRunsMillis) * time.Millisecond
			clusterOpts.PoolTimeout = time.Duration(s.ConnectionPool.MaxWaitMillis) * time.Millisecond
		}
		s.ClusterOptions = clusterOpts
		return
	} else if s.Type == ServerTypeSentinel {
		s.FailoverOptions = s.getFailoverOptions()
		return
	}
	opts := &redis.Options{
		Addr: s.Hosts,
	}
	if len(s.Password) > 0 {
		opts.Password = password.GetDecipher().Decode(s.Password)
	}
	opts.DB = s.Db
	opts.DialTimeout = timeout
	opts.WriteTimeout = timeout
	opts.ReadTimeout = timeout
	if s.ConnectionPool != nil && *s.ConnectionPool != (ServerConnectionPoolConfiguration{}) {
		opts.PoolSize = s.ConnectionPool.MaxTotal
		opts.MinIdleConns = s.ConnectionPool.MinIdle
		opts.IdleCheckFrequency = time.Duration(s.ConnectionPool.TimeBetweenEvictionRunsMillis) * time.Millisecond
		opts.PoolTimeout = time.Duration(s.ConnectionPool.MaxWaitMillis) * time.Millisecond
	}
	s.Options = opts
}

func (s *ServerConfiguration) getFailoverOptions() *redis.FailoverOptions {
	opts := &redis.FailoverOptions{
		SentinelAddrs: util.ConvertAddressStrToSlice(s.Hosts, false),
	}
	if len(s.Password) > 0 {
		opts.Password = password.GetDecipher().Decode(s.Password)
	}
	if len(s.SentinelPassword) > 0 {
		opts.SentinelPassword = password.GetDecipher().Decode(s.Password)
	}
	if s.MasterName == "" {
		s.MasterName = "mymaster"
	}
	opts.MasterName = s.MasterName
	opts.MaxRetries = s.MaxAttempts
	timeout := time.Millisecond * time.Duration(s.Timeout)
	opts.DialTimeout = timeout
	opts.WriteTimeout = timeout
	opts.ReadTimeout = timeout
	opts.DB = s.Db
	if s.ConnectionPool != nil && *s.ConnectionPool != (ServerConnectionPoolConfiguration{}) {
		opts.PoolSize = s.ConnectionPool.MaxTotal
		opts.MinIdleConns = s.ConnectionPool.MinIdle
		opts.IdleCheckFrequency = time.Duration(s.ConnectionPool.TimeBetweenEvictionRunsMillis) * time.Millisecond
		opts.PoolTimeout = time.Duration(s.ConnectionPool.MaxWaitMillis) * time.Millisecond
		opts.PoolFIFO = s.ConnectionPool.Lifo
	}
	return opts
}

func newDefaultConnectionPool() *ServerConnectionPoolConfiguration {
	return &ServerConnectionPoolConfiguration{
		MaxTotal:                      100,
		MaxIdle:                       10,
		MinIdle:                       0,
		MaxWaitMillis:                 10000,
		TimeBetweenEvictionRunsMillis: 60000,
		Lifo:                          true,
	}
}
