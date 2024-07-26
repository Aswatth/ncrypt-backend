package services

import (
	"encoding/json"
	"log"
	"ncrypt/models"
	"ncrypt/utils/file_handler"
)

type DataService struct {
	file_name string
}

func (obj *DataService) Init(file_name string) *DataService {
	obj.file_name = file_name
	return obj
}

func (obj *DataService) AddData(new_data models.Data) error {

	existing_data, err := obj.GetAllData()

	if err != nil {
		log.Println(err.Error())
		return err
	}

	existing_data = append(existing_data, new_data)

	data_to_save, err := json.Marshal(existing_data)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return file_handler.Save(obj.file_name, data_to_save)

}

func (obj *DataService) GetAllData() ([]models.Data, error) {
	fetched_data, err := file_handler.Read(obj.file_name)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if(len(fetched_data)== 0){
		return []models.Data{}, nil
	}

	var data_list []models.Data

	err = json.Unmarshal(fetched_data, &data_list)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data_list, nil
}
