package database

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
)

type BadgerDb struct {
	database_name string
}

func (obj *BadgerDb) SetDatabase(database_name string) error {
	obj.database_name = "data/" + database_name

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
func (obj *BadgerDb) GetAllData(params ...string) ([]interface{}, error) {
	db, err := badger.Open(badger.DefaultOptions(obj.database_name))

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var result_list []interface{}
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		// opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var result interface{}

				err := json.Unmarshal(v, &result)

				if err != nil {
					return err
				}

				result_list = append(result_list, result)

				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result_list, nil
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
	db, err := badger.Open(badger.DefaultOptions(obj.database_name))

	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(table_name))
		return err
	})

	return err
}
