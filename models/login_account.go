package models

type Account struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
