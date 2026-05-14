package config 

const configFileName = ".gatorconfig.json"

type Config struct{ 
	DbURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) (error) { 

	cfg.CurrentUserName = username
	err := write(*cfg)
	if err!= nil{
		return err
	}
	return nil
}

