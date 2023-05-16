package config

import (
	"log"
	"os"
)

type Config struct {
	JWTSecret string
}

var Cfg *Config

func InitializeConfigs() {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("ERROR env var JWT_SECRET not set")
	}
	Cfg = &Config{
		JWTSecret: jwtSecret,
	}
}
