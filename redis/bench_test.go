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
 * 2021.11.15-Changed modify the constructor Client
 * 			Huawei Technologies Co., Ltd.
 */

package redis

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

func benchmarkDevsporeRedisClient(ctx context.Context, poolSize int) *DevsporeClient {
	configuration := benchmarkDevsporeClientConfiguration(config.ServerTypeNormal, poolSize)
	client := NewDevsporeClient(configuration)
	if err := client.FlushDB(ctx).Err(); err != nil {
		panic(err)
	}
	return client
}

func benchmarkDevsporeClusterClient(ctx context.Context) *DevsporeClient {

	configuration := benchmarkDevsporeClientConfiguration(config.ServerTypeCluster, 0)
	client := NewDevsporeClient(configuration)
	if err := client.FlushDB(ctx).Err(); err != nil {
		panic(err)
	}
	return client
}

func benchmarkRedisOptions(poolSize int) *redis.Options {
	return &redis.Options{
		Addr:         ":6379",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     poolSize,
	}
}

func benchmarkCLusterOptions() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:        []string{":6383", ":6384", ":6385", ":6386", ":6387", ":6388"},
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,

		MaxRedirects: 8,

		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func benchmarkDevsporeClientConfiguration(serverType string, poolSize int) *config.Configuration {
	configuration := &config.Configuration{
		RedisConfig:    &config.RedisConfiguration{},
		RouteAlgorithm: "single-read-write",
		Active:         "server1",
	}
	if serverType == config.ServerTypeCluster {
		servers := map[string]*config.ServerConfiguration{
			"server1": {
				Type:           config.ServerTypeCluster,
				Cloud:          "huawei cloud",
				Region:         "beijing",
				Azs:            "az0",
				ClusterOptions: benchmarkCLusterOptions(),
			},
		}
		configuration.RedisConfig.Servers = servers
	} else {
		servers := map[string]*config.ServerConfiguration{
			"server1": {
				Type:    config.ServerTypeNormal,
				Cloud:   "huawei cloud",
				Region:  "beijing",
				Azs:     "az0",
				Options: benchmarkRedisOptions(poolSize),
			},
		}
		configuration.RedisConfig.Servers = servers
	}
	return configuration
}

func BenchmarkRedisPing(b *testing.B) {
	ctx := context.Background()
	rdb := benchmarkDevsporeRedisClient(ctx, 10)
	defer rdb.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := rdb.Ping(ctx).Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSetGoroutines(b *testing.B) {
	ctx := context.Background()
	rdb := benchmarkDevsporeRedisClient(ctx, 10)
	defer rdb.Close()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := rdb.Set(ctx, "hello", "world", 0).Err()
				if err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSet(b *testing.B) {
	ctx := context.Background()
	rdb := benchmarkDevsporeRedisClient(ctx, 10)
	defer rdb.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := rdb.Set(ctx, "hello", "test", 0).Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRedisGetNil(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.Get(ctx, "key").Err(); err != redis.Nil {
				b.Fatal(err)
			}
		}
	})
}

type setStringBenchmark struct {
	poolSize  int
	valueSize int
}

func (bm setStringBenchmark) String() string {
	return fmt.Sprintf("pool=%d value=%d", bm.poolSize, bm.valueSize)
}

func BenchmarkRedisSetString(b *testing.B) {
	benchmarks := []setStringBenchmark{
		{10, 64},
		{10, 1024},
		{10, 64 * 1024},
		{10, 1024 * 1024},
		{10, 10 * 1024 * 1024},

		{100, 64},
		{100, 1024},
		{100, 64 * 1024},
		{100, 1024 * 1024},
		{100, 10 * 1024 * 1024},
	}
	for _, bm := range benchmarks {
		b.Run(bm.String(), func(b *testing.B) {
			ctx := context.Background()
			client := benchmarkDevsporeRedisClient(ctx, bm.poolSize)
			defer client.Close()

			value := strings.Repeat("1", bm.valueSize)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					err := client.Set(ctx, "key", value, 0).Err()
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkRedisSetGetBytes(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	value := bytes.Repeat([]byte{'1'}, 10000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.Set(ctx, "key", value, 0).Err(); err != nil {
				b.Fatal(err)
			}

			got, err := client.Get(ctx, "key").Bytes()
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(got, value) {
				b.Fatalf("got != value")
			}
		}
	})
}

func BenchmarkRedisMGet(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	if err := client.MSet(ctx, "key1", "hello1", "key2", "hello2").Err(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.MGet(ctx, "key1", "key2").Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSetExpire(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.Set(ctx, "key", "hello", 0).Err(); err != nil {
				b.Fatal(err)
			}
			if err := client.Expire(ctx, "key", time.Second).Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPipeline(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, "key", "hello", 0)
				pipe.Expire(ctx, "key", time.Second)
				return nil
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkZAdd(b *testing.B) {
	ctx := context.Background()
	client := benchmarkDevsporeRedisClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := client.ZAdd(ctx, "key", &redis.Z{
				Score:  float64(1),
				Member: "hello",
			}).Err()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

var ringSink *redis.Ring

func BenchmarkRingWithContext(b *testing.B) {
	ctx := context.Background()
	rdb := redis.NewRing(&redis.RingOptions{})
	defer rdb.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ringSink = rdb.WithContext(ctx)
	}
}

func BenchmarkClusterPing(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	ctx := context.Background()
	client := benchmarkDevsporeClusterClient(ctx)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := client.Ping(ctx).Err()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkClusterSetString(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	ctx := context.Background()
	client := benchmarkDevsporeClusterClient(ctx)
	defer client.Close()

	value := string(bytes.Repeat([]byte{'1'}, 10000))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := client.Set(ctx, "key", value, 0).Err()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
