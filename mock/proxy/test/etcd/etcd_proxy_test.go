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

// Package etcd proxy test cases
package etcd

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/mock/proxy"
)

func TestEtcdMock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcd")
}

var _ = Describe("Etcd", func() {
	var (
		etcdProxy *proxy.Proxy
		etcdMock  *mock.MockEtcd

		addrs      = []string{"127.0.0.1:2382"}
		dataDir    = "etcd_data"
		proxyAddrs = []string{"127.0.0.1:3382"}
		client     *clientv3.Client
		err        error

		key = "key"
		val = "val"
	)

	BeforeSuite(func() {
		metadata := mock.NewEtcdMetadata()
		metadata.ClientAddrs = addrs
		metadata.DataDir = dataDir
		etcdMock = &mock.MockEtcd{}
		etcdMock.StartMockEtcd(metadata)

		etcdProxy = proxy.NewProxy(addrs[0], proxyAddrs[0], proxy.Etcd)
		err := etcdProxy.StartProxy()
		if err != nil {
			log.Fatalln(err)
		}
		client, err = clientv3.New(clientv3.Config{Endpoints: proxyAddrs, Username: "root", Password: "root"})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Println(err)
			}
		}(dataDir)
		defer etcdMock.StopMockEtcd()
		defer etcdProxy.StopProxy()

		defer func(client *clientv3.Client) {
			err = client.Close()
			if err != nil {
				log.Println(err)
			}
		}(client)
	})

	AfterEach(func() {
		etcdProxy.DeleteAllRule()
	})

	It("Put", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = client.Put(ctx, key, val, clientv3.WithPrevKV())
		cancel()
		Expect(err).NotTo(HaveOccurred())
	})

	It("Get", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, key)
		cancel()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Count).To(Equal(int64(1)))
		Expect(string(resp.Kvs[0].Value)).To(Equal(val))
	})

	It("GetDelay", func() {
		_ = etcdProxy.AddDelay("delay", 200, 0, "", "")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, key)
		cancel()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Count).To(Equal(int64(1)))
		Expect(string(resp.Kvs[0].Value)).To(Equal(val))
	})

	It("GetJitter", func() {
		_ = etcdProxy.AddJitter("jitter", 200, 0, "", "")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, key)
		cancel()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Count).To(Equal(int64(1)))
		Expect(string(resp.Kvs[0].Value)).To(Equal(val))
	})
})
