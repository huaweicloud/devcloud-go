# devcloud-go/redis

### Quickstart：
1. use yaml configuartion file
```bigquery
package main

import (
    "context"
    "time"

    "github.com/huaweicloud/devcloud-go/redis"
)

func main()  {
    ctx := context.Background()
    client := redis.NewDevsporeClientWithYaml("./config_with_password.yaml")
    client.Set(ctx, "test_key", "test_val", time.Hour)
    client.Get(ctx, "test_key")
}
```
2. use code(recommend)
```bigquery
package main

import (
    "context"
    "time"

    goredis "github.com/go-redis/redis/v8"
    "github.com/huaweicloud/devcloud-go/redis"
    "github.com/huaweicloud/devcloud-go/redis/config"
)
func main() {
    servers := map[string]*config.ServerConfiguration{
        "server1": {
            Type:   config.ServerTypeNormal,
            Cloud:  "huawei cloud",
            Region: "beijing",
            Azs:    "az0",
            Options: &goredis.Options{
                Addr:     "127.0.0.0:6379",
                Password: "123456",
            },
        },
    }
    configuration := &config.Configuration{
        RedisConfig: &config.RedisConfiguration{
            Servers: servers,
        },
        RouteAlgorithm: "single-read-write",
        Active:         "server1",  
    }

    client := redis.NewDevsporeClient(configuration)
    ctx := context.Background()
    client.Set(ctx, "test_key", "test_val", time.Hour)
    client.Get(ctx, "test_key")
}
```
### Yaml configuration file

```bigquery
props:
  version: v1
  appId: xxx
  monitorId: sdk_test
  cloud: huaweicloud
  region: cn-north-4
  azs: az1 
etcd: # Optional
  address: 127.0.0.1:2379,127.0.0.2:2379,127.0.0.3:2379
  apiVersion: v3
  username: root
  password: root
  httpsEnable: false
redis:
  redisGroupName: xxx-redis-group
  username: xxx # for redis 6.0
  password: yyy
  nearest: dc1
  connectionPool:
    enable: true
  servers:
    dc1:
      hosts: 127.0.0.1:6379
      password: password
      type: normal  # cluster, master-slave, normal
      cloud: huaweicloud  # cloud
      region: cn-north-4  # region id
      azs: az1  # azs
      pool: # Optional
        maxTotal: 100
        maxIdle: 8
        minIdle: 0
        maxWaitMillis: 10000
        timeBetweenEvictionRunsMillis: 1000
    dc2:
      hosts: 127.0.0.1:6380
      password: password
      type: master-slave  # cluster, master-slave, normal
      cloud: huaweicloud  # cloud
      region: cn-north-4  # region id
      azs: az1  # azs
      pool: # Optional
        maxTotal: 100
        maxIdle: 8
        minIdle: 0
        maxWaitMillis: 10000
        timeBetweenEvictionRunsMillis: 1000
routeAlgorithm: local-read-single-write  # local-read-single-write, single-read-write, double-write
active: dc1
```
### Testing
package commands_test needs redis 6.2.0+, so if your redis is redis 5.0+, you need to execute 
```bigquery
ginkgo -skip="redis6" 
```
See more usages of ginkgo in **https://github.com/onsi/ginkgo**