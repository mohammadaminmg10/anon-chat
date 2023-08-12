package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
}

type JWTConfig struct {
	JWTKey        string `json:"JWTKey"`
	JWTCookieName string `json:"JWTCookieName"`
}

type CookieConfig struct {
	Name string `json:"name"`
}

type Configuration struct {
	Database DBConfig     `json:"database"`
	Jwt      JWTConfig    `json:"jwt"`
	Cookie   CookieConfig `json:"cookie"`
}

func LoadConfig(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
