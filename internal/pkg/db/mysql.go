package db

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

type MysqlImpl struct {
	client *sqlx.DB
}

func (m MysqlImpl) Query(sql []byte) ([]byte, error) {
	s := string(sql)
	split := strings.Split(s, " ")
	if !strings.Contains(split[0], "SELECT") || !strings.Contains(split[0], "select") {
		panic(errors.New("查询接口不执行删除"))
	}

	rows, err := m.client.Query(s)
	if err != nil {
		return nil, err
	}
	// 读出查询出的字段
	cols, _ := rows.Columns()
	// 每列的数据
	values := make([][]byte, len(cols))
	// query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}

	results := make([]map[string]string, 0, 10)
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]string, 10)
		for k, v := range values {
			key := cols[k]
			value := string(v)
			row[key] = value
		}
		results = append(results, row)
	}

	dto := MysqlDataDTO{
		Columns: cols,
		Values:  results,
	}
	marshal, err := json.Marshal(&dto)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func NewMysqlImpl(dbType DBType) *MysqlImpl {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbType.Username,
		dbType.Password,
		dbType.Address,
		dbType.Port,
		dbType.Database)
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	return &MysqlImpl{client: db}
}
