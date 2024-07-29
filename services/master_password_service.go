package services

import (
	"ncrypt/utils/encryptor"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type MasterPasswordService struct {
}

func (obj *MasterPasswordService) Init() {
	godotenv.Load("../.env")
}

func (obj *MasterPasswordService) SetMasterPassword(password string) error {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("MASTER_PASSWORD_DB_NAME")))

	if err != nil {
		return err
	}
	defer db.Close()

	password = encryptor.CreateHash(password)

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(os.Getenv("MASTER_PASSOWRD_KEY")), []byte(password))

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (obj *MasterPasswordService) ValidateMasterPassword(password string) (bool, error) {
	stored_password, err := obj.GetMasterPassword()

	if err != nil {
		return false, err
	}

	return stored_password == encryptor.CreateHash(password), nil
}

func (obj *MasterPasswordService) GetMasterPassword() (string, error) {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("MASTER_PASSWORD_DB_NAME")))

	if err != nil {
		return "", err
	}
	defer db.Close()

	var stored_password string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(os.Getenv("MASTER_PASSOWRD_KEY")))

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			temp_data := append([]byte{}, val...)

			stored_password = string(temp_data[:])

			return nil
		})

		return err
	})

	if err != nil {
		return "", err
	}

	return stored_password, nil
}
