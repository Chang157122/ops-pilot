package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	sqlLiteFilePath = `db`
	SqlLiteDB       *sql.DB
)

func init() {
	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {

		panic(err)
	}

	SqlLiteDB = db
}
