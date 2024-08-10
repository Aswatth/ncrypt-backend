package services

import (
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"os"

	"github.com/joho/godotenv"
)

type MasterPasswordService struct {
	database database.IDatabase
}

// Initialize Master password service
func (obj *MasterPasswordService) Init() {
	//Load env
	godotenv.Load("../.env")

	//Initialize database
	obj.database = database.InitBadgerDb()
	obj.database.SetDatabase(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}

// helper function to set master_password
func (obj *MasterPasswordService) setMasterPassword(password string) error {

	password = encryptor.CreateHash(password)

	err := obj.database.AddData(os.Getenv("MASTER_PASSWORD_KEY"), password)

	return err
}

// Sets up master_password for the very first time and also creates system->login_info if not found
func (obj *MasterPasswordService) SetMasterPassword(password string) error {
	err := obj.setMasterPassword(password)

	return err
}

/*
	Update master_password

1. Decrypt all encrypted content using old master_password
2. Encrypt all encrpyed content using new master_passwordl̥
*/
func (obj *MasterPasswordService) UpdateMasterPassword(password string) error {

	key, err := obj.GetMasterPassword()

	if err != nil {
		return nil
	}

	//Decrypt all encrpyted data using old password
	login_service := new(LoginService)
	login_service.Init()
	login_list, err := login_service.decryptAllData(key)

	if err != nil {
		return nil
	}

	err = obj.setMasterPassword(password)

	if err != nil {
		return nil
	}

	//Encrypt all login data using new password
	err = login_service.encrytAllData(login_list)

	if err != nil {
		return nil
	}

	return nil
}

func (obj *MasterPasswordService) Validate(password string) (bool, error) {

	stored_password, err := obj.GetMasterPassword()

	if err != nil {
		return false, err
	}

	result := stored_password == encryptor.CreateHash(password)

	return result, nil
}

func (obj *MasterPasswordService) GetMasterPassword() (string, error) {
	fetched_data, err := obj.database.GetData(os.Getenv("MASTER_PASSWORD_KEY"))

	if err != nil {
		return "", err
	}

	return fetched_data.(string), err
}
