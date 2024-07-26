package services

import (
	"ncrypt/models"
	"os"
	"testing"
)

func TestAddData(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("add_data_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	err := data_service.AddData(*data)

	if err != nil {
		t.Errorf(err.Error())
	}

	os.Remove("add_data_test.txt")
}

func TestGetAllData(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("get_all_test.txt")

	expected_data_list := []models.Data{{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}, {Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}}

	for i := range len(expected_data_list) {
		data_service.AddData(expected_data_list[i])
	}

	actual_data_list, err := data_service.GetAllData()

	if err != nil {
		t.Errorf(err.Error())
	}

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

	os.Remove("get_all_test.txt")
}
