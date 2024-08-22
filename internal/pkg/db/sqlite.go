package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var sqlLiteFilePath = `db`

func InitSqlLiteDB() {
	open, err := sql.Open("sqlite3", ":memory:")
	if err != nil {

		panic(err)
	}
	defer open.Close()

}
