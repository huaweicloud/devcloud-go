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

	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterDataSource_EmptyConfiguration(t *testing.T) {
	clusterDatasource, err := NewClusterDataSource(nil)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "clusterConfiguration cannot be nil", err.Error())
}

func TestNewClusterDataSource_EmptyDataSource(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.DataSource = nil
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "datasource config cannot be nil", err.Error())
}

func TestNewClusterDataSource_EmptyRouter(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.RouterConfig = nil
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "router config cannot be nil", err.Error())
}

func TestNewClusterDataSource_EmptyProps(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.Props = nil
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "props cannot be nil", err.Error())
}

func TestNewClusterDataSource_EmptyAppID(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.Props.AppID = ""
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "appId is required", err.Error())
}

func TestNewClusterDataSource_EmptyMonitorID(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.Props.MonitorID = ""
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "monitorId is required", err.Error())
}

func TestNewClusterDataSource_EmptyDatabaseName(t *testing.T) {
	clusterConfiguration, _ := config.Unmarshal("../resources/config.yaml")
	clusterConfiguration.Props.DatabaseName = ""
	clusterDatasource, err := NewClusterDataSource(clusterConfiguration)
	assert.Nil(t, clusterDatasource)
	assert.Equal(t, "databaseName is required", err.Error())
}
