# devcloud-go/mock

### Quickstart：

##### etcd

```bigquery
import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    
    clientv3 "go.etcd.io/etcd/client/v3"
    "github.com/huaweicloud/devcloud-go/mock"
)

func main()  {
    addrs := []string{"127.0.0.1:2382"}
    dataDir := "etcd_data"
    defer func(path string) {
        err := os.RemoveAll(path)
        if err != nil {
            log.Println("ERROR: remove data dir failed, %v", err)
        }
    }(dataDir)
    metadata := mock.NewEtcdMetadata()
    metadata.ClientAddrs = addrs
    metadata.DataDir = dataDir
    mockEtcd := &mock.MockEtcd{}
    mockEtcd.StartMockEtcd(metadata)
    defer mockEtcd.StopMockEtcd()
    
    client, err := clientv3.New(clientv3.Config{Endpoints: addrs, Username: "root", Password: "root"})
    defer func(client *clientv3.Client) {
        err = client.Close()
        if err != nil {
            log.Println("ERROR: close client failed, %v", err)
        }
    }(client)
    
    key := "key"
    val := "val"
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    client.Put(ctx, key, val, clientv3.WithPrevKV())
    cancel()
    
    ctx, cancel = context.WithTimeout(context.Background(), time.Second)
    resp, _ := client.Get(ctx, key)
    cancel()
    fmt.Println(string(resp.Kvs[0].Value))
}

```

##### mysql

```bigquery
import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    "github.com/dolthub/go-mysql-server/memory"
    mocksql "github.com/dolthub/go-mysql-server/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/huaweicloud/devcloud-go/mock"
)

func main() {
    metadata := mock.MysqlMock{
        User:         "root",
        Password:     "root",
        Address:      "127.0.0.1:3318",
        Databases:    []string{"mydb"},
        MemDatabases: []*memory.Database{createTestDatabase("mydb", "user")},
    }
    metadata.StartMockMysql()
    defer metadata.StopMockMysql()
    
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3318)/mydb")
    defer db.Close()
        if err != nil {
            log.Println(err)
            return 
        }
    var name, email string
    err = db.QueryRow("SELECT name, email FROM user WHERE id=?", 1).Scan(&name, &email)
    fmt.Println(name, email)
}

func createTestDatabase(dbName, tableName string) *memory.Database {
    db := memory.NewDatabase(dbName)
    table := memory.NewTable(tableName, mocksql.Schema{
        {Name: "id", Type: mocksql.Int64, Nullable: false, AutoIncrement: true, PrimaryKey: true, Source: tableName},
        {Name: "name", Type: mocksql.Text, Nullable: false, Source: tableName},
        {Name: "email", Type: mocksql.Text, Nullable: false, Source: tableName},
        {Name: "phone_numbers", Type: mocksql.JSON, Nullable: false, Source: tableName},
        {Name: "created_at", Type: mocksql.Timestamp, Nullable: false, Source: tableName},
    })
    
    db.AddTable(tableName, table)
    ctx := mocksql.NewEmptyContext()
    
    rows := []mocksql.Row{
        mocksql.NewRow(1, "John Doe", "jasonkay@doe.com", []string{"555-555-555"}, time.Now()),
        mocksql.NewRow(2, "John Doe", "johnalt@doe.com", []string{}, time.Now()),
        mocksql.NewRow(3, "Jane Doe", "jane@doe.com", []string{}, time.Now()),
        mocksql.NewRow(4, "Evil Bob", "jasonkay@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
    }
    
    for _, row := range rows {
        _ = table.Insert(ctx, row)
    }
    return db
}
```

##### redis

```bigquery
import (
    "context"
    "fmt"
    
    goredis "github.com/go-redis/redis/v8"
    "github.com/huaweicloud/devcloud-go/mock"
)

func main() {
    redisMock := mock.RedisMock{Addr: "127.0.0.1:16379"}
    redisMock.StartMockRedis()
    defer redisMock.StopMockRedis()
    cluster := goredis.NewClusterClient(&goredis.ClusterOptions{
        Addrs: []string{"127.0.0.1:16379"},
    })
    
    ctx := context.Background()
    cluster.Set(ctx, "key", "val", 0)
    res := cluster.Get(ctx, "key")
    fmt.Println(res.Val())
}
```

### Fault injection
Fault injection through TCP proxy

![image](../img/proxy.png)

Faults such as delay, fluctuation, disconnection, null value, and error can be injected.

##### Delay
```bigquery
func (p *Proxy) AddDelay(name string, delay, percentage int, clientAddr, command string) error
```
Add a delay fault named name, filter the clientAddr trustlist, intercept command commands, set the delay time to delay, and set the trigger probability to percentage.

##### Jitter
```bigquery
func (p *Proxy) AddJitter(name string, jitter, percentage int, clientAddr, command string) error
```
Add a fluctuating fault named name, filter the clientAddr trustlist, intercept command commands, set the fluctuating duration to jitter, and set the triggering probability to percentage.

##### Drop
```bigquery
func (p *Proxy) AddDrop(name string, percentage int, clientAddr, command string) error
```
Add a disconnection fault named name. The fault filters out the clientAddr trustlist and intercepts command commands. The triggering probability is percentage.

##### ReturnEmpty
```bigquery
func (p *Proxy) AddReturnEmpty(name string, percentage int, clientAddr, command string) error
```
Add a null fault named name. The fault filters the clientAddr trustlist and intercepts command commands. The triggering probability is percentage.

##### ReturnErr
```bigquery
func (p *Proxy) AddReturnErr(name string, returnErr error, percentage int, clientAddr, command string) error
```
Add an error fault named name, filter the clientAddr trustlist, intercept command commands, and set the fault information to returnErr and trigger probability to percentage.

##### Redis is used as an example.
```bigquery
import (
    "context"
    "fmt"
    
    proxyredis "github.com/huaweicloud/devcloud-go/mock/proxy/proxy-redis"
    goredis "github.com/go-redis/redis/v8"
    "github.com/huaweicloud/devcloud-go/mock"
)

func main() {
    redisMock := mock.RedisMock{Addr: "127.0.0.1:16379"}
    redisMock.StartMockRedis()
    defer redisMock.StopMockRedis()
    redisProxy := proxyredis.NewProxy(redisMock.Addr, "127.0.0.1:26379")
    redisProxy.StartProxy()
    defer redisProxy.StopProxy()
    client := goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:26379"})
    ctx := context.Background()
    client.Set(ctx, "key", "val", 0)
    time1 := time.Now().Unix() / 1e6
    res1 := client.Get(ctx, "key")
    time2 := time.Now().Unix() / 1e6
    fmt.Println(res1.Val(), time2-time1)
    redisProxy.AddDelay("delay", 1500, 0, "", "")
    //redisProxy.AddJitter("jitter", 3500, 0, "", "")
    //redisProxy.AddDrop("drop", 0, "", "")
    //redisProxy.AddReturnEmpty("returnEmpty", 0, "", "")
    //redisProxy.AddReturnErr("returnErr", proxyredis.UnknownError, 0, "", "")
    time3 := time.Now().UnixNano() / 1e6
    res2 := client.Get(ctx, "key")
    time4 := time.Now().UnixNano() / 1e6
    fmt.Println(res2.Val(), time4-time3)
}
```