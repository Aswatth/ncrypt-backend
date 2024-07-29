package services

import (
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
}
