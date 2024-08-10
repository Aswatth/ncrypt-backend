package services

import (
	"ncrypt/utils"
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

// Function to set master_password
func (obj *MasterPasswordService) SetMasterPassword(password string) error {

	password = encryptor.CreateHash(password)

	err := obj.database.AddData(os.Getenv("MASTER_PASSWORD_KEY"), password)

	return err
}

/*
	Update master_password

1. Decrypt all encrypted content using old master_password
2. Encrypt all encrpyed content using new master_passwordl̥
*/
func (obj *MasterPasswordService) UpdateMasterPassword(password string) error {
	old_password, err := obj.GetMasterPassword()

	if err != nil {
		return nil
	}

	err = obj.SetMasterPassword(password)

	new_password, err := obj.GetMasterPassword()

	if err != nil {
		return nil
	}

	//Setup broadcast to update encrypted data across services
	broadcast := utils.NewBroadcast()

	login_service := InitBadgerLoginService()
	login_service.Init()
	broadcast.Subscribe("UPDATE_MASTER_PASSWORD", login_service.recryptData)

	event_data := utils.Event{
		Type: "UPDATE_MASTER_PASSWORD",
		Data: map[string]string{
			"old_password": old_password,
			"new_password": new_password,
		},
	}

	broadcast.Publish(event_data)

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
