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
 *  specific language governing permissions and limitations under the License.
 *
 */

package mock

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestMockEtcd(t *testing.T) {
	// start mock etcd
	addrs := []string{"127.0.0.1:2382"}
	dataDir := "etcd_data"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Errorf("ERROR: remove data dir failed, %v", err)
		}
	}(dataDir)
	metadata := NewEtcdMetadata()
	metadata.ClientAddrs = addrs
	metadata.DataDir = dataDir
	mockEtcd := &MockEtcd{}
	mockEtcd.StartMockEtcd(metadata)
	defer mockEtcd.StopMockEtcd()

	client, err := clientv3.New(clientv3.Config{Endpoints: addrs, Username: "root", Password: "root"})
	assert.Nil(t, err)
	defer func(client *clientv3.Client) {
		err = client.Close()
		if err != nil {
			t.Errorf("ERROR: close client failed, %v", err)
		}
	}(client)

	testKey := "key"
	testVal := "val"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, testKey, testVal, clientv3.WithPrevKV())
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, testKey)
	cancel()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), resp.Count)
	assert.Equal(t, testVal, string(resp.Kvs[0].Value))
}
