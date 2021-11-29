# Devcloud-go

Devcloud-go provides a sql-driver for mysql which named devspore driver and a redis client which named devspore client, you can use them with MAS or use them separately. 
The devspore driver is developed based on [go-sql-driver/mysql v1.6.0](https://github.com/go-sql-driver/mysql), the devspore client is developed based on [go-redis v8.11.3](https://github.com/go-redis/redis).  
This document introduces how to obtain and use Devcloud-go.
***
## Requirements
* To use devcloud-go multi datasource disaster recovery capability, you need to create an MAS application in huaweicloud.
* Devcloud-go requires go 1.14.6 or later, run command `go version` to check the version of Go.
***
## Install
Run the following command to install Devcloud-go:
```bigquery
go get github.com/huaweicloud/devcloud-go
```
***
## Code Example
* **sql-driver** : see details in [sql-driver/mysql/README.md](sql-driver/mysql/README.md)
* **redis** : see details in [redis/README.md](redis/README.md)
***
## ChangeLog
Detailed changes for each released version are documented in the [CHANGELOG.md](CHANGELOG.md).

***
## License
This project is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.