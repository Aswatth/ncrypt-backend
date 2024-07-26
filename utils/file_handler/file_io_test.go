package file_handler

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	file_name := "save_test.txt"

	data_to_save := "Hi there. This is a test"

	err := Save(file_name, []byte(data_to_save))

	if err != nil {
		t.Errorf("Saving failed")
		return
	}

	os.Remove("save_test.txt")
}

func TestRead(t *testing.T) {
	file_name := "read_test.txt"

	data_to_save := "Hi there. This is a test"

	Save(file_name, []byte(data_to_save))

	fetched_data, err := Read(file_name)

	if err != nil {
		t.Errorf("Reading failed")
	}

	if string(fetched_data[:]) != data_to_save {
		t.Errorf("Incorrect data")
	}

	os.Remove("read_test.txt")
}
