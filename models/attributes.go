package models

type Attributes struct {
	IsFavourite           bool `json:"isFavourite" bson:"isFavourite"`
	RequireMasterPassword bool `json:"requireMasterPassword" bson:"requireMasterPassword"`
}
