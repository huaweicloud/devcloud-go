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
	"fmt"
	"log"

	"github.com/huaweicloud/devcloud-go/common/etcd"
	"github.com/huaweicloud/devcloud-go/mas"
	"go.etcd.io/etcd/client/v3"
)

const (
	serversPrefix        = "/mas-monitor/conf/dcs/services/%s/%s/servers"
	routeAlgorithmPrefix = "/mas-monitor/conf/dcs/services/%s/%s/route-algorithm"
	activePrefix         = "/mas-monitor/status/dcs/services/%s/%s/active"
)

// RemoteConfigurationLoader to load remote configuration from etcd
type RemoteConfigurationLoader struct {
	etcdClient         etcd.EtcdClient
	routerAlgorithmKey string
	activeKey          string
	serversKey         string
	listeners          []Listener
}

// NewRemoteConfigurationLoader create a loader to load remote configuration
func NewRemoteConfigurationLoader(props *mas.PropertiesConfiguration,
	etcdConfiguration *etcd.EtcdConfiguration) *RemoteConfigurationLoader {
	var appID, monitorID string
	if props != nil {
		appID = props.AppID
		monitorID = props.MonitorID
	}
	loader := &RemoteConfigurationLoader{
		routerAlgorithmKey: fmt.Sprintf(routeAlgorithmPrefix, appID, monitorID),
		activeKey:          fmt.Sprintf(activePrefix, appID, monitorID),
		serversKey:         fmt.Sprintf(serversPrefix, appID, monitorID),
	}
	if etcdConfiguration != nil && etcdConfiguration.Address != "" {
		loader.etcdClient = etcd.CreateEtcdClient(etcdConfiguration)
	}
	return loader
}

// GetConfiguration from etcd, which contains route algorithm, active server and all redis servers.
func (l *RemoteConfigurationLoader) GetConfiguration() *RemoteRedisConfiguration {
	if l.etcdClient == nil {
		return nil
	}
	routeAlgorithm, err := l.etcdClient.Get(l.routerAlgorithmKey)
	if err != nil {
		log.Printf("ERROR: get remote routerConfig failed, %v", err)
		return nil
	}

	active, err := l.etcdClient.Get(l.activeKey)
	if err != nil {
		log.Printf("ERROR: get remote active failed, %v", err)
		return nil
	}

	serversStr, err := l.etcdClient.Get(l.serversKey)
	if err != nil {
		log.Printf("ERROR: get remote serversConfig failed, %v", err)
		return nil
	}

	return NewRemoteRedisConfiguration(routeAlgorithm, active, serversStr)
}

// AddRouterListener add a listener
func (l *RemoteConfigurationLoader) AddRouterListener(listener Listener) {
	if l.listeners == nil {
		l.listeners = []Listener{}
	}
	l.listeners = append(l.listeners, listener)
}

// onChanged listening for etcd activeKey changes
func (l *RemoteConfigurationLoader) onChanged(event *clientv3.Event) {
	if string(event.Kv.Key) == l.activeKey && event.Type == clientv3.EventTypePut {
		for _, listener := range l.listeners {
			listener.OnChanged(string(event.Kv.Value))
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
