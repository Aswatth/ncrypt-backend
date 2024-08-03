package services

import (
	"encoding/json"
	"ncrypt/models"

	"github.com/dgraph-io/badger/v4"
)

type SystemService struct {
	system_db_name string
}

func (obj *SystemService) Init() {
	obj.system_db_name = "SYSTEM"
}

func (obj *SystemService) setSystemData(system_data models.SystemData) error {
	db, err := badger.Open(badger.DefaultOptions(obj.system_db_name))

	if err != nil {
		return err
	}
	defer db.Close()

	system_data_bytes, err := json.Marshal(system_data)

	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(obj.system_db_name), system_data_bytes)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (obj *SystemService) GetSystemData() (*models.SystemData, error) {
	db, err := badger.Open(badger.DefaultOptions(obj.system_db_name))

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var system_data models.SystemData
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(obj.system_db_name))

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			temp_data := append([]byte{}, val...)

			err := json.Unmarshal(temp_data[:], &system_data)

			if err != nil {
				return err
			}

			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return &system_data, nil
}
