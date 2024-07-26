package models

type Content struct {
	Name string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}