package bootstrap

import (
	"os"
)

type Config struct {
	Database
	KFactorMapping
	Sentry
	StatisticoDataService
	SupportedCompetitions []uint64
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type KFactorMapping map[uint64]float64

type Sentry struct {
	DSN string
}

type StatisticoDataService struct {
	Host string
	Port string
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

	config.KFactorMapping = map[uint64]float64{
		8: 15,
		9: 12,
		12: 10,
		14: 8,
	}

	config.Sentry = Sentry{DSN: os.Getenv("SENTRY_DSN")}

	config.StatisticoDataService = StatisticoDataService{
		Host: os.Getenv("STATISTICO_DATA_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_DATA_SERVICE_PORT"),
	}

	config.SupportedCompetitions = []uint64{8}

	return &config
}
