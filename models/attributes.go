package models

type Attributes struct {
	IsFavourite           bool `json:"is_favourite" bson:"is_favourite"`
	RequireMasterPassword bool `json:"require_master_password" bson:"require_master_password"`
}

func (obj *Attributes) fromMap(data map[string]interface{}) *Attributes {
	obj.IsFavourite = data["is_favourite"].(bool)
	obj.RequireMasterPassword = data["require_master_password"].(bool)

	return obj
}
