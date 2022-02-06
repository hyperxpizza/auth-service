package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Host     string `json:"host"`
	} `json:"database"`
	AuthService struct {
		JWTSecret           string `json:"jwtSecret"`
		Issuer              string `json:"issuer"`
		ExpirationTimeHours int64  `json:"expirationTimeHours"`
		Audience            string `json:"audience"`
		Host                string `json:"host"`
		Port                int64  `json:"port"`
	} `json:"authService"`
}

func NewConfig(pathToFile string) (*Config, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var c Config

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) PrettyPrint() {

}
