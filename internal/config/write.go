package config 

import (
	"os"
	"encoding/json"
)

func write(cfg Config) error{ 

	jsonData, err := json.Marshal(cfg)
	if err!=nil{
		return err
	}

	configFile, err := getConfigFilePath()
	if err!= nil{
		return err
	}
	if err:= os.WriteFile(configFile, jsonData, 0644); err!= nil{
		return err
	}

	return nil
}