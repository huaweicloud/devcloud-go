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
 *
 * Package password defines Decipher interface, which is used to decode password,
 * user can set customize decipher by SetDecipher function.
 */

package password

import "sync"

// Decipher to decode password
type Decipher interface {
	Decode(string) string
}

var (
	actualDecipher Decipher
	lock           = &sync.Mutex{}
)

// SetDecipher for set a customized Decipher
func SetDecipher(decipher Decipher) {
	lock.Lock()
	actualDecipher = decipher
	lock.Unlock()
}

// GetDecipher return a password Decipher
func GetDecipher() Decipher {
	lock.Lock()
	defer lock.Unlock()
	if actualDecipher != nil {
		return actualDecipher
	}
	return &defaultDecipher{}
}

type defaultDecipher struct {
}

func (dd *defaultDecipher) Decode(password string) string {
	return password
}
