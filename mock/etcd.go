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

package mock

import (
	"log"
	"net/url"
	"time"

	"go.etcd.io/etcd/api/v3/authpb"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/server/v3/embed"
)

type MockEtcd struct {
	e   *embed.Etcd
	cfg *embed.Config
}

func (m *MockEtcd) StartMockEtcd(metadata *EtcdMetadata) {
	cfg := embed.NewConfig()
	cfg.LCUrls = convertAddrsToUrls(metadata.ClientAddrs)
	cfg.LPUrls = convertAddrsToUrls(metadata.PeerAddrs)
	cfg.Dir = metadata.DataDir
	cfg.LogLevel = "warn"
	m.cfg = cfg

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Printf("ERROR: start embed etcd failed, %v", err)
		return
	}
	m.e = e
	if metadata.AuthEnable {
		m.initRootRole()
		m.AddUser(defaultUser, defaultPassword)
		err = m.e.Server.AuthStore().AuthEnable()
		if err != nil {
			log.Printf("ERROR: enable auth failed, %v", err)
			return
		}
	}
	if len(metadata.UserName) > 0 && len(metadata.Password) > 0 {
		m.AddUser(metadata.UserName, metadata.Password)
	}
	select {
	case <-m.e.Server.ReadyNotify():
		log.Printf("Start mock etcd!")
	case <-time.After(60 * time.Second):
		m.e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
}

func (m *MockEtcd) StopMockEtcd() {
	m.e.Server.Stop()
	m.e.Close()
	log.Println("Stop mock etcd!")
}

func (m *MockEtcd) initRootRole() {
	authStore := m.e.Server.AuthStore()
	_, err := authStore.RoleAdd(&pb.AuthRoleAddRequest{Name: defaultRole})
	if err != nil {
		log.Printf("ERROR: add 'root' role failed, %+v", err)
	}

	_, err = authStore.RoleGrantPermission(&pb.AuthRoleGrantPermissionRequest{
		Name: defaultRole,
		Perm: &authpb.Permission{
			PermType: 2,
		}})
	if err != nil {
		log.Printf("ERROR: RoleGrantPermission failed, %+v", err)
	}
}

const (
	defaultRole     = "root"
	defaultUser     = "root"
	defaultPassword = "root"
)

func (m *MockEtcd) AddUser(user, password string) {
	authStore := m.e.Server.AuthStore()
	_, err := authStore.UserAdd(&pb.AuthUserAddRequest{Name: user, Password: password})
	if err != nil {
		log.Printf("WARNING: add user %s failed, %v", user, err)
	}
	_, err = authStore.UserGrantRole(&pb.AuthUserGrantRoleRequest{User: user, Role: defaultRole})
	if err != nil {
		log.Printf("WARNING: user %s grant root role failed, %v", user, err)
	}
}

type EtcdMetadata struct {
	ClientAddrs []string
	PeerAddrs   []string // optional
	AuthEnable  bool
	UserName    string
	Password    string
	DataDir     string
}

func NewEtcdMetadata() *EtcdMetadata {
	return &EtcdMetadata{
		ClientAddrs: []string{"127.0.0.1:2379"},
		PeerAddrs:   []string{"127.0.0.1:2380"},
		AuthEnable:  true,
		DataDir:     "default.etcd",
	}
}

func convertAddrsToUrls(addrs []string) []url.URL {
	var urls []url.URL
	for _, addr := range addrs {
		tUrl, err := url.Parse("http://" + addr)
		if err != nil {
			log.Println(err)
		}
		urls = append(urls, *tUrl)
	}
	return urls
}
