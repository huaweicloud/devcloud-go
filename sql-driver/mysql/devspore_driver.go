/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2021.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 *
 */

/*
Package mysql provides a DevsporeDriver for Go's database/sql package.
Which can automatically switch data sources, read and write separation.
The driver should be used with "github.com/go-sql-driver/mysql":

import (
"database/sql"
"github.com/huaweicloud/devcloud-go/common/password"
_ "github.com/huaweicloud/devcloud-go/sql-driver/mysql"
)

password.SetDecipher(&MyDecipher{})
db, err := sql.Open("devspore_mysql", yamlConfigPath)

See README.md for more details.
*/
package mysql

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/huaweicloud/devcloud-go/common/util"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
)

// DevsporeDriver is exported to make the driver directly accessible.
type DevsporeDriver struct {
}

// Open new connection according to the dsn
func (d DevsporeDriver) Open(dsn string) (driver.Conn, error) {
	return actualDriver.Open(dsn)
}

func init() {
	sql.Register("devspore_mysql", &DevsporeDriver{})
	actualDriver = mysql.MySQLDriver{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	var err error
	idGenerator, err = util.NewNode(util.GetWorkerIDByIp())
	if err != nil {
		log.Printf("WARNING: create snowflake node failed, err: %v", err)
	}
}

var (
	actualDriver driver.Driver
	idGenerator  *util.Node
)

// OpenConnector implements driver.DriverContext
func (d DevsporeDriver) OpenConnector(yamlFilePath string) (driver.Connector, error) {
	configuration, err := getClusterConfiguration(yamlFilePath)
	if err != nil {
		log.Printf("ERROR: getConfiuration failed, err is %v", err)
	}
	clusterDataSource, err := datasource.NewClusterDataSource(configuration)
	if err != nil {
		log.Printf("ERROR: create clusterdataSource failed, %v", err)
		return nil, err
	}
	actualExecutor := newExecutor(clusterDataSource.RouterConfiguration.Retry)
	return &devsporeConnector{clusterDataSource: clusterDataSource, executor: actualExecutor}, nil
}

var clusterConfiguration *config.ClusterConfiguration

func SetClusterConfiguration(cfg *config.ClusterConfiguration) {
	clusterConfiguration = cfg
}

func getClusterConfiguration(yamlFilePath string) (*config.ClusterConfiguration, error) {
	if len(yamlFilePath) == 0 && clusterConfiguration != nil {
		return clusterConfiguration, nil
	}
	realPath, err := filepath.Abs(yamlFilePath)
	if err != nil {
		return nil, err
	}
	// parse yaml config
	configuration, err := config.Unmarshal(realPath)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
