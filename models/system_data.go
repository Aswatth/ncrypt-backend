package models

type SystemData struct {
	Login_count int    `json:"login_count" bson:"login_count"`
	Last_login  string `json:"last_login" bson:"last_login"`
}
