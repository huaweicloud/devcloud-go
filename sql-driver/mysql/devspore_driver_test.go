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
	"database/sql"
	"testing"

	"github.com/huaweicloud/devcloud-go/mock"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/config"
	"github.com/huaweicloud/devcloud-go/sql-driver/rds/datasource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mysql")
}

var _ = Describe("CRUD", func() {
	var (
		devsporeDB *sql.DB
		masterDB   *sql.DB
		err        error
		activeNode *datasource.NodeDataSource
	)
	metadata := mock.MysqlMock{
		User:      "XXXX",
		Password:  "XXXX",
		Address:   "127.0.0.1:13306",
		Databases: []string{"ds0", "ds0-slave0", "ds0-slave1", "ds1", "ds1-slave0", "ds1-slave1"},
	}
	metadata.StartMockMysql()

	BeforeEach(func() {
		devsporeDB, err = sql.Open("devspore_mysql", "../rds/resources/driver_test_config.yaml")
		Expect(err).NotTo(HaveOccurred())
		activeNode, err = initDB()
		Expect(err).NotTo(HaveOccurred())
		masterDB, err = sql.Open("mysql", activeNode.MasterDataSource.Dsn)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(devsporeDB.Close()).NotTo(HaveOccurred())
		Expect(masterDB.Close()).NotTo(HaveOccurred())
	})

	It("Test Query", func() {
		var (
			val  string
			flag bool
		)
		err = devsporeDB.QueryRow("SELECT val FROM foo WHERE id=?", id1).Scan(&val)
		Expect(err).NotTo(HaveOccurred())
		for _, slave := range activeNode.SlavesDatasource {
			if slave.Name == val {
				flag = true
			}
		}
		Expect(flag).To(Equal(true))
	})

	It("Test Insert", func() {
		var val string
		_, err = devsporeDB.Exec(`INSERT INTO foo (id, val) VALUES (?, ?)`, id2, "insert")
		Expect(err).NotTo(HaveOccurred())
		err = masterDB.QueryRow("SELECT val FROM foo WHERE id=?", id2).Scan(&val)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("insert"))
	})

	It("Test Update", func() {
		var val string
		_, err = devsporeDB.Exec(`UPDATE foo set val=? where id=?`, "update", id1)
		Expect(err).NotTo(HaveOccurred())
		err = masterDB.QueryRow("SELECT val FROM foo WHERE id=?", id1).Scan(&val)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("update"))
	})

	It("Test Delete", func() {
		var val string
		_, err = devsporeDB.Exec(`DELETE FROM foo where id=?`, id1)
		Expect(err).NotTo(HaveOccurred())
		err = masterDB.QueryRow("SELECT val FROM foo WHERE id=?", id1).Scan(&val)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("sql: no rows in result set"))
	})
})

func initDB() (*datasource.NodeDataSource, error) {
	yamlConfigPath := "../rds/resources/driver_test_config.yaml"
	// parse yaml config
	clusterConfiguration, err := config.Unmarshal(yamlConfigPath)
	if err != nil {
		return nil, err
	}
	clusterDataSource, _ := datasource.NewClusterDataSource(clusterConfiguration)
	for _, nodeDataSource := range clusterDataSource.DataSources {
		if err = createTable(nodeDataSource.MasterDataSource); err != nil {
			return nil, err
		}
		for _, slave := range nodeDataSource.SlavesDatasource {
			if err = createTable(slave); err != nil {
				return nil, err
			}
		}
	}
	return clusterDataSource.DataSources[clusterDataSource.Active], nil
}

var (
	id1 = 1
	id2 = 2
)

func createTable(actualDataSource *datasource.ActualDataSource) error {
	db, err := sql.Open("mysql", actualDataSource.Dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec("DROP TABLE IF EXISTS foo"); err != nil {
		return err
	}

	if _, err = db.Exec("CREATE TABLE foo (id INT PRIMARY KEY, val CHAR(50))"); err != nil {
		return err
	}

	if _, err = db.Exec(`INSERT INTO foo VALUES (?, ?)`, id1, actualDataSource.Name); err != nil {
		return err
	}
	return nil
}
