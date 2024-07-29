package models

type Note struct {
	Title      string     `json:"title" bson:"title"`
	Content    string     `json:"content" bson:"content"`
	Attributes Attributes `json:"attributes" bson:"attributes"`
}
