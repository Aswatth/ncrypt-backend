package services

import (
	"errors"
	"ncrypt/models"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"ncrypt/utils/logger"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type LoginDataService struct {
	database                database.IDatabase
	master_password_service IMasterPasswordService
}

func (obj *LoginDataService) Init() {
	logger.Log.Printf("Initializing login service")
	logger.Log.Printf("Loading .env variables")
	godotenv.Load("../.env")

	logger.Log.Printf("Setting up database")
	obj.database = database.InitBadgerDb()
	obj.database.SetDatabase(os.Getenv("LOGIN_DB_NAME"))

	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()

	logger.Log.Printf("DONE")
}

func (obj *LoginDataService) GetLoginData(login_data_name string) (models.Login, error) {
	logger.Log.Printf("Getting login data")
	login_data_name = strings.ToUpper(login_data_name)

	fetched_data, err := obj.database.GetData(login_data_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return models.Login{}, err
	}

	var login_data models.Login
	login_data.FromMap(fetched_data.(map[string]interface{}))

	logger.Log.Printf("DONE")
	return login_data, err
}

func (obj *LoginDataService) GetDecryptedAccountPassword(login_data_name string, account_username string) (string, error) {
	fetched_login_data, err := obj.GetLoginData(login_data_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	logger.Log.Printf("Decrypting login data")
	var decrypted_password string
	for _, account := range fetched_login_data.Accounts {
		if account.Username == account_username {
			master_password_hash, err := obj.master_password_service.GetMasterPassword()

			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
				return "", err
			}

			decrypted_password, err = encryptor.Decrypt(account.Password, master_password_hash+login_data_name+account_username)

			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
				return "", err
			}

			break
		}
	}

	if decrypted_password == "" {
		logger.Log.Printf("ERROR: account username not found")
		return "", errors.New("account username not found")
	}

	logger.Log.Printf("DONE")
	return decrypted_password, nil
}

func (obj *LoginDataService) GetAllLoginData() ([]models.Login, error) {
	logger.Log.Printf("Getting all login data")
	var login_data_list []models.Login

	result_list, err := obj.database.GetAllData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil, err
	}

	for _, result := range result_list {
		var login_data models.Login
		login_data.FromMap(result.(map[string]interface{}))

		login_data_list = append(login_data_list, login_data)
	}

	logger.Log.Printf("DONE")
	return login_data_list, nil
}

func (obj *LoginDataService) setLoginData(login_data *models.Login) error {

	logger.Log.Printf("Checking for duplicate accounts")
	//Check for duplicate account-username
	account_username_map := make(map[string]bool)

	for _, account := range login_data.Accounts {
		if !account_username_map[account.Username] {
			account_username_map[account.Username] = true
		} else {
			return errors.New("duplicate username " + account.Username)
		}
	}

	logger.Log.Printf("Encrypting data")
	//Encrypt login_data - account_passwords
	// Get master password
	master_password_hash, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		if strings.ToUpper(err.Error()) == "KEY NOT FOUND" {
			err = errors.New("master_password not set")
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
		return err
	}

	for index := range len(login_data.Accounts) {
		login_data.Accounts[index].Password, _ = encryptor.Encrypt(login_data.Accounts[index].Password, master_password_hash+login_data.Name+login_data.Accounts[index].Username)
	}

	err = obj.database.AddData(strings.ToUpper(login_data.Name), login_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	logger.Log.Printf("Saved to database!")
	return err
}

func (obj *LoginDataService) AddLoginData(login_data *models.Login) error {
	logger.Log.Printf("Adding login data")
	logger.Log.Printf("Checking for duplicate data")
	existing_data, err := obj.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}
	if existing_data.Name != "" {
		logger.Log.Printf("ERROR: %s", "CONFLICTING NAMES")
		return errors.New(login_data.Name + " already exists")
	}

	err = obj.setLoginData(login_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	logger.Log.Printf("DONE")
	return err
}

func (obj *LoginDataService) UpdateLoginData(old_login_data_name string, login_data *models.Login) error {
	logger.Log.Printf("Updating login data")

	logger.Log.Printf("Checking for name conflicts")
	if old_login_data_name != login_data.Name {
		existing_data, err := obj.GetLoginData(login_data.Name)

		if err != nil && err != badger.ErrKeyNotFound {
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
		if existing_data.Name != "" {
			err = errors.New(login_data.Name + " already exists")
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}

		err = obj.DeleteLoginData(old_login_data_name)

		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
	}

	logger.Log.Printf("Decrypting data")

	key, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	for index := range len(login_data.Accounts) {
		decrypted_data, err := encryptor.Decrypt(login_data.Accounts[index].Password, key+old_login_data_name+login_data.Accounts[index].Username)

		// login_data.Accounts[index].Password = decrypted_data
		if err == nil {
			login_data.Accounts[index].Password = decrypted_data
		}
	}

	err = obj.setLoginData(login_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	logger.Log.Printf("DONE")

	return err
}

func (obj *LoginDataService) DeleteLoginData(login_data_name string) error {
	logger.Log.Printf("Deleting login data")
	err := obj.database.DeleteData(strings.ToUpper(login_data_name))
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}
	logger.Log.Printf("DONE")
	return err
}

func (obj *LoginDataService) recryptData(password_data map[string]string) error {
	logger.Log.Printf("Re-crpyting login data")

	//Get all login data
	login_list, err := obj.GetAllLoginData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	//Decrypt all login data
	old_password := password_data["OLD_PASSWORD"]
	new_password := password_data["NEW_PASSWORD"]

	for i := range len(login_list) {
		for j := range len(login_list[i].Accounts) {
			login_list[i].Accounts[j].Password, err = encryptor.Decrypt(login_list[i].Accounts[j].Password, old_password+login_list[i].Name+login_list[i].Accounts[j].Username)
			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
				return err
			}
		}
	}

	//Save updated data
	for i := range len(login_list) {
		for j := range len(login_list[i].Accounts) {
			login_list[i].Accounts[j].Password, err = encryptor.Encrypt(login_list[i].Accounts[j].Password, new_password+login_list[i].Name+login_list[i].Accounts[j].Username)

			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
				return err
			}

			obj.database.UpdateData(login_list[i].Name, login_list[i])
		}
	}

	logger.Log.Printf("DONE")
	return nil
}

func (obj *LoginDataService) importData(login_data_list []models.Login) error {

	os.RemoveAll("data/" + os.Getenv("LOGIN_DB_NAME"))

	for _, login_data := range login_data_list {
		err := obj.database.AddData(strings.ToUpper(login_data.Name), login_data)

		if err != nil {
			return err
		}
	}

	return nil
}
