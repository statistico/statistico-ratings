package bootstrap

import (
	"os"
)

type Config struct {
	AwsConfig
	Database
	KFactorMapping
	Sentry
	StatisticoDataService
	SupportedCompetitions []uint64
}

type AwsConfig struct {
	Key      string
	Region   string
	Secret   string
	S3Bucket string
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

	config.AwsConfig = AwsConfig{
		Key:      os.Getenv("AWS_KEY"),
		Region:   os.Getenv("AWS_REGION"),
		Secret:   os.Getenv("AWS_SECRET"),
		S3Bucket: os.Getenv("AWS_S3_BUCKET"),
	}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	config.KFactorMapping = map[uint64]float64{
		8: 5,
		9: 4,
		12: 3,
		14: 2,
	}

	config.Sentry = Sentry{DSN: os.Getenv("SENTRY_DSN")}

	config.StatisticoDataService = StatisticoDataService{
		Host: os.Getenv("STATISTICO_DATA_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_DATA_SERVICE_PORT"),
	}

	config.SupportedCompetitions = []uint64{8}

	return &config
}
