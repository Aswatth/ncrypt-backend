package models

type AutoBackupSetting struct {
	IsEnabled      bool   `json:"is_enabled" bson:"is_enabled"`
	BackupLocation string `json:"backup_location" bson:"backup_location"`
	BackupFileName string `json:"backup_file_name" bson:"backup_file_name"`
}

func (obj *AutoBackupSetting) FromMap(data map[string]interface{}) *AutoBackupSetting {
	obj.IsEnabled = data["is_enabled"].(bool)
	obj.BackupLocation = data["backup_location"].(string)
	obj.BackupFileName = data["backup_file_name"].(string)

	return obj
}
