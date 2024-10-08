package models

type Note struct {
	CreatedDateTime string     `json:"created_date_time" bson:"created_date_time"`
	Title           string     `json:"title" bson:"title"`
	Content         string     `json:"content" bson:"content"`
	Attributes      Attributes `json:"attributes" bson:"attributes"`
}

func (obj *Note) FromMap(data map[string]interface{}) *Note {
	obj.CreatedDateTime = data["created_date_time"].(string)
	obj.Title = data["title"].(string)
	obj.Content = data["content"].(string)
	obj.Attributes = *new(Attributes).fromMap(data["attributes"].(map[string]interface{}))

	return obj
}
