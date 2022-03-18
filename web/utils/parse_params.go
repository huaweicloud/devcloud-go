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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryConditions struct {
	Query  map[string]string
	Fields []string
	Order  []string
	Limit  int
	Offset int
}

func ParseQueryCond(ctx *gin.Context) (QueryConditions, error) {
	query, err := GetQuery(ctx)
	if err != nil {
		return QueryConditions{}, err
	}
	fields := GetFields(ctx)
	order := GetOrder(ctx)
	limit := GetLimit(ctx)
	offset := GetOffset(ctx)
	return QueryConditions{
		Query:  query,
		Fields: fields,
		Order:  order,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func GetQuery(ctx *gin.Context) (map[string]string, error) {
	var query = make(map[string]string)
	if v := ctx.Query("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				return query, errors.New("invalid query key/value pair")
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	return query, nil
}

func GetFields(ctx *gin.Context) []string {
	var fields []string
	if v := ctx.Query("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	return fields
}

func GetOrder(ctx *gin.Context) []string {
	var order []string
	if v := ctx.Query("order"); v != "" {
		order = strings.Split(v, ",")
	}
	return order
}

func GetLimit(ctx *gin.Context) int {
	limit := 10
	if v := ctx.Query("limit"); v != "" {
		if vm, err := strconv.Atoi(v); err == nil {
			limit = vm
		}
	}
	return limit
}

func GetOffset(ctx *gin.Context) int {
	offset := 0
	if v := ctx.Query("offset"); v != "" {
		if vm, err := strconv.Atoi(v); err == nil {
			offset = vm
		}
	}
	return offset
}

func GetStringParam(ctx *gin.Context, param string) string {
	return ctx.Param(param)
}

func GetIntParam(ctx *gin.Context, param string) (int, error) {
	return strconv.Atoi(ctx.Param(param))
}

func GetInt64Param(ctx *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(ctx.Param(param), 10, 64)
}
