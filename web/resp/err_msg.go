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

type ErrorMsg struct {
	Errno  int
	Reason string
}

func Err2Json(code int, err error) *ErrorMsg {
	return &ErrorMsg{
		Errno:  code,
		Reason: err.Error(),
	}
}

func InternalServerErr2Json(err error) *ErrorMsg {
	return &ErrorMsg{
		Errno:  http.StatusInternalServerError,
		Reason: err.Error(),
	}
}

func BadRequestErr2Json(err error) *ErrorMsg {
	return &ErrorMsg{
		Errno:  http.StatusBadRequest,
		Reason: err.Error(),
	}
}

func NotFoundErr2Json(err error) *ErrorMsg {
	return &ErrorMsg{
		Errno:  http.StatusNotFound,
		Reason: err.Error(),
	}
}
