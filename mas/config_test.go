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

package mas

import "testing"

func TestPropertiesConfiguration_CalHashCode(t *testing.T) {
	type fields struct {
		Version           string
		AppID             string
		MonitorID         string
		DatabaseName      string
		DecipherClassName string
		Region            string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal_case",
			fields: fields{
				Version:      "v3",
				AppID:        "xxxappId",
				MonitorID:    "xxxmonitorId",
				DatabaseName: "xxxdatabase",
			},
			want: "v3_xxxappId_xxxmonitorId_xxxdatabase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PropertiesConfiguration{
				Version:      tt.fields.Version,
				AppID:        tt.fields.AppID,
				MonitorID:    tt.fields.MonitorID,
				DatabaseName: tt.fields.DatabaseName,
				Region:       tt.fields.Region,
			}
			if got := p.CalHashCode(); got != tt.want {
				t.Errorf("CalHashCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
