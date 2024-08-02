package controllers

import (
	"bytes"
	"encoding/json"
	"ncrypt/services"
	"ncrypt/utils/encryptor"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetPassword(t *testing.T) {
	password_data := make(map[string]string)

	password_data["master_password"] = "12345"

	hashed_password := encryptor.CreateHash(password_data["master_password"])

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init(master_password_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password", master_password_controller.SetPassword)
	req, _ := http.NewRequest("POST", "/master_password", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	} else {
		//Validate set password
		data, err := master_password_service.GetMasterPassword()

		if err != nil {
			t.Error(err.Error())
		}

		if data != hashed_password {
			t.Errorf("Mistmatch in password\nExpected: %s\nActual:%s", hashed_password, data)
		}
	}

	os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}

func TestSetPassword_RESET(t *testing.T) {
	password_data := make(map[string]string)

	password_data["master_password"] = "12345"

	hashed_password := encryptor.CreateHash(password_data["master_password"])

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init(master_password_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password", master_password_controller.SetPassword)
	req, _ := http.NewRequest("POST", "/master_password", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	} else {
		//Validate set password
		data, err := master_password_service.GetMasterPassword()

		if err != nil {
			t.Error(err.Error())
		}

		if data != hashed_password {
			t.Errorf("Mistmatch in password\nExpected: %s\nActual:%s", hashed_password, data)
			return
		}

		password_data["master_password"] = "123"
		req, _ := http.NewRequest("POST", "/master_password", bytes.NewReader(password_data_bytes))

		server.ServeHTTP(test, req)

		if test.Code != 200 {
			var data interface{}

			json.Unmarshal(test.Body.Bytes(), &data)

			t.Error(data)
		} else {
			//Validate set password
			data, err := master_password_service.GetMasterPassword()

			if err != nil {
				t.Error(err.Error())
			}

			if data != hashed_password {
				t.Errorf("Mistmatch in password\nExpected: %s\nActual:%s", hashed_password, data)
				return
			}
		}
	}

	os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}

func TestValidatePassword_PASS(t *testing.T) {

	password := "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword(password)

	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init(master_password_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data := make(map[string]string)
	password_data["master_password"] = password
	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password/validate", master_password_controller.Validate)
	req, _ := http.NewRequest("POST", "/master_password/validate", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}

func TestValidatePassword_FAIL(t *testing.T) {

	password := "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword(password)

	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init(master_password_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data := make(map[string]string)
	password_data["master_password"] = "random"
	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password/validate", master_password_controller.Validate)
	req, _ := http.NewRequest("POST", "/master_password/validate", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 400 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	os.RemoveAll(os.Getenv("MASTER_PASSWORD_DB_NAME"))
}
