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
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/stretchr/testify/assert"
)

const (
	appId        = "xxxappId"
	monitorId    = "xxx-monitor-id"
	databaseName = "xxx-database"
)

var (
	etcdConfiguration = &etcd.EtcdConfiguration{
		Address:     "127.0.0.1:2379",
		Username:    "root",
		Password:    "123456",
		HTTPSEnable: false,
	}
	props = &mas.PropertiesConfiguration{
		AppID:        appId,
		MonitorID:    monitorId,
		DatabaseName: databaseName,
	}
	wrongEtcdAddress = "127.0.0.1:2380"
)

// TestRemoteConfigurationLoader_GetConfiguration need an actual etcd address.
func TestRemoteConfigurationLoader_GetConfiguration(t *testing.T) {
	loader := NewRemoteConfigurationLoader(props, etcdConfiguration)
	createRemoteConfiguration(loader)
	loader.Init()
	defer func() {
		err := loader.Close()
		if err != nil {
			t.Errorf("close remote configuration loader err, %v", err)
		}
	}()

	remoteConfiguration := loader.GetConfiguration(props.CalHashCode())
	assert := assert.New(t)
	assert.NotNil(remoteConfiguration)
	assert.NotNil(remoteConfiguration.DataSources)
	assert.Equal(len(remoteConfiguration.DataSources), 6)
	assert.NotNil(remoteConfiguration.RouterConfig)
	_, ok := remoteConfiguration.DataSources["ds0"]
	assert.True(ok)
	return
}

func TestGetConfigurationFromCache(t *testing.T) {
	handler := NewConfigurationFileHandler()
	// remove defaultCacheConfigFile if exists
	if _, err := os.Stat(handler.cacheFilePath); err == nil {
		os.Remove(handler.cacheFilePath)
	}
	handler.Save(&config.RemoteClusterConfiguration{
		DataSources:  dataSources,
		RouterConfig: routerConfig}, props.CalHashCode())

	wrongEtcdConfiguration := etcdConfiguration
	wrongEtcdConfiguration.Address = wrongEtcdAddress
	loader := NewRemoteConfigurationLoader(props, wrongEtcdConfiguration)
	loader.Init()
	defer func() {
		err := loader.Close()
		if err != nil {
			t.Errorf("close remote configuration loader err, %v", err)
		}
	}()

	localConfiguration := loader.GetConfiguration(props.CalHashCode())
	assert := assert.New(t)
	assert.NotNil(localConfiguration)
	assert.NotNil(localConfiguration.DataSources)
	assert.Equal(len(localConfiguration.DataSources), 6)
	assert.NotNil(localConfiguration.RouterConfig)
	_, ok := localConfiguration.DataSources["ds0"]
	assert.True(ok)

}

func createRemoteConfiguration(loader *RemoteConfigurationLoader) {
	client := loader.etcdClient
	datasourceStr, err := json.Marshal(dataSources)
	if err != nil {
		log.Printf("json marshal datasources failed, err %v", err)
	}
	_, err = client.Put(loader.dataSourceKey, string(datasourceStr))
	if err != nil {
		log.Printf("etcd put datasource failed, err %v", err)
	}

	routerConfigStr, err := json.Marshal(routerConfig)
	if err != nil {
		log.Printf("json marshal routerConfig failed, err %v", err)
	}
	_, err = client.Put(loader.routerKey, string(routerConfigStr))
	if err != nil {
		log.Printf("etcd put routerConfig failed, err %v", err)
	}

	_, err = client.Put(loader.activeKey, "c1")
	if err != nil {
		log.Printf("etcd put active failed, err %v", err)
	}
}

func TestListener(t *testing.T) {
	loader := NewRemoteConfigurationLoader(props, etcdConfiguration)
	loader.Init()

	loader.AddRouterListener(&mockListener{})
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		if err := modifyRouterConfig(); err != nil {
			t.Errorf("modify router config failed, err %s", err)
		}
		time.Sleep(time.Second)
	}
	time.Sleep(100 * time.Second)
}

func modifyRouterConfig() error {
	loader := NewRemoteConfigurationLoader(props, etcdConfiguration)
	client := loader.etcdClient
	val, err := client.Get(loader.activeKey)
	if err != nil {
		return err
	}
	var newVal string
	if val == "c0" {
		newVal = "c1"
	} else {
		newVal = "c0"
	}
	resp, err := client.Put(loader.activeKey, newVal)
	if err != nil {
		return err
	}
	log.Printf("previous value is %v\n", resp)
	return nil
}

type mockListener struct {
}

func (m *mockListener) OnChanged(config *config.RouterConfiguration) {
	log.Printf("active node:%v", config.Active)
	println("mockListener onchanged")
}
