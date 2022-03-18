/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2022. All rights reserved.
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

package utils

import (
	"errors"
	"reflect"
	"strings"
)

type ModelInfo struct {
	Model  reflect.Type
	Models reflect.Type
	PKName string
	PKJson string
}

func NewModel(model, models interface{}, pkName, pkJson string) *ModelInfo {
	modelInfo := &ModelInfo{
		Model:  reflect.TypeOf(model),
		Models: reflect.TypeOf(models),
		PKName: pkName,
		PKJson: pkJson,
	}
	if modelInfo.Model.Kind() == reflect.Ptr {
		modelInfo.Model = modelInfo.Model.Elem()
	}
	if modelInfo.Models.Kind() == reflect.Ptr {
		modelInfo.Models = modelInfo.Models.Elem()
	}
	return modelInfo
}

func (m *ModelInfo) GetModel() reflect.Value {
	return reflect.New(m.Model)
}

func (m *ModelInfo) GetModels() reflect.Value {
	return reflect.New(m.Models)
}

func GetPKName(i interface{}) (pkName, pkJson string, err error) {
	t := reflect.TypeOf(i).Elem()
	for n := 0; n < t.NumField(); n++ {
		tf := t.Field(n)
		if strings.Index(tf.Tag.Get("gorm"), "primaryKey") >= 0 {
			return tf.Name, tf.Tag.Get("json"), nil
		}
	}
	return "", "", errors.New("the model format is abnormal")
}
