package config

import (
	"log"
	"os"
)

type Config struct {
	PostgresConnStr string
	JWTSecret       string
}

var Cfg *Config

func InitializeConfigs() {
	pgConn, ok := os.LookupEnv("POSTGRES_URI")
	if !ok {
		log.Fatal("ERROR env var POSTGRES_URI not set")
	}
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("ERROR env var JWT_SECRET not set")
	}
	Cfg = &Config{
		PostgresConnStr: pgConn,
		JWTSecret:       jwtSecret,
	}
}
