package services

import (
	"ncrypt/models"
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

func TestAddLoginData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, *fetched_data)

	//Clean up
	login_service.DeleteLogin(login_data.Name)
}

func TestAddLoginData_DuplicateData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Duplicate key with updated data
	login_data = &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.AddLoginData(login_data)

	if err == nil {
		t.Error("Should not allow duplicate keys/ overwriting existing keys")
	}

	//Clean up
	login_service.DeleteLogin(login_data.Name)
}

func TestGetLoginData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, *fetched_data)

	//Clean up
	login_service.DeleteLogin(login_data.Name)
}

func TestGetAllLoginData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	fetched_data_list, err := login_service.GetAllLoginData()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != 1 {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", 1, len(fetched_data_list))
	}

	compareLoginData(t, *login_data, fetched_data_list[0])

	//Clean up
	login_service.DeleteLogin(login_data.Name)
}

func TestDeleteLoginData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	login_service.DeleteLogin(login_data.Name)

	_, err := login_service.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		t.Error(err.Error())
	}
}

func TestUpdateLoginData(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Duplicate key with updated data
	login_data = &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.UpdateLoginData(login_data.Name, login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *login_data, *fetched_data)

	//Clean up
	login_service.DeleteLogin(login_data.Name)
}

func TestUpdateLoginData_ChangeName(t *testing.T) {
	login_data := &models.Login{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}

	login_service := new(LoginService)
	login_service.Init()

	login_service.AddLoginData(login_data)

	//Duplicate key with updated data
	updated_login_data := &models.Login{Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.UpdateLoginData(login_data.Name, updated_login_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := login_service.GetLoginData(updated_login_data.Name)

	if err != nil {
		t.Error(err.Error())
	}

	compareLoginData(t, *updated_login_data, *fetched_data)

	//Data with old name should be deleted
	_, err = login_service.GetLoginData(login_data.Name)

	if err != nil && err != badger.ErrKeyNotFound {
		t.Error(err.Error())
	}

	//Clean up
	login_service.DeleteLogin(updated_login_data.Name)
}

func TestUpdateLoginData_Fail(t *testing.T) {
	login_data_list := []models.Login{{Name: "github", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}, {Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "abc", Password: "123"}, {Username: "pqr", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}}

	login_service := new(LoginService)
	login_service.Init()

	for _, login_data := range login_data_list {
		login_service.AddLoginData(&login_data)
	}

	//Duplicate key with updated data
	updated_login_data := &models.Login{Name: "email", URL: "https://github.com", Accounts: []models.Account{{Username: "ABC", Password: "123"}, {Username: "PQR", Password: "456"}}, Attributes: &models.Attributes{IsFavourite: true, RequireMasterPassword: false}}
	err := login_service.UpdateLoginData(login_data_list[0].Name, updated_login_data)

	if err == nil {
		t.Error("Should fail to update due to conflicting keys")
	}

	//Clean up
	for _, login_data := range login_data_list {
		login_service.DeleteLogin(login_data.Name)
	}
}
