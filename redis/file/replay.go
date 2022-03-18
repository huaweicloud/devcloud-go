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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"

	"github.com/go-redis/redis/v8"
)

var Pattern = regexp.MustCompile("(.*/)(.*)?-(.*)?-(.*)?\\.dat")

// BatchReplay Run the following command to filter the files to be executed
func BatchReplay(clients map[string]redis.UniversalClient, fileNames []string) {
	redis2file := make(map[string][]string)
	for _, fileName := range fileNames {
		fileNameInfo, err := Parse(fileName)
		if err != nil {
			log.Println(err)
			continue
		}
		redis2file[fileNameInfo.redisName] = append(redis2file[fileNameInfo.redisName], fileName)
	}
	for key, value := range redis2file {
		sort.Strings(value)
		for _, filename := range value {
			if !Replay(filename, clients[key]) {
				break
			}
		}
	}
}

// ReplayExec Run the command on the client
func ReplayExec(client redis.UniversalClient, interrupted, lineIndex *int64, line []byte) {
	var fileItem Item
	if err := json.Unmarshal(line, &fileItem); err != nil {
		*interrupted++
		log.Println(err)
	} else {
		if c := client.Do(context.Background(), fileItem.Args...); c.Err() != nil {
			*interrupted++
			log.Println(fileItem.Args, c.Err())
		} else {
			*lineIndex++
		}
	}
}

// Replay Traverse all commands in the file and execute them
func Replay(filename string, client redis.UniversalClient) bool {
	var lineIndex, interrupted int64 = 0, 0
	file, err := os.OpenFile(filename, os.O_APPEND, CacheFilePerm)
	if err != nil {
		interrupted++
		log.Println("ERROR: OpenFile " + filename + " failed")
	} else {
		br := bufio.NewReader(file)
		for {
			line, _, err := br.ReadLine()
			if err != nil {
				break
			}
			ReplayExec(client, &interrupted, &lineIndex, line)
		}
		file.Close()
	}
	if interrupted > 0 {
		oldFilenameInfo, err := Parse(filename)
		if err != nil {
			log.Println(err)
		} else {
			oldFilenameInfo.increaseVersion()
			failDispose(filename, oldFilenameInfo.joining(Delimiter)+Suffix, lineIndex)
		}
	}
	err = os.Remove(filename)
	if err == nil {
		log.Println("success delete file: ", filename)
	} else {
		log.Println("ERROR: delete file " + filename + " failed")
	}
	return interrupted == 0
}

// failDispose Write the execution exception and subsequent contents to the new version number file
func failDispose(srcPath, dstPath string, startLine int64) {
	srcFile, err := os.OpenFile(srcPath, os.O_APPEND, CacheFilePerm)
	if err != nil {
		log.Println("ERROR: OpenFile " + srcPath + " failed")
		return
	}
	defer srcFile.Close()
	br := bufio.NewReader(srcFile)

	dstFile, err := os.OpenFile(dstPath, os.O_CREATE, CacheFilePerm)
	if err != nil {
		log.Println("ERROR: OpenFile " + dstPath + " failed")
		return
	}
	defer dstFile.Close()
	bw := bufio.NewWriter(dstFile)
	var curLine int64 = 0
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		if curLine >= startLine {
			_, err = bw.Write(line)
			if err != nil {
				log.Println(err)
			}
			_, err = bw.WriteString("\r\n")
			if err != nil {
				log.Println(err)
			}
		}
		curLine++
	}
	err = bw.Flush()
	if err != nil {
		log.Println(err)
	}
}

// Parse Check whether the parsed file name meets requirements
func Parse(fileName string) (*NameInfo, error) {
	matcher := Pattern.FindStringSubmatch(fileName)
	if len(matcher) == 0 || len(matcher)-1 != reflect.TypeOf(NameInfo{}).NumField() {
		return nil, errors.New(fmt.Sprintf("fail format not match, expect: "+
			"${dir}${redisName}-${fileIndex}-${version}.dat actual: %s", fileName))
	}
	return &NameInfo{matcher[1], matcher[2], matcher[3], matcher[4]}, nil
}
