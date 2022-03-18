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
 */

/*
Package loader provides ConfigurationFileHandler and RemoteConfigurationLoader,
the ConfigurationFileHandler is used for load config from cache and save remote config to cache,
the RemoteConfigurationLoader is used for load config from remote etcd.
*/
package loader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

const (
	cacheConfigFilePrefix  = "remote-config"
	cacheConfigFileSuffix  = ".json"
	defaultCacheConfigFile = cacheConfigFilePrefix + cacheConfigFileSuffix
	maxRetryTimes          = 5
	retryDuration          = 100 * time.Millisecond
	cacheFilePerm          = 0666
)

// ConfigurationFileHandler for cacheConfig load and save
type ConfigurationFileHandler struct {
	cacheFilePath string
	cacheFileLock *sync.Mutex
	cacheFileDir  string
}

func NewConfigurationFileHandler() *ConfigurationFileHandler {
	homeDir := getHomeDir() + string(os.PathSeparator) + ".devspore" + string(os.PathSeparator)
	cacheFilePath := homeDir + defaultCacheConfigFile
	return &ConfigurationFileHandler{
		cacheFilePath: cacheFilePath,
		cacheFileLock: &sync.Mutex{},
		cacheFileDir:  homeDir,
	}
}

// Save remote config to local
func (h *ConfigurationFileHandler) Save(config *config.RemoteClusterConfiguration, hashCode string) {
	h.cacheFileLock.Lock()
	defer h.cacheFileLock.Unlock()
	for i := 0; i < maxRetryTimes; i++ {
		err := h.doWrite(config, hashCode)
		if err == nil {
			break
		}
		if i == maxRetryTimes {
			log.Printf("WARNING: save config to local failed, err %v", err)
			return
		}
		time.Sleep(retryDuration)
	}
	log.Print("INFO: save config success")
}

// write remote cluster configuration to local according to the hashCode
func (h *ConfigurationFileHandler) doWrite(config *config.RemoteClusterConfiguration, hashCode string) error {
	cfgStr, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(h.cacheFileDir, cacheFilePerm); err != nil {
		return err
	}
	if err = ioutil.WriteFile(h.getCompleteCacheFilePath(hashCode), cfgStr, cacheFilePerm); err != nil {
		return err
	}
	return nil
}

// Load local cache config when get remote config from etcd failed
func (h *ConfigurationFileHandler) Load(hashCode string) *config.RemoteClusterConfiguration {
	h.cacheFileLock.Lock()
	defer h.cacheFileLock.Unlock()
	// if cacheConfigFile not exists, then return
	cacheConfigPath := h.getCompleteCacheFilePath(hashCode)
	if _, err := os.Stat(cacheConfigPath); err != nil {
		log.Printf("WARNING: %v", err)
		return nil
	}
	content, err := ioutil.ReadFile(cacheConfigPath)
	if err != nil {
		log.Printf("WARNING: read config from local cache file failed, err %v", err)
		return nil
	}
	remoteClusterConfiguration := &config.RemoteClusterConfiguration{}
	if err = json.Unmarshal(content, remoteClusterConfiguration); err != nil {
		log.Printf("WARNING: unmarshal local config failed, err %v", err)
		return nil
	}
	return remoteClusterConfiguration
}

// get user current home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// getCompleteCacheFilePath returns cacheFilePath according to the hashCode
func (h *ConfigurationFileHandler) getCompleteCacheFilePath(hashCode string) string {
	if hashCode == "" {
		return h.cacheFilePath
	}
	fileName := cacheConfigFilePrefix + "-" + hashCode + cacheConfigFileSuffix
	return strings.Replace(h.cacheFilePath, defaultCacheConfigFile, fileName, 1)
}
