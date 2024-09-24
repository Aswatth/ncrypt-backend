package services

type IMasterPasswordService interface {
	Init()
	GetMasterPassword() (string, error)
	SetMasterPassword(master_password string) error
	UpdateMasterPassword(old_master_password string, new_master_password string) error
	Validate(password string) (bool, error)
	importData(password string) error
}

func InitBadgerMasterPasswordService() *MasterPasswordService {
	return &MasterPasswordService{}
}
