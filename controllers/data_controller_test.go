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

	server.GET("/data", data_controller.GetAllData)
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
		var result string
		json.Unmarshal(test.Body.Bytes(), &result)
		t.Errorf(result)
	}

	os.Remove("get_all_data_test.txt")
}

func TestGetData_PASS(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("get_data_pass_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	expected_data := models.Data{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(expected_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/data", data_controller.GetData)
	req, _ := http.NewRequest("GET", "/data?name="+expected_data.Name, bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var actual_data models.Data

		json.Unmarshal(test.Body.Bytes(), &actual_data)

		if actual_data.Name == expected_data.Name {
			for j := range expected_data.Contents {
				if (expected_data.Contents[j].Name != actual_data.Contents[j].Name) || (expected_data.Contents[j].Value != actual_data.Contents[j].Value) {
					t.Errorf("Mismatch in content")
				}
			}
		}
	} else {
		var result string
		json.Unmarshal(test.Body.Bytes(), &result)
		t.Errorf(result)
	}

	os.Remove("get_data_pass_test.txt")
}

func TestGetData_FAIL(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("get_data_fail_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	expected_data := models.Data{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(expected_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.GET("/data", data_controller.GetData)
	req, _ := http.NewRequest("GET", "/data?name=random_name", bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var actual_data string

		json.Unmarshal(test.Body.Bytes(), &actual_data)

		if actual_data != "DATA NOT FOUND" {
			t.Errorf("INCORRECT ERROR")
		}
	} else {
		var result string
		json.Unmarshal(test.Body.Bytes(), &result)
		t.Errorf(result)
	}

	os.Remove("get_data_fail_test.txt")
}

func TestDeleteData_PASS(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("delete_data_pass_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	expected_data := models.Data{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(expected_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.DELETE("/data/:name", data_controller.DeleteData)
	req, _ := http.NewRequest("DELETE", "/data/"+expected_data.Name, bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var error_message string
		json.Unmarshal(test.Body.Bytes(), &error_message)
		t.Errorf(error_message)
	}

	os.Remove("delete_data_pass_test.txt")
}

func TestDeleteData_FAIL(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("delete_data_fail_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	expected_data := models.Data{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(expected_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.DELETE("/data/:name", data_controller.DeleteData)
	req, _ := http.NewRequest("DELETE", "/data/random_name", bytes.NewReader([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var error_message string

		json.Unmarshal(test.Body.Bytes(), &error_message)

		if error_message != "DATA NOT FOUND" {
			t.Errorf(error_message)
		}
	}

	os.Remove("delete_data_fail_test.txt")
}

func TestUpdateData_PASS(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("update_data_pass_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	data := models.Data{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(data)

	updated_data := models.Data{Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}}}
	updated_data_bytes, _ := json.Marshal(updated_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/data/:name", data_controller.UpdateData)
	req, _ := http.NewRequest("PUT", "/data/"+data.Name, bytes.NewReader(updated_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var error_message string

		json.Unmarshal(test.Body.Bytes(), &error_message)

		t.Errorf(error_message)
	}

	os.Remove("update_data_pass_test.txt")
}

func TestUpdateData_FAIL(t *testing.T) {
	data_service := new(services.DataService)
	data_service.Init("update_data_fail_test.txt")

	data_controller := new(DataController)
	data_controller.Init(data_service)

	data_list := []models.Data{{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}, {Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}}

	for _, data := range data_list {
		data_service.AddData(data)
	}

	updated_data := models.Data{Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}}}
	updated_data_bytes, _ := json.Marshal(updated_data)

	server := gin.Default()
	test := httptest.NewRecorder()

	server.PUT("/data/:name", data_controller.UpdateData)
	req, _ := http.NewRequest("PUT", "/data/"+data_list[0].Name, bytes.NewReader(updated_data_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var error_message string
		json.Unmarshal(test.Body.Bytes(), &error_message)

		if error_message != "DATA NAME ALREADY EXISTS" {
			t.Errorf(error_message)
		}
	} else {
		t.Errorf("SHOULD NOT UPDATE")
		t.Fail()
	}

	os.Remove("update_data_fail_test.txt")
}