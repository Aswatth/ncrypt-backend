package services

import (
	"encoding/json"
	"errors"
	"ncrypt/models"
	"ncrypt/utils/encryptor"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type LoginService struct {
	master_password_service MasterPasswordService
}

func (obj *LoginService) Init() {
	godotenv.Load("../.env")
	master_password_service := new(MasterPasswordService)

	obj.master_password_service = *master_password_service
	obj.master_password_service.Init()
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

func (obj *LoginService) GetDecryptedAccountPassword(login_data_name string, account_username string) (string, error) {
	fetched_login_data, err := obj.GetLoginData(login_data_name)

	if err != nil {
		return "", err
	}

	var decrypted_password string
	for _, account := range fetched_login_data.Accounts {
		if account.Username == account_username {
			master_password_hash, err := obj.master_password_service.GetMasterPassword()

			if err != nil {
				return "", err
			}

			decrypted_password, err = encryptor.Decrypt(account.Password, master_password_hash)

			if err != nil {
				return "", err
			}

			break
		}
	}

	if decrypted_password == "" {
		return "", errors.New("account username not found")
	}

	return decrypted_password, nil
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

	//Check for duplicate account-username
	account_username_map := make(map[string]bool)

	for _, account := range login_data.Accounts {
		if !account_username_map[account.Username] {
			account_username_map[account.Username] = true
		} else {
			return errors.New("duplicate username " + account.Username)
		}
	}

	//Encrypt login_data - account_passwords
	// Get master password
	master_password_hash, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		if strings.ToUpper(err.Error()) == "KEY NOT FOUND" {
			return errors.New("master_password not set")
		}
		return err
	}

	for index := range len(login_data.Accounts) {
		login_data.Accounts[index].Password, _ = encryptor.Encrypt(login_data.Accounts[index].Password, master_password_hash)
	}

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

	//Decrypt data
	master_password_service := new(MasterPasswordService)
	master_password_service.Init()

	key, err := master_password_service.GetMasterPassword()

	if err != nil {
		return err
	}

	for index := range len(login_data.Accounts) {
		decrypted_data, err := encryptor.Decrypt(login_data.Accounts[index].Password, key)

		if err == nil {
			login_data.Accounts[index].Password = decrypted_data
		}
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

func (obj *LoginService) decryptData(login_data models.Login, key string) *models.Login {
	for index := range len(login_data.Accounts) {
		login_data.Accounts[index].Password, _ = encryptor.Decrypt(login_data.Accounts[index].Password, key)
	}

	return &login_data
}

func (obj *LoginService) encryptData(login_data models.Login, key string) *models.Login {
	for index := range len(login_data.Accounts) {
		login_data.Accounts[index].Password, _ = encryptor.Encrypt(login_data.Accounts[index].Password, key)
	}

	return &login_data
}

func (obj *LoginService) decryptAllData(key string) ([]models.Login, error) {

	login_data_list, err := obj.GetAllLoginData()

	if err != nil {
		return nil, err
	}

	for index := range len(login_data_list) {
		login_data_list[index] = *obj.decryptData(login_data_list[index], key)
	}

	return login_data_list, nil
}

func (obj *LoginService) encrytAllData(login_data_list []models.Login) error {

	for index := range len(login_data_list) {
		err := obj.setLoginData(&login_data_list[index])
		if err != nil {
			return err
		}
	}

	return nil
}
