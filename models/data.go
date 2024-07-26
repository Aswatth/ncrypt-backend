package models

type Data struct {
	Name string `json:"name" bson:"name"`
	Contents []Content `json:"contents" bson:"contents"`
}