package services

import (
	"ncrypt/utils"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"ncrypt/utils/logger"
	"os"

	"github.com/joho/godotenv"
)

type MasterPasswordService struct {
	database database.IDatabase
}

// Initialize Master password service
func (obj *MasterPasswordService) Init() {
	logger.Log.Printf("Initiazling Master password service")

	logger.Log.Printf("Loading .env variables")
	//Load env
	godotenv.Load("../.env")

	logger.Log.Printf("Setting up db")
	//Initialize database
	obj.database = database.InitBadgerDb()
	obj.database.SetDatabase(os.Getenv("MASTER_PASSWORD_DB_NAME"))
	logger.Log.Printf("Master password service initialized")
}

// Function to set master_password
func (obj *MasterPasswordService) SetMasterPassword(password string) error {
	logger.Log.Printf("Setting master pasword")
	password = encryptor.CreateHash(password)
	logger.Log.Printf("Created hash")

	err := obj.database.AddData(os.Getenv("MASTER_PASSWORD_KEY"), password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	logger.Log.Printf("Saved to database!")
	return err
}

/*
	Update master_password

1. Decrypt all encrypted content using old master_password
2. Encrypt all encrpyed content using new master_passwordl̥
*/
func (obj *MasterPasswordService) UpdateMasterPassword(password string) error {
	logger.Log.Printf("Updating master pasword")
	old_password, err := obj.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil
	}

	err = obj.SetMasterPassword(password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil
	}

	new_password, err := obj.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil
	}

	logger.Log.Printf("Creating new broadcast")
	//Setup broadcast to update encrypted data across services
	broadcast := utils.NewBroadcast()

	logger.Log.Printf("Subscribing login service to listen for changes")
	login_service := InitBadgerLoginService()
	login_service.Init()
	broadcast.Subscribe("UPDATE_MASTER_PASSWORD", login_service.recryptData)

	logger.Log.Printf("Creating event")
	event_data := utils.Event{
		Type: "UPDATE_MASTER_PASSWORD",
		Data: map[string]string{
			"old_password": old_password,
			"new_password": new_password,
		},
	}

	logger.Log.Printf("Broadcasting event")
	broadcast.Publish(event_data)

	logger.Log.Printf("Master password updated!")

	return nil
}

func (obj *MasterPasswordService) Validate(password string) (bool, error) {

	logger.Log.Printf("Validating master password")
	stored_password, err := obj.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return false, err
	}

	logger.Log.Printf("Comparing password...")
	result := stored_password == encryptor.CreateHash(password)

	logger.Log.Printf("Validation completed!")
	return result, nil
}

func (obj *MasterPasswordService) GetMasterPassword() (string, error) {
	logger.Log.Printf("Getting master password")
	fetched_data, err := obj.database.GetData(os.Getenv("MASTER_PASSWORD_KEY"))

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	return fetched_data.(string), err
}

func (obj *MasterPasswordService) importData(password string) error {
	err := obj.database.AddData(os.Getenv("MASTER_PASSWORD_KEY"), password)

	return err
}
