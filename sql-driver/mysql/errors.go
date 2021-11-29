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

import (
	"github.com/go-sql-driver/mysql"
)

var recoverableErrorCodesMap = map[int]string{
	// sqlState start with '08' which means connection error, refer to https://mariadb.com/kb/en/mariadb-error-codes/#shared-mariadbmysql-error-codes
	1040: "ER_CON_COUNT_ERROR",
	1042: "ER_BAD_HOST_ERROR",
	1043: "ER_HANDSHAKE_ERROR",
	1047: "ER_UNKNOWN_COM_ERROR",
	1053: "ER_SERVER_SHUTDOWN",
	1080: "ER_FORCING_CLOSE",
	1081: "ER_IPSOCK_ERROR",
	1152: "ER_ABORTING_CONNECTION",
	1153: "ER_NET_PACKET_TOO_LARGE",
	1154: "ER_NET_READ_ERROR_FROM_PIPE",
	1155: "ER_NET_FCNTL_ERROR",
	1156: "ER_NET_PACKETS_OUT_OF_ORDER",
	1157: "ER_NET_UNCOMPRESS_ERROR",
	1158: "ER_NET_READ_ERROR",
	1159: "ER_NET_READ_INTERRUPTED",
	1160: "ER_NET_ERROR_ON_WRITE",
	1161: "ER_NET_WRITE_INTERRUPTED",
	1184: "ER_NEW_ABORTING_CONNECTION",
	1189: "ER_MASTER_NET_READ",
	1190: "ER_MASTER_NET_WRITE",
	1218: "ER_CONNECT_TO_MASTER",

	// Communications Errors
	1129: "ER_HOST_IS_BLOCKED",
	1130: "ER_HOST_NOT_PRIVILEGED",

	// Authentication Errors
	1045: "ER_ACCESS_DENIED_ERROR",

	// Resource Errors
	1004: "ER_CANT_CREATE_FILE",
	1005: "ER_CANT_CREATE_TABLE",
	1015: "ER_CANT_LOCK",
	1021: "ER_DISK_FULL",
	1041: "ER_OUT_OF_RESOURCES",

	// Out-of-memory errors
	1037: "ER_OUTOFMEMORY",
	1038: "ER_OUT_OF_SORTMEMORY",

	// Access denied
	1142: "ER_TABLEACCESS_DENIED_ERROR",
	1227: "ER_SPECIFIC_ACCESS_DENIED_ERROR",

	1023: "ER_ERROR_ON_CLOSE",

	1290: "ER_OPTION_PREVENTS_STATEMENT",
}

func isErrorRecoverable(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if _, exist := recoverableErrorCodesMap[int(mysqlError.Number)]; exist {
			return true
		}
	}
	return false
}
