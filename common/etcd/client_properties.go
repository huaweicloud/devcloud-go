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

package etcd

// ClientProperties properties for etcd
type ClientProperties struct {
	Endpoints          []string
	UserName           string
	Password           string
	NeedAuthentication bool
	ClientCert         string
	ClientKey          string
	CaCert             string
}

// KeyValue is etcd-Kv Simplified version
type KeyValue struct {
	Key           string
	Val           string
	ModifiedIndex int64
}

// EtcdConfiguration yaml etcd configuration entity
type EtcdConfiguration struct {
	APIVersion  string `yaml:"apiVersion"`
	Address     string `yaml:"address"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	HTTPSEnable bool   `yaml:"httpsEnable"`
	ClientCert  string `yaml:"clientCert"` // etcd cert file
	ClientKey   string `yaml:"clientKey"`  // etcd cert-key file
	CaCert      string `yaml:"caCert"`     // etcd ca file
}
