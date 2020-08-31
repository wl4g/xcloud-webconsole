/**
 * Copyright 2017 ~ 2035 the original author or authors<Wanglsir@gmail.com, 983708408@qq.com>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"xcloud-webconsole/pkg/config"
	utils "xcloud-webconsole/pkg/utils"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlStore ...
type MysqlStore struct {
	mysqlDB *sql.DB
}

// NewMysqlStore ...
func NewMysqlStore() (*MysqlStore, error) {
	mysqlConfig := config.GlobalConfig.DataSource.Mysql
	log.Print("Connecting to MySQL of configuration: " + utils.ToJSONString(mysqlConfig))

	mydb, err := utils.OpenMysqlConnection(
		mysqlConfig.DbConnectStr,
		mysqlConfig.MaxOpenConns,
		mysqlConfig.MaxIdleConns,
		time.Duration(mysqlConfig.ConnMaxLifetimeSec) * time.Second,
	)
	//defer mysqlDB.Close(); // @see #Close()

	if err != nil {
		return nil, errors.New("Cannot connect to mysql. " + err.Error())
	}

	return &MysqlStore{
		mysqlDB: mydb,
	}, nil
}

// GetSessionByID find session info by id
func (that MysqlStore) GetSessionByID(id int64) *SessionBean {
	session := new(SessionBean)
	row := that.mysqlDB.QueryRow("select id,name,address,username,IFNULL(password, ''),IFNULL(ssh_key, '') from webconsole_session where id=?", id)
	if err := row.Scan(&session.ID, &session.Name, &session.Address, &session.Username, &session.Password, &session.SSHPrivateKey); err != nil {
		log.Fatal("GetSessionById", err)
	}
	fmt.Println(session.ID, session.Name, session.Username)
	return session
}

// QuerySessionList ...
func (that MysqlStore) QuerySessionList() []SessionBean {
	// 通过切片存储
	sessions := make([]SessionBean, 0)
	rows, _ := that.mysqlDB.Query("SELECT id,name,address,username,IFNULL(password, ''),IFNULL(ssh_key, '') FROM `webconsole_session` limit ?", 100)

	// 遍历
	var session SessionBean
	for rows.Next() {
		rows.Scan(&session.ID, &session.Name, &session.Address, &session.Username, &session.Password, &session.SSHPrivateKey)
		sessions = append(sessions, session)
	}

	fmt.Println(sessions)
	return sessions
}

// SaveSession ...
func (that MysqlStore) SaveSession(session *SessionBean) int64 {
	ret, e := that.mysqlDB.Exec("insert INTO webconsole_session(name,address,username,password,ssh_key) values(?,?,?,?,?)", session.Name, session.Address, session.Username, session.Password, session.SSHPrivateKey)
	if nil != e {
		log.Print("add fail", e)
		return 0
	}
	//影响行数
	rowsaffected, _ := ret.RowsAffected()
	id, _ := ret.LastInsertId()
	log.Printf("RowsAffected: %d", rowsaffected)
	return id
}

// DeleteSession ...
func (that MysqlStore) DeleteSession(ID int64) int64 {
	result, _ := that.mysqlDB.Exec("delete from webconsole_session where id=?", ID)
	rowsaffected, _ := result.RowsAffected()
	log.Printf("RowsAffected: %d", rowsaffected)

	return rowsaffected
}

// Close ...
func (that MysqlStore) Close() error {
	return that.mysqlDB.Close()
}