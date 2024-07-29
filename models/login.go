package models

type Login struct {
	Name       string      `json:"name" bson:"name"`
	URL        string      `json:"url" bson:"url"`
	Attributes *Attributes `json:"attributes" bson:"attributes"`
	Accounts   []Account   `json:"accounts" bson:"accounts"`
}
