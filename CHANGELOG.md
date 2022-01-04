#### 2022-01-04
version: github.com/huaweicloud/devcloud-go v0.1.1  
feature:
1. mock package add etcd mock, replace test cases that rely on real etcd.

#### 2021-12-27
version: github.com/huaweicloud/devcloud-go v0.1.0  
feature:
1. add dms package, which is a high performance and high reliability kafka consumer.
2. add mock package, which can mock redis, mysql and interface.
#### 2021-12-25
1. dms: persist the first N continuous offsets in offsetNode to the database and kafka broker, this will reduce repeated consumption of messages.
2. change dms/method.go BizHandler from interface to function types.
#### 2021-12-24
1. add mock package, which contains interface mock, redis mock and mysql mock.

#### 2021-12-16
version: github.com/huaweicloud/devcloud-go v0.0.1  
feature:
1. add dms which is a kafka consumer.

#### 2021-12-03
version: github.com/huaweicloud/devcloud-go v0.0.1  
feature：
1. fix bug: sql-driver transaction will panic in high concurrency situations.