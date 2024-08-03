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

func system_controller_test_cleanup() {
	os.RemoveAll("SYSTEM")
}
func TestGetSystem_Fail(t *testing.T) {
	system_service := new(services.SystemService)
	system_service.Init()
	system_controller := new(SystemController)
	system_controller.Init(*system_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/system/login_info", system_controller.GetLoginInfo)
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
	system_controller.Init(*system_service)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/system/login_info", system_controller.GetLoginInfo)
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

		if system_data.Login_count != 1 {
			t.Errorf("Incorrect login count\nExpected: %d\nActual: %d", 1, system_data.Login_count)
		}
	}

	t.Cleanup(system_controller_test_cleanup)
}
