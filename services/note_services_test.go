package services

import (
	"ncrypt/models"
	"os"
	"strings"
	"testing"

	"github.com/dgraph-io/badger/v4"
)

func note_service_test_init() {
	master_password_service := new(MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword("12345")
}

func compareNote(t *testing.T, expected_note models.Note, actual_note models.Note) {
	if !(expected_note.CreatedDateTime == actual_note.CreatedDateTime && expected_note.Title == actual_note.Title && expected_note.Attributes.IsFavourite == actual_note.Attributes.IsFavourite && expected_note.Attributes.RequireMasterPassword == actual_note.Attributes.RequireMasterPassword) {
		t.Errorf("Mismatch in data\nExpected: %v\nActual: %v", expected_note, actual_note)
	}
}

func TestAddNote_With_Master_Password(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_note, err := note_service.GetNote(note.CreatedDateTime)

	if err != nil {
		t.Error(err.Error())
	}

	compareNote(t, note, *fetched_note)

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestAddNote_Without_Master_Password(t *testing.T) {

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)

	if err != nil {
		if strings.ToUpper(err.Error()) != "MASTER_PASSWORD NOT SET" {
			t.Error(err.Error())
		}
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestAddNote_DuplicateData(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	duplicate_note_data := make(map[string]interface{})
	duplicate_note_data["created_date_time"] = "testing"
	duplicate_note_data["title"] = "test"
	duplicate_note_data["content"] = "this is a test"
	duplicate_note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	note_service := new(NoteService)
	note_service.Init()

	note_service.AddNote(note_data)

	//Duplicate key with updated data
	err := note_service.AddNote(duplicate_note_data)

	if err == nil {
		t.Error("Should not allow duplicate keys/ overwriting existing keys")
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestGetNote(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	fetched_note, err := note_service.GetNote(note.CreatedDateTime)

	if err != nil {
		t.Error(err.Error())
	}

	compareNote(t, note, *fetched_note)

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestGetAllNote_With_Single_Record(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_note_list, err := note_service.GetAllNotes()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_note_list) != 1 {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", 1, len(fetched_note_list))
	}

	compareNote(t, note, fetched_note_list[0])

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestGetAllNote_With_Multiple_Record(t *testing.T) {
	note_service_test_init()

	var note_data_list []map[string]interface{}

	note_data_1 := make(map[string]interface{})
	note_data_1["created_date_time"] = "testing1"
	note_data_1["title"] = "test1"
	note_data_1["content"] = "this is a test"
	note_data_1["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}
	note_data_list = append(note_data_list, note_data_1)

	note_data_2 := make(map[string]interface{})
	note_data_2["created_date_time"] = "testing2"
	note_data_2["title"] = "test2"
	note_data_2["content"] = "this is a test"
	note_data_2["attributes"] = map[string]interface{}{"is_favourite": false, "require_master_password": true}
	note_data_list = append(note_data_list, note_data_2)

	note_service := new(NoteService)
	note_service.Init()

	for _, note_data := range note_data_list {
		err := note_service.AddNote(note_data)

		if err != nil {
			t.Error(err.Error())
		}
	}

	fetched_note_list, err := note_service.GetAllNotes()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_note_list) != len(note_data_list) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(note_data_list), len(fetched_note_list))
	}

	for index := range note_data_list {
		var note models.Note
		note.FromMap(note_data_list[index])
		compareNote(t, note, fetched_note_list[index])
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestDeleteNote(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	note_service.AddNote(note_data)

	note_service.DeleteNote(note.CreatedDateTime)

	_, err := note_service.GetNote(note.CreatedDateTime)

	if err != nil && err != badger.ErrKeyNotFound {
		t.Error(err.Error())
	}

	t.Cleanup(note_service_test_cleanup)
}

func TestUpdateNote(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	note_service.AddNote(note_data)

	//Updating accounts only
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test_update"
	note_data["content"] = "this is a test update"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}
	note.FromMap(note_data)

	err := note_service.UpdateNote(note.CreatedDateTime, note_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_note, err := note_service.GetNote(note.CreatedDateTime)

	if err != nil {
		t.Error(err.Error())
	}

	compareNote(t, note, *fetched_note)

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestUpdateNote_ChangeAll(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	note_service.AddNote(note_data)

	//Updating entire data
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test_update"
	note_data["content"] = "this is a test updated"
	note_data["attributes"] = map[string]interface{}{"is_favourite": false, "require_master_password": false}

	note.FromMap(note_data)

	err := note_service.UpdateNote(note.CreatedDateTime, note_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data, err := note_service.GetNote(note.CreatedDateTime)

	if err != nil {
		t.Error(err.Error())
	}

	compareNote(t, note, *fetched_data)

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestGetDecryptedContent_ValidCreatedDateTime(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	expected_content := note.Content

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_content, err := note_service.GetDecryptedContent(note.CreatedDateTime)

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_content != expected_content {
		t.Errorf("Expected: %s\nActual: %s", expected_content, fetched_content)
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestGetDecryptedContent_InvalidCreatedDateTime(t *testing.T) {
	note_service_test_init()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "testing"
	note_data["title"] = "test"
	note_data["content"] = "this is a test"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.AddNote(note_data)

	if err != nil {
		t.Error(err.Error())
	}

	_, err = note_service.GetDecryptedContent("random_value") // invalid created date time

	if err == nil {
		t.Error("Should result in an error")
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestNoteImport(t *testing.T) {
	note_service_test_init()

	note_datas := []models.Note{
		{CreatedDateTime: "testing1", Title: "test1", Content: "this is a test1", Attributes: models.Attributes{IsFavourite: true, RequireMasterPassword: true}},
		{CreatedDateTime: "testing2", Title: "test2", Content: "this is a test2", Attributes: models.Attributes{IsFavourite: false, RequireMasterPassword: false}},
	}

	note_service := new(NoteService)
	note_service.Init()

	err := note_service.importData(note_datas)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_note_list, err := note_service.GetAllNotes()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_note_list) != len(note_datas) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(note_datas), len(fetched_note_list))
	}

	for index := range note_datas {
		compareNote(t, note_datas[index], fetched_note_list[index])
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func TestNoteRecrpyt(t *testing.T) {
	master_password_service := new(MasterPasswordService)
	master_password_service.Init()

	master_password_service.SetMasterPassword("12345")

	old_password, err := master_password_service.GetMasterPassword()

	if err != nil {
		t.Error(err.Error())
	}

	var note_data_list []map[string]interface{}

	note_data_1 := make(map[string]interface{})
	note_data_1["created_date_time"] = "testing1"
	note_data_1["title"] = "test1"
	note_data_1["content"] = "this is a test"
	note_data_1["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}
	note_data_list = append(note_data_list, note_data_1)

	note_data_2 := make(map[string]interface{})
	note_data_2["created_date_time"] = "testing2"
	note_data_2["title"] = "test2"
	note_data_2["content"] = "this is a test"
	note_data_2["attributes"] = map[string]interface{}{"is_favourite": false, "require_master_password": true}
	note_data_list = append(note_data_list, note_data_2)

	note_service := new(NoteService)
	note_service.Init()

	for _, note_data := range note_data_list {
		note_service.AddNote(note_data)
	}

	data := make(map[string]string)
	data["OLD_PASSWORD"] = old_password
	data["NEW_PASSWORD"] = "123"

	err = note_service.recryptData(data)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_data_list, err := note_service.GetAllNotes()

	if err != nil {
		t.Error(err.Error())
	}

	if len(fetched_data_list) != len(note_data_list) {
		t.Errorf("Mismatch in count\nExpected:\t%d\nActual:\t%d", len(note_data_list), len(fetched_data_list))
	}

	//Clean up
	t.Cleanup(note_service_test_cleanup)
}

func note_service_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
	os.RemoveAll("logs")
}
