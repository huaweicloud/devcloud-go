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

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/common/etcd/mocks"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/stretchr/testify/assert"
)

const (
	appId        = "xxxappId"
	monitorId    = "xxx-monitor-id"
	databaseName = "xxx-database"
)

var (
	props = &mas.PropertiesConfiguration{
		AppID:        appId,
		MonitorID:    monitorId,
		DatabaseName: databaseName,
	}
	wrongEtcdAddress = "127.0.0.1:2381"
)

func TestRemoteConfigurationLoader_GetConfiguration(t *testing.T) {
	loader := NewRemoteConfigurationLoader(props, nil)
	mockClient := &mocks.EtcdClient{}
	loader.etcdClient = mockClient
	createRemoteConfiguration(mockClient, loader)
	remoteConfiguration := loader.GetConfiguration(props.CalHashCode())

	assert.NotNil(t, remoteConfiguration)
	assert.NotNil(t, remoteConfiguration.DataSources)
	assert.Equal(t, len(remoteConfiguration.DataSources), 6)
	assert.NotNil(t, remoteConfiguration.RouterConfig)
	_, ok := remoteConfiguration.DataSources["ds0"]
	assert.True(t, ok)

	err := removeTempFile(props.CalHashCode())
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
}

func removeTempFile(hashCode string) error {
	fileHandler := NewConfigurationFileHandler()
	filePath := fileHandler.getCompleteCacheFilePath(hashCode)
	return os.RemoveAll(filePath)
}

func TestGetConfigurationFromCache(t *testing.T) {
	handler := NewConfigurationFileHandler()
	// remove defaultCacheConfigFile if exists
	if _, err := os.Stat(handler.cacheFilePath); err == nil {
		if err := os.Remove(handler.cacheFilePath); err != nil {
			t.Log("remove cache file failed")
			return
		}
	}
	handler.Save(&config.RemoteClusterConfiguration{
		DataSources:  dataSources,
		RouterConfig: routerConfig}, props.CalHashCode())

	wrongEtcdConfiguration := getEtcdConfiguration()
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
	assert.NotNil(t, localConfiguration)
	assert.NotNil(t, localConfiguration.DataSources)
	assert.Equal(t, len(localConfiguration.DataSources), 6)
	assert.NotNil(t, localConfiguration.RouterConfig)
	_, ok := localConfiguration.DataSources["ds0"]
	assert.True(t, ok)

	err := removeTempFile(props.CalHashCode())
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
}

func createRemoteConfiguration(mockClient *mocks.EtcdClient, loader *RemoteConfigurationLoader) {
	datasourceStr, err := json.Marshal(dataSources)
	if err != nil {
		log.Printf("json marshal datasources failed, err %v", err)
	}
	mockClient.On("Get", loader.dataSourceKey).Return(string(datasourceStr), nil).Once()
	routerConfigStr, err := json.Marshal(routerConfig)
	if err != nil {
		log.Printf("json marshal routerConfig failed, err %v", err)
	}
	mockClient.On("Get", loader.routerKey).Return(string(routerConfigStr), nil).Once()
	mockClient.On("Get", loader.activeKey).Return("c1", nil).Once()
}

func TestListener(t *testing.T) {
	dataDir := "etcd_data"
	defer func() {
		err := os.RemoveAll(dataDir)
		if err != nil {
			log.Println("remove failed")
		}
	}()
	metadata := mock.NewEtcdMetadata()
	metadata.DataDir = dataDir
	mockEtcd := &mock.MockEtcd{}
	mockEtcd.StartMockEtcd(metadata)
	defer mockEtcd.StopMockEtcd()

	loader := NewRemoteConfigurationLoader(props, getEtcdConfiguration())
	loader.Init()
	defer func() {
		err := loader.Close()
		log.Println(err)
	}()

	loader.AddRouterListener(&mockListener{})
	err := modifyRouterConfig()
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)

	active, err := loader.etcdClient.Get(loader.activeKey)
	assert.Nil(t, err)
	assert.NotNil(t, active)
}

func modifyRouterConfig() error {
	loader := NewRemoteConfigurationLoader(props, getEtcdConfiguration())
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
	_, err = client.Put(loader.activeKey, newVal)
	if err != nil {
		return err
	}
	return err
}

type mockListener struct {
}

func (m *mockListener) OnChanged(config *config.RouterConfiguration) {
	log.Printf("change active node to:%v", config.Active)
}

func getEtcdConfiguration() *etcd.EtcdConfiguration {
	return &etcd.EtcdConfiguration{
		Address:     "127.0.0.1:2379",
		Username:    "root",
		Password:    "root",
		HTTPSEnable: false,
	}
}
