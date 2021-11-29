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
	"testing"

	"github.com/huaweicloud/devcloud-go/common/password"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/stretchr/testify/assert"
)

func Test_convertDataSourceToDSN(t *testing.T) {
	dataSource := &config.DataSourceConfiguration{
		URL:      "tcp(127.0.0.1:3306)/ds0?timeout=200ms",
		Username: "root",
		Password: "123456",
	}
	// test with default password decipher
	noDecipherDsn := convertDataSourceToDSN(dataSource)
	assert.Equal(t, "root:123456@tcp(127.0.0.1:3306)/ds0?timeout=200ms", noDecipherDsn)

	// test datasource with server and schema
	dataSource.Server = "127.0.0.1:3307"
	dataSource.Schema = "ds1"
	noDecipherWithServerSchemaDsn := convertDataSourceToDSN(dataSource)
	assert.Equal(t, "root:123456@tcp(127.0.0.1:3307)/ds1?timeout=200ms", noDecipherWithServerSchemaDsn)

	// test datasource without options
	dataSource.URL = "tcp(127.0.0.1:3306)/ds0"
	noDecipherWithoutOptDsn := convertDataSourceToDSN(dataSource)
	assert.Equal(t, "root:123456@tcp(127.0.0.1:3307)/ds1", noDecipherWithoutOptDsn)

	// test with password decipher
	password.SetDecipher(&myDecipher{})
	withDecipherDsn := convertDataSourceToDSN(dataSource)
	assert.Equal(t, "root:123456_test@tcp(127.0.0.1:3307)/ds1", withDecipherDsn)
}

type myDecipher struct {
}

func (d *myDecipher) Decode(input string) string {
	return input + "_test"
}
