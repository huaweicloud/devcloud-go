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

package util

import (
	"reflect"
	"testing"
)

func TestValidateHostPort(t *testing.T) {
	type args struct {
		hostPort string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal_case",
			args: args{
				hostPort: "127.0.0.1:2379",
			},
			wantErr: false,
		},
		{
			name: "invalid_port",
			args: args{
				hostPort: "127.0.0.1:237999",
			},
			wantErr: true,
		},
		{
			name: "case1,with white space",
			args: args{
				hostPort: "127.0.0.1: 2379",
			},
			wantErr: true,
		},
		{
			name: "case2,with white space",
			args: args{
				hostPort: " 127.0.0.1:2379",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateHostPort(tt.args.hostPort); (err != nil) != tt.wantErr {
				t.Errorf("ValidateHostPort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertAddressStrToSlice(t *testing.T) {
	type args struct {
		addressStr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no space",
			args: args{
				addressStr: "127.0.0.1:6379,127.0.0.1:6380",
			},
			want: []string{"127.0.0.1:6379", "127.0.0.1:6380"},
		},
		{
			name: "with space",
			args: args{
				addressStr: "127.0.0.1:6379, 127.0.0.1:6380",
			},
			want: []string{"127.0.0.1:6379", "127.0.0.1:6380"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertAddressStrToSlice(tt.args.addressStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAddressStrToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
