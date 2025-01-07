package db

import (
	"fmt"
	"testing"
)

func TestNewMysql(t *testing.T) {
	impl := NewMysqlImpl(DBType{Schema: `mysql`, Username: `im`, Password: `xxxx`, Database: `shinemo_im`, Port: `3333`, Address: `1xxxx1`})
	query, err := impl.client.Query("show tables")
	if err != nil {
		t.Fatal(err)
	}

	columns, err := query.Columns()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(columns)

}
