package config

import "os"

type DatabaseConfig struct {
	PostgresConnectionString string
}

var DB *DatabaseConfig

func initDatabaseConfig() {
	postgresConnectionString, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING")
	if !ok {
		panic("POSTGRES_CONNECTION_STRING not set")
	}
	DB = &DatabaseConfig{
		PostgresConnectionString: postgresConnectionString,
	}
}
