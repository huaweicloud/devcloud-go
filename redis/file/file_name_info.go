/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2022.
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

package file

import (
	"log"
	"strconv"
	"strings"
)

type NameInfo struct {
	dir       string
	redisName string
	fileIndex string
	version   string
}

// increaseVersion Failed to add the version number
func (f *NameInfo) increaseVersion() {
	version, err := strconv.Atoi(f.version)
	if err != nil {
		log.Println(err)
	}
	f.version = strconv.Itoa(version + 1)
}

func (f *NameInfo) joining(delimiter string) string {
	return f.dir + strings.Join([]string{f.redisName, f.fileIndex, f.version}, delimiter)
}

type Item struct {
	Args []interface{}
}
