package services

import (
	"errors"
	"ncrypt/models"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type LoginService struct {
	database                database.IDatabase
	master_password_service IMasterPasswordService
}

func (obj *LoginService) Init() {
	godotenv.Load("../.env")

	obj.database = database.InitBadgerDb()
	obj.database.SetDatabase(os.Getenv("LOGIN_DB_NAME"))

	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()
}

func (obj *LoginService) GetLoginData(name string) (models.Login, error) {
	name = strings.ToUpper(name)

	fetched_data, err := obj.database.GetData(name)

	if err != nil {
		return models.Login{}, err
	}

	var login_data models.Login
	login_data.FromMap(fetched_data.(map[string]interface{}))

	return login_data, err
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
	var login_data_list []models.Login

	result_list, err := obj.database.GetAllData()

	if err != nil {
		return nil, err
	}

	for _, result := range result_list {
		var login_data models.Login
		login_data.FromMap(result.(map[string]interface{}))

		login_data_list = append(login_data_list, login_data)
	}

	return login_data_list, nil
}

func (obj *LoginService) setLoginData(login_data *models.Login) error {

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

	err = obj.database.AddData(strings.ToUpper(login_data.Name), login_data)

	return err
}

func (obj *LoginService) AddLoginData(login_data *models.Login) error {

	existing_data, err := obj.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if existing_data.Name != "" {
		return errors.New(login_data.Name + " already exists")
	}

	return obj.setLoginData(login_data)
}

func (obj *LoginService) UpdateLoginData(name string, login_data *models.Login) error {
	if name != login_data.Name {
		existing_data, err := obj.GetLoginData(login_data.Name)

		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if existing_data.Name != "" {
			return errors.New(login_data.Name + " already exists")
		}

		err = obj.DeleteLoginData(name)

		if err != nil {
			return err
		}
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

func (obj *LoginService) DeleteLoginData(name string) error {
	err := obj.database.DeleteData(strings.ToUpper(name))

	return err
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
