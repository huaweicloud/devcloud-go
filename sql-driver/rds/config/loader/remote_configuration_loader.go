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
	"fmt"
	"log"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	datasourcePrefix = "/mas-monitor/conf/db/services/%s/%s/database/%s/datasource"
	routerPrefix     = "/mas-monitor/conf/db/services/%s/%s/database/%s/router"
	activePrefix     = "/mas-monitor/status/db/services/%s/%s/database/%s/active"
)

// RemoteConfigurationLoader to load remote configuration from etcd
type RemoteConfigurationLoader struct {
	etcdClient    etcd.EtcdClient
	dataSourceKey string
	routerKey     string
	activeKey     string
	listeners     []config.RouterConfigurationListener
}

//  @param props is yaml properties configuration entity
//  @param etcdConfiguration is yaml etcd configuration entity
func NewRemoteConfigurationLoader(props *mas.PropertiesConfiguration,
	etcdConfiguration *etcd.EtcdConfiguration) *RemoteConfigurationLoader {
	var appID, monitorID, databaseTag string
	if props != nil {
		appID = props.AppID
		monitorID = props.MonitorID
		databaseTag = props.DatabaseName
	}
	loader := &RemoteConfigurationLoader{
		dataSourceKey: fmt.Sprintf(datasourcePrefix, appID, monitorID, databaseTag),
		routerKey:     fmt.Sprintf(routerPrefix, appID, monitorID, databaseTag),
		activeKey:     fmt.Sprintf(activePrefix, appID, monitorID, databaseTag),
	}
	if etcdConfiguration != nil && etcdConfiguration.Address != "" {
		loader.etcdClient = etcd.CreateEtcdClient(etcdConfiguration)
	}
	return loader
}

// GetConfiguration form etcd, or from local cache
func (l *RemoteConfigurationLoader) GetConfiguration(hashCode string) *config.RemoteClusterConfiguration {
	handler := NewConfigurationFileHandler()
	if l.etcdClient == nil {
		return handler.Load(hashCode)
	}

	dataSourceConfig, err := l.etcdClient.Get(l.dataSourceKey)
	if err != nil || dataSourceConfig == "" {
		log.Printf("ERROR: get remote datasourceConfig failed, %v", err)
		return handler.Load(hashCode)
	}

	routerConfig, err := l.etcdClient.Get(l.routerKey)
	if err != nil || routerConfig == "" {
		log.Printf("ERROR: get remote routerConfig failed, %v", err)
		return handler.Load(hashCode)
	}

	remoteClusterConfiguration := config.NewRemoteClusterConfiguration(dataSourceConfig, routerConfig)
	active, err := l.etcdClient.Get(l.activeKey)
	if err != nil {
		log.Printf("ERROR: get remote active failed, %v", err)
		return handler.Load(hashCode)
	}

	remoteClusterConfiguration.RouterConfig.Active = active
	// save file to local
	handler.Save(remoteClusterConfiguration, hashCode)
	return remoteClusterConfiguration
}

// AddRouterListener add a router configuration listener
func (l *RemoteConfigurationLoader) AddRouterListener(listener config.RouterConfigurationListener) {
	if l.listeners == nil {
		l.listeners = []config.RouterConfigurationListener{}
	}
	l.listeners = append(l.listeners, listener)
}

// onChanged listening for etcd activeKey changes
func (l *RemoteConfigurationLoader) onChanged(event *clientv3.Event) {
	if string(event.Kv.Key) == l.activeKey && event.Type == clientv3.EventTypePut {
		for _, listener := range l.listeners {
			newRouterConfiguration := &config.RouterConfiguration{Active: string(event.Kv.Value)}
			listener.OnChanged(newRouterConfiguration)
		}
	}
}

// Init etcd start watch activeKey
func (l *RemoteConfigurationLoader) Init() {
	if l.etcdClient == nil {
		return
	}
	go l.etcdClient.Watch(l.activeKey, 0, l.onChanged)
}

// Close loader's etcdClient and set loader's listeners nil
func (l *RemoteConfigurationLoader) Close() error {
	if l.etcdClient == nil {
		return nil
	}
	err := l.etcdClient.Close()
	l.etcdClient = nil
	l.listeners = nil
	return err
}
