package config

import (
	"encoding/json"
	"fmt"
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
		JWTSecret             string `json:"jwtSecret"`
		AccessIssuer          string `json:"accessIssuer"`
		RefreshIssuer         string `json:"refreshIssuer"`
		ExpirationTimeAccess  int64  `json:"expirationTimeAccesss"`
		ExpirationTimeRefresh int64  `json:"expirationTimeRefresh"`
		AutoLogoff            int64  `json:"autoLogOff"`
		Audience              string `json:"audience"`
		Host                  string `json:"host"`
		Port                  int64  `json:"port"`
	} `json:"authService"`
	Redis struct {
		Host     string `json:"host"`
		Port     int64  `json:"port"`
		Network  string `json:"network"`
		Password string `json:"password"`
		DB       int64  `json:"db"`
	} `json:"redis"`
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
	data, _ := json.MarshalIndent(c, "", " ")
	fmt.Println(string(data))
}
