package bootstrap

import (
	"os"
)

type Config struct {
	CompetitionScoreMapping
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

type KFactorMapping map[uint64]uint8

type CompetitionScoreMapping map[uint64]uint16

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

	config.KFactorMapping = map[uint64]uint8{
		8: 20,
		9: 15,
		12: 10,
		14: 5,
	}

	config.CompetitionScoreMapping = map[uint64]uint16{
		8: 1500,
		9: 1450,
		12: 1400,
		14: 1350,
	}

	config.Sentry = Sentry{DSN: os.Getenv("SENTRY_DSN")}

	config.StatisticoDataService = StatisticoDataService{
		Host: os.Getenv("STATISTICO_DATA_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_DATA_SERVICE_PORT"),
	}

	config.SupportedCompetitions = []uint64{8}

	return &config
}
