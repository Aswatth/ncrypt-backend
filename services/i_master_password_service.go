package services

type IMasterPasswordService interface {
	Init()
	GetMasterPassword() (string, error)
	SetMasterPassword(password string) error
	UpdateMasterPassword(new_password string) error
	Validate(password string) (bool, error)
}

func InitBadgerMasterPasswordService() *MasterPasswordService {
	return &MasterPasswordService{}
}
