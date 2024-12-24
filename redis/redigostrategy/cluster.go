package redigostrategy

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/huaweicloud/devcloud-go/redis/config"
)

type RedigoClusterConfig struct {
	MaxIdle   int
	MaxActive int
	Password  string
	Timeout   time.Duration
}

func newRedigoClusterConfig(serverConfig *config.ServerConfiguration) *RedigoClusterConfig {
	return &RedigoClusterConfig{
		MaxIdle:   serverConfig.ConnectionPool.MaxIdle,
		MaxActive: serverConfig.ConnectionPool.MaxTotal,
		Password:  serverConfig.Password,
		Timeout:   serverConfig.ClusterOptions.DialTimeout,
	}
}

func (r RedigoClusterConfig) createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxActive,
		Wait:        true,
		IdleTimeout: time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr,
				redis.DialPassword(r.Password),
				redis.DialConnectTimeout(time.Duration(r.Timeout)),
				redis.DialWriteTimeout(time.Duration(r.Timeout)),
				redis.DialReadTimeout(time.Duration(r.Timeout)))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func createPool(r *RedigoClusterConfig) func(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return func(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
		return r.createPool(addr, opts...)
	}
}
