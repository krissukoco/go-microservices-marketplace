package config

import "log"

func InitializeConfigs() {
	initApiConfig()
	initDatabaseConfig()
	log.Println("Successfully initialized Configs")
}
