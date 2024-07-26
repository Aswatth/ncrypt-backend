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

func TestGetData_PASS(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("get_data_pass_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(*data)

	fetched_data, err := data_service.GetData(data.Name)

	if err != nil {
		t.Errorf(err.Error())
	}

	if fetched_data.Name == data.Name {
		for j := range data.Contents {
			if (data.Contents[j].Name != fetched_data.Contents[j].Name) || (data.Contents[j].Value != fetched_data.Contents[j].Value) {
				t.Errorf("Mismatch in content")
			}
		}
	}

	os.Remove("get_data_pass_test.txt")
}

func TestGetData_FAIL(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("get_data_fail_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(*data)

	_, err := data_service.GetData("random_name")

	if err != nil {
		if err.Error() != "DATA NOT FOUND" {
			t.Errorf(err.Error())
		}
	}

	os.Remove("get_data_fail_test.txt")
}

func TestDeleteData_PASS(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("delete_data_pass_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(*data)

	err := data_service.DeleteData(data.Name)

	if err != nil {
		t.Errorf(err.Error())
	}

	os.Remove("delete_data_pass_test.txt")
}

func TestDeleteData_FAIL(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("delete_data_fail_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(*data)

	err := data_service.DeleteData("random_name")

	if err != nil {
		if err.Error() != "DATA NOT FOUND" {
			t.Errorf(err.Error())
		}
	}

	os.Remove("delete_data_fail_test.txt")
}

func TestUpdateData_PASS(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("update_data_pass_test.txt")

	data := &models.Data{Name: "test", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}

	data_service.AddData(*data)

	updated_data := &models.Data{Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}}}

	err := data_service.UpdateData(data.Name, *updated_data)

	if err != nil {
		t.Errorf(err.Error())
	}

	fetched_data, err := data_service.GetData(updated_data.Name)

	if err != nil {
		t.Errorf(err.Error())
	}

	if fetched_data.Name == updated_data.Name {
		for j := range updated_data.Contents {
			if (updated_data.Contents[j].Name != fetched_data.Contents[j].Name) || (updated_data.Contents[j].Value != fetched_data.Contents[j].Value) {
				t.Errorf("Mismatch in content")
			}
		}
	}

	os.Remove("update_data_pass_test.txt")
}

func TestUpdateData_FAIL(t *testing.T) {
	data_service := new(DataService)

	data_service.Init("update_data_fail_test.txt")

	data_list := []models.Data{{Name: "test1", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}, {Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}, {Name: "password", Value: "12345"}}}}

	for _, data := range data_list {
		data_service.AddData(data)
	}

	updated_data := &models.Data{Name: "test2", Contents: []models.Content{{Name: "id", Value: "test_id"}}}

	err := data_service.UpdateData(data_list[0].Name, *updated_data)

	if err != nil {
		if err.Error() != "DATA NAME ALREADY EXISTS" {
			t.Errorf(err.Error())
		}
	}

	os.Remove("update_data_fail_test.txt")
}
