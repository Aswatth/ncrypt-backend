package database

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
)

type BadgerDb struct {
	database_name string
}

func (obj *BadgerDb) SetDatabase(database_name string) error {
	obj.database_name = database_name

	return nil
}
func (obj *BadgerDb) GetData(table_name string, params ...string) (interface{}, error) {
	db, err := badger.Open(badger.DefaultOptions(obj.database_name))

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var fetched_data interface{}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(table_name))

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &fetched_data)

			return err
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return fetched_data, nil
}
func (obj *BadgerDb) GetAllData(table_name string, params ...string) ([]interface{}, error) {
	return nil, nil
}
func (obj *BadgerDb) AddData(table_name string, data interface{}) error {
	db, err := badger.Open(badger.DefaultOptions(obj.database_name))

	if err != nil {
		return err
	}
	defer db.Close()

	data_bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(table_name), data_bytes)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}
func (obj *BadgerDb) UpdateData(table_name string, data interface{}, params ...string) error {
	return nil
}
func (obj *BadgerDb) DeleteData(table_name string, params ...string) error {
	return nil
}
