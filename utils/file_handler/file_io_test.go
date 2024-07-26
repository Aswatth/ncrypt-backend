package file_handler

import (
	"encoding/json"
	"ncrypt/models"
	"testing"
)

func TestSave(t *testing.T) {

	file_name := "test.txt"

	contents := []models.Content{{Name: "item_2", Value: "secret"}}
	data_list := []models.Data{{Name: "test_data", Contents: contents}}

	data_as_bytes, _ := json.Marshal(data_list)

	err := Save(file_name, data_as_bytes)

	if err != nil {
		t.Errorf("Saving failed")
		return
	}
}

func TestRead(t *testing.T) {

	TestSave(t)

	file_name := "test.txt"

	contents := []models.Content{{Name: "item_2", Value: "secret"}}
	expected_data_list := []models.Data{{Name: "test_data", Contents: contents}}

	fetched_data, err := Read(file_name)

	if err != nil {
		t.Errorf("Reading failed")
		return
	}

	var actual_data_list []models.Data

	json.Unmarshal(fetched_data, &actual_data_list)

	if len(expected_data_list) != len(actual_data_list) {
		t.Errorf("Record count mismatch")
		return
	}

	for i := range len(expected_data_list) {
		if expected_data_list[i].Name != actual_data_list[i].Name {
			t.Errorf("Mismatch in data")
			return
		}
		if len(expected_data_list[i].Contents) != len(actual_data_list[i].Contents) {
			t.Errorf("Content count mismatch")
			return
		}

		for j := range expected_data_list[i].Contents {
			if (expected_data_list[i].Contents[j].Name != actual_data_list[i].Contents[j].Name) || (expected_data_list[i].Contents[j].Value != actual_data_list[i].Contents[j].Value) {
				t.Errorf("Mismatch in content")
				return
			}
		}
	}
}
