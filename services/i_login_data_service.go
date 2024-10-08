package services

import "ncrypt/models"

type ILoginDataService interface {
	Init()
	GetLoginData(login_data_name string) (models.Login, error)
	GetAllLoginData() ([]models.Login, error)
	GetDecryptedAccountPassword(login_data_name string, account_username string) (string, error)
	AddLoginData(login_data map[string]interface{}) error
	UpdateLoginData(old_login_data_name string, login_data map[string]interface{}) error
	DeleteLoginData(login_data_name string) error
	recryptData(password_data map[string]string) error
	importData(login_datas []models.Login) error
}

func InitBadgerLoginService() *LoginDataService {
	return &LoginDataService{}
}
