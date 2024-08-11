package services

import (
	"errors"
	"ncrypt/models"
	"ncrypt/utils/database"
	"ncrypt/utils/logger"
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
		err = obj.setSystemData(models.SystemData{LoginCount: 0, LastLoginDateTime: "", CurrentLoginDateTime: "", IsLoggedIn: false, AutomaticBackup: false, AutomaticBackupLocation: "", BackupFileName: ""})
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

func (obj *SystemService) Login(password string) error {
	logger.Log.Printf("Logging in")
	result, err := obj.master_password_service.Validate(password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	if !result {
		logger.Log.Printf("ERROR: invalid password")
		return errors.New("invalid password")
	}

	system_data, err := obj.GetSystemData()
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	system_data.IsLoggedIn = true
	system_data.LoginCount += 1
	system_data.CurrentLoginDateTime = time.Now().Format(time.RFC3339)

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Logged in")

	return err
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
