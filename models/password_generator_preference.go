package models

type PasswordGeneratorPreference struct {
	HasDigits      bool `json:"has_digits" bson:"has_digits"`
	HasUpperCase   bool `json:"has_uppercase" bson:"has_uppercase"`
	HasSpecialChar bool `json:"has_special_char" bson:"has_special_char"`
	Length         int  `json:"length" bson:"length"`
}

func (obj *PasswordGeneratorPreference) FromMap(data map[string]interface{}) *PasswordGeneratorPreference {
	obj.HasDigits = data["has_digits"].(bool)
	obj.HasUpperCase = data["has_uppercase"].(bool)
	obj.HasSpecialChar = data["has_special_char"].(bool)
	obj.Length = (int)(data["length"].(float64))

	return obj
}
