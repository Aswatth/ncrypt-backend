package models

type SystemData struct {
	LoginCount                  int                         `json:"login_count" bson:"login_count"`
	LastLoginDateTime           string                      `json:"last_login" bson:"last_login"`
	IsLoggedIn                  bool                        `json:"is_logged_in" bson:"is_logged_in"`
	CurrentLoginDateTime        string                      `json:"current_login_date_time" bson:"current_login_date_time"`
	AutomaticBackup             bool                        `json:"automatic_backup" bson:"automatic_backup"`
	AutomaticBackupLocation     string                      `json:"automatic_backup_location" bson:"automatic_backup_location"`
	BackupFileName              string                      `json:"backup_file_name" bson:"backup_file_name"`
	SessionDurationInMinutes    int                         `json:"session_duration_in_minutes" bson:"session_duration_in_minutes"`
	PasswordGeneratorPreference PasswordGeneratorPreference `json:"password_generator_preference" bson:"password_generator_preference"`
}

func (obj *SystemData) FromMap(data map[string]interface{}) *SystemData {
	obj.LoginCount = int(data["login_count"].(float64))
	obj.LastLoginDateTime = data["last_login"].(string)
	obj.IsLoggedIn = data["is_logged_in"].(bool)
	obj.CurrentLoginDateTime = data["current_login_date_time"].(string)
	obj.AutomaticBackup = data["automatic_backup"].(bool)
	obj.AutomaticBackupLocation = data["automatic_backup_location"].(string)
	obj.BackupFileName = data["backup_file_name"].(string)
	obj.SessionDurationInMinutes = int(data["session_duration_in_minutes"].(float64))
	obj.PasswordGeneratorPreference = *new(PasswordGeneratorPreference).FromMap(data["password_generator_preference"].(map[string]interface{}))

	return obj
}
