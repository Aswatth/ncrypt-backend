package models

type Attributes struct {
	IsFavourite           bool `json:"isFavourite" bson:"isFavourite"`
	RequireMasterPassword bool `json:"requireMasterPassword" bson:"requireMasterPassword"`
}

func (obj *Attributes) fromMap(data map[string]interface{}) *Attributes {
	obj.IsFavourite = data["isFavourite"].(bool)
	obj.RequireMasterPassword = data["requireMasterPassword"].(bool)

	return obj
}
