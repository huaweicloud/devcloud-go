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
	"errors"
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
	model, errMsg := c.domain.GetOneByPK(pk)
	if errMsg != nil {
		ctx.JSON(errMsg.Errno, errMsg)
		return
	}
	ctx.JSON(http.StatusOK, model)
}

func (c *AbstractController) Add(ctx *gin.Context) {
	model := c.domain.GetModel()
	if err := ctx.BindJSON(model.Interface()); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.BadRequestErr2Json(err))
		return
	}
	curModel, errMsg := c.domain.Add(model.Interface())
	if errMsg != nil {
		ctx.JSON(errMsg.Errno, errMsg)
		return
	}
	ctx.JSON(http.StatusOK, curModel)
}

func (c *AbstractController) Update(ctx *gin.Context) {
	pk := utils.GetStringParam(ctx, c.domain.PKJson)
	model := c.domain.GetModel()
	var err error
	if err = ctx.BindJSON(model.Interface()); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.BadRequestErr2Json(err))
		return
	}
	if model.Elem().FieldByName(c.domain.PKName).String() != pk {
		ctx.JSON(http.StatusBadRequest, resp.BadRequestErr2Json(errors.New("different "+c.domain.PKJson+" between path and body")))
		return
	}
	curModel, errMsg := c.domain.Update(model.Interface(), pk)
	if errMsg != nil {
		ctx.JSON(errMsg.Errno, errMsg)
		return
	}
	ctx.JSON(http.StatusOK, curModel)
}

func (c *AbstractController) Delete(ctx *gin.Context) {
	pk := utils.GetStringParam(ctx, c.domain.PKJson)
	if errMsg := c.domain.Delete(pk); errMsg != nil {
		ctx.JSON(errMsg.Errno, errMsg)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func (c *AbstractController) GetAll(ctx *gin.Context) {
	queryCond, err := utils.ParseQueryCond(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resp.BadRequestErr2Json(err))
		return
	}
	models, errMsg := c.domain.GetList(queryCond)
	if errMsg != nil {
		ctx.JSON(errMsg.Errno, errMsg)
		return
	}
	ctx.JSON(http.StatusOK, models)
}
