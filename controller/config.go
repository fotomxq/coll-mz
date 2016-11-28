package controller

import (
	"encoding/json"
)

//Read the configuration file
func LoadConfigFile(configSrc string) (map[string]interface{}, error) {
	fileContent, err := LoadFile(configSrc)
	var data map[string]interface{}
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(fileContent, &data)
	return data, err
}

//Write the configuration file
func SaveConfigFile(configSrc string, configData interface{}) error {
	contentJson, err := json.Marshal(configData)
	if err != nil {
		return err
	}
	return WriteFile(configSrc, contentJson)
}
