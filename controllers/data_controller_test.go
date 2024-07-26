package controllers

import (
	"bytes"
	"encoding/json"
	"ncrypt/models"
	"ncrypt/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAddData(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("add_data_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	data_to_add := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_as_bytes, _ := json.Marshal(data_to_add)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.POST("/data", data_controller.AddData)
	req, _ := http.NewRequest("POST", "/data", bytes.NewReader(data_as_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		t.Errorf("Failed to add data")
	}

	os.Remove("add_data_test.txt")
}

func TestGetAllData(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("get_all_data_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	expected_data_list := []models.Data{{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}, {Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}}

	for i := range len(expected_data_list) {
		data_service.AddData(expected_data_list[i])
	}

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/data", data_controller.GetData)
	req, _ := http.NewRequest("GET", "/data", bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var actual_data_list []models.Data

		json.Unmarshal(test.Body.Bytes(), &actual_data_list)

		if len(actual_data_list) != len(expected_data_list) {
			t.Errorf("Record count does not match")
		}

		for i := range len(expected_data_list) {
			if expected_data_list[i].Name != actual_data_list[i].Name {
				t.Errorf("Mismatch in data")
			}
			if len(expected_data_list[i].Contents) != len(actual_data_list[i].Contents) {
				t.Errorf("Content count mismatch")
			}

			for j := range expected_data_list[i].Contents {
				if (expected_data_list[i].Contents[j].Name != actual_data_list[i].Contents[j].Name) || (expected_data_list[i].Contents[j].Value != actual_data_list[i].Contents[j].Value) {
					t.Errorf("Mismatch in content")
				}
			}
		}
	} else {
		t.Errorf("Failed to fetch all data")
	}

	os.Remove("get_all_data_test.txt")
}
