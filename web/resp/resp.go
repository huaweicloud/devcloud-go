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

package resp

import "net/http"

type ResponseInfo struct {
	Code  int         `json:"code,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

var FuncResp func(info *ResponseInfo) interface{}

func FormatResp(funcResp func(info *ResponseInfo) interface{}) {
	FuncResp = funcResp
}

func GetResp(info *ResponseInfo) interface{} {
	if FuncResp != nil {
		return FuncResp(info)
	}
	return info
}

func Success() *ResponseInfo {
	return &ResponseInfo{
		Code: http.StatusOK,
	}
}

func CreateData(data interface{}) *ResponseInfo {
	return &ResponseInfo{
		Code: http.StatusCreated,
		Data: data,
	}
}

func SuccessData(data interface{}) *ResponseInfo {
	return &ResponseInfo{
		Code: http.StatusOK,
		Data: data,
	}
}

func Failure(msg string) *ResponseInfo {
	return &ResponseInfo{
		Code:  http.StatusInternalServerError,
		Error: msg,
	}
}

func FailureStatus(status int, msg string) *ResponseInfo {
	return &ResponseInfo{
		Code:  status,
		Error: msg,
	}
}
