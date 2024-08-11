package services

import "ncrypt/models"

type ILoginService interface {
	Init()
	GetLoginData(name string) (models.Login, error)
	GetAllLoginData() ([]models.Login, error)
	GetDecryptedAccountPassword(login_data_name string, account_username string) (string, error)
	AddLoginData(login_data *models.Login) error
	UpdateLoginData(name string, login_data *models.Login) error
	DeleteLoginData(name string) error
	recryptData(data interface{}) error
	importData(login_data_list []models.Login) error
}

func InitBadgerLoginService() *LoginService {
	return &LoginService{}
}
