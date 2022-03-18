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
Package etcd defines EtcdClient interface, and use "go.etcd.io/etcd/client/v3"
implements the interface.
*/
package etcd

import (
	"log"

	"github.com/huaweicloud/devcloud-go/common/password"
	"github.com/huaweicloud/devcloud-go/common/util"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//go:generate mockery -name=EtcdClient
type EtcdClient interface {
	Get(key string) (string, error)
	Put(key, value string) (string, error)
	List(prefix string) ([]*KeyValue, error)
	Del(key string) (int64, error)
	Watch(prefix string, startIndex int64, onEvent func(event *clientv3.Event))
	Close() error
}

// CreateEtcdClient according to yaml etcdConfiguration
func CreateEtcdClient(etcdConfiguration *EtcdConfiguration) EtcdClient {
	properties := &ClientProperties{
		Endpoints: util.ConvertAddressStrToSlice(etcdConfiguration.Address),
	}
	if etcdConfiguration.Username != "" {
		properties.UserName = etcdConfiguration.Username
		properties.Password = password.GetDecipher().Decode(etcdConfiguration.Password)
		properties.NeedAuthentication = true
	}
	client, err := NewEtcdV3Client(properties)
	if err != nil || client == nil {
		log.Printf("ERROR: create etcd client failed, err %v", err)
		return nil
	}
	return client
}
