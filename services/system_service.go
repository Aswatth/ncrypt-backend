package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"ncrypt/models"
	"ncrypt/utils"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type SystemService struct {
	database                database.IDatabase
	database_name           string
	master_password_service IMasterPasswordService
}

func (obj *SystemService) Init() {
	logger.Log.Printf("Initializing system service")
	logger.Log.Printf("Setting up database")
	obj.database = &database.BadgerDb{}
	obj.database_name = "SYSTEM"
	obj.database.SetDatabase(obj.database_name)

	//Initialize system
	logger.Log.Printf("Setting up intial data")
	obj.initSystem()

	logger.Log.Printf("Setting up master password service")
	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()
	logger.Log.Printf("System service initialized")
}

func (obj *SystemService) initSystem() {
	_, err := obj.GetSystemData()

	if err != nil && err == badger.ErrKeyNotFound {
		err = obj.setSystemData(models.SystemData{LoginCount: 0, LastLoginDateTime: "", CurrentLoginDateTime: "", IsLoggedIn: false, AutomaticBackup: false, AutomaticBackupLocation: "", BackupFileName: "", SessionTimeInMinutes: 20})
		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
		}
	}
}

func (obj *SystemService) setSystemData(system_data models.SystemData) error {
	logger.Log.Printf("Setting system data")
	err := obj.database.AddData(obj.database_name, system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	return err
}

func (obj *SystemService) GetSystemData() (*models.SystemData, error) {
	logger.Log.Printf("Getting system data")
	fetched_data, err := obj.database.GetData(obj.database_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil, err
	}

	var system_data models.SystemData
	system_data.FromMap(fetched_data.(map[string]interface{}))

	return &system_data, err
}

func (obj *SystemService) Login(password string) (string, error) {
	logger.Log.Printf("Logging in")
	result, err := obj.master_password_service.Validate(password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		if err.Error() == "Key not found" {
			token, err := jwt.ShortLivedToken()

			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
				return "", err
			}
			return token, nil
		}
		return "", err
	}

	if !result {
		logger.Log.Printf("ERROR: invalid password")
		return "", errors.New("invalid password")
	}

	system_data, err := obj.GetSystemData()
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	system_data.IsLoggedIn = true
	system_data.LoginCount += 1
	system_data.CurrentLoginDateTime = time.Now().Format(time.RFC3339)

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	logger.Log.Printf("Logged in")

	token, err := jwt.GenerateToken()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	return token, nil
}

func (obj *SystemService) Logout() error {
	logger.Log.Printf("Logging out")

	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	system_data.IsLoggedIn = false
	system_data.LastLoginDateTime = system_data.CurrentLoginDateTime

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Logged out")
	return err
}

func (obj *SystemService) Export(file_name string, file_path string) error {
	logger.Log.Println("Exporting data...")
	//Get system data
	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Fetching master password")
	//Get master password data
	master_password, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Fetching login data")
	//Get login data
	login_service := InitBadgerLoginService()
	login_service.Init()
	login_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	export_data := new(ExportData)

	export_data.SYSTEM_DATA = *system_data
	export_data.LOGIN_DATA = login_data_list
	export_data.MASTER_PASSWORD = master_password

	logger.Log.Println("Exporting to " + file_path + "\\" + file_name)
	file, err := os.Create(file_path + "\\" + file_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	defer file.Close()

	export_data_bytes, err := json.Marshal(export_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Encrypting export data")
	//Encrpyt data using master_password
	encrypted_export_data, err := encryptor.Encrypt(base64.StdEncoding.EncodeToString(export_data_bytes), master_password)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	encrypted_export_data_bytes, err := base64.StdEncoding.DecodeString(encrypted_export_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Saving to file")
	_, err = file.Write(encrypted_export_data_bytes)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Export complete!")
	return nil
}

func (obj *SystemService) Import(file_name string, file_path string, master_password string) error {
	logger.Log.Println("Importing data")
	file, err := os.Open(file_path + "\\" + file_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}
	defer file.Close()

	logger.Log.Println("Reading import file")
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Decrypting import content")
	decrypted_data, err := encryptor.Decrypt(base64.StdEncoding.EncodeToString(data), encryptor.CreateHash(master_password))
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	decrypted_data_bytes, err := base64.StdEncoding.DecodeString(decrypted_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return errors.New("incorrect master password or corrupted file")
	}

	logger.Log.Println("Importing system data")
	imported_data := new(ExportData)
	json.Unmarshal(decrypted_data_bytes, &imported_data)

	//Import system data
	logger.Log.Println("Importing system data")
	err = obj.setSystemData(imported_data.SYSTEM_DATA)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	//Import master password
	logger.Log.Println("Importing master password")
	err = obj.master_password_service.importData(imported_data.MASTER_PASSWORD)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	//Import login data
	logger.Log.Println("Importing login data")
	login_service := InitBadgerLoginService()
	login_service.Init()
	login_service.importData(imported_data.LOGIN_DATA)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("DONE")
	return nil
}

type ExportData struct {
	SYSTEM_DATA     models.SystemData `json:"SYSTEM" bson:"SYSTEM"`
	LOGIN_DATA      []models.Login    `json:"LOGIN_DATA" bson:"LOGIN_DATA"`
	MASTER_PASSWORD string            `json:"MASTER_PASSWORD" bson:"MASTER_PASSWORD"`
}

func (obj *SystemService) GeneratePassword(has_digits bool, has_upper_case bool, has_special_char bool, length int) string {
	return utils.GeneratePassword(has_digits, has_upper_case, has_special_char, length)
}
