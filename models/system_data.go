package models

type SystemData struct {
	LoginCount                  int                         `json:"login_count" bson:"login_count"`
	LastLoginDateTime           string                      `json:"last_login" bson:"last_login"`
	IsLoggedIn                  bool                        `json:"is_logged_in" bson:"is_logged_in"`
	CurrentLoginDateTime        string                      `json:"current_login_date_time" bson:"current_login_date_time"`
	SessionDurationInMinutes    int                         `json:"session_duration_in_minutes" bson:"session_duration_in_minutes"`
	AutoBackupSetting           AutoBackupSetting           `json:"auto_backup_setting" bson:"auto_backup_setting"`
	PasswordGeneratorPreference PasswordGeneratorPreference `json:"password_generator_preference" bson:"password_generator_preference"`
}

func (obj *SystemData) FromMap(data map[string]interface{}) *SystemData {
	obj.LoginCount = int(data["login_count"].(float64))
	obj.LastLoginDateTime = data["last_login"].(string)
	obj.IsLoggedIn = data["is_logged_in"].(bool)
	obj.CurrentLoginDateTime = data["current_login_date_time"].(string)
	obj.AutoBackupSetting = *new(AutoBackupSetting).FromMap(data["auto_backup_setting"].(map[string]interface{}))
	obj.SessionDurationInMinutes = int(data["session_duration_in_minutes"].(float64))
	obj.PasswordGeneratorPreference = *new(PasswordGeneratorPreference).FromMap(data["password_generator_preference"].(map[string]interface{}))

	return obj
}
