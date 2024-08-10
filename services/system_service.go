package services

import (
	"errors"
	"ncrypt/models"
	"ncrypt/utils/database"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type SystemService struct {
	database                database.IDatabase
	database_name           string
	master_password_service IMasterPasswordService
}

func (obj *SystemService) Init() {
	obj.database = &database.BadgerDb{}
	obj.database_name = "SYSTEM"
	obj.database.SetDatabase(obj.database_name)

	//Initialize system
	obj.initSystem()

	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()
}

func (obj *SystemService) initSystem() {
	_, err := obj.GetSystemData()

	if err != nil && err == badger.ErrKeyNotFound {
		err = obj.setSystemData(models.SystemData{LoginCount: 0, LastLoginDateTime: "", CurrentLoginDateTime: "", IsLoggedIn: false, AutomaticBackup: false, AutomaticBackupLocation: "", BackupFileName: ""})
	}
}

func (obj *SystemService) setSystemData(system_data models.SystemData) error {
	err := obj.database.AddData(obj.database_name, system_data)

	return err
}

func (obj *SystemService) GetSystemData() (*models.SystemData, error) {
	fetched_data, err := obj.database.GetData(obj.database_name)

	if err != nil {
		return nil, err
	}

	var system_data models.SystemData
	system_data.FromMap(fetched_data.(map[string]interface{}))

	return &system_data, err
}

func (obj *SystemService) Login(password string) error {
	result, err := obj.master_password_service.Validate(password)

	if err != nil {
		return err
	}

	if !result {
		return errors.New("invalid password")
	}

	system_data, err := obj.GetSystemData()
	if err != nil {
		return err
	}

	system_data.IsLoggedIn = true
	system_data.LoginCount += 1
	system_data.CurrentLoginDateTime = time.Now().Format(time.RFC3339)

	err = obj.setSystemData(*system_data)

	return err
}

func (obj *SystemService) Logout() error {
	system_data, err := obj.GetSystemData()
	if err != nil {
		return err
	}

	system_data.IsLoggedIn = false
	system_data.LastLoginDateTime = system_data.CurrentLoginDateTime

	err = obj.setSystemData(*system_data)

	return err
}
