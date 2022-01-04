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

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// need an actual etcd properties for test
var (
	props = &ClientProperties{
		Endpoints:          []string{"127.0.0.1:2379"},
		UserName:           "root",
		Password:           "root",
		NeedAuthentication: true,
	}
	key = "etcd_test_key"
	val = "etcd_test_value"
)

// TestEtcdV3Client test etcd put, get, del and watch operations
func TestEtcdV3Client(t *testing.T) {
	dataDir := "etcd_data"
	defer os.RemoveAll(dataDir)
	metadata := mock.NewEtcdMetadata()
	metadata.DataDir = dataDir
	mockEtcd := &mock.MockEtcd{}
	mockEtcd.StartMockEtcd(metadata)
	defer mockEtcd.StopMockEtcd()

	client, err := NewEtcdV3Client(props)
	if err != nil {
		t.Errorf("create etcd client err, %v", err.Error())
		return
	}
	println("create client success!")
	// watch key
	go client.Watch(key, 0, watchOnEvent)
	time.Sleep(time.Second)
	// test put
	putResp, err := client.Put(key, val)
	if err != nil {
		t.Errorf("put to etcd failed, err:%v\n", err)
		return
	}
	if putResp != "" {
		t.Logf("put etcd k:%s, v:%s, previous value is :%s", key, val, putResp)
	}
	// test get
	getResp, err := client.Get(key)
	if err != nil {
		t.Errorf("get etcd key:%s failed, err:%v\n", key, err)
		return
	}
	assert.Equal(t, val, getResp)
	// test list
	listResp, err := client.List("")
	assert.Equal(t, 1, len(listResp))
	// test delete
	_, err = client.Del(key)
	assert.Nil(t, err)

	err = client.Close()
	assert.Nil(t, err)
}

func watchOnEvent(event *clientv3.Event) {
	s := fmt.Sprintf("k:%v, v:%v, modversion:%v, type:%v", string(event.Kv.Key),
		string(event.Kv.Value), event.Kv.ModRevision, event.Type)
	log.Print(s)
}
