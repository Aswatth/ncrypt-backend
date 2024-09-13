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

	t.Cleanup(master_password_service_test_cleanup)
}

func TestUpdateMasterPassword(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateMasterPassword("123")

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.Validate("123")

	if err != nil {
		t.Error(err.Error())
	}

	if !result {
		t.Errorf("Expected: %t\nActual: %t", true, result)
	}

	t.Cleanup(master_password_service_test_cleanup)
}

func TestValidate_PASS(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.Validate("12345")

	if err != nil {
		t.Error(err.Error())
	}

	if !result {
		t.Errorf("Expected: %t\nActual: %t", true, result)
	}

	t.Cleanup(master_password_service_test_cleanup)
}

func TestValidate_FAIL(t *testing.T) {
	new_password := map[string]string{"master_password": "12345"}

	service := new(MasterPasswordService)
	service.Init()

	err := service.SetMasterPassword(new_password["master_password"])

	if err != nil {
		t.Error(err.Error())
	}

	result, err := service.Validate("123")

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

func TestImport(t *testing.T) {
	service := new(MasterPasswordService)
	service.Init()

	imported_password := "12345"
	err := service.importData(imported_password)

	if err != nil {
		t.Error(err.Error())
	}

	stored_password, err := service.GetMasterPassword()

	if err != nil {
		t.Error(err.Error())
	}

	if stored_password != imported_password {
		t.Error("Password mismatch")
	}


	t.Cleanup(master_password_service_test_cleanup)
}

func master_password_service_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}