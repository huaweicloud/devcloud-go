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

package dao

import (
	"gorm.io/gorm"

	"github.com/huaweicloud/devcloud-go/web/utils"
)

type AbstractDao struct {
	DB *gorm.DB
	*utils.ModelInfo
}

func NewAbstractDao(db *gorm.DB, modelInfo *utils.ModelInfo) *AbstractDao {
	return &AbstractDao{
		db,
		modelInfo,
	}
}

func (d *AbstractDao) Add(model interface{}) (interface{}, error) {
	if err := d.DB.Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *AbstractDao) AddBatch(models interface{}) (interface{}, error) {
	if err := d.DB.Create(models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (d *AbstractDao) Update(model interface{}, pk string) (interface{}, error) {
	if _, err := d.GetOneByPrimaryKey(pk); err != nil {
		return nil, err
	}
	if err := d.DB.Where(d.PKJson+" = ?", pk).Updates(model).Error; err != nil {
		return nil, err
	}
	return d.GetOneByPrimaryKey(pk)
}

func (d *AbstractDao) DeleteByPrimaryKey(primaryKey string) error {
	model := d.GetModel().Interface()
	return d.DB.Where(d.PKJson+" = ?", primaryKey).Delete(model).Error
}

func (d *AbstractDao) DeleteByPrimaryKeys(primaryKeys []string) error {
	model := d.GetModel().Interface()
	return d.DB.Where(d.PKJson+" IN ?", primaryKeys).Delete(model).Error
}

func (d *AbstractDao) GetOneByPrimaryKey(primaryKey string) (interface{}, error) {
	model := d.GetModel().Interface()
	if err := d.DB.Where(d.PKJson+" = ?", primaryKey).First(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *AbstractDao) GetListByPrimaryKeys(primaryKeys []string) (interface{}, error) {
	models := d.GetModels().Interface()
	if err := d.DB.Where(d.PKJson+" IN ?", primaryKeys).Find(models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (d *AbstractDao) GetList(queryCond utils.QueryConditions) (interface{}, error) {
	models := d.GetModels().Interface()
	db := d.DB.Model(d.GetModel().Interface())
	db = getCondGormDB(db, queryCond)
	if err := db.Find(models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func getCondGormDB(db *gorm.DB, queryCond utils.QueryConditions) *gorm.DB {
	if len(queryCond.Query) > 0 {
		db = db.Where(queryCond.Query)
	}
	if len(queryCond.Fields) > 0 {
		db = db.Select(queryCond.Fields)
	}
	if len(queryCond.Order) > 0 {
		for _, order := range queryCond.Order {
			db = db.Order(order)
		}
	}
	if queryCond.Offset >= 0 {
		db = db.Offset(queryCond.Offset)
	}
	if queryCond.Limit > 0 {
		db = db.Limit(queryCond.Limit)
	}
	return db
}
