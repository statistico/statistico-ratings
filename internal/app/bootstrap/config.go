package bootstrap

import (
	"os"
)

type Config struct {
	Database
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func BuildConfig() *Config {
	config := Config{}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	return &config
}
