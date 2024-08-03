package services

import (
	"ncrypt/utils/encryptor"
	"os"
	"testing"
)

func master_password_service_test_cleanup() {
	os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}

func TestSetMasterPassword(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(master_password_service_test_cleanup)
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

	t.Cleanup(master_password_service_test_cleanup)
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

	t.Cleanup(master_password_service_test_cleanup)
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

	t.Cleanup(master_password_service_test_cleanup)
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

	t.Cleanup(master_password_service_test_cleanup)
}
