# Devcloud-go

Devcloud-go provides a sql-driver for mysql, a redis client and a kafka consumer which named dms, you can use the driver and redis client with MAS or use them separately,
at the same time, they also support use cases with injection faults for scenario simulation.
The driver is developed based on [go-sql-driver/mysql v1.6.0](https://github.com/go-sql-driver/mysql), the redis client is developed based on [go-redis v8.11.3](https://github.com/go-redis/redis).
The kafka consumer is developed based on [github.com/Shopify/sarama v1.29.1](https://github.com/Shopify/sarama).
The mock package provides the simulation of MySQL, redis and etcd services, and realizes the fault injection function of MySQL and redis through TCP.
This document introduces how to obtain and use Devcloud-go.

## Requirements
* To use devcloud-go multi datasource disaster recovery capability, you need to create an MAS application in huaweicloud.
* Devcloud-go requires go 1.14.6 or later, run command `go version` to check the version of Go.

## Install
Run the following command to install Devcloud-go:
```bigquery
go get github.com/huaweicloud/devcloud-go
```

## Code Example
* **sql-driver** : see details in [sql-driver/mysql/README.md](sql-driver/mysql/README.md)
* **redis** : see details in [redis/README.md](redis/README.md)
* **dms**: see details in [dms/README.md](dms/README.md)
* **mock**: see details in [mock/README.md](mock/README.md)

## ChangeLog
Detailed changes for each released version are documented in the [CHANGELOG.md](CHANGELOG.md).


## License
This project is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.