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

package domain

import (
	"github.com/huaweicloud/devcloud-go/web/dao"
	"github.com/huaweicloud/devcloud-go/web/resp"
	"github.com/huaweicloud/devcloud-go/web/utils"
)

type AbstractDomain struct {
	abstractDao *dao.AbstractDao
	*utils.ModelInfo
}

func NewAbstractDomain(abstractDao *dao.AbstractDao) *AbstractDomain {
	return &AbstractDomain{
		abstractDao,
		abstractDao.ModelInfo,
	}
}

func (d *AbstractDomain) Add(model interface{}) (interface{}, *resp.ErrorMsg) {
	model, err := d.abstractDao.Add(model)
	if err != nil {
		return nil, resp.InternalServerErr2Json(err)
	}
	return model, nil
}

func (d *AbstractDomain) GetOneByPK(pk string) (interface{}, *resp.ErrorMsg) {
	model, err := d.abstractDao.GetOneByPrimaryKey(pk)
	if err != nil {
		return nil, resp.InternalServerErr2Json(err)
	}
	return model, nil
}

func (d *AbstractDomain) Update(model interface{}, pk string) (interface{}, *resp.ErrorMsg) {
	model, err := d.abstractDao.Update(model, pk)
	if err != nil {
		return nil, resp.InternalServerErr2Json(err)
	}
	return model, nil
}

func (d *AbstractDomain) Delete(pk string) *resp.ErrorMsg {
	err := d.abstractDao.DeleteByPrimaryKey(pk)
	if err != nil {
		return resp.InternalServerErr2Json(err)
	}
	return nil
}

func (d *AbstractDomain) GetList(queryCond utils.QueryConditions) (interface{}, *resp.ErrorMsg) {
	models, err := d.abstractDao.GetList(queryCond)
	if err != nil {
		return nil, resp.InternalServerErr2Json(err)
	}
	return models, nil
}
