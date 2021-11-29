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
 */

package datasource

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/huaweicloud/devcloud-go/common/password"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
)

// DataSource interface
type DataSource interface {
	// extend is only for implement
	extend()
}

// ActualDataSource an actual datasource which contains actual dsn.
type ActualDataSource struct {
	Available     bool
	Dsn           string
	Name          string
	RetryTimes    int
	LastRetryTime int64 // latest retry timestamp, ms
}

// NewActualDataSource create an actual datasource through dsn
func NewActualDataSource(name string, datasource *config.DataSourceConfiguration) *ActualDataSource {
	return &ActualDataSource{
		Available: true,
		Name:      name,
		Dsn:       convertDataSourceToDSN(datasource),
	}
}

func (ad *ActualDataSource) extend() {
}

// DsnFmt mysql dsn Example: username:password@protocol(address)/dbname?param=value,
// see details https://github.com/go-sql-driver/mysql#dsn-data-source-name
const DsnFmt = "%s:%s@%s"

// convert datasourceConfiguration to mysql dsn.
func convertDataSourceToDSN(dataSource *config.DataSourceConfiguration) string {
	if dataSource == nil {
		return ""
	}
	if dataSource.Server != "" && dataSource.Schema != "" {
		server := dataSource.Server
		dbName := dataSource.Schema
		serverReg := regexp.MustCompile(`\((.*?)\)`)
		dataSource.URL = serverReg.ReplaceAllString(dataSource.URL, "("+server+")")

		if ok, err := regexp.MatchString(`/(.*?)\?`, dataSource.URL); ok && err == nil {
			schemaReg := regexp.MustCompile(`/(.*?)\?`)
			dataSource.URL = schemaReg.ReplaceAllString(dataSource.URL, "/"+dbName+"?")
		} else {
			strs := strings.Split(dataSource.URL, "/")
			dataSource.URL = strs[0] + "/" + dbName
		}
	}
	pwd := password.GetDecipher().Decode(dataSource.Password)
	dsn := fmt.Sprintf(DsnFmt, dataSource.Username, pwd, dataSource.URL)
	return dsn
}
