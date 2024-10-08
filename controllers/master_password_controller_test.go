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

func TestUpdatePassword(t *testing.T) {
	password_data := make(map[string]string)

	password_data["master_password"] = "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init()

	master_password_service.SetMasterPassword(password_data["master_password"])

	//Updating password
	password_data = make(map[string]string)
	password_data["old_master_password"] = "12345"
	password_data["new_master_password"] = "123"

	hashed_password := encryptor.CreateHash(password_data["new_master_password"])

	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/master_password", master_password_controller.UpdatePassword)
	req, _ := http.NewRequest("PUT", "/master_password", bytes.NewReader(password_data_bytes))

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

	t.Cleanup(master_password_controller_test_cleanup)
}

func TestUpdatePassword_IncorrectOldPassword(t *testing.T) {
	password_data := make(map[string]string)

	password_data["master_password"] = "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init()

	master_password_service.SetMasterPassword(password_data["master_password"])

	//Updating password
	password_data = make(map[string]string)
	password_data["old_master_password"] = "456"
	password_data["new_master_password"] = "123"

	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/master_password", master_password_controller.UpdatePassword)
	req, _ := http.NewRequest("PUT", "/master_password", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error("Should result in an error")
	}

	t.Cleanup(master_password_controller_test_cleanup)
}
func TestUpdatePassword_NewPasswordSameAsOld(t *testing.T) {
	password_data := make(map[string]string)

	password_data["master_password"] = "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init()

	master_password_service.SetMasterPassword(password_data["master_password"])

	//Updating password
	password_data = make(map[string]string)
	password_data["old_master_password"] = "12345"
	password_data["new_master_password"] = "12345"

	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/master_password", master_password_controller.UpdatePassword)
	req, _ := http.NewRequest("PUT", "/master_password", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error("Should result in an error")
	}

	t.Cleanup(master_password_controller_test_cleanup)
}

func TestValidatePassword_PASS(t *testing.T) {

	password := "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword(password)

	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data := make(map[string]string)
	password_data["master_password"] = password
	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password/validate", master_password_controller.ValidatePassword)
	req, _ := http.NewRequest("POST", "/master_password/validate", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(master_password_controller_test_cleanup)
}

func TestValidatePassword_FAIL(t *testing.T) {

	password := "12345"

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword(password)

	master_password_controller := new(MasterPasswordController)
	master_password_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	password_data := make(map[string]string)
	password_data["master_password"] = "random"
	password_data_bytes, err := json.Marshal(password_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/master_password/validate", master_password_controller.ValidatePassword)
	req, _ := http.NewRequest("POST", "/master_password/validate", bytes.NewReader(password_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 400 {
		var data interface{}

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(master_password_controller_test_cleanup)
}

func master_password_controller_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
	os.RemoveAll("logs")
}
