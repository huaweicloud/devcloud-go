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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateClusterConfiguration(t *testing.T) {
	type args struct {
		configuration *ClusterConfiguration
	}
	// normal_case
	normalConfigPath, err := filepath.Abs("../resources/config.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	clusterConfiguration, err := Unmarshal(normalConfigPath)
	if err != nil {
		t.Errorf("unmarshal configuration failed, err: %v", err)
	}

	// no_datasources
	noDataSourceConfigPath, err := filepath.Abs("../resources/no_datasources.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	noDataSourceConfiguration, err := Unmarshal(noDataSourceConfigPath)
	if err != nil {
		t.Errorf("unmarshal configuration failed, err: %v", err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// normal clusterConfiguration
		{
			name: "normal_case",
			args: args{
				configuration: clusterConfiguration,
			},
			wantErr: false,
		},
		// no datasources clusterConfiguration
		{
			name: "no_datasources_case",
			args: args{
				configuration: noDataSourceConfiguration,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateClusterConfiguration(tt.args.configuration); (err != nil) != tt.wantErr {
				t.Errorf("ValidateClusterConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshalWithEnvVariables(t *testing.T) {
	path, err := filepath.Abs("../resources/config_env.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	const (
		etcd_address  = "127.0.0.1:8081"
		etcd_username = "root_env"
		etcd_password = "password_env"
	)

	os.Clearenv()
	if err := os.Setenv("etcd_address", etcd_address); err != nil {
		t.Errorf("set env 'etcd_address' failed, err %v", err)
		return
	}
	if err := os.Setenv("etcd_username", etcd_username); err != nil {
		t.Errorf("set env 'etcd_username' failed, err %v", err)
		return
	}
	if err := os.Setenv("etcd_password", etcd_password); err != nil {
		t.Errorf("set env 'etcd_password' failed, err %v", err)
		return
	}
	configuration, err := Unmarshal(path)
	if err != nil {
		t.Errorf("Unmarshal config_env.yaml failed, err %v", err)
		return
	}

	assert.NotNil(t, configuration)
	assert.NotNil(t, configuration.EtcdConfig)
	assert.Equal(t, etcd_address, configuration.EtcdConfig.Address)
	assert.Equal(t, etcd_username, configuration.EtcdConfig.Username)
	assert.Equal(t, etcd_password, configuration.EtcdConfig.Password)

	assert.Len(t, configuration.DataSource, 1)
	assert.True(t, len(configuration.DataSource["ds0"].URL) != 0)
}
