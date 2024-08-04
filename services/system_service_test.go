package services

import (
	"ncrypt/models"
	"os"
	"testing"
	"time"
)

func system_service_test_cleanup() {
	os.RemoveAll("SYSTEM")
}

func TestSetSystemData(t *testing.T) {
	service := new(SystemService)
	service.Init()

	now := time.Now().Format(time.RFC3339)
	err := service.setSystemData(models.SystemData{Login_count: 1, Last_login: now})

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestGetSystemData(t *testing.T) {
	service := new(SystemService)
	service.Init()

	now := time.Now().Format(time.RFC3339)
	initial_data := models.SystemData{Login_count: 1, Last_login: now}

	err := service.setSystemData(initial_data)

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	if system_data.Login_count != initial_data.Login_count || system_data.Last_login != initial_data.Last_login {
		t.Errorf("Mismatch in data:Expected: %v\nActual: %v", initial_data, system_data)
	}

	t.Cleanup(system_service_test_cleanup)
}
