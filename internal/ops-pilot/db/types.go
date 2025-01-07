package db

type DBType struct {
	Schema   string `json:"schema"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type MysqlDataDTO struct {
	Columns []string   `json:"columns"`
	Values  [][]string `json:"Values"`
}
