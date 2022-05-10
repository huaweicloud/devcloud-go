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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	contextTimeOut = 1 * time.Second
	retryDelay     = 1000 * time.Millisecond // for etcd client watch loop delay
)

// EtcdV3Client implements EtcdClient interface.
type EtcdV3Client struct {
	*clientv3.Client
}

// NewEtcdV3Client create an *EtcdV3Client based on "go.etcd.io/etcd/client/v3"
func NewEtcdV3Client(props *ClientProperties) (*EtcdV3Client, error) {
	if props == nil || len(props.Endpoints) == 0 {
		return nil, errors.New("etcd endpoints can not be null")
	}
	config := &clientv3.Config{
		Endpoints:   props.Endpoints,
		DialTimeout: 1 * time.Second,
	}
	if props.NeedAuthentication {
		config.Username = props.UserName
		config.Password = props.Password
	}
	if props.CaCert != "" && props.ClientCert != "" && props.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(props.ClientCert, props.ClientKey)
		if err != nil {
			return nil, err
		}
		caData, err := ioutil.ReadFile(props.CaCert)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)

		config.TLS = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
	}

	client, err := clientv3.New(*config)

	if err != nil {
		return nil, err
	}

	return &EtcdV3Client{
		client,
	}, nil
}

// Get return the val corresponding to the key in etcd
func (c *EtcdV3Client) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	resp, err := c.Client.Get(ctx, key)
	cancel()
	if err != nil {
		log.Printf("ERROR: etcd get '%s' failed, err %v", key, err.Error())
		return "", err
	}
	if resp.Count <= 0 {
		log.Printf("ERROR: etcd get '%s' resp count <= 0", key)
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// Put store key-value in etcd
func (c *EtcdV3Client) Put(key, value string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	putResp, err := c.Client.Put(ctx, key, value, clientv3.WithPrevKV())
	cancel()
	if err != nil {
		log.Printf("ERROR: etcd put '%s'-'%s' failed, err %v", key, value, err)
		return "", err
	}
	if putResp.PrevKv != nil {
		return string(putResp.PrevKv.Value), nil
	}
	return "", nil
}

// List return a list key-val from etcd with prefix
func (c *EtcdV3Client) List(prefix string) ([]*KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	listResp, err := c.Client.Get(ctx, prefix, clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Printf("ERROR: etcd get prefix '%s' failed, err %v", prefix, err)
		return nil, err
	}
	var kvList []*KeyValue
	for _, kv := range listResp.Kvs {
		kvList = append(kvList, &KeyValue{
			Key:           string(kv.Key),
			Val:           string(kv.Value),
			ModifiedIndex: kv.ModRevision,
		})
	}
	return kvList, nil
}

// Del delete the key which in etcd
func (c *EtcdV3Client) Del(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	delResp, err := c.Client.Delete(ctx, key)
	cancel()
	if err != nil {
		log.Printf("ERROR: etcd delete '%s' failed, err %v", key, err)
		return 0, err
	}
	return delResp.Deleted, nil
}

// Watch monitor active_key changes
func (c *EtcdV3Client) Watch(prefix string, startIndex int64, onEvent func(event *clientv3.Event)) {
	var nextIndex int64 = 0
	if startIndex > 0 {
		nextIndex = startIndex + 1
	}
	for {
		if c.Client == nil {
			return
		}
		watchRespChan := c.Client.Watch(context.Background(), prefix, clientv3.WithRev(nextIndex), clientv3.WithPrefix())
		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				log.Printf("INFO: watch event type:%v, key:%s, v:%s, modversion:%v", event.Type,
					string(event.Kv.Key), string(event.Kv.Value), event.Kv.ModRevision)
				onEvent(event)
			}
		}
		time.Sleep(retryDelay)
	}
}

func (c *EtcdV3Client) Close() error {
	if c.Client == nil {
		log.Print("WARNING: etcd client is already nil")
		return nil
	}
	return c.Client.Close()
}
