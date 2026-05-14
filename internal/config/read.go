package config

import (
	"os"
	"encoding/json"
)

func Read() (Config, error) {
	
	configFile, err := getConfigFilePath()
	if err!=nil{
		return Config{}, err
	}
	data, err := os.ReadFile(configFile)
	if err!= nil { 
		return Config{}, err
	}

	var configuration Config 
	if err := json.Unmarshal(data, &configuration); err!=nil{
		return Config{}, err
	}

	return configuration, nil
}

func getConfigFilePath() (string, error){

	home_directory, err := os.UserHomeDir()
	if err!=nil{
		return "", err
	}

	configFilePath := home_directory + "/" + configFileName
	
	return configFilePath, nil
}