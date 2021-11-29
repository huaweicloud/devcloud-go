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

package loader

import (
	"os"
	"testing"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/stretchr/testify/assert"
)

var (
	dataSources = map[string]*config.RemoteDataSourceConfiguration{
		"ds0": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-1",
			Schema:   "ds0",
			Username: "root",
		},
		"ds0-slave0": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-1",
			Schema:   "ds0-slave0",
			Username: "root",
		},
		"ds0-slave1": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-1",
			Schema:   "ds0-slave1",
			Username: "root",
		},
		"ds1": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-2",
			Schema:   "ds1",
			Username: "root",
		},
		"ds1-slave0": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-2",
			Schema:   "ds1-slave0",
			Username: "root",
		},
		"ds1-slave1": {
			Server:   "127.0.0.1:3306",
			Cloud:    "huaweicloud",
			Region:   "cn-north-2",
			Schema:   "ds1-slave1",
			Username: "root",
		},
	}
	routerConfig = &config.RouterConfiguration{
		Active:         "c0",
		RouteAlgorithm: "single-read-write",
		Retry: &config.RetryConfiguration{
			Times: "10",
			Delay: "50",
		},
		Nodes: map[string]*config.NodeConfiguration{
			"c0": {
				Master:      "ds0",
				LoadBalance: "RANDOM",
				Slaves:      []string{"ds0-slave0", "ds0-slave1"},
			},
			"c1": {
				Master:      "ds1",
				LoadBalance: "ROUND_ROBIN",
				Slaves:      []string{"ds1-slave0", "ds1-slave1"},
			},
		},
	}
	remoteClusterConfiguration = &config.RemoteClusterConfiguration{
		DataSources:  dataSources,
		RouterConfig: routerConfig,
	}
	testHashCode = "test_hashCode"
)

func TestConfigurationFileHandler_Save(t *testing.T) {
	handler := NewConfigurationFileHandler()
	// remove defaultCacheConfigFile if exists
	if _, err := os.Stat(handler.cacheFilePath); err == nil {
		os.Remove(handler.cacheFilePath)
	}
	// test save
	handler.Save(remoteClusterConfiguration, testHashCode)

	// check if the cacheFile exists
	cacheFilePath := handler.getCompleteCacheFilePath(testHashCode)
	if _, err := os.Stat(cacheFilePath); err != nil && os.IsNotExist(err) {
		t.Error(err)
	}

	assertions := assert.New(t)

	configuration := handler.Load(testHashCode)
	assertions.Equal(6, len(configuration.DataSources))
	assertions.Equal(2, len(configuration.RouterConfig.Nodes))
	os.Remove(cacheFilePath)
	t.Log("OK")
}

// test with the cacheFile not exist
func TestConfigurationFileHandler_LoadFailed(t *testing.T) {
	homeDir := getHomeDir() + string(os.PathSeparator) + ".devspore" + string(os.PathSeparator)
	cacheFilePath := homeDir + "remote-config-test_hashCode.json"
	if _, err := os.Stat(cacheFilePath); err == nil {
		os.Remove(cacheFilePath)
	}
	handler := NewConfigurationFileHandler()
	cacheConfig := handler.Load(testHashCode)
	assert.Nil(t, cacheConfig)
	t.Log("OK")
}
