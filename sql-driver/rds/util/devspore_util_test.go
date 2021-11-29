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
	"testing"
)

func TestIsOnlyRead(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "read_case",
			args: args{
				sql: "select * from user",
			},
			want: true,
		},
		{
			name: "update case",
			args: args{
				sql: "UPDATE Websites SET alexa='5000', country='USA' WHERE name='test'",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOnlyRead(tt.args.sql); got != tt.want {
				t.Errorf("IsOnlyRead() = %v, want %v", got, tt.want)
			}
		})
	}
}
