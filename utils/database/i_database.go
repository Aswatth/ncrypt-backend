package database

type IDatabase interface {
	SetDatabase(database_name string) error
	GetData(table_name string, params ...string) (interface{}, error)
	GetAllData(params ...string) ([]interface{}, error)
	AddData(table_name string, data interface{}) error
	UpdateData(table_name string, data interface{}, params ...string) error
	DeleteData(table_name string, params ...string) error
}

func InitBadgerDb() IDatabase {
	return &BadgerDb{}
}
