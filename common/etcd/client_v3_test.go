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
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// need an actual etcd properties for test
var (
	props = &ClientProperties{
		Endpoints:          []string{"127.0.0.1:2379"},
		UserName:           "root",
		Password:           "123456",
		NeedAuthentication: true,
	}
	key = "etcd_test_key"
	val = "etcd_test_value"
)

// TestEtcdV3Client test etcd put, get, del and watch operations
func TestEtcdV3Client(t *testing.T) {
	client, err := NewEtcdV3Client(props)
	if err != nil {
		t.Errorf("create etcd client err, %v", err.Error())
		return
	}

	t.Logf("create client success, client info: %+v", client)

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
	if getResp != val {
		t.Errorf("want value is %s, return value is :%s", val, getResp)
	}

	// test delete
	delResp, err := client.Del(key)
	if err != nil {
		t.Errorf("del etcd key:%s failed, err:%v\n", key, err)
	}
	t.Logf("del resp is:%v", delResp)

	listResp, err := client.List("/a")
	for i, resp := range listResp {
		t.Logf("list %d, resp:%v", i, *resp)
	}
	time.Sleep(time.Second)

	client.Close()

	time.Sleep(time.Second)
}

func watchOnEvent(event *clientv3.Event) {
	s := fmt.Sprintf("k:%v, v:%v, modversion:%v, type:%v", string(event.Kv.Key),
		string(event.Kv.Value), event.Kv.ModRevision, event.Type)
	log.Print(s)
}
