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

func login_controller_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}

func compateData(expected models.Login, actual models.Login) bool {
	if expected.Name != actual.Name || expected.URL != actual.URL || expected.Attributes.IsFavourite != actual.Attributes.IsFavourite || expected.Attributes.RequireMasterPassword != actual.Attributes.RequireMasterPassword {
		return false
	}
	if len(expected.Accounts) != len(actual.Accounts) {
		return false
	}

	for index := range len(expected.Accounts) {
		if expected.Accounts[index].Username != actual.Accounts[index].Username || expected.Accounts[index].Password != actual.Accounts[index].Password {
			return false
		}
	}

	return true
}

func TestCreateLogin_Without_Master_Password(t *testing.T) {

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_data_bytes, err := json.Marshal(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.POST("/login", login_controller.AddLoginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 400 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if strings.ToUpper(data) != "MASTER_PASSWORD NOT SET" {
			t.Error(data)
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestCreateLogin_With_Master_Password(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_data_bytes, err := json.Marshal(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.POST("/login", login_controller.AddLoginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestCreateLogin_With_Duplicate_Account_Username(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "abc", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_data_bytes, err := json.Marshal(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.POST("/login", login_controller.AddLoginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 400 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if strings.ToUpper(data) != "DUPLICATE USERNAME ABC" {
			t.Error(data)
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestGetLoginData_All(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)
	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/login", login_controller.GetLoginData)
	req, _ := http.NewRequest("GET", "/login", bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		t.Error(data)
	} else {
		var data_list []models.Login

		err := json.Unmarshal(test.Body.Bytes(), &data_list)

		if err != nil {
			t.Error(err.Error())
		}

		if len(data_list) != 1 {
			t.Errorf("invalid data count\nExpected: %d\nActual: %d", 1, len(data_list))
		} else {
			if !compateData(*login_data, data_list[0]) {
				t.Errorf("Mismatch in data\nExpected:%v\nActual:%v", login_data, data_list[0])
			}
		}

	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestGetLoginData_PASS(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)
	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/login", login_controller.GetLoginData)
	req, _ := http.NewRequest("GET", "/login?name="+login_data.Name, bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		t.Error(data)
	} else {
		var data models.Login

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if !compateData(*login_data, data) {
			t.Errorf("Mismatch in data\nExpected:%v\nActual:%v", login_data, data)
		}

	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestGetLoginData_FAIL(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)
	if err != nil {
		t.Error(err.Error())
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/login", login_controller.GetLoginData)
	req, _ := http.NewRequest("GET", "/login?name=random_name", bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if strings.ToUpper(data) != "KEY NOT FOUND" {
			t.Error(data)
		}
	} else {
		var data *models.Login

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if data != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestUpdateLoginData_Without_Duplicate_Account_Username(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	updated_login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "efg", Password: "123"}, {Username: "EFG", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	updated_login_data_bytes, err := json.Marshal(updated_login_data)

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/login/:name", login_controller.UpdateLoginData)
	req, _ := http.NewRequest("PUT", "/login/"+login_data.Name, bytes.NewReader(updated_login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	} else {
		fetched_data, err := login_service.GetLoginData(updated_login_data.Name)

		if err != nil {
			t.Error(err.Error())
		}

		for _, account := range fetched_data.Accounts {
			if account.Username != "efg" && account.Username != "EFG" {
				t.Error("incorrect username")
			}
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestUpdateLoginData_With_Duplicate_Account_Username(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	updated_login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "ttt", Password: "123"}, {Username: "ttt", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	updated_login_data_bytes, err := json.Marshal(updated_login_data)

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/login/:name", login_controller.UpdateLoginData)
	req, _ := http.NewRequest("PUT", "/login/"+login_data.Name, bytes.NewReader(updated_login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if strings.ToUpper(data) != "DUPLICATE USERNAME TTT" {
			t.Error(data)
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestUpdateLoginData_Without_Conflicting_Names(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	updated_login_data := &models.Login{Name: "email", URL: "https://email.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	updated_login_data_bytes, err := json.Marshal(updated_login_data)

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/login/:name", login_controller.UpdateLoginData)
	req, _ := http.NewRequest("PUT", "/login/"+login_data.Name, bytes.NewReader(updated_login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)

	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestUpdateLoginData_With_Conflicting_Names(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data_list := []models.Login{{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}, {Name: "email", URL: "https://email.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}}

	for _, login_data := range login_data_list {
		err := login_service.AddLoginData(&login_data)

		if err != nil {
			t.Error(err.Error())
		}
	}

	updated_login_data := &models.Login{Name: "email", URL: "https://email.com", Accounts: []models.Account{{Username: "ttt", Password: "123"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	updated_login_data_bytes, err := json.Marshal(updated_login_data)

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/login/:name", login_controller.UpdateLoginData)
	req, _ := http.NewRequest("PUT", "/login/"+login_data_list[0].Name, bytes.NewReader(updated_login_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		if strings.ToUpper(data) != strings.ToUpper(updated_login_data.Name+" already exists") {
			t.Error(data)
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestDeleteLoginData(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.DELETE("/login/:name", login_controller.DeleteLoginData)
	req, _ := http.NewRequest("DELETE", "/login/"+login_data.Name, bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		t.Error(data)
	}

	t.Cleanup(login_controller_test_cleanup)
}

func TestGetAccountPassword(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	login_service := new(services.LoginDataService)
	login_service.Init()

	login_controller := new(LoginDataController)
	login_controller.Init()

	password := "12345"
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: password}}, Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	if err != nil {
		t.Error(err)
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/login/:name", login_controller.GetAccountPassword)
	req, _ := http.NewRequest("GET", "/login/"+login_data.Name+"?username="+login_data.Accounts[0].Username, bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}

		t.Error(data)
	} else {
		var decrypted_password string

		err := json.Unmarshal(test.Body.Bytes(), &decrypted_password)

		if err != nil {
			t.Error(err.Error())
		}

		if decrypted_password != password {
			t.Errorf("Mismatch in password\nExpected: %s\nActual: %s", password, decrypted_password)
		}
	}

	t.Cleanup(login_controller_test_cleanup)
}
