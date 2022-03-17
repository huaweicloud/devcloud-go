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

// Package file Double-write of files
package file

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CacheFilePerm                      = 0666
	Suffix                             = ".dat"
	Delimiter                          = "-"
	DefaultVersion                     = "0"
	FileTimestampGapMillions           = 60000
	FileTimestampCloseMillions         = 70000
	FileCloseCheckTimestampGapMillions = 5000
	FlushGapMillions                   = 10000
	RelativePattern                    = "(dc[12])-([0-9]{13})-([0-9]*)\\.dat"
	MaxLine                            = 1 << 20
)

type Operation struct {
	lastCreateTime int64
	lastFlushTime  int64
	lineIndex      int64
	file           *os.File
	fileWriter     *bufio.Writer
	close          func()
	mutex          *sync.Mutex
}

func NewFileOperation() *Operation {
	fileOperation := &Operation{mutex: &sync.Mutex{}}
	go fileOperation.CloseWriter()
	return fileOperation
}

// CloseWriter Polling disables the write of files that are not operated within a specified period of time
func (f *Operation) CloseWriter() {
	ticker := time.NewTicker(time.Millisecond * FileCloseCheckTimestampGapMillions)
	defer ticker.Stop()
	for range ticker.C {
		f.mutex.Lock()
		if f.fileWriter == nil {
			continue
		}
		if time.Now().UnixNano()/1e6 > f.lastFlushTime+FileTimestampGapMillions {
			err := f.fileWriter.Flush()
			if err != nil {
				log.Println(err)
			}
			f.fileWriter = nil
			err = f.file.Close()
			if err != nil {
				log.Println(err)
			}
		}
		f.mutex.Unlock()
	}
}

// WriteFile Command Write
func (f *Operation) WriteFile(path string, content Item) {
	// Maximum line or time reached
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.fileWriter == nil || f.isShouldNewFile() {
		if f.fileWriter != nil {
			err := f.fileWriter.Flush()
			if err != nil {
				log.Println(err)
			}
			f.fileWriter = nil
			err = f.file.Close()
			if err != nil {
				log.Println(err)
			}
		}
		f.lastCreateTime = time.Now().UnixNano() / 1e6
		f.lastFlushTime = time.Now().UnixNano() / 1e6
		f.lineIndex = 0
		var err error
		f.file, err = f.CreateFile(path)
		if err != nil {
			log.Println("CreateFile failed")
			return
		}
		f.fileWriter = bufio.NewWriter(f.file)
	}
	data, err := json.Marshal(content)
	if err != nil {
		log.Println(err)
	} else {
		_, err = f.fileWriter.Write(data)
		if err != nil {
			log.Println(err)
		}
		_, err = f.fileWriter.WriteString("\r\n")
		if err != nil {
			log.Println(err)
		}
		f.lineIndex = (f.lineIndex + 1) & (0x000003ff)
	}
	if f.lineIndex == 0x000003ff || time.Now().UnixNano()/1e6 > f.lastFlushTime+FlushGapMillions {
		err = f.fileWriter.Flush()
		if err != nil {
			log.Println(err)
		}
		f.lastFlushTime = time.Now().UnixNano() / 1e6
	}
}

// isShouldNewFile Whether to re-create the write file
func (f *Operation) isShouldNewFile() bool {
	return f.lineIndex > MaxLine || time.Now().UnixNano()/1e6 > f.lastCreateTime+FileTimestampGapMillions
}

func (f *Operation) CreateFile(originPath string) (*os.File, error) {
	nameItem := make([]string, 0)
	nameItem = append(nameItem, originPath)
	nameItem = append(nameItem, strconv.FormatInt(f.lastCreateTime, 10))
	nameItem = append(nameItem, DefaultVersion)
	path := strings.Join(nameItem, Delimiter) + Suffix
	return os.OpenFile(path, os.O_CREATE, CacheFilePerm)
}

// traversal Traverse and check whether the execution requirements are met
func traversal(info os.FileInfo, matchFile *[]string, nameMap map[string]string, nowTime int64) {
	if matched, err := regexp.MatchString(RelativePattern, info.Name()); matched {
		if err != nil {
			log.Println(err)
		}
		*matchFile = append(*matchFile, info.Name())
		startIndex := strings.Index(info.Name(), Delimiter) + 1
		endIndex := strings.LastIndex(info.Name(), "-")
		ux, err := strconv.ParseInt(info.Name()[startIndex:endIndex], 10, 64)
		if err == nil && ux > nowTime-FileTimestampCloseMillions {
			if info.Name() > nameMap[info.Name()[0:3]] {
				nameMap[info.Name()[0:3]] = info.Name()
			}
		}
	}
}

// FileListNeedReplay Obtain the list of files to be executed
func FileListNeedReplay(dir string) []string {
	totalFilenames := make([]string, 0)
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	if rd != nil && len(rd) > 0 {
		matchFile := make([]string, 0)
		nameMap := make(map[string]string)
		nowTime := time.Now().UnixNano() / 1e6
		for _, info := range rd {
			traversal(info, &matchFile, nameMap, nowTime)
		}
		for _, name := range matchFile {
			if nameMap[name[0:3]] != name {
				totalFilenames = append(totalFilenames, dir+name)
			}
		}
	}
	return totalFilenames
}

func MkDirs(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, CacheFilePerm)
		if err != nil {
			log.Println("mkdir fail ", dir, err)
		} else {
			log.Println("mkdir success ", dir)
		}
	}
}
