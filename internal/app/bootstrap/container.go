package bootstrap

import (
	"database/sql"
	"fmt"
)

type Container struct {
	Config   *Config
	Database *sql.DB
}

func BuildContainer(config *Config) Container {
	c := Container{
		Config: config,
	}

	c.Database = databaseConnection(config)

	return c
}

func databaseConnection(config *Config) *sql.DB {
	db := config.Database

	dsn := "host=%s port=%s user=%s " +
		"password=%s dbname=%s sslmode=disable"

	psqlInfo := fmt.Sprintf(dsn, db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(25)

	return conn
}
