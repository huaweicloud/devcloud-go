# devcloud-go/sql-driver/mysql

### Introduction
Currently, MySQL supports two modes.single-read-write and local-read-single-write.
In addition, read/write separation is supported, which can be configured as random or RoundRobin.
##### single-read-write
![image](../../mas/img/mysql-single-read-write.png)
##### local-read-single-write
![image](../../mas/img/mysql-local-read-single-write.png)
### Quickstart：
```bigquery
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/huaweicloud/devcloud-go/common/etcd"
    "github.com/huaweicloud/devcloud-go/mas"
    devspore "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
    "github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

func main()  {
    devspore.SetClusterConfiguration(devsporeConfiguration())
    var err error
    db, err = gorm.Open(mysql.New(mysql.Config{
    	DriverName: "devspore_mysql",
    	DSN:        "",
    }), &gorm.Config{})
    log.Printf("create db failed, %v", err)
}
func devsporeConfiguration() *config.ClusterConfiguration {
    return &config.ClusterConfiguration{
        Props: &mas.PropertiesConfiguration{
            AppID:        "xxx",
            MonitorID:    "xxx",
            DatabaseName: "xx",
        },
    	EtcdConfig: &etcd.EtcdConfiguration{
            Address:     "xxx.xxx.xxx.xxx:xxxx,xxx.xxx.xxx.xxx:xxxx,xxx.xxx.xxx.xxx:xxxx",
            Username:    "XXXX",
            Password:    "XXXX",
            HTTPSEnable: false,
    	},
    	RouterConfig: &config.RouterConfiguration{
            Nodes: map[string]*config.NodeConfiguration{
                "dc1": {
                    Master: "ds1",
                },
                "dc2": {
                    Master: "ds2",
                },
            },
            Active: "dc1",
    	},
    	DataSource: map[string]*config.DataSourceConfiguration{
            "ds1": {
                URL:      "tcp(xxx.xxx.xxx.xxx:xxxx)/ds0?charset=utf8&parseTime=true",
                Username: "XXXX",
                Password: "XXXX",
            },
            "ds2": {
                URL:      "tcp(xxx.xxx.xxx.xxx:xxxx)/ds0?charset=utf8&parseTime=true",
                Username: "XXXX",
                Password: "XXXX",
            },
    	},
    }
}

```
you also can use yaml file.
```bigquery
1.sql
import (
    "database/sql"
    "fmt"

    "github.com/huaweicloud/devcloud-go/common/password"
    _ "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

func main() {
    password.SetDecipher(&MyDecipher{}) //MyDecipher implements password.Decipher interface
    yamlConfigPath := "xxx/config_with_password.yaml"
    db, err := sql.Open("devspore_mysql", yamlConfigPath)
    if err != nil {
        fmt.Errorf(err.Error())
    }
    ......THEN 
}

2.gorm
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    _ "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

func main() {
    db, err := gorm.Open(mysql.New(mysql.Config{
        DriverName: "devspore_mysql",
        DSN:        "xxx/config_with_password.yaml",
    }), &gorm.Config{})
    ......THEN 
}

3.beego-orm
import (
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

func main() {
    // register devspore_mysql
    err = orm.RegisterDriver("devspore_mysql", orm.DRMySQL)
    if err != nil {
        log.Fatalln(err)
    }
    // register model
    orm.RegisterModel(new(interface{}),new(interface{}))
    
    err = orm.RegisterDataBase("default", "devspore_mysql", "xxx/config_with_password.yaml")
    if err != nil {
        log.Fatalln(err)
    }
    db:= orm.NewOrm()
    ......THEN 
}

```
**Version requirements：go1.14.6 and above**

### Configuration file format：
you can just configure datasource and router if you don't use mas.
```bigquery
props: # Optional
  version: v1  // project version
  appId: xxxxx  // mas appId
  monitorId: xxxxx  // mas monitorId
  databaseName: xxxxx  // dbName

etcd: # Optional
  address: 127.0.0.2:2379,127.0.0.2:2379,127.0.0.2:2379  
  apiVersion: v3  // etcd version
  username: etcduser  
  password: etcdpwd  
  httpsEnable: false  
  
datasource: # Require
  ds0:
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds0 
    username: datasourceuser 
    password: datasourcepwd  
  ds0-slave0:
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds0_slave0
    username: datasourceuser
    password: datasourcepwd
  ds0-slave1: 
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds0_slave1
    username: datasourceuser
    password: datasourcepwd
  ds1:
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds1
    username: datasourceuser
    password: datasourcepwd
  ds1-slave0:
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds1_slave0
    username: datasourceuser
    password: datasourcepwd
  ds1-slave1:
    url: tcp(xxx.xxx.xxx.xxx:xxxx)/ds1_slave1
    username: datasourceuser
    password: datasourcepwd

router: # Require
  active: c0 
  routeAlgorithm: single-read-write  // single-read-write(default), local-read-single-write
  retry:
    times: 3  
    delay: 50  // ms
  nodes:
    c0:  
      weight: ""  // not yet used
      master: ds0  // 
      loadBalance: ROUND_ROBIN  // ROUND_ROBIN(default),RANDOM
      slaves:  
        - ds0-slave0
        - ds0-slave1
    c1:
      weight: ""
      master: ds1
      loadBalance: ROUND_ROBIN
      slaves:
        - ds1-slave0
        - ds1-slave1

```

### Description of Configuration Parameters

<table width="100%">
<thead><b>ClusterConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>props</td><td>PropertiesConfiguration</td><td>For details,see the description of the data structure of PropertiesConfiguration</td><td>Mas monitoring configuration,which is used together with etcd</td></tr>
<tr><td>etcd</td><td>EtcdConfiguration</td><td>For details,see the description of the data structure of EtcdConfiguration</td><td>Etcd configuration.If it is configured, it will be pulled from the remote end</td></tr>
<tr><td>datasource</td><td>map[string]DataSourceConfiguration</td><td>The key is customized,for details about a single dimension,see the description of the data structure of DataSourceConfiguration</td><td>DataSource</td></tr>
<tr><td>router</td><td>RouterConfiguration</td><td>For details,see the description of the data structure of RouterConfiguration</td><td>Route-related configuration</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>PropertiesConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>version</td><td>string</td><td>-</td><td>Project version number</td></tr>
<tr><td>appId</td><td>string</td><td>-</td><td>Project name</td></tr>
<tr><td>monitorId</td><td>string</td><td>-</td><td>Monitoring group name</td></tr>
<tr><td>databaseName</td><td>string</td><td>-</td><td>Database name</td></tr>
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
<thead><b>DataSourceConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>url</td><td>string</td><td>protocol(address)/dbname?param=value</td><td>Data Source Name</td></tr>
<tr><td>username</td><td>string</td><td>-</td><td>Username</td></tr>
<tr><td>password</td><td>string</td><td>-</td><td>Password</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>RouterConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>active</td><td>string</td><td>Key of the node</td><td>Activating Nodes</td></tr>
<tr><td>routeAlgorithm</td><td>string</td><td>single-read-write,local-read-single-write</td><td>Routing algorithm</td></tr>
<tr><td>retry.times</td><td>string</td><td>-</td><td>Failed Retry Times</td></tr>
<tr><td>retry.delay</td><td>string</td><td>-</td><td>Retry interval,in milliseconds</td></tr>
<tr><td>nodes</td><td>map[string]NodeConfiguration</td><td>The key is customized,for details about a single dimension,see the description of the data structure of NodeConfiguration</td><td>Node-related configuration</td></tr>
</tbody>
</table>

<table width="100%">
<thead><b>NodeConfiguration</b></thead>
<tbody>
<tr><th>Parameter Name</th><th>Parameter Type</th><th>Value range</th><th>Description</th></tr>
<tr><td>master</td><td>string</td><td>key of the datasource</td><td>Master node datasource</td></tr>
<tr><td>loadBalance</td><td>string</td><td>RANDOM,ROUND_ROBIN</td><td>Load balancing algorithm for read/write separation</td></tr>
<tr><td>slaves</td><td>[]string</td><td>key of the datasource</td><td>Slave node datasource</td></tr>
</tbody>
</table>
