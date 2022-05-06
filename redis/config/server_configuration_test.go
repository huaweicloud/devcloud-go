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
	"testing"
	"time"

	"github.com/huaweicloud/devcloud-go/common/password"
	"github.com/stretchr/testify/assert"
)

func TestServerConfigurationConvertOptions(t *testing.T) {
	configuration, err := LoadConfiguration("../resources/config_with_password.yaml")
	if err != nil {
		t.Errorf("load configuration from failed, err: %v", err)
		return
	}
	password.SetDecipher(&myDecipher{})
	configuration.ConvertServerConfiguration()
	assert.Equal(t, len(configuration.RedisConfig.Servers), 2)
	for serverName, serverConfig := range configuration.RedisConfig.Servers {
		if serverName == "dc1" {
			assert.NotNil(t, serverConfig.Options)
			assert.Equal(t, "127.0.0.1:6379", serverConfig.Options.Addr)
			assert.Equal(t, "XXXX!!!", serverConfig.Options.Password)
			assert.Equal(t, 100, serverConfig.Options.PoolSize)
			assert.Equal(t, 0, serverConfig.Options.MinIdleConns)
			assert.Equal(t, time.Duration(1000)*time.Millisecond, serverConfig.Options.IdleCheckFrequency)
			assert.Equal(t, time.Duration(10000)*time.Millisecond, serverConfig.Options.PoolTimeout)
		} else if serverName == "dc2" {
			assert.NotNil(t, serverConfig.ClusterOptions)
			assert.Equal(t, []string{"127.0.0.1:6380", "127.0.0.1:6381"}, serverConfig.ClusterOptions.Addrs)
			assert.Equal(t, "XXXX!!!", serverConfig.ClusterOptions.Password)
			assert.Equal(t, 100, serverConfig.ClusterOptions.PoolSize)
			assert.Equal(t, 0, serverConfig.ClusterOptions.MinIdleConns)
			assert.Equal(t, time.Duration(1000)*time.Millisecond, serverConfig.ClusterOptions.IdleCheckFrequency)
			assert.Equal(t, time.Duration(10000)*time.Millisecond, serverConfig.ClusterOptions.PoolTimeout)
		}
	}
}

func TestServerConfigurationConvertOptions_DisablePool(t *testing.T) {
	configuration, err := LoadConfiguration("../resources/config_disable_pool.yaml")
	if err != nil {
		t.Errorf("load configuration from config_disable_pool.yaml failed, err: %v", err)
		return
	}
	configuration.ConvertServerConfiguration()
	for _, serverConfig := range configuration.RedisConfig.Servers {
		assert.NotNil(t, serverConfig.Options)
		assert.Equal(t, 0, len(serverConfig.Options.Password))
		assert.Equal(t, 0, serverConfig.Options.PoolSize)
		assert.Equal(t, 0, serverConfig.Options.MinIdleConns)
		assert.Equal(t, time.Duration(0), serverConfig.Options.IdleCheckFrequency)
		assert.Equal(t, time.Duration(0), serverConfig.Options.PoolTimeout)
	}
}

func TestServerConfigurationConvertOptions_DefaultPool(t *testing.T) {
	configuration, err := LoadConfiguration("../resources/config_no_pool.yaml")
	if err != nil {
		t.Errorf("load configuration from config_no_pool.yaml failed, err: %v", err)
		return
	}
	configuration.ConvertServerConfiguration()
	assert.Equal(t, len(configuration.RedisConfig.Servers), 2)
	for _, serverConfig := range configuration.RedisConfig.Servers {
		assert.Equal(t, 0, len(serverConfig.Options.Password))
		assert.Equal(t, 100, serverConfig.Options.PoolSize)
		assert.Equal(t, 0, serverConfig.Options.MinIdleConns)
		assert.Equal(t, time.Duration(60000)*time.Millisecond, serverConfig.Options.IdleCheckFrequency)
		assert.Equal(t, time.Duration(10000)*time.Millisecond, serverConfig.Options.PoolTimeout)
	}
}

type myDecipher struct {
}

func (d *myDecipher) Decode(password string) string {
	return password + "!!!"
}
