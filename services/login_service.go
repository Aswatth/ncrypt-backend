package services

import (
	"encoding/json"
	"errors"
	"ncrypt/models"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type LoginService struct {
}

func (obj *LoginService) Init() {
	godotenv.Load("../.env")
}

func (obj *LoginService) GetLoginData(name string) (*models.Login, error) {
	name = strings.ToUpper(name)
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("LOGIN_DB_NAME")))

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var login_data models.Login
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(name))

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			temp_data := append([]byte{}, val...)

			return json.Unmarshal(temp_data, &login_data)
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return &login_data, nil
}

func (obj *LoginService) GetAllLoginData() ([]models.Login, error) {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("LOGIN_DB_NAME")))

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var login_data_list []models.Login
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		// opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var login_data models.Login

				err := json.Unmarshal(v, &login_data)

				if err != nil {
					return err
				}

				login_data_list = append(login_data_list, login_data)

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
	return login_data_list, nil
}

func (obj *LoginService) setLoginData(login_data *models.Login) error {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("LOGIN_DB_NAME")))

	if err != nil {
		return err
	}
	defer db.Close()

	login_bytes, err := json.Marshal(login_data)

	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(strings.ToUpper(login_data.Name)), login_bytes)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (obj *LoginService) AddLoginData(new_login_data *models.Login) error {

	existing_data, err := obj.GetLoginData(new_login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if existing_data != nil {
		return errors.New(new_login_data.Name + " already exists")
	}

	return obj.setLoginData(new_login_data)
}

func (obj *LoginService) UpdateLoginData(name string, login_data *models.Login) error {
	if name != login_data.Name {
		existing_data, err := obj.GetLoginData(login_data.Name)

		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if existing_data != nil {
			return errors.New(login_data.Name + " already exists")
		}

		obj.DeleteLogin(name)
	}

	return obj.setLoginData(login_data)
}

func (obj *LoginService) DeleteLogin(name_to_delete string) error {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("LOGIN_DB_NAME")))

	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(strings.ToUpper(name_to_delete)))
		return err
	})

	if err != nil {
		return err
	}

	return nil
}
