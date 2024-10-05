package services

import (
	"ncrypt/models"
	"os"
	"strings"
	"testing"

	"github.com/dgraph-io/badger/v4"
)

func compareLoginData(t *testing.T, expected_login_data models.Login, actual_login_data models.Login) {
	if expected_login_data.Name == actual_login_data.Name && expected_login_data.URL == actual_login_data.URL && expected_login_data.Attributes.IsFavourite == actual_login_data.Attributes.IsFavourite && expected_login_data.Attributes.RequireMasterPassword == actual_login_data.Attributes.RequireMasterPassword {
		if len(expected_login_data.Accounts) == len(actual_login_data.Accounts) {
			for i := range len(expected_login_data.Accounts) {
				if expected_login_data.Accounts[i].Username != actual_login_data.Accounts[i].Username || expected_login_data.Accounts[i].Password != actual_login_data.Accounts[i].Password {
					t.Error("Mismatch in account data")
					t.Errorf("Expected: %v\n Actual: %v", expected_login_data.Accounts[i], actual_login_data.Accounts[i])
				}
			}
		} else {
			t.Error("Mismatch in accounts")
			t.Errorf("Expected: %d\n Actual: %d", len(expected_login_data.Accounts), len(actual_login_data.Accounts))
		}
	} else {
		t.Error("Mismatch in data")
		t.Errorf("Expected: %v\n Actual: %v", expected_login_data, actual_login_data)
	}
}

func login_service_test_init() {
	master_password_service := new(MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword("12345")
}

func TestAddLoginData_With_Master_Password(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{
		Name: "github",
		URL:  "https://github.com",
		Accounts: []models.Account{
			{Username: "abc", Password: "123"},
			{Username: "pqr", Password: "456"}},
		Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false},
	}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, fetched_data)

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestAddLoginData_Without_Master_Password(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		if strings.ToUpper(err.Error()) != "MASTER_PASSWORD NOT SET" {
			t.Error(err.Error())
		}
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestAddLoginData_Duplicate_Account_Username(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{
		Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "abc", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		if strings.ToUpper(err.Error()) != "DUPLICATE USERNAME ABC" {
			t.Error(err.Error())
		}
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestAddLoginData_DuplicateData(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	duplicate_login_data := &models.Login{Name: "GITHUB", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Duplicate key with updated data
	err := login_service.AddLoginData(duplicate_login_data)

	if err == nil {
		t.Error("Should not allow duplicate keys/ overwriting existing keys")
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestGetLoginData(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)
	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, fetched_data)

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestGetAllLoginData_With_Single_Record(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != 1 {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", 1, len(fetched_data_list))
	}

	compareLoginData(t, *login_data, fetched_data_list[0])

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestGetAllLoginData_With_Multiple_Record(t *testing.T) {
	login_service_test_init()

	login_datas := []models.Login{
		{Name: "github1", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
		{Name: "github2", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
	}

	login_service := new(LoginDataService)
	login_service.Init()

	for _, login_data := range login_datas {
		err := login_service.AddLoginData(&login_data)

		if err != nil {
			t.Error(err.Error())
		}
	}

	fetched_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != len(login_datas) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(login_datas), len(fetched_data_list))
	}

	for index := range login_datas {
		compareLoginData(t, login_datas[index], fetched_data_list[index])
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestDeleteLoginData(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	login_service.DeleteLoginData(login_data.Name)

	_, err := login_service.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		t.Error(err.Error())
	}

	t.Cleanup(login_service_test_cleanup)
}

func TestUpdateLoginData_ChangingAccounts(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Updating accounts only
	login_data = &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.UpdateLoginData(login_data.Name, login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, fetched_data)

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestUpdateLoginData_ChangeAll(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Updating entire data
	updated_login_data := &models.Login{Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: true}}
	err := login_service.UpdateLoginData(login_data.Name, updated_login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(updated_login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *updated_login_data, fetched_data)

	//Data with old name should be deleted
	_, err = login_service.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		t.Error(err.Error())
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestUpdateLoginData_ConflictingName(t *testing.T) {
	login_service_test_init()

	login_datas := []models.Login{{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
		{Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
	}

	login_service := new(LoginDataService)
	login_service.Init()

	for _, login_data := range login_datas {
		login_service.AddLoginData(&login_data)
	}

	//Duplicate key with updated data
	updated_login_data := &models.Login{Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.UpdateLoginData(login_datas[0].Name, updated_login_data)

	if err == nil {
		t.Error("Should fail to update due to conflicting keys")
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestGetDecryptedAccountPassword_ValidUsername(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	expected_password := login_data.Accounts[0].Password

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_password, err := login_service.GetDecryptedAccountPassword(login_data.Name, login_data.Accounts[0].Username)

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_password != expected_password {
		t.Errorf("Expected: %s\nActual: %s", expected_password, fetched_password)
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestGetDecryptedAccountPassword_InvalidUsername(t *testing.T) {
	login_service_test_init()

	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.AddLoginData(login_data)

	if err != nil {
		t.Error(err.Error())
	}

	_, err = login_service.GetDecryptedAccountPassword(login_data.Name, "ttt") //ttt - invalid username

	if err != nil {
		if strings.ToUpper(err.Error()) != "ACCOUNT USERNAME NOT FOUND" {
			t.Error(err.Error())
		}
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestLoginDataImport(t *testing.T) {
	login_service_test_init()

	login_datas := []models.Login{
		{Name: "github1", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
		{Name: "github2", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
	}

	login_service := new(LoginDataService)
	login_service.Init()

	err := login_service.importData(login_datas)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != len(login_datas) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(login_datas), len(fetched_data_list))
	}

	for index := range login_datas {
		compareLoginData(t, login_datas[index], fetched_data_list[index])
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func TestLoginDataRecrpyt(t *testing.T) {
	master_password_service := new(MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword("12345")

	old_password, err := master_password_service.GetMasterPassword()

	if err != nil {
		t.Error(err.Error())
	}

	login_datas := []models.Login{
		{Name: "github1", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
		{Name: "github2", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}},
	}

	login_service := new(LoginDataService)
	login_service.Init()

	for _, login_data := range login_datas {
		login_service.AddLoginData(&login_data)
	}

	data := make(map[string]string)
	data["OLD_PASSWORD"] = old_password
	data["NEW_PASSWORD"] = "123"

	err = login_service.recryptData(data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != len(login_datas) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(login_datas), len(fetched_data_list))
	}

	//Clean up
	t.Cleanup(login_service_test_cleanup)
}

func login_service_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
	os.RemoveAll("logs")
}
