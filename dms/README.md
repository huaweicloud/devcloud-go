# devcloud-go/dms

### Feature
1. support asynchronous consume kafka message and ensure message not lost.
2. support consumption speed-limiting.

### QuickStart
1. First you need implement the OffsetPersist interface which is defined in [offset_persist.go](offset_persist.go), the create_table sql see [example/create_table.sql](example/create_table.sql). 
```go
type OffsetPersist interface {
	Find(groupId, topic string, partition int) (int64, error)
	Save(groupId, topic string, partition int, offset int64) error
}
```
2. Then you need implement the message Handler which is defined in [method_info.go#L30](method_info.go)
```go
type BizHandler func(msg *sarama.ConsumerMessage) error
```
3. Create a props for dms consumer, there are several modes of props, async and sync, you also can specify how to commit offset, interval or quantitative by set CommitInterval or CommitSize.

- async: consume messages asynchronous
- sync: consume messages synchronous
4. Create a dms consumer to consume kafka messages.

See details in package example.

### Note
1. when using async mode, the pool size should be larger than topic*partition numbers.