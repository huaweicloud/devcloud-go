/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2025.
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

package redigostrategy

import (
	"context"
	"fmt"
	"log"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/huaweicloud/devcloud-go/mas"
	"github.com/huaweicloud/devcloud-go/redis/config"
	"github.com/mna/redisc"
)

// nomal & sebtinel by redis.Pool , cluster by redisc.Cluster
type RedigoUniversalClient struct {
	*redis.Pool
	*redisc.Cluster
}

func (r RedigoUniversalClient) Close() error {
	if r.Pool != nil {
		return r.Pool.Close()
	} else if r.Cluster != nil {
		return r.Cluster.Close()
	}
	return fmt.Errorf("close no available pool")
}

func (r RedigoUniversalClient) Get() redis.Conn {
	if r.Pool != nil {
		return r.Pool.Get()
	} else if r.Cluster != nil {
		return r.Cluster.Get()
	}
	log.Println("Error: get no available pool")
	return nil
}

func (r RedigoUniversalClient) Dial() (redis.Conn, error) {
	if r.Pool != nil {
		return r.Pool.Dial()
	} else if r.Cluster != nil {
		return r.Cluster.Dial()
	}
	return nil, fmt.Errorf("dial no available pool")
}

func (r RedigoUniversalClient) Stats() redis.PoolStats {
	if r.Pool != nil {
		return r.Pool.Stats()
	}
	return redis.PoolStats{}
}

func (r RedigoUniversalClient) ClusterStats() map[string]redis.PoolStats {
	if r.Cluster != nil {
		return r.Cluster.Stats()
	}
	return nil
}

func (r RedigoUniversalClient) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (r RedigoUniversalClient) DoContext(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.Get()
	defer conn.Close()
	return redis.DoContext(conn, ctx, commandName, args...)
}

type RedigoCommandArgs struct {
	CommandName string
	Args        []interface{}
}

type RedigoPipeineArgs struct {
	CommandName string
	Args        []interface{}
}

func (r RedigoUniversalClient) Pipeline(transactions bool, cmds interface{}) (reply []interface{}, err error) {
	args, err := covertPipelineCmds(cmds)
	if err != nil {
		return nil, err
	}
	if transactions {
		return r.Transactions(args)
	}
	conn := r.Get()
	defer conn.Close()
	reviceCount := 0
	for k, cmd := range args {
		if cmd.CommandName == "" {
			log.Printf("pipeline cmds[%d] CommandName is empty\n", k)
			continue
		}
		conn.Send(cmd.CommandName, cmd.Args...)
		reviceCount++
	}
	conn.Flush()
	err = nil
	reply = make([]interface{}, reviceCount)
	for i := 0; i < reviceCount; i++ {
		r, err1 := conn.Receive()
		if err1 != nil {
			reply[i] = err1
			err = err1
		} else {
			reply[i] = r
		}
	}
	return reply, err
}

func covertPipelineCmds(cmds interface{}) ([]*RedigoCommandArgs, error) {
	args := make([]*RedigoCommandArgs, 0)
	switch cmds := cmds.(type) {
	case [][]string:
		for k, cmd := range cmds {
			if len(cmd) == 0 {
				log.Printf("pipeline cmds[%d] len is 0\n", k)
				continue
			}
			commandName := cmd[0]
			if len(cmd) == 1 {
				args = append(args, &RedigoCommandArgs{
					CommandName: commandName,
					Args:        nil,
				})
				continue
			}
			arglist := make([]interface{}, 0, len(cmd))
			for i := 1; i < len(cmd); i++ {
				arglist = append(arglist, cmd[i])
			}
			args = append(args, &RedigoCommandArgs{
				CommandName: commandName,
				Args:        arglist,
			})
		}
	case []*RedigoCommandArgs:
		args = cmds
	default:
		log.Fatalf("unsupported command type: %T", cmds)
		return nil, fmt.Errorf("unsupported command type: %T", cmds)
	}
	return args, nil
}

func (r RedigoUniversalClient) Transactions(args []*RedigoCommandArgs) (reply []interface{}, err error) {
	conn := r.Get()
	defer conn.Close()
	conn.Send("MULTI")
	for k, cmd := range args {
		if cmd.CommandName == "" {
			log.Printf("pipeline cmds[%d] CommandName is empty\n", k)
			continue
		}
		conn.Send(cmd.CommandName, cmd.Args...)
	}

	return redis.Values(conn.Do("EXEC"))
}

type abstractRedigoStrategy struct {
	ClientPool          map[string]*RedigoUniversalClient
	Configuration       *config.Configuration
	injectionManagement *mas.InjectionManagement
}

func newAbstractStrategy(configuration *config.Configuration) abstractRedigoStrategy {
	strategy := abstractRedigoStrategy{
		Configuration: configuration,
		ClientPool:    map[string]*RedigoUniversalClient{}}
	strategy.initClients(false)

	return strategy
}

func (a *abstractRedigoStrategy) initClients(chaos bool) {
	for name, serverConfig := range a.Configuration.RedisConfig.Servers {
		client := newClient(serverConfig)
		if chaos {
			log.Println("Info: redigo client no support chaos")
		}
		a.ClientPool[name] = client
	}
}

func (a *abstractRedigoStrategy) activeClient() *RedigoUniversalClient {
	activeServer := a.Configuration.Active
	return a.getClientByServerName(activeServer)
}

func (a *abstractRedigoStrategy) noActiveClient() *RedigoUniversalClient {
	activeServer := a.Configuration.Active
	for name, _ := range a.Configuration.RedisConfig.Servers {
		if name != activeServer {
			return a.getClientByServerName(name)
		}
	}
	log.Println("info: 'single-read-async-double-write' need another redis server for double write!")
	return nil
}

func (a *abstractRedigoStrategy) nearestClient() *RedigoUniversalClient {
	nearest := a.Configuration.RedisConfig.Nearest
	return a.getClientByServerName(nearest)
}

func (a *abstractRedigoStrategy) remoteClient() *RedigoUniversalClient {
	nearest := a.Configuration.RedisConfig.Nearest
	for name, _ := range a.Configuration.RedisConfig.Servers {
		if name != nearest {
			return a.getClientByServerName(name)
		}
	}
	log.Println("ERROR: routeAlgorithm 'local-read-async-double-write' need another redis server for double write!")
	return &RedigoUniversalClient{}
}

func (a *abstractRedigoStrategy) getClientByServerName(serverName string) *RedigoUniversalClient {
	if client, ok := a.ClientPool[serverName]; ok {
		return client
	}
	if serverConfig, ok := a.Configuration.RedisConfig.Servers[serverName]; ok && serverConfig != nil {
		a.ClientPool[serverName] = newClient(serverConfig)
		return a.ClientPool[serverName]
	}
	log.Printf("ERROR: server '%s' has no config!", serverName)
	return &RedigoUniversalClient{}
}

func (a *abstractRedigoStrategy) Close() error {
	var err error
	for _, client := range a.ClientPool {
		err = client.Close()
	}
	return err
}

func newClient(serverConfig *config.ServerConfiguration) *RedigoUniversalClient {
	var client *RedigoUniversalClient
	switch serverConfig.Type {
	case config.ServerTypeCluster:
		client = newClusterClient(serverConfig)
	case config.ServerTypeNormal:
		client = newNormalClient(serverConfig)
	case config.ServerTypeSentinel:
		client = newSentinelClient(serverConfig)
	default:
		log.Printf("WARNING: invalid server type '%s'", serverConfig.Type)
	}
	return client
}

func newNormalClient(serverConfig *config.ServerConfiguration) *RedigoUniversalClient {
	pool := &redis.Pool{

		MaxIdle:     serverConfig.ConnectionPool.MaxIdle,
		MaxActive:   serverConfig.ConnectionPool.MaxTotal,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", serverConfig.Hosts,
				redis.DialPassword(serverConfig.Password),
				redis.DialConnectTimeout(time.Duration(serverConfig.Options.DialTimeout)),
				redis.DialWriteTimeout(time.Duration(serverConfig.Options.DialTimeout)),
				redis.DialReadTimeout(time.Duration(serverConfig.Options.DialTimeout)),
				redis.DialDatabase(serverConfig.Options.DB),
			)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedigoUniversalClient{
		Pool: pool,
	}
}

func newSentinelClient(serverConfig *config.ServerConfiguration) *RedigoUniversalClient {
	sntnl := &sentinel.Sentinel{
		Addrs:      serverConfig.FailoverOptions.SentinelAddrs,
		MasterName: serverConfig.FailoverOptions.MasterName,
		Dial: func(addr string) (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr,
				redis.DialConnectTimeout(time.Duration(serverConfig.FailoverOptions.DialTimeout)))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	pool := &redis.Pool{
		MaxIdle:     serverConfig.FailoverOptions.PoolSize,
		MaxActive:   serverConfig.FailoverOptions.PoolSize,
		Wait:        true,
		IdleTimeout: serverConfig.FailoverOptions.PoolTimeout,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr,
				redis.DialPassword(serverConfig.Password),
				redis.DialConnectTimeout(time.Duration(serverConfig.FailoverOptions.DialTimeout)),
				redis.DialWriteTimeout(time.Duration(serverConfig.FailoverOptions.DialTimeout)),
				redis.DialReadTimeout(time.Duration(serverConfig.FailoverOptions.DialTimeout)),
				redis.DialDatabase(serverConfig.FailoverOptions.DB))
			if err != nil {
				return nil, err
			}
			isMaster, err := sentinel.TestRole(c, "master")
			if err != nil {
				c.Close()
				return nil, err
			}
			if !isMaster {
				c.Close()
				return nil, fmt.Errorf("%s is not redis master", masterAddr)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedigoUniversalClient{
		Pool: pool,
	}
}

func newClusterClient(serverConfig *config.ServerConfiguration) *RedigoUniversalClient {
	config := newRedigoClusterConfig(serverConfig)
	cluster := &redisc.Cluster{
		StartupNodes: serverConfig.ClusterOptions.Addrs,
		DialOptions:  []redis.DialOption{redis.DialConnectTimeout(5 * time.Second)},
		CreatePool:   createPool(config),
	}
	return &RedigoUniversalClient{
		Cluster: cluster,
	}
}

func Watch(client redis.Conn, keys ...string) error {
	if len(keys) > 0 {
		args := make([]interface{}, len(keys))
		for i, key := range keys {
			args[i] = key
		}
		_, err := client.Do("Watch", args...)
		return err
	}
	return nil
}

func (a *abstractRedigoStrategy) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	err := a.injectionManagement.Inject()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (a *abstractRedigoStrategy) AfterProcess(ctx context.Context, cmd goredis.Cmder) error {
	return nil
}

func (a *abstractRedigoStrategy) BeforeProcessPipeline(ctx context.Context, cmds []goredis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (a *abstractRedigoStrategy) AfterProcessPipeline(ctx context.Context, cmds []goredis.Cmder) error {
	return nil
}
