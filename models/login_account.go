package models

type Account struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func (obj *Account) fromMap(data map[string]interface{}) *Account {
	obj.Username = data["username"].(string)
	obj.Password = data["password"].(string)

	return obj
}