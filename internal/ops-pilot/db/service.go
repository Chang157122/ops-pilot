package db

type DBService interface {
	Query(sql []byte) ([]byte, error)
}

func NewDBService(dbType DBType) DBService {
	switch dbType.Schema {
	case "mysql":
		return NewMysqlImpl(dbType)
	default:
		return nil
	}
}
