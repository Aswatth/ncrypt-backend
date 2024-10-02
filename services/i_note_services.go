package services

import "ncrypt/models"

type INoteService interface {
	Init()
	GetNote(created_date_time string) (*models.Note, error)
	GetAllNotes() ([]models.Note, error)
	GetDecryptedContent(created_date_time string) (string, error)
	AddNote(note *models.Note) error
	UpdateNote(created_date_time string, updated_note models.Note) error
	DeleteNote(created_date_time string) error
	recryptData(password_data map[string]string) error
	importData(notes []models.Note) error
}

func InitBadgerNoteService() *NoteService {
	return &NoteService{}
}
