package controllers

import (
	"bytes"
	"encoding/json"
	"ncrypt/models"
	"ncrypt/services"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetSystem_Fail(t *testing.T) {
	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/system/login_info", system_controller.GetSystemData)
	req, _ := http.NewRequest("GET", "/system/login_info", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		if strings.ToUpper(data) != "KEY NOT FOUND" {
			t.Error(data)
		}

	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestGetSystem_Pass(t *testing.T) {

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()

	master_password_service.Init()

	system_service := new(services.SystemService)
	system_service.Init()
	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/system/login_info", system_controller.GetSystemData)
	req, _ := http.NewRequest("GET", "/system/login_info", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		if strings.ToUpper(data) != "KEY NOT FOUND" {
			t.Error(data)
		}

	} else {
		var system_data models.SystemData

		err := json.Unmarshal(test.Body.Bytes(), &system_data)

		if err != nil {
			t.Error(err.Error())
		}

		if system_data.LoginCount != 1 {
			t.Errorf("Incorrect login count\nExpected: %d\nActual: %d", 1, system_data.LoginCount)
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetup_WithNoBackup(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	request_data := make(map[string]interface{})

	request_data["master_password"] = "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""
	request_data["auto_backup_setting"] = auto_backup_setting

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/setup", system_controller.Setup)
	req, _ := http.NewRequest("POST", "/system/setup", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetup_WithBackup(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	request_data := make(map[string]interface{})

	request_data["master_password"] = "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = "test_backup"
	request_data["auto_backup_setting"] = auto_backup_setting

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/setup", system_controller.Setup)
	req, _ := http.NewRequest("POST", "/system/setup", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSignin(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["master_password"] = "12345"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/sigin", system_controller.SignIn)
	req, _ := http.NewRequest("POST", "/system/sigin", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSignIn_WithIncorrectPassword(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["master_password"] = "123"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/sigin", system_controller.SignIn)
	req, _ := http.NewRequest("POST", "/system/sigin", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestLogout(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	system_service.SignIn("12345")

	server.POST("/system/logout", system_controller.Logout)
	req, _ := http.NewRequest("POST", "/system/logout", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestLogout_WithoutSignin(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/logout", system_controller.Logout)
	req, _ := http.NewRequest("POST", "/system/logout", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestExport_CurrentFolder(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.ncrypt"
	request_data["path"] = ""

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/export", system_controller.Export)
	req, _ := http.NewRequest("POST", "/system/export", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	if _, err := os.Stat("test_export.ncrypt"); os.IsNotExist(err) {
		t.Error("Exported file not found")
	} else {
		err = os.Remove("test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestExport_WithCustomPath(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.ncrypt"
	request_data["path"] = "..\\models"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/export", system_controller.Export)
	req, _ := http.NewRequest("POST", "/system/export", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	if _, err := os.Stat("..\\models\\test_export.ncrypt"); os.IsNotExist(err) {
		t.Error("Exported file not found")
	} else {
		err = os.Remove("..\\models\\test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestExport_WithInvalidPath(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.ncrypt"
	request_data["path"] = "..\\test"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/export", system_controller.Export)
	req, _ := http.NewRequest("POST", "/system/export", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestExport_IncorrectFormat(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.txt"
	request_data["path"] = ""

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/export", system_controller.Export)
	req, _ := http.NewRequest("POST", "/system/export", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_Default(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = false
	request_data["has_uppercase"] = false
	request_data["has_special_char"] = false
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlyDigits(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = true
	request_data["has_uppercase"] = false
	request_data["has_special_char"] = false
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlyUppercase(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = false
	request_data["has_uppercase"] = true
	request_data["has_special_char"] = false
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlySpecialChar(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = false
	request_data["has_uppercase"] = false
	request_data["has_special_char"] = true
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlyDigitsAndUppercase(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = true
	request_data["has_uppercase"] = true
	request_data["has_special_char"] = false
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlyDigitsAndSpecialChar(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = true
	request_data["has_uppercase"] = false
	request_data["has_special_char"] = true
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_OnlyUppercaseAndSpecialChar(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = false
	request_data["has_uppercase"] = true
	request_data["has_special_char"] = true
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_All(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = true
	request_data["has_uppercase"] = true
	request_data["has_special_char"] = true
	request_data["length"] = float64(8)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_Length(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["has_digits"] = true
	request_data["has_uppercase"] = true
	request_data["has_special_char"] = true
	request_data["length"] = float64(16)

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/password_generator_preference", system_controller.UpdatePasswordGeneratorPreference)
	req, _ := http.NewRequest("PUT", "/system/password_generator_preference", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestSetPasswordPreference_GeneratePassword(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	passwordPreference := make(map[string]interface{})
	passwordPreference["has_digits"] = true
	passwordPreference["has_uppercase"] = true
	passwordPreference["has_special_char"] = true
	passwordPreference["length"] = float64(16)

	err = system_service.UpdatePasswordGeneratorPreference(passwordPreference)

	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/system/generate_password", system_controller.GeneratePassword)
	req, _ := http.NewRequest("GET", "/system/generate_password", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	} else {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		if len(data) == 0 {
			t.Errorf("Generated password cannot be empty")
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestImport(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	err = system_service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.ncrypt"
	request_data["path"] = ""
	request_data["master_password"] = "12345"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/import", system_controller.Import)
	req, _ := http.NewRequest("POST", "/system/import", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	err = os.Remove("test_export.ncrypt")
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestImport_WithIncorrectMasterPassword(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	err = system_service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["file_name"] = "test_export.ncrypt"
	request_data["path"] = ""
	request_data["master_password"] = "123"

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/import", system_controller.Import)
	req, _ := http.NewRequest("POST", "/system/import", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	err = os.Remove("test_export.ncrypt")
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestBackup(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = "test_backup"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/system/backup", system_controller.Backup)
	req, _ := http.NewRequest("POST", "/system/backup", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	// Remove backup up file
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err.Error())
		return
	}

	files, err := os.ReadDir(cwd)
	if err != nil {
		t.Error(err.Error())
	}
	var backup_file_name string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "test_backup") {
			backup_file_name = file.Name()
			break
		}
	}

	if _, err := os.Stat(backup_file_name); os.IsNotExist(err) {
		t.Error("Backup file not found")
	} else {
		err = os.Remove(backup_file_name)
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToTrue(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting = make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	auto_backup_setting_bytes, err := json.Marshal(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/automatic_backup_setting", system_controller.UpdateAutomaticBackup)
	req, _ := http.NewRequest("PUT", "/system/automatic_backup_setting", bytes.NewBuffer(auto_backup_setting_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToFalse(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["is_enabled"] = false

	auto_backup_setting_bytes, err := json.Marshal(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/automatic_backup_setting", system_controller.UpdateAutomaticBackup)
	req, _ := http.NewRequest("PUT", "/system/automatic_backup_setting", bytes.NewBuffer(auto_backup_setting_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpadteAutomaticBackupData_ChangingFileName(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["backup_file_name"] = "backup.ncrypt"

	auto_backup_setting_bytes, err := json.Marshal(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/automatic_backup_setting", system_controller.UpdateAutomaticBackup)
	req, _ := http.NewRequest("PUT", "/system/automatic_backup_setting", bytes.NewBuffer(auto_backup_setting_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpadteAutomaticBackupData_EmptyFileName(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["backup_file_name"] = ""

	auto_backup_setting_bytes, err := json.Marshal(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/automatic_backup_setting", system_controller.UpdateAutomaticBackup)
	req, _ := http.NewRequest("PUT", "/system/automatic_backup_setting", bytes.NewBuffer(auto_backup_setting_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpdateSessionDuration(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["session_duration_in_minutes"] = 20

	request_data_bytes, err := json.Marshal(request_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/session_duration", system_controller.UpdateSessionDuration)
	req, _ := http.NewRequest("PUT", "/system/session_duration", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestExtendSessionDuration(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/system/session_duration", system_controller.ExtendSession)
	req, _ := http.NewRequest("GET", "/system/session_duration", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func TestUpdateTheme(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()

	system_controller := new(SystemController)
	system_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := system_service.Setup("12345", auto_backup_setting)
	if err != nil {
		t.Error(err.Error())
	}

	request_data := make(map[string]interface{})
	request_data["theme"] = "DARK"

	request_data_bytes, err := json.Marshal(request_data)
	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/system/theme", system_controller.UpdateTheme)
	req, _ := http.NewRequest("PUT", "/system/theme", bytes.NewBuffer(request_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		json.Unmarshal(test.Body.Bytes(), &data)

		t.Error(data)
	}

	t.Cleanup(system_controller_test_cleanup)
}

func system_controller_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}
