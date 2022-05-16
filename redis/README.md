 # devcloud-go/redis

### Introduction
Currently, the Redis supports three modes.single-read-write,local-read-single-write and double-write
##### single-read-write
![image](../mas/img/redis-single-read-write.png)
##### local-read-single-write
![image](../mas/img/redis-local-read-single-write.png)
##### double-write
![image](../mas/img/redis-double-write.png)
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
            Hosts:      "127.0.0.0:6379",
            Password:   "XXXX",
            Type:       config.ServerTypeNormal,
            Cloud:      "huawei cloud",
            Region:     "beijing",
            Azs:        "az0",
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
  address: xxx.xxx.xxx.xxx:xxxx,xxx.xxx.xxx.xxx:xxxx,xxx.xxx.xxx.xxx:xxxx
  apiVersion: v3
  username: XXXX
  password: XXXX
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
      hosts: xxx.xxx.xxx.xxx:xxxx
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
      hosts: xxx.xxx.xxx.xxx:xxxx
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
routeAlgorithm: single-read-write  # local-read-single-write, single-read-write, double-write
active: dc1
```
### Double-write
Redis also supports double-write modes, including memory double-write and file double-write, 
depending on asyncRemotePool.persist. true: file double-write; false: memory double-write
```bigquery
redis:
  redisGroupName: xxx-redis-group
  username: xxx # for redis 6.0
  password: yyy
  nearest: dc1
  asyncRemoteWrite:
    retryTimes: 4
  connectionPool:
    enable: true
  asyncRemotePool:
    persist: true
    threadCoreSize: 10
    taskQueueSize: 5
    persistDir: dataDir/
  servers:
    dc1:
      hosts: xxx.xxx.xxx.xxx:xxxx
      password:
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
      hosts: xxx.xxx.xxx.xxx:xxxx
      password:
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
routeAlgorithm: double-write  # local-read-single-write, single-read-write, double-write
active: dc2
```
### Testing
package commands_test needs redis 6.2.0+, so if your redis is redis 5.0+, you need to execute 
```bigquery
ginkgo -skip="redis6" 
```
See more usages of ginkgo in **https://github.com/onsi/ginkgo**

### Description of Configuration Parameters

<table width="100%">
<thead><b>Configuration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>props</td><td>PropertiesConfiguration</td><td>For details,see the description of the data structure of PropertiesConfiguration</td><td>Mas monitoring configuration,which is used together with etcd</td></tr>
<tr><td>etcd</td><td>EtcdConfiguration</td><td>For details,see the description of the data structure of EtcdConfiguration</td><td>Etcd configuration.If it is configured, it will be pulled from the remote end</td></tr>
<tr><td>redis</td><td>RedisConfiguration</td><td>For details,see the description of the data structure of RedisConfiguration</td><td>RedisServer configuration</td></tr>
<tr><td>routeAlgorithm</td><td>string</td><td>single-read-write,local-read-single-write,double-write</td><td>Routing algorithm</td></tr>
<tr><td>active</td><td>string</td><td>The value can only be dc1 or dc2</td><td>Activated Redis</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>PropertiesConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>version</td><td>string</td><td>-</td><td>Project version number</td></tr>
<tr><td>appId</td><td>string</td><td>-</td><td>Project name</td></tr>
<tr><td>monitorId</td><td>string</td><td>-</td><td>Monitoring group name</td></tr>
<tr><td>cloud</td><td>string</td><td>-</td><td>Project deployment cloud group</td></tr>
<tr><td>region</td><td>string</td><td>-</td><td>Project deployment region</td></tr>
<tr><td>azs</td><td>string</td><td>-</td><td>Project deployment AZ</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>EtcdConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>address</td><td>string</td><td>-</td><td>Etcd address</td></tr>
<tr><td>apiVersion</td><td>string</td><td>-</td><td>Etcd interface Version</td></tr>
<tr><td>username</td><td>string</td><td>-</td><td>Etcd username</td></tr>
<tr><td>password</td><td>string</td><td>-</td><td>Etcd password</td></tr>
<tr><td>httpEnable</td><td>bool</td><td>-</td><td>Specifies whether HTTPS is enabled for Etcd</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>RedisConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>nearest</td><td>string</td><td>The value can only be dc1 or dc2</td><td>Indicates the local Redis</td></tr>
<tr><td>asyncRemoteWrite.retryTimes</td><td>int</td><td>-</td><td>Number of retries of asynchronous remote write operations</td></tr>
<tr><td>connectionPool.enable</td><td>bool</td><td>true/false</td><td>Indicates whether to enable the connection pool</td></tr>
<tr><td>asyncRemotePool</td><td>AsyncRemotePoolConfiguration</td><td>For details,see the description of the data structure of AsyncRemotePoolConfiguration</td><td>Configure the asynchronous write thread pool</td></tr>
<tr><td>servers</td><td>map[string]ServerConfiguration</td><td>The key is dc1/dc2.for details about a single dimension,see the description of the data structure of ServerConfiguration</td><td>RedisServer connection configuration of dc1 and dc2</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>AsyncRemotePoolConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>threadCoreSize</td><td>int</td><td>-</td><td>Basic size of the thread pool</td></tr>
<tr><td>persist</td><td>bool</td><td>true/false</td><td>Indicates whether the command is persistent.No:The command is fast.Yes:The speed is lower than that of non-persistent</td></tr>
<tr><td>taskQueueSize</td><td>int</td><td>-</td><td>Number of buffer queues</td></tr>
<tr><td>persistDir</td><td>string</td><td>Default root directory "/"</td><td>Redis persistent file directory</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>ServerConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>hosts</td><td>string</td><td>-</td><td>RedisServer IP address</td></tr>
<tr><td>password</td><td>string</td><td>-</td><td>RedisServer password</td></tr>
<tr><td>type</td><td>string</td><td>cluster,master-slave,normal</td><td>RedisServer Type</td></tr>
<tr><td>cloud</td><td>string</td><td>-</td><td>RedisServer cloud</td></tr>
<tr><td>region</td><td>string</td><td>-</td><td>Region to which the RedisServer belongs</td></tr>
<tr><td>azs</td><td>string</td><td>-</td><td>AZ to which RedisServer belongs</td></tr>
<tr><td>pool</td><td>ServerConnectionPoolConfiguration</td><td>For details,see the description of the data structure of ServerConnectionPoolConfiguration</td><td>Connection pool configuration</td></tr>
</tbody>
</table>


<table width="100%">
<thead><b>ServerConnectionPoolConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>maxTotal</td><td>int</td><td>-</td><td>Maximum number of active objects</td></tr>
<tr><td>maxIdle</td><td>int</td><td>-</td><td>Maximum number of objects that can remain in the idle state</td></tr>
<tr><td>minIdle</td><td>int</td><td>-</td><td>Minimum number of objects that can remain in the idle state</td></tr>
<tr><td>maxWaitMillis</td><td>int</td><td>-</td><td>Maximum wait time when no object is returned in the pool</td></tr>
<tr><td>timeBetweenEvictionRunsMillis</td><td>int</td><td>-</td><td>Idle link detection thread,detection interval,in milliseconds.A negative value indicates that the detection thread is not running</td></tr>
</tbody>
</table>
