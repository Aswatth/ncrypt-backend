package services

import (
	"errors"
	"ncrypt/models"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"ncrypt/utils/logger"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type NoteService struct {
	database                database.IDatabase
	master_password_service IMasterPasswordService
}

func (obj *NoteService) Init() {
	logger.Log.Printf("Initializing notes service")
	logger.Log.Printf("Loading .env variables")
	godotenv.Load("../.env")

	logger.Log.Printf("Setting up database")
	obj.database = database.InitBadgerDb()
	obj.database.SetDatabase(os.Getenv("NOTE_DB_NAME"))

	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()

	logger.Log.Printf("DONE")
}

func (obj *NoteService) GetNote(created_date_time string) (*models.Note, error) {
	logger.Log.Printf("Getting note from databse")
	var note models.Note
	fetched_note, err := obj.database.GetData(created_date_time)

	if err != nil {
		return &models.Note{}, err
	}
	note.FromMap(fetched_note.(map[string]interface{}))

	return &note, err
}

func (obj *NoteService) GetAllNotes() ([]models.Note, error) {
	logger.Log.Printf("Getting all notes")
	fetched_datas, err := obj.database.GetAllData()

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return nil, err
	}

	var notes []models.Note
	for _, fetched_data := range fetched_datas {
		notes = append(notes, *new(models.Note).FromMap(fetched_data.(map[string]interface{})))
	}

	return notes, nil
}

func (obj *NoteService) GetDecryptedContent(created_date_time string) (string, error) {
	logger.Log.Printf("Decrypting note content")
	master_password, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return "", err
	}

	fetched_note, err := obj.GetNote(created_date_time)

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return "", err
	}

	decrypted_content, err := encryptor.Decrypt(fetched_note.Content, master_password+fetched_note.CreatedDateTime)

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return "", err
	}

	return decrypted_content, err
}

func (obj *NoteService) AddNote(note *models.Note) error {
	logger.Log.Printf("Adding note")
	master_password_hash, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		if strings.ToUpper(err.Error()) == "KEY NOT FOUND" {
			err = errors.New("master_password not set")
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
		return err
	}

	existing_note, err := obj.GetNote(note.CreatedDateTime)
	if err != nil && err != badger.ErrKeyNotFound {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}
	if existing_note.CreatedDateTime != "" {
		logger.Log.Printf("ERROR: %s", "CONFLICTING NAMES")
		return errors.New(existing_note.CreatedDateTime + " already exists")
	}

	logger.Log.Printf("Encrypting content")
	encrypted_content, err := encryptor.Encrypt(note.Content, master_password_hash+note.CreatedDateTime)

	if err != nil {
		return err
	}

	note.Content = encrypted_content

	logger.Log.Printf("Storing to DB")
	err = obj.database.AddData(note.CreatedDateTime, note)

	return err
}
func (obj *NoteService) UpdateNote(created_date_time string, updated_note models.Note) error {
	logger.Log.Printf("Updating note")
	fetched_note, err := obj.GetNote(created_date_time)

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return err
	}

	master_password, err := obj.master_password_service.GetMasterPassword()
	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return err
	}

	decrypted_content, err := encryptor.Decrypt(updated_note.Content, master_password+created_date_time)

	if err == nil {
		updated_note.Content = decrypted_content
	}

	updated_note.CreatedDateTime = fetched_note.CreatedDateTime

	return obj.database.AddData(updated_note.CreatedDateTime, updated_note)
}

func (obj *NoteService) DeleteNote(created_date_time string) error {
	logger.Log.Printf("Deleting note")
	err := obj.database.DeleteData(created_date_time)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}
	return err
}
func (obj *NoteService) recryptData(password_data map[string]string) error {
	logger.Log.Printf("Recrypting notes content")
	notes, err := obj.GetAllNotes()

	if err != nil {
		logger.Log.Printf("ERROR: " + err.Error())
		return err
	}

	logger.Log.Printf("Decrypting all notes content")
	old_password := password_data["OLD_PASSWORD"]
	new_password := password_data["NEW_PASSWORD"]

	for i := range len(notes) {
		decrypted_content, err := encryptor.Decrypt(notes[i].Content, old_password+notes[i].CreatedDateTime)

		if err != nil {
			logger.Log.Printf("ERROR: " + err.Error())
			return err
		}

		notes[i].Content = decrypted_content
	}

	logger.Log.Printf("Recrypting content")
	for i := range len(notes) {
		notes[i].Content, err = encryptor.Encrypt(notes[i].Content, new_password+notes[i].CreatedDateTime)

		if err != nil {
			logger.Log.Printf("ERROR: " + err.Error())
			return err
		}
		obj.database.AddData(notes[i].CreatedDateTime, notes[i])
	}

	return nil
}

func (obj *NoteService) importData(notes []models.Note) error {
	logger.Log.Printf("Importing notes")
	os.RemoveAll("data/" + os.Getenv("NOTE_DB_NAME"))

	for _, note := range notes {
		err := obj.database.AddData(note.CreatedDateTime, note)
		if err != nil {
			logger.Log.Printf("ERROR: " + err.Error())
			return err
		}
	}

	return nil
}
