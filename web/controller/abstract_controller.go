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

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/huaweicloud/devcloud-go/web/domain"
	"github.com/huaweicloud/devcloud-go/web/resp"
	"github.com/huaweicloud/devcloud-go/web/utils"
)

type AbstractController struct {
	domain *domain.AbstractDomain
}

func NewAbstractController(domain *domain.AbstractDomain) *AbstractController {
	return &AbstractController{
		domain,
	}
}

func (c *AbstractController) Get(ctx *gin.Context) {
	pk := utils.GetStringParam(ctx, c.domain.PKJson)
	response := c.domain.GetOneByPK(pk)
	ctx.JSON(response.Code, resp.GetResp(response))
}

func (c *AbstractController) Add(ctx *gin.Context) {
	model := c.domain.GetModel()
	response := new(resp.ResponseInfo)
	if err := ctx.BindJSON(model.Interface()); err != nil {
		response = resp.FailureStatus(http.StatusBadRequest, err.Error())
	} else {
		response = c.domain.Add(model.Interface())
	}
	ctx.JSON(response.Code, resp.GetResp(response))
}

func (c *AbstractController) Update(ctx *gin.Context) {
	pk := utils.GetStringParam(ctx, c.domain.PKJson)
	model := c.domain.GetModel()
	response := new(resp.ResponseInfo)
	if err := ctx.BindJSON(model.Interface()); err != nil {
		response = resp.FailureStatus(http.StatusBadRequest, err.Error())
	} else if model.Elem().FieldByName(c.domain.PKName).String() != pk {
		response = resp.FailureStatus(http.StatusBadRequest, "different "+c.domain.PKJson+" between path and body")
	} else {
		response = c.domain.Update(model.Interface(), pk)
	}
	ctx.JSON(response.Code, resp.GetResp(response))
}

func (c *AbstractController) Delete(ctx *gin.Context) {
	pk := utils.GetStringParam(ctx, c.domain.PKJson)
	response := c.domain.Delete(pk)
	ctx.JSON(response.Code, resp.GetResp(response))
}

func (c *AbstractController) GetAll(ctx *gin.Context) {
	queryCond, err := utils.ParseQueryCond(ctx)
	response := new(resp.ResponseInfo)
	if err != nil {
		response = resp.FailureStatus(http.StatusBadRequest, err.Error())
	} else {
		response = c.domain.GetList(queryCond)
	}
	ctx.JSON(response.Code, resp.GetResp(response))
}
