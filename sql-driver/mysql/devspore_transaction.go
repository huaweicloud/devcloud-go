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
 */

package mysql

import "database/sql/driver"

// devsporeTx implements driver.Tx interface.
type devsporeTx struct {
	dc       *devsporeConn
	actualTx driver.Tx
}

// Commit implements driver.Tx interface.
// Commit send transactionChan to clear transactionHolder before return.
func (tx *devsporeTx) Commit() (err error) {
	tx.dc.inTransaction = false
	return tx.actualTx.Commit()
}

// Rollback implements driver.Tx interface.
// Rollback send transactionChan to clear transactionHolder before return.
func (tx *devsporeTx) Rollback() (err error) {
	tx.dc.inTransaction = false
	return tx.actualTx.Rollback()
}
