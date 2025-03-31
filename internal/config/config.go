package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB_URL          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const filename = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't get user's home directory: %s", err)
	}
	filePath := homeDir + "/" + filename
	return filePath, nil
}

func write(c Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("couldn't get filepath: %s", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open file: %s", err)
	}

	defer file.Close()

	encodedData, err := json.Marshal(&c)
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %s", err)
	}

	file.Write(encodedData)

	return nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("couldn't get filepath: %s", err)
	}

	userConf := Config{}
	body, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't read file: %s", err)
	}

	err = json.Unmarshal(body, &userConf)
	if err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal JSON: %s", err)
	}

	return userConf, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(*c)
	if err != nil {
		return fmt.Errorf("couldn't write data: %s", err)
	}
	return nil
}
