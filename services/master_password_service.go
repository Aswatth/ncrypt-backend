package services

import (
	"ncrypt/models"
	"ncrypt/utils/encryptor"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type MasterPasswordService struct {
}

func (obj *MasterPasswordService) Init() {
	godotenv.Load("../.env")
}

// helper function to set master_password
func (obj *MasterPasswordService) setMasterPassword(password string) error {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("MASTER_PASSWORD_DB_NAME")))

	if err != nil {
		return err
	}
	defer db.Close()

	password = encryptor.CreateHash(password)

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(os.Getenv("MASTER_PASSWORD_KEY")), []byte(password))

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

// Sets up master_password for the very first time and also creates system->login_info if not found
func (obj *MasterPasswordService) SetMasterPassword(password string) error {
	err := obj.setMasterPassword(password)

	if err != nil {
		return err
	}

	system_service := new(SystemService)
	system_service.Init()

	_, err = system_service.GetSystemData()

	if err != nil {
		if err == badger.ErrKeyNotFound {
			err := system_service.setSystemData(models.SystemData{Login_count: 1, Last_login: time.Now().Format(time.RFC3339)})
			if err != nil {
				return nil
			}
		} else {
			return nil
		}

	}

	return nil
}

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

func (obj *MasterPasswordService) ValidateMasterPassword(password string, is_login ...bool) (bool, error) {

	stored_password, err := obj.GetMasterPassword()

	if err != nil {
		return false, err
	}

	result := stored_password == encryptor.CreateHash(password)

	//If it is a validation for logging in, then update system's last_login_date_time
	if len(is_login) == 1 {
		if is_login[0] && result {
			system_service := new(SystemService)
			system_service.Init()

			system_data, err := system_service.GetSystemData()

			if err != nil {
				return false, err
			}

			system_data.Login_count += 1
			now := time.Now()
			system_data.Last_login = now.Format(time.RFC3339)

			err = system_service.setSystemData(*system_data)

			if err != nil {
				return false, err
			}
		}
	}

	return result, nil
}

func (obj *MasterPasswordService) GetMasterPassword() (string, error) {
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("MASTER_PASSWORD_DB_NAME")))

	if err != nil {
		return "", err
	}
	defer db.Close()

	var stored_password string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(os.Getenv("MASTER_PASSWORD_KEY")))

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
