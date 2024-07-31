package services

import (
	"ncrypt/utils/encryptor"
	"os"
	"testing"
)

func TestSetMasterPassword(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	err = os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))

	if err != nil {
		t.Error(err.Error())
	}
}

func TestSetMasterPassword_RESET(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	err = service.SetMasterPassword("123")

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.ValidateMasterPassword("123")

	if err != nil {
		t.Error(err.Error())
	}

	if !result {
		t.Errorf("Expected: %t\nActual: %t", true, result)
	}

	err = os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))

	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidateMasterPassword_PASS(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.ValidateMasterPassword("12345")

	if err != nil {
		t.Error(err.Error())
	}

	if !result {
		t.Errorf("Expected: %t\nActual: %t", true, result)
	}

	err = os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))

	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidateMasterPassword_FAIL(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.ValidateMasterPassword("123")

	if err != nil {
		if err.Error() != "invalid passowrd" {
			t.Error(err.Error())
		}
	}

	if result {
		t.Errorf("Expected: %t\nActual: %t", false, result)
	}

	err = os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))

	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetMasterPassword(t *testing.T) {
	password := "12345"

	hashed_password := encryptor.CreateHash(password)

	service := new(MasterPasswordService)
	service.Init()

	service.SetMasterPassword(password)

	stored_password, err := service.GetMasterPassword()

	if err != nil {
		t.Error(err.Error())
	}

	if stored_password != hashed_password {
		t.Error("Password mismatch")
	}

	err = os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))

	if err != nil {
		t.Error(err.Error())
	}
}
