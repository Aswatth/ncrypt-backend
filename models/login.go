package models

type Login struct {
	Name       string     `json:"name" bson:"name"`
	URL        string     `json:"url" bson:"url"`
	Attributes Attributes `json:"attributes" bson:"attributes"`
	Accounts   []Account  `json:"accounts" bson:"accounts"`
}

func (obj *Login) FromMap(data map[string]interface{}) *Login {

	obj.Name = data["name"].(string)
	obj.URL = data["url"].(string)
	obj.Attributes = *new(Attributes).fromMap(data["attributes"].(map[string]interface{}))

	account_data_list := data["accounts"].([]interface{})
	for _, account_data := range account_data_list {
		obj.Accounts = append(obj.Accounts, *new(Account).fromMap(account_data.(map[string]interface{})))
	}

	return obj
}
