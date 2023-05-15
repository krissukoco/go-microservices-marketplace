package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
)

var PG *sql.DB

func InitializePostgres() {
	pg, err := sql.Open("postgres", config.Cfg.PostgresConnStr)
	if err != nil {
		log.Fatal("ERROR connecting to PostgreSQL: ", err)
	}
	PG = pg
	err = PG.Ping()
	if err != nil {
		log.Fatal("ERROR pinging PostgreSQL: ", err)
	}
	log.Println("Connected to PostgreSQL")
}

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
