package database

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func errEnvVar(key string) error {
	return fmt.Errorf("env var %s not set", key)
}

func NewPostgresGorm() (*gorm.DB, error) {
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return nil, errEnvVar("POSTGRES_HOST")
	}
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, errEnvVar("POSTGRES_USER")
	}
	pwd, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, errEnvVar("POSTGRES_PASSWORD")
	}
	dbName, ok := os.LookupEnv("POSTGRES_DBNAME")
	if !ok {
		return nil, errEnvVar("POSTGRES_DBNAME")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, errEnvVar("POSTGRES_PORT")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=Asia/Jakarta",
		host, user, pwd, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	var err error
	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
