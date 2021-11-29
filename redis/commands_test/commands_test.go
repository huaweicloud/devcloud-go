/*
 * Copyright (c) 2013 The github.com/go-redis/redis Authors.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *    * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * 2021.11.15-Changed modify the constructor Client, split a file into
 * multiple small files, and extract duplicate code.
 * 			Huawei Technologies Co., Ltd.
 */

package commands_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	devspore "github.com/huaweicloud/devcloud-go/redis"
	"github.com/huaweicloud/devcloud-go/redis/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	redisPort          = "6379"
	redisAddr          = ":" + redisPort
	redisSecondaryPort = "6388"
)

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-redis")
}

func redisOptions() *redis.Options {
	return &redis.Options{
		Addr:         redisAddr,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,

		MaxRetries: -1,

		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func configuration() *config.Configuration {
	servers := map[string]*config.ServerConfiguration{
		"server1": {
			Type:    config.ServerTypeNormal,
			Cloud:   "huaweiCloud",
			Region:  "beijing",
			Azs:     "az0",
			Options: redisOptions(),
		},
	}
	return &config.Configuration{
		RedisConfig: &config.RedisConfiguration{
			Servers: servers,
		},
		RouteAlgorithm: "single-read-write",
		Active:         "server1",
	}
}

// "Commands" contains some commands which are available in redis 6.2.0+, so if your redis version is 5.0+, you need to
// skip those commands in redis 6.2.0+, execute 'ginkgo -skip="redis6"' in the terminal.
var _ = Describe("Commands", func() {
	ctx := context.TODO()
	var client *devspore.DevsporeClient

	conf := configuration()
	BeforeEach(func() {
		client = devspore.NewDevsporeClient(conf)
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	Describe("server", func() {
		It("should Auth", func() {
			cmds, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Auth(ctx, "password")
				pipe.Auth(ctx, "")
				return nil
			})
			Expect(err).To(HaveOccurred())
			Expect(cmds[0].Err()).To(HaveOccurred())
			Expect(cmds[1].Err()).To(HaveOccurred())

			stats := client.PoolStats()
			Expect(stats.Hits).To(Equal(uint32(1)))
			Expect(stats.Misses).To(Equal(uint32(1)))
			Expect(stats.Timeouts).To(Equal(uint32(0)))
			Expect(stats.TotalConns).To(Equal(uint32(1)))
			Expect(stats.IdleConns).To(Equal(uint32(1)))
		})

		It("should Echo", func() {
			pipe := client.Pipeline()
			echo := pipe.Echo(ctx, "hello")
			_, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(echo.Err()).NotTo(HaveOccurred())
			Expect(echo.Val()).To(Equal("hello"))
		})

		It("should Ping", func() {
			ping := client.Ping(ctx)
			Expect(ping.Err()).NotTo(HaveOccurred())
			Expect(ping.Val()).To(Equal("PONG"))
		})

		It("should Wait", func() {
			const wait = 3 * time.Second

			// assume testing on single redis instance
			start := time.Now()
			val, err := client.Wait(ctx, 1, wait).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(int64(0)))
			Expect(time.Now()).To(BeTemporally("~", start.Add(wait), 3*time.Second))
		})

		It("should Select", func() {
			pipe := client.Pipeline()
			sel := pipe.Select(ctx, 1)
			_, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(sel.Err()).NotTo(HaveOccurred())
			Expect(sel.Val()).To(Equal("OK"))
		})

		It("should SwapDB", func() {
			pipe := client.Pipeline()
			sel := pipe.SwapDB(ctx, 1, 2)
			_, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(sel.Err()).NotTo(HaveOccurred())
			Expect(sel.Val()).To(Equal("OK"))
		})

		It("should BgRewriteAOF", func() {
			Skip("flaky test")

			val, err := client.BgRewriteAOF(ctx).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(ContainSubstring("Background append only file rewriting"))
		})

		It("should BgSave", func() {
			Skip("flaky test")

			// workaround for "ERR Can't BGSAVE while AOF log rewriting is in progress"
			Eventually(func() string {
				return client.BgSave(ctx).Val()
			}, "30s").Should(Equal("Background saving started"))
		})

		It("should ClientKill", func() {
			r := client.ClientKill(ctx, "1.1.1.1:1111")
			Expect(r.Err()).To(MatchError("ERR No such client"))
			Expect(r.Val()).To(Equal(""))
		})

		It("should ClientKillByFilter", func() {
			r := client.ClientKillByFilter(ctx, "TYPE", "test")
			Expect(r.Err()).To(MatchError("ERR Unknown client type 'test'"))
			Expect(r.Val()).To(Equal(int64(0)))
		})

		It("should ClientID", func() {
			err := client.ClientID(ctx).Err()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.ClientID(ctx).Val()).To(BeNumerically(">=", 0))
		})

		It("should ClientUnblock", func() {
			id := client.ClientID(ctx).Val()
			r, err := client.ClientUnblock(ctx, id).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal(int64(0)))
		})

		It("should ClientUnblockWithError", func() {
			id := client.ClientID(ctx).Val()
			r, err := client.ClientUnblockWithError(ctx, id).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal(int64(0)))
		})

		It("should ClientPause", func() {
			err := client.ClientPause(ctx, time.Second).Err()
			Expect(err).NotTo(HaveOccurred())

			start := time.Now()
			err = client.Ping(ctx).Err()
			Expect(err).NotTo(HaveOccurred())
			Expect(time.Now()).To(BeTemporally("~", start.Add(time.Second), 800*time.Millisecond))
		})

		It("should ClientSetName and ClientGetName", func() {
			pipe := client.Pipeline()
			set := pipe.ClientSetName(ctx, "theclientname")
			get := pipe.ClientGetName(ctx)
			_, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(set.Err()).NotTo(HaveOccurred())
			Expect(set.Val()).To(BeTrue())

			Expect(get.Err()).NotTo(HaveOccurred())
			Expect(get.Val()).To(Equal("theclientname"))
		})

		It("should ConfigGet", func() {
			val, err := client.ConfigGet(ctx, "*").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeEmpty())
		})

		It("should ConfigResetStat", func() {
			r := client.ConfigResetStat(ctx)
			Expect(r.Err()).NotTo(HaveOccurred())
			Expect(r.Val()).To(Equal("OK"))
		})

		It("should ConfigSet", func() {
			configGet := client.ConfigGet(ctx, "maxmemory")
			Expect(configGet.Err()).NotTo(HaveOccurred())
			Expect(configGet.Val()).To(HaveLen(2))
			Expect(configGet.Val()[0]).To(Equal("maxmemory"))

			configSet := client.ConfigSet(ctx, "maxmemory", configGet.Val()[1].(string))
			Expect(configSet.Err()).NotTo(HaveOccurred())
			Expect(configSet.Val()).To(Equal("OK"))
		})

		It("should ConfigRewrite", func() {
			configRewrite := client.ConfigRewrite(ctx)
			Expect(configRewrite.Err()).NotTo(HaveOccurred())
			Expect(configRewrite.Val()).To(Equal("OK"))
		})

		It("should DBSize", func() {
			size, err := client.DBSize(ctx).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("should Info", func() {
			info := client.Info(ctx)
			Expect(info.Err()).NotTo(HaveOccurred())
			Expect(info.Val()).NotTo(Equal(""))
		})

		It("should Info cpu", func() {
			info := client.Info(ctx, "cpu")
			Expect(info.Err()).NotTo(HaveOccurred())
			Expect(info.Val()).NotTo(Equal(""))
			Expect(info.Val()).To(ContainSubstring(`used_cpu_sys`))
		})

		It("should LastSave", func() {
			lastSave := client.LastSave(ctx)
			Expect(lastSave.Err()).NotTo(HaveOccurred())
			Expect(lastSave.Val()).NotTo(Equal(0))
		})

		It("should Save", func() {
			// workaround for "ERR Background save already in progress"
			Eventually(func() string {
				return client.Save(ctx).Val()
			}, "10s").Should(Equal("OK"))
		})

		It("should SlaveOf", func() {
			slaveOf := client.SlaveOf(ctx, "localhost", "8888")
			Expect(slaveOf.Err()).NotTo(HaveOccurred())
			Expect(slaveOf.Val()).To(Equal("OK"))

			slaveOf = client.SlaveOf(ctx, "NO", "ONE")
			Expect(slaveOf.Err()).NotTo(HaveOccurred())
			Expect(slaveOf.Val()).To(Equal("OK"))
		})

		It("should Time", func() {
			tm, err := client.Time(ctx).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(tm).To(BeTemporally("~", time.Now(), 3*time.Second))
		})

		It("should Command", func() {
			cmds, err := client.Command(ctx).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(cmds)).To(BeNumerically("~", 200, 25))

			cmd := cmds["mget"]
			Expect(cmd.Name).To(Equal("mget"))
			Expect(cmd.Arity).To(Equal(int8(-2)))
			Expect(cmd.Flags).To(ContainElement("readonly"))
			Expect(cmd.FirstKeyPos).To(Equal(int8(1)))
			Expect(cmd.LastKeyPos).To(Equal(int8(-1)))
			Expect(cmd.StepCount).To(Equal(int8(1)))

			cmd = cmds["ping"]
			Expect(cmd.Name).To(Equal("ping"))
			Expect(cmd.Arity).To(Equal(int8(-1)))
			Expect(cmd.Flags).To(ContainElement("stale"))
			Expect(cmd.Flags).To(ContainElement("fast"))
			Expect(cmd.FirstKeyPos).To(Equal(int8(0)))
			Expect(cmd.LastKeyPos).To(Equal(int8(0)))
			Expect(cmd.StepCount).To(Equal(int8(0)))
		})
	})

	Describe("debugging", func() {
		It("should DebugObject", func() {
			err := client.DebugObject(ctx, "foo").Err()
			Expect(err).To(MatchError("ERR no such key"))

			err = client.Set(ctx, "foo", "bar", 0).Err()
			Expect(err).NotTo(HaveOccurred())

			s, err := client.DebugObject(ctx, "foo").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(ContainSubstring("serializedlength:4"))
		})

		It("should MemoryUsage", func() {
			err := client.MemoryUsage(ctx, "foo").Err()
			Expect(err).To(Equal(redis.Nil))

			err = client.Set(ctx, "foo", "bar", 0).Err()
			Expect(err).NotTo(HaveOccurred())

			n, err := client.MemoryUsage(ctx, "foo").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).NotTo(BeZero())

			n, err = client.MemoryUsage(ctx, "foo", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).NotTo(BeZero())
		})
	})

	Describe("scanning", func() {
		It("should Scan", func() {
			for i := 0; i < 1000; i++ {
				set := client.Set(ctx, fmt.Sprintf("key%d", i), "hello", 0)
				Expect(set.Err()).NotTo(HaveOccurred())
			}

			keys, cursor, err := client.Scan(ctx, 0, "", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeZero())
		})

		It("[redis6] should ScanType", func() {
			for i := 0; i < 1000; i++ {
				set := client.Set(ctx, fmt.Sprintf("key%d", i), "hello", 0)
				Expect(set.Err()).NotTo(HaveOccurred())
			}

			keys, cursor, err := client.ScanType(ctx, 0, "", 0, "string").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeZero())
		})

		It("should SScan", func() {
			for i := 0; i < 1000; i++ {
				sadd := client.SAdd(ctx, "myset", fmt.Sprintf("member%d", i))
				Expect(sadd.Err()).NotTo(HaveOccurred())
			}

			keys, cursor, err := client.SScan(ctx, "myset", 0, "", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeZero())
		})

		It("should HScan", func() {
			for i := 0; i < 1000; i++ {
				sadd := client.HSet(ctx, "myhash", fmt.Sprintf("key%d", i), "hello")
				Expect(sadd.Err()).NotTo(HaveOccurred())
			}

			keys, cursor, err := client.HScan(ctx, "myhash", 0, "", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeZero())
		})

		It("should ZScan", func() {
			for i := 0; i < 1000; i++ {
				err := client.ZAdd(ctx, "myset", &redis.Z{
					Score:  float64(i),
					Member: fmt.Sprintf("member%d", i),
				}).Err()
				Expect(err).NotTo(HaveOccurred())
			}

			keys, cursor, err := client.ZScan(ctx, "myset", 0, "", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeZero())
		})
	})

	Describe("hashes", func() {
		It("should HDel", func() {
			hSet := client.HSet(ctx, "hash", "key", "hello")
			Expect(hSet.Err()).NotTo(HaveOccurred())

			hDel := client.HDel(ctx, "hash", "key")
			Expect(hDel.Err()).NotTo(HaveOccurred())
			Expect(hDel.Val()).To(Equal(int64(1)))

			hDel = client.HDel(ctx, "hash", "key")
			Expect(hDel.Err()).NotTo(HaveOccurred())
			Expect(hDel.Val()).To(Equal(int64(0)))
		})

		It("should HExists", func() {
			hSet := client.HSet(ctx, "hash", "key", "hello")
			Expect(hSet.Err()).NotTo(HaveOccurred())

			hExists := client.HExists(ctx, "hash", "key")
			Expect(hExists.Err()).NotTo(HaveOccurred())
			Expect(hExists.Val()).To(Equal(true))

			hExists = client.HExists(ctx, "hash", "key1")
			Expect(hExists.Err()).NotTo(HaveOccurred())
			Expect(hExists.Val()).To(Equal(false))
		})

		It("should HGet", func() {
			hSet := client.HSet(ctx, "hash", "key", "hello")
			Expect(hSet.Err()).NotTo(HaveOccurred())

			hGet := client.HGet(ctx, "hash", "key")
			Expect(hGet.Err()).NotTo(HaveOccurred())
			Expect(hGet.Val()).To(Equal("hello"))

			hGet = client.HGet(ctx, "hash", "key1")
			Expect(hGet.Err()).To(Equal(redis.Nil))
			Expect(hGet.Val()).To(Equal(""))
		})

		It("should HGetAll", func() {
			err := client.HSet(ctx, "hash", "key1", "hello1").Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.HSet(ctx, "hash", "key2", "hello2").Err()
			Expect(err).NotTo(HaveOccurred())

			m, err := client.HGetAll(ctx, "hash").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(m).To(Equal(map[string]string{"key1": "hello1", "key2": "hello2"}))
		})

		It("should scan", func() {
			err := client.HMSet(ctx, "hash", "key1", "hello1", "key2", 123).Err()
			Expect(err).NotTo(HaveOccurred())

			res := client.HGetAll(ctx, "hash")
			Expect(res.Err()).NotTo(HaveOccurred())

			type data struct {
				Key1 string `redis:"key1"`
				Key2 int    `redis:"key2"`
			}
			var d data
			Expect(res.Scan(&d)).NotTo(HaveOccurred())
			Expect(d).To(Equal(data{Key1: "hello1", Key2: 123}))
		})

		It("should HIncrBy", func() {
			hSet := client.HSet(ctx, "hash", "key", "5")
			Expect(hSet.Err()).NotTo(HaveOccurred())

			hIncrBy := client.HIncrBy(ctx, "hash", "key", 1)
			Expect(hIncrBy.Err()).NotTo(HaveOccurred())
			Expect(hIncrBy.Val()).To(Equal(int64(6)))

			hIncrBy = client.HIncrBy(ctx, "hash", "key", -1)
			Expect(hIncrBy.Err()).NotTo(HaveOccurred())
			Expect(hIncrBy.Val()).To(Equal(int64(5)))

			hIncrBy = client.HIncrBy(ctx, "hash", "key", -10)
			Expect(hIncrBy.Err()).NotTo(HaveOccurred())
			Expect(hIncrBy.Val()).To(Equal(int64(-5)))
		})

		It("should HIncrByFloat", func() {
			hSet := client.HSet(ctx, "hash", "field", "10.50")
			Expect(hSet.Err()).NotTo(HaveOccurred())
			Expect(hSet.Val()).To(Equal(int64(1)))

			hIncrByFloat := client.HIncrByFloat(ctx, "hash", "field", 0.1)
			Expect(hIncrByFloat.Err()).NotTo(HaveOccurred())
			Expect(hIncrByFloat.Val()).To(Equal(10.6))

			hSet = client.HSet(ctx, "hash", "field", "5.0e3")
			Expect(hSet.Err()).NotTo(HaveOccurred())
			Expect(hSet.Val()).To(Equal(int64(0)))

			hIncrByFloat = client.HIncrByFloat(ctx, "hash", "field", 2.0e2)
			Expect(hIncrByFloat.Err()).NotTo(HaveOccurred())
			Expect(hIncrByFloat.Val()).To(Equal(float64(5200)))
		})

		It("should HKeys", func() {
			hkeys := client.HKeys(ctx, "hash")
			Expect(hkeys.Err()).NotTo(HaveOccurred())
			Expect(hkeys.Val()).To(Equal([]string{}))

			hset := client.HSet(ctx, "hash", "key1", "hello1")
			Expect(hset.Err()).NotTo(HaveOccurred())
			hset = client.HSet(ctx, "hash", "key2", "hello2")
			Expect(hset.Err()).NotTo(HaveOccurred())

			hkeys = client.HKeys(ctx, "hash")
			Expect(hkeys.Err()).NotTo(HaveOccurred())
			Expect(hkeys.Val()).To(Equal([]string{"key1", "key2"}))
		})

		It("should HLen", func() {
			hSet := client.HSet(ctx, "hash", "key1", "hello1")
			Expect(hSet.Err()).NotTo(HaveOccurred())
			hSet = client.HSet(ctx, "hash", "key2", "hello2")
			Expect(hSet.Err()).NotTo(HaveOccurred())

			hLen := client.HLen(ctx, "hash")
			Expect(hLen.Err()).NotTo(HaveOccurred())
			Expect(hLen.Val()).To(Equal(int64(2)))
		})

		It("should HMGet", func() {
			err := client.HSet(ctx, "hash", "key1", "hello1", "key2", "hello2").Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := client.HMGet(ctx, "hash", "key1", "key2", "_").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]interface{}{"hello1", "hello2", nil}))
		})

		It("should HSet", func() {
			ok, err := client.HSet(ctx, "hash", map[string]interface{}{
				"key1": "hello1",
				"key2": "hello2",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(Equal(int64(2)))

			v, err := client.HGet(ctx, "hash", "key1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal("hello1"))

			v, err = client.HGet(ctx, "hash", "key2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal("hello2"))

			keys, err := client.HKeys(ctx, "hash").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).To(ConsistOf([]string{"key1", "key2"}))
		})

		It("should HSet", func() {
			hSet := client.HSet(ctx, "hash", "key", "hello")
			Expect(hSet.Err()).NotTo(HaveOccurred())
			Expect(hSet.Val()).To(Equal(int64(1)))

			hGet := client.HGet(ctx, "hash", "key")
			Expect(hGet.Err()).NotTo(HaveOccurred())
			Expect(hGet.Val()).To(Equal("hello"))
		})

		It("should HSetNX", func() {
			hSetNX := client.HSetNX(ctx, "hash", "key", "hello")
			Expect(hSetNX.Err()).NotTo(HaveOccurred())
			Expect(hSetNX.Val()).To(Equal(true))

			hSetNX = client.HSetNX(ctx, "hash", "key", "hello")
			Expect(hSetNX.Err()).NotTo(HaveOccurred())
			Expect(hSetNX.Val()).To(Equal(false))

			hGet := client.HGet(ctx, "hash", "key")
			Expect(hGet.Err()).NotTo(HaveOccurred())
			Expect(hGet.Val()).To(Equal("hello"))
		})

		It("should HVals", func() {
			err := client.HSet(ctx, "hash", "key1", "hello1").Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.HSet(ctx, "hash", "key2", "hello2").Err()
			Expect(err).NotTo(HaveOccurred())

			v, err := client.HVals(ctx, "hash").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal([]string{"hello1", "hello2"}))

			var slice []string
			err = client.HVals(ctx, "hash").ScanSlice(&slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(slice).To(Equal([]string{"hello1", "hello2"}))
		})

		It("[redis6] should HRandField", func() {
			err := client.HSet(ctx, "hash", "key1", "hello1").Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.HSet(ctx, "hash", "key2", "hello2").Err()
			Expect(err).NotTo(HaveOccurred())

			v := client.HRandField(ctx, "hash", 1, false)
			Expect(v.Err()).NotTo(HaveOccurred())
			Expect(v.Val()).To(Or(Equal([]string{"key1"}), Equal([]string{"key2"})))

			v = client.HRandField(ctx, "hash", 0, false)
			Expect(v.Err()).NotTo(HaveOccurred())
			Expect(v.Val()).To(HaveLen(0))

			var slice []string
			err = client.HRandField(ctx, "hash", 1, true).ScanSlice(&slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(slice).To(Or(Equal([]string{"key1", "hello1"}), Equal([]string{"key2", "hello2"})))
		})
	})

	Describe("hyperloglog", func() {
		It("should PFMerge", func() {
			pfAdd := client.PFAdd(ctx, "hll1", "1", "2", "3", "4", "5")
			Expect(pfAdd.Err()).NotTo(HaveOccurred())

			pfCount := client.PFCount(ctx, "hll1")
			Expect(pfCount.Err()).NotTo(HaveOccurred())
			Expect(pfCount.Val()).To(Equal(int64(5)))

			pfAdd = client.PFAdd(ctx, "hll2", "a", "b", "c", "d", "e")
			Expect(pfAdd.Err()).NotTo(HaveOccurred())

			pfMerge := client.PFMerge(ctx, "hllMerged", "hll1", "hll2")
			Expect(pfMerge.Err()).NotTo(HaveOccurred())

			pfCount = client.PFCount(ctx, "hllMerged")
			Expect(pfCount.Err()).NotTo(HaveOccurred())
			Expect(pfCount.Val()).To(Equal(int64(10)))

			pfCount = client.PFCount(ctx, "hll1", "hll2")
			Expect(pfCount.Err()).NotTo(HaveOccurred())
			Expect(pfCount.Val()).To(Equal(int64(10)))
		})
	})

	Describe("sets", func() {
		Describe("use same sadd", func() {
			BeforeEach(func() {
				sAdd := client.SAdd(ctx, "set1", "a")
				Expect(sAdd.Err()).NotTo(HaveOccurred())
				sAdd = client.SAdd(ctx, "set1", "b")
				Expect(sAdd.Err()).NotTo(HaveOccurred())
				sAdd = client.SAdd(ctx, "set1", "c")
				Expect(sAdd.Err()).NotTo(HaveOccurred())

				sAdd = client.SAdd(ctx, "set2", "c")
				Expect(sAdd.Err()).NotTo(HaveOccurred())
				sAdd = client.SAdd(ctx, "set2", "d")
				Expect(sAdd.Err()).NotTo(HaveOccurred())
				sAdd = client.SAdd(ctx, "set2", "e")
				Expect(sAdd.Err()).NotTo(HaveOccurred())
			})

			It("should SDiff", func() {
				sDiff := client.SDiff(ctx, "set1", "set2")
				Expect(sDiff.Err()).NotTo(HaveOccurred())
				Expect(sDiff.Val()).To(ConsistOf([]string{"a", "b"}))
			})

			It("should SDiffStore", func() {
				sDiffStore := client.SDiffStore(ctx, "set", "set1", "set2")
				Expect(sDiffStore.Err()).NotTo(HaveOccurred())
				Expect(sDiffStore.Val()).To(Equal(int64(2)))

				sMembers := client.SMembers(ctx, "set")
				Expect(sMembers.Err()).NotTo(HaveOccurred())
				Expect(sMembers.Val()).To(ConsistOf([]string{"a", "b"}))
			})

			It("should SInter", func() {
				sInter := client.SInter(ctx, "set1", "set2")
				Expect(sInter.Err()).NotTo(HaveOccurred())
				Expect(sInter.Val()).To(Equal([]string{"c"}))
			})

			It("should SInterStore", func() {
				sInterStore := client.SInterStore(ctx, "set", "set1", "set2")
				Expect(sInterStore.Err()).NotTo(HaveOccurred())
				Expect(sInterStore.Val()).To(Equal(int64(1)))

				sMembers := client.SMembers(ctx, "set")
				Expect(sMembers.Err()).NotTo(HaveOccurred())
				Expect(sMembers.Val()).To(Equal([]string{"c"}))
			})

			It("should SUnion", func() {
				sUnion := client.SUnion(ctx, "set1", "set2")
				Expect(sUnion.Err()).NotTo(HaveOccurred())
				Expect(sUnion.Val()).To(HaveLen(5))
			})

			It("should SUnionStore", func() {
				sUnionStore := client.SUnionStore(ctx, "set", "set1", "set2")
				Expect(sUnionStore.Err()).NotTo(HaveOccurred())
				Expect(sUnionStore.Val()).To(Equal(int64(5)))

				sMembers := client.SMembers(ctx, "set")
				Expect(sMembers.Err()).NotTo(HaveOccurred())
				Expect(sMembers.Val()).To(HaveLen(5))
			})
		})

		It("should SAdd", func() {
			sAdd := client.SAdd(ctx, "set", "Hello")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(1)))

			sAdd = client.SAdd(ctx, "set", "World")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(1)))

			sAdd = client.SAdd(ctx, "set", "World")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(0)))

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
		})

		It("should SAdd strings", func() {
			set := []string{"Hello", "World", "World"}
			sAdd := client.SAdd(ctx, "set", set)
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(2)))

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
		})

		It("should SCard", func() {
			sAdd := client.SAdd(ctx, "set", "Hello")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(1)))

			sAdd = client.SAdd(ctx, "set", "World")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			Expect(sAdd.Val()).To(Equal(int64(1)))

			sCard := client.SCard(ctx, "set")
			Expect(sCard.Err()).NotTo(HaveOccurred())
			Expect(sCard.Val()).To(Equal(int64(2)))
		})

		It("should IsMember", func() {
			sAdd := client.SAdd(ctx, "set", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sIsMember := client.SIsMember(ctx, "set", "one")
			Expect(sIsMember.Err()).NotTo(HaveOccurred())
			Expect(sIsMember.Val()).To(Equal(true))

			sIsMember = client.SIsMember(ctx, "set", "two")
			Expect(sIsMember.Err()).NotTo(HaveOccurred())
			Expect(sIsMember.Val()).To(Equal(false))
		})

		It("[redis6] should SMIsMember", func() {
			sAdd := client.SAdd(ctx, "set", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sMIsMember := client.SMIsMember(ctx, "set", "one", "two")
			Expect(sMIsMember.Err()).NotTo(HaveOccurred())
			Expect(sMIsMember.Val()).To(Equal([]bool{true, false}))
		})

		It("should SMembers", func() {
			sAdd := client.SAdd(ctx, "set", "Hello")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "World")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
		})

		It("should SMembersMap", func() {
			sAdd := client.SAdd(ctx, "set", "Hello")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "World")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sMembersMap := client.SMembersMap(ctx, "set")
			Expect(sMembersMap.Err()).NotTo(HaveOccurred())
			Expect(sMembersMap.Val()).To(Equal(map[string]struct{}{"Hello": {}, "World": {}}))
		})

		It("should SMove", func() {
			sAdd := client.SAdd(ctx, "set1", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set1", "two")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sAdd = client.SAdd(ctx, "set2", "three")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sMove := client.SMove(ctx, "set1", "set2", "two")
			Expect(sMove.Err()).NotTo(HaveOccurred())
			Expect(sMove.Val()).To(Equal(true))

			sMembers := client.SMembers(ctx, "set1")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(Equal([]string{"one"}))

			sMembers = client.SMembers(ctx, "set2")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(ConsistOf([]string{"three", "two"}))
		})

		It("should SPop", func() {
			sAdd := client.SAdd(ctx, "set", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "two")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "three")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sPop := client.SPop(ctx, "set")
			Expect(sPop.Err()).NotTo(HaveOccurred())
			Expect(sPop.Val()).NotTo(Equal(""))

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(HaveLen(2))
		})

		It("should SPopN", func() {
			sAdd := client.SAdd(ctx, "set", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "two")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "three")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "four")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sPopN := client.SPopN(ctx, "set", 1)
			Expect(sPopN.Err()).NotTo(HaveOccurred())
			Expect(sPopN.Val()).NotTo(Equal([]string{""}))

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(HaveLen(3))

			sPopN = client.SPopN(ctx, "set", 4)
			Expect(sPopN.Err()).NotTo(HaveOccurred())
			Expect(sPopN.Val()).To(HaveLen(3))

			sMembers = client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(HaveLen(0))
		})

		It("should SRandMember and SRandMemberN", func() {
			err := client.SAdd(ctx, "set", "one").Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.SAdd(ctx, "set", "two").Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.SAdd(ctx, "set", "three").Err()
			Expect(err).NotTo(HaveOccurred())

			members, err := client.SMembers(ctx, "set").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(members).To(HaveLen(3))

			member, err := client.SRandMember(ctx, "set").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(member).NotTo(Equal(""))

			members, err = client.SRandMemberN(ctx, "set", 2).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(members).To(HaveLen(2))
		})

		It("should SRem", func() {
			sAdd := client.SAdd(ctx, "set", "one")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "two")
			Expect(sAdd.Err()).NotTo(HaveOccurred())
			sAdd = client.SAdd(ctx, "set", "three")
			Expect(sAdd.Err()).NotTo(HaveOccurred())

			sRem := client.SRem(ctx, "set", "one")
			Expect(sRem.Err()).NotTo(HaveOccurred())
			Expect(sRem.Val()).To(Equal(int64(1)))

			sRem = client.SRem(ctx, "set", "four")
			Expect(sRem.Err()).NotTo(HaveOccurred())
			Expect(sRem.Val()).To(Equal(int64(0)))

			sMembers := client.SMembers(ctx, "set")
			Expect(sMembers.Err()).NotTo(HaveOccurred())
			Expect(sMembers.Val()).To(ConsistOf([]string{"three", "two"}))
		})
	})

	Describe("streams", func() {
		BeforeEach(func() {
			id, err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				ID:     "1-0",
				Values: map[string]interface{}{"uno": "un"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal("1-0"))

			// Values supports []interface{}.
			id, err = client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				ID:     "2-0",
				Values: []interface{}{"dos", "deux"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal("2-0"))

			// Value supports []string.
			id, err = client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				ID:     "3-0",
				Values: []string{"tres", "troix"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal("3-0"))
		})

		It("should XTrim", func() {
			n, err := client.XTrim(ctx, "stream", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("should XTrimApprox", func() {
			n, err := client.XTrimApprox(ctx, "stream", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		// TODO XTrimMaxLenApprox/XTrimMinIDApprox There is a bug in the limit parameter.
		// TODO Don't test it for now.
		// TODO link: https://github.com/redis/redis/issues/9046
		It("should XTrimMaxLen", func() {
			n, err := client.XTrimMaxLen(ctx, "stream", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("should XTrimMaxLenApprox", func() {
			n, err := client.XTrimMaxLenApprox(ctx, "stream", 0, 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("[redis6] should XTrimMinID", func() {
			n, err := client.XTrimMinID(ctx, "stream", "4-0").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("[redis6] should XTrimMinIDApprox", func() {
			n, err := client.XTrimMinIDApprox(ctx, "stream", "4-0", 0).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("should XAdd", func() {
			id, err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				Values: map[string]interface{}{"quatro": "quatre"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())

			vals, err := client.XRange(ctx, "stream", "-", "+").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				{ID: id, Values: map[string]interface{}{"quatro": "quatre"}},
			}))
		})

		// TODO XAdd There is a bug in the limit parameter.
		// TODO Don't test it for now.
		// TODO link: https://github.com/redis/redis/issues/9046
		It("should XAdd with MaxLen", func() {
			id, err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				MaxLen: 1,
				Values: map[string]interface{}{"quatro": "quatre"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())

			vals, err := client.XRange(ctx, "stream", "-", "+").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.XMessage{
				{ID: id, Values: map[string]interface{}{"quatro": "quatre"}},
			}))
		})

		It("[redis6] should XAdd with MinID", func() {
			id, err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream",
				MinID:  "5-0",
				ID:     "4-0",
				Values: map[string]interface{}{"quatro": "quatre"},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal("4-0"))

			vals, err := client.XRange(ctx, "stream", "-", "+").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(HaveLen(0))
		})

		It("should XDel", func() {
			n, err := client.XDel(ctx, "stream", "1-0", "2-0", "3-0").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("should XLen", func() {
			n, err := client.XLen(ctx, "stream").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(3)))
		})

		It("should XRange", func() {
			msgs, err := client.XRange(ctx, "stream", "-", "+").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
			}))

			msgs, err = client.XRange(ctx, "stream", "2", "+").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
			}))

			msgs, err = client.XRange(ctx, "stream", "-", "2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
			}))
		})

		It("should XRangeN", func() {
			msgs, err := client.XRangeN(ctx, "stream", "-", "+", 2).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
			}))

			msgs, err = client.XRangeN(ctx, "stream", "2", "+", 1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
			}))

			msgs, err = client.XRangeN(ctx, "stream", "-", "2", 1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
			}))
		})

		It("should XRevRange", func() {
			msgs, err := client.XRevRange(ctx, "stream", "+", "-").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
				{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
			}))

			msgs, err = client.XRevRange(ctx, "stream", "+", "2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
			}))
		})

		It("should XRevRangeN", func() {
			msgs, err := client.XRevRangeN(ctx, "stream", "+", "-", 2).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
			}))

			msgs, err = client.XRevRangeN(ctx, "stream", "+", "2", 1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(msgs).To(Equal([]redis.XMessage{
				{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
			}))
		})

		var streamRes = []redis.XStream{
			{
				Stream: "stream",
				Messages: []redis.XMessage{
					{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
					{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
					{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				},
			},
		}

		It("should XRead", func() {
			res, err := client.XReadStreams(ctx, "stream", "0").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(streamRes))

			_, err = client.XReadStreams(ctx, "stream", "3").Result()
			Expect(err).To(Equal(redis.Nil))
		})

		var streamRes2 = []redis.XStream{
			{
				Stream: "stream",
				Messages: []redis.XMessage{
					{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
					{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
				},
			},
		}

		It("should XRead", func() {
			res, err := client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{"stream", "0"},
				Count:   2,
				Block:   100 * time.Millisecond,
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(streamRes2))

			_, err = client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{"stream", "3"},
				Count:   1,
				Block:   100 * time.Millisecond,
			}).Result()
			Expect(err).To(Equal(redis.Nil))
		})

		Describe("group", func() {
			BeforeEach(func() {
				err := client.XGroupCreate(ctx, "stream", "group", "0").Err()
				Expect(err).NotTo(HaveOccurred())

				res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "group",
					Consumer: "consumer",
					Streams:  []string{"stream", ">"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(streamRes))
			})

			AfterEach(func() {
				n, err := client.XGroupDestroy(ctx, "stream", "group").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))
			})

			It("should XReadGroup skip empty", func() {
				n, err := client.XDel(ctx, "stream", "2-0").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))

				res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "group",
					Consumer: "consumer",
					Streams:  []string{"stream", "0"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal([]redis.XStream{
					{
						Stream: "stream",
						Messages: []redis.XMessage{
							{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
							{ID: "2-0", Values: nil},
							{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
						},
					},
				}))
			})

			It("should XGroupCreateMkStream", func() {
				err := client.XGroupCreateMkStream(ctx, "stream2", "group", "0").Err()
				Expect(err).NotTo(HaveOccurred())

				err = client.XGroupCreateMkStream(ctx, "stream2", "group", "0").Err()
				Expect(err.Error()).To(Equal("BUSYGROUP Consumer Group name already exists"))

				n, err := client.XGroupDestroy(ctx, "stream2", "group").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))

				n, err = client.Del(ctx, "stream2").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))
			})

			It("[redis6] should XPending", func() {
				info, err := client.XPending(ctx, "stream", "group").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(info).To(Equal(&redis.XPending{
					Count:     3,
					Lower:     "1-0",
					Higher:    "3-0",
					Consumers: map[string]int64{"consumer": 3},
				}))
				args := &redis.XPendingExtArgs{
					Stream:   "stream",
					Group:    "group",
					Start:    "-",
					End:      "+",
					Count:    10,
					Consumer: "consumer",
				}
				infoExt, err := client.XPendingExt(ctx, args).Result()
				Expect(err).NotTo(HaveOccurred())
				for i := range infoExt {
					infoExt[i].Idle = 0
				}
				Expect(infoExt).To(Equal([]redis.XPendingExt{
					{ID: "1-0", Consumer: "consumer", Idle: 0, RetryCount: 1},
					{ID: "2-0", Consumer: "consumer", Idle: 0, RetryCount: 1},
					{ID: "3-0", Consumer: "consumer", Idle: 0, RetryCount: 1},
				}))

				args.Idle = 72 * time.Hour
				infoExt, err = client.XPendingExt(ctx, args).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(infoExt).To(HaveLen(0))
			})

			It("[redis6] should XGroup Create Delete Consumer", func() {
				n, err := client.XGroupCreateConsumer(ctx, "stream", "group", "c1").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))

				n, err = client.XGroupDelConsumer(ctx, "stream", "group", "consumer").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(3)))
			})

			It("[redis6] should XAutoClaim", func() {
				xca := &redis.XAutoClaimArgs{
					Stream:   "stream",
					Group:    "group",
					Consumer: "consumer",
					Start:    "-",
					Count:    2,
				}
				msgs, start, err := client.XAutoClaim(ctx, xca).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(start).To(Equal("3-0"))
				Expect(msgs).To(Equal([]redis.XMessage{{
					ID:     "1-0",
					Values: map[string]interface{}{"uno": "un"},
				}, {
					ID:     "2-0",
					Values: map[string]interface{}{"dos": "deux"},
				}}))

				xca.Start = start
				msgs, start, err = client.XAutoClaim(ctx, xca).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(start).To(Equal("0-0"))
				Expect(msgs).To(Equal([]redis.XMessage{{
					ID:     "3-0",
					Values: map[string]interface{}{"tres": "troix"},
				}}))

				ids, start, err := client.XAutoClaimJustID(ctx, xca).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(start).To(Equal("0-0"))
				Expect(ids).To(Equal([]string{"3-0"}))
			})

			It("should XClaim", func() {
				msgs, err := client.XClaim(ctx, &redis.XClaimArgs{
					Stream:   "stream",
					Group:    "group",
					Consumer: "consumer",
					Messages: []string{"1-0", "2-0", "3-0"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(msgs).To(Equal([]redis.XMessage{{
					ID:     "1-0",
					Values: map[string]interface{}{"uno": "un"},
				}, {
					ID:     "2-0",
					Values: map[string]interface{}{"dos": "deux"},
				}, {
					ID:     "3-0",
					Values: map[string]interface{}{"tres": "troix"},
				}}))

				ids, err := client.XClaimJustID(ctx, &redis.XClaimArgs{
					Stream:   "stream",
					Group:    "group",
					Consumer: "consumer",
					Messages: []string{"1-0", "2-0", "3-0"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(ids).To(Equal([]string{"1-0", "2-0", "3-0"}))
			})

			It("should XAck", func() {
				n, err := client.XAck(ctx, "stream", "group", "1-0", "2-0", "4-0").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(2)))
			})
		})

		Describe("xinfo", func() {
			BeforeEach(func() {
				err := client.XGroupCreate(ctx, "stream", "group1", "0").Err()
				Expect(err).NotTo(HaveOccurred())

				res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "group1",
					Consumer: "consumer1",
					Streams:  []string{"stream", ">"},
					Count:    2,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(streamRes2))

				res, err = client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "group1",
					Consumer: "consumer2",
					Streams:  []string{"stream", ">"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal([]redis.XStream{
					{
						Stream: "stream",
						Messages: []redis.XMessage{
							{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
						},
					},
				}))

				err = client.XGroupCreate(ctx, "stream", "group2", "1-0").Err()
				Expect(err).NotTo(HaveOccurred())

				res, err = client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "group2",
					Consumer: "consumer1",
					Streams:  []string{"stream", ">"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal([]redis.XStream{
					{
						Stream: "stream",
						Messages: []redis.XMessage{
							{ID: "2-0", Values: map[string]interface{}{"dos": "deux"}},
							{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
						},
					},
				}))
			})

			AfterEach(func() {
				n, err := client.XGroupDestroy(ctx, "stream", "group1").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))
				n, err = client.XGroupDestroy(ctx, "stream", "group2").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(1)))
			})

			It("should XINFO STREAM", func() {
				res, err := client.XInfoStream(ctx, "stream").Result()
				Expect(err).NotTo(HaveOccurred())
				res.RadixTreeKeys = 0
				res.RadixTreeNodes = 0

				Expect(res).To(Equal(&redis.XInfoStream{
					Length:          3,
					RadixTreeKeys:   0,
					RadixTreeNodes:  0,
					Groups:          2,
					LastGeneratedID: "3-0",
					FirstEntry:      redis.XMessage{ID: "1-0", Values: map[string]interface{}{"uno": "un"}},
					LastEntry:       redis.XMessage{ID: "3-0", Values: map[string]interface{}{"tres": "troix"}},
				}))

				// stream is empty
				n, err := client.XDel(ctx, "stream", "1-0", "2-0", "3-0").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(3)))

				res, err = client.XInfoStream(ctx, "stream").Result()
				Expect(err).NotTo(HaveOccurred())
				res.RadixTreeKeys = 0
				res.RadixTreeNodes = 0

				Expect(res).To(Equal(&redis.XInfoStream{
					Length:          0,
					RadixTreeKeys:   0,
					RadixTreeNodes:  0,
					Groups:          2,
					LastGeneratedID: "3-0",
					FirstEntry:      redis.XMessage{},
					LastEntry:       redis.XMessage{},
				}))
			})

			It("should XINFO GROUPS", func() {
				res, err := client.XInfoGroups(ctx, "stream").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal([]redis.XInfoGroup{
					{Name: "group1", Consumers: 2, Pending: 3, LastDeliveredID: "3-0"},
					{Name: "group2", Consumers: 1, Pending: 2, LastDeliveredID: "3-0"},
				}))
			})

			It("should XINFO CONSUMERS", func() {
				res, err := client.XInfoConsumers(ctx, "stream", "group1").Result()
				Expect(err).NotTo(HaveOccurred())
				for i := range res {
					res[i].Idle = 0
				}
				Expect(res).To(Equal([]redis.XInfoConsumer{
					{Name: "consumer1", Pending: 2, Idle: 0},
					{Name: "consumer2", Pending: 1, Idle: 0},
				}))
			})
		})
	})

	Describe("Geo add and radius search", func() {
		BeforeEach(func() {
			n, err := client.GeoAdd(
				ctx,
				"Sicily",
				&redis.GeoLocation{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
				&redis.GeoLocation{Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"},
			).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(2)))
		})

		It("should not add same geo location", func() {
			geoAdd := client.GeoAdd(
				ctx,
				"Sicily",
				&redis.GeoLocation{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
			)
			Expect(geoAdd.Err()).NotTo(HaveOccurred())
			Expect(geoAdd.Val()).To(Equal(int64(0)))
		})

		It("should search geo radius", func() {
			res, err := client.GeoRadius(ctx, "Sicily", 15, 37, &redis.GeoRadiusQuery{
				Radius: 200,
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(HaveLen(2))
			Expect(res[0].Name).To(Equal("Palermo"))
			Expect(res[1].Name).To(Equal("Catania"))
		})

		It("should geo radius and store the result", func() {
			n, err := client.GeoRadiusStore(ctx, "Sicily", 15, 37, &redis.GeoRadiusQuery{
				Radius: 200,
				Store:  "result",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(2)))

			res, err := client.ZRangeWithScores(ctx, "result", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(ContainElement(redis.Z{
				Score:  3.479099956230698e+15,
				Member: "Palermo",
			}))
			Expect(res).To(ContainElement(redis.Z{
				Score:  3.479447370796909e+15,
				Member: "Catania",
			}))
		})

		It("should geo radius and store dist", func() {
			n, err := client.GeoRadiusStore(ctx, "Sicily", 15, 37, &redis.GeoRadiusQuery{
				Radius:    200,
				StoreDist: "result",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(2)))

			res, err := client.ZRangeWithScores(ctx, "result", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(ContainElement(redis.Z{
				Score:  190.44242984775784,
				Member: "Palermo",
			}))
			Expect(res).To(ContainElement(redis.Z{
				Score:  56.4412578701582,
				Member: "Catania",
			}))
		})

		var geoRadius = &redis.GeoRadiusQuery{
			Radius:      200,
			Unit:        "km",
			WithGeoHash: true,
			WithCoord:   true,
			WithDist:    true,
			Count:       2,
			Sort:        "ASC",
		}

		It("should search geo radius with options", func() {
			res, err := client.GeoRadius(ctx, "Sicily", 15, 37, geoRadius).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(HaveLen(2))
			Expect(res[1].Name).To(Equal("Palermo"))
			Expect(res[1].Dist).To(Equal(190.4424))
			Expect(res[1].GeoHash).To(Equal(int64(3479099956230698)))
			Expect(res[1].Longitude).To(Equal(13.361389338970184))
			Expect(res[1].Latitude).To(Equal(38.115556395496299))
			Expect(res[0].Name).To(Equal("Catania"))
			Expect(res[0].Dist).To(Equal(56.4413))
			Expect(res[0].GeoHash).To(Equal(int64(3479447370796909)))
			Expect(res[0].Longitude).To(Equal(15.087267458438873))
			Expect(res[0].Latitude).To(Equal(37.50266842333162))
		})

		It("should search geo radius with WithDist=false", func() {
			res, err := client.GeoRadius(ctx, "Sicily", 15, 37, &redis.GeoRadiusQuery{
				Radius:      200,
				Unit:        "km",
				WithGeoHash: true,
				WithCoord:   true,
				Count:       2,
				Sort:        "ASC",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(HaveLen(2))
			Expect(res[1].Name).To(Equal("Palermo"))
			Expect(res[1].Dist).To(Equal(float64(0)))
			Expect(res[1].GeoHash).To(Equal(int64(3479099956230698)))
			Expect(res[1].Longitude).To(Equal(13.361389338970184))
			Expect(res[1].Latitude).To(Equal(38.115556395496299))
			Expect(res[0].Name).To(Equal("Catania"))
			Expect(res[0].Dist).To(Equal(float64(0)))
			Expect(res[0].GeoHash).To(Equal(int64(3479447370796909)))
			Expect(res[0].Longitude).To(Equal(15.087267458438873))
			Expect(res[0].Latitude).To(Equal(37.50266842333162))
		})

		It("should search geo radius by member with options", func() {
			res, err := client.GeoRadiusByMember(ctx, "Sicily", "Catania", geoRadius).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(HaveLen(2))
			Expect(res[0].Name).To(Equal("Catania"))
			Expect(res[0].Dist).To(Equal(0.0))
			Expect(res[0].GeoHash).To(Equal(int64(3479447370796909)))
			Expect(res[0].Longitude).To(Equal(15.087267458438873))
			Expect(res[0].Latitude).To(Equal(37.50266842333162))
			Expect(res[1].Name).To(Equal("Palermo"))
			Expect(res[1].Dist).To(Equal(166.2742))
			Expect(res[1].GeoHash).To(Equal(int64(3479099956230698)))
			Expect(res[1].Longitude).To(Equal(13.361389338970184))
			Expect(res[1].Latitude).To(Equal(38.115556395496299))
		})

		It("should search geo radius with no results", func() {
			res, err := client.GeoRadius(ctx, "Sicily", 99, 37, &redis.GeoRadiusQuery{
				Radius:      200,
				Unit:        "km",
				WithGeoHash: true,
				WithCoord:   true,
				WithDist:    true,
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(HaveLen(0))
		})

		It("should get geo distance with unit options", func() {
			// From Redis CLI, note the difference in rounding in m vs
			// km on Redis itself.
			//
			// GEOADD Sicily 13.361389 38.115556 "Palermo" 15.087269 37.502669 "Catania"
			// GEODIST Sicily Palermo Catania m
			// "166274.15156960033"
			// GEODIST Sicily Palermo Catania km
			// "166.27415156960032"
			dist, err := client.GeoDist(ctx, "Sicily", "Palermo", "Catania", "km").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(dist).To(BeNumerically("~", 166.27, 0.01))

			dist, err = client.GeoDist(ctx, "Sicily", "Palermo", "Catania", "m").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(dist).To(BeNumerically("~", 166274.15, 0.01))
		})

		It("should get geo hash in string representation", func() {
			hashes, err := client.GeoHash(ctx, "Sicily", "Palermo", "Catania").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(hashes).To(ConsistOf([]string{"sqc8b49rny0", "sqdtr74hyu0"}))
		})

		It("should return geo position", func() {
			pos, err := client.GeoPos(ctx, "Sicily", "Palermo", "Catania", "NonExisting").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(pos).To(ConsistOf([]*redis.GeoPos{
				{
					Longitude: 13.361389338970184,
					Latitude:  38.1155563954963,
				},
				{
					Longitude: 15.087267458438873,
					Latitude:  37.50266842333162,
				},
				nil,
			}))
		})

		var geoLocation = []redis.GeoLocation{
			{
				Name:      "Catania",
				Longitude: 15.08726745843887329,
				Latitude:  37.50266842333162032,
				Dist:      56.4413,
				GeoHash:   3479447370796909,
			},
			{
				Name:      "Palermo",
				Longitude: 13.36138933897018433,
				Latitude:  38.11555639549629859,
				Dist:      190.4424,
				GeoHash:   3479099956230698,
			},
		}

		It("[redis6] should geo search with options", func() {
			q := &redis.GeoSearchLocationQuery{
				GeoSearchQuery: redis.GeoSearchQuery{
					Longitude:  15,
					Latitude:   37,
					Radius:     200,
					RadiusUnit: "km",
					Sort:       "asc",
				},
				WithHash:  true,
				WithDist:  true,
				WithCoord: true,
			}
			val, err := client.GeoSearchLocation(ctx, "Sicily", q).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(geoLocation))
		})

		It("[redis6] should geo search store", func() {
			q := &redis.GeoSearchStoreQuery{
				GeoSearchQuery: redis.GeoSearchQuery{
					Longitude:  15,
					Latitude:   37,
					Radius:     200,
					RadiusUnit: "km",
					Sort:       "asc",
				},
				StoreDist: false,
			}

			val, err := client.GeoSearchStore(ctx, "Sicily", "key1", q).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(int64(2)))

			q.StoreDist = true
			val, err = client.GeoSearchStore(ctx, "Sicily", "key2", q).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(int64(2)))

			loc, err := client.GeoSearchLocation(ctx, "key1", &redis.GeoSearchLocationQuery{
				GeoSearchQuery: q.GeoSearchQuery,
				WithCoord:      true,
				WithDist:       true,
				WithHash:       true,
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(loc).To(Equal(geoLocation))

			v, err := client.ZRangeWithScores(ctx, "key2", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal([]redis.Z{
				{
					Score:  56.441257870158204,
					Member: "Catania",
				},
				{
					Score:  190.44242984775784,
					Member: "Palermo",
				},
			}))
		})
	})

	Describe("Eval", func() {
		It("returns keys and values", func() {
			vals, err := client.Eval(
				ctx,
				"return {KEYS[1],ARGV[1]}",
				[]string{"key"},
				"hello",
			).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]interface{}{"key", "hello"}))
		})

		It("returns all values after an error", func() {
			vals, err := client.Eval(
				ctx,
				`return {12, {err="error"}, "abc"}`,
				nil,
			).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(vals.([]interface{}))).To(Equal(3))
			Expect(vals.([]interface{})[0]).To(Equal(int64(12)))
			Expect(vals.([]interface{})[2]).To(Equal("abc"))
		})
	})

	Describe("SlowLogGet", func() {
		It("returns slow query result", func() {
			const key = "slowlog-log-slower-than"

			old := client.ConfigGet(ctx, key).Val()
			client.ConfigSet(ctx, key, "0")
			defer client.ConfigSet(ctx, key, old[1].(string))

			err := client.Do(ctx, "slowlog", "reset").Err()
			Expect(err).NotTo(HaveOccurred())

			client.Set(ctx, "test", "true", 0)

			result, err := client.SlowLogGet(ctx, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).NotTo(BeZero())
		})
	})
})
