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

func compareNote(t *testing.T, expected_note models.Note, actual_note models.Note) {
	if !(expected_note.CreatedDateTime == actual_note.CreatedDateTime && expected_note.Title == actual_note.Title && expected_note.Attributes.IsFavourite == actual_note.Attributes.IsFavourite && expected_note.Attributes.RequireMasterPassword == actual_note.Attributes.RequireMasterPassword) {
		t.Errorf("Mismatch in data\nExpected: %v\nActual: %v", expected_note, actual_note)
	}
}

func TestAddNote_With_Master_Password(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note := &models.Note{
		CreatedDateTime: "123",
		Title:           "abc",
		Content:         "my content",
		Attributes:      models.Attributes{IsFavourite: true, RequireMasterPassword: false},
	}

	note_bytes, err := json.Marshal(note)
	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/note", note_controller.AddNote)
	req, _ := http.NewRequest("POST", "/note", bytes.NewBuffer(note_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestAddNote_Without_Master_Password(t *testing.T) {

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note := &models.Note{
		CreatedDateTime: "123",
		Title:           "abc",
		Content:         "my content",
		Attributes:      models.Attributes{IsFavourite: true, RequireMasterPassword: false},
	}

	note_bytes, err := json.Marshal(note)
	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/note", note_controller.AddNote)
	req, _ := http.NewRequest("POST", "/note", bytes.NewBuffer(note_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestAddNote_DuplicateData(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	note_service.AddNote(note_data)

	note_bytes, err := json.Marshal(note)
	if err != nil {
		t.Error(err.Error())
	}

	server.POST("/note", note_controller.AddNote)
	req, _ := http.NewRequest("POST", "/note", bytes.NewBuffer(note_bytes))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestGetNote(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/note", note_controller.GetNote)
	req, _ := http.NewRequest("GET", "/note?created_date_time=123", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	} else {
		var fetched_note models.Note

		err := json.Unmarshal(test.Body.Bytes(), &fetched_note)

		if err != nil {
			t.Error(err.Error())
		}

		compareNote(t, note, fetched_note)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestGetAllNote_With_Single_Record(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/note", note_controller.GetNote)
	req, _ := http.NewRequest("GET", "/note", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	} else {
		var fetched_note []models.Note

		err := json.Unmarshal(test.Body.Bytes(), &fetched_note)

		if err != nil {
			t.Error(err.Error())
		}

		compareNote(t, note, fetched_note[0])
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestGetAllNote_With_Multiple_Record(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

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

	for _, note_data := range note_data_list {
		err := note_service.AddNote(note_data)

		if err != nil {
			t.Error(err.Error())
		}
	}

	server.GET("/note", note_controller.GetNote)
	req, _ := http.NewRequest("GET", "/note", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	} else {
		var fetched_notes []models.Note

		err := json.Unmarshal(test.Body.Bytes(), &fetched_notes)

		if err != nil {
			t.Error(err.Error())
		}

		for index := range note_data_list {
			var note models.Note
			note.FromMap(note_data_list[index])
			compareNote(t, note, fetched_notes[index])
		}
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestDeleteNote(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)
	err := note_service.AddNote(note_data)

	if err != nil {
		t.Error(err.Error())
	}

	server.DELETE("/note/:created_date_time", note_controller.DeleteNote)
	req, _ := http.NewRequest("DELETE", "/note/"+note.CreatedDateTime, bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestUpdateNote(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	note.Title = "updated_title"
	note.Content = "updated_content"

	notes_bytes, err := json.Marshal(note)
	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/note/:created_date_time", note_controller.UpdateNote)
	req, _ := http.NewRequest("PUT", "/note/"+note.CreatedDateTime, bytes.NewBuffer(notes_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestUpdateNote_ChangeAll(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	note.Title = "updated_title"
	note.Content = "updated_content"
	note.Attributes.IsFavourite = false
	note.Attributes.RequireMasterPassword = true

	notes_bytes, err := json.Marshal(note)
	if err != nil {
		t.Error(err.Error())
	}

	server.PUT("/note/:created_date_time", note_controller.UpdateNote)
	req, _ := http.NewRequest("PUT", "/note/"+note.CreatedDateTime, bytes.NewBuffer(notes_bytes))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestGetDecryptedContent_ValidCreatedDateTime(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	expected_content := note.Content

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/note/:created_date_time", note_controller.GetContent)
	req, _ := http.NewRequest("GET", "/note/"+note.CreatedDateTime, bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code != 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	} else {
		var content string

		err := json.Unmarshal(test.Body.Bytes(), &content)

		if err != nil {
			t.Error(err.Error())
		}

		if expected_content != content {
			t.Errorf("Mismatch in content\nExpeceted: %s\nActual: %s", expected_content, content)
		}
	}

	t.Cleanup(note_controller_test_cleanup)
}

func TestGetDecryptedContent_InvalidCreatedDateTime(t *testing.T) {
	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_service.SetMasterPassword("12345")

	note_service := new(services.NoteService)
	note_service.Init()

	note_controller := new(NoteController)
	note_controller.Init()

	server := gin.Default()
	test := httptest.NewRecorder()

	note_data := make(map[string]interface{})
	note_data["created_date_time"] = "123"
	note_data["title"] = "abc"
	note_data["content"] = "my content"
	note_data["attributes"] = map[string]interface{}{"is_favourite": true, "require_master_password": false}

	var note models.Note
	note.FromMap(note_data)

	err := note_service.AddNote(note_data)
	if err != nil {
		t.Error(err.Error())
	}

	server.GET("/note/:created_date_time", note_controller.GetContent)
	req, _ := http.NewRequest("GET", "/note/112", bytes.NewBuffer([]byte{}))

	server.ServeHTTP(test, req)

	if test.Code == 200 {
		var data string

		err := json.Unmarshal(test.Body.Bytes(), &data)

		if err != nil {
			t.Error(err.Error())
		}
		t.Error(data)
	}

	t.Cleanup(note_controller_test_cleanup)
}

func note_controller_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}
