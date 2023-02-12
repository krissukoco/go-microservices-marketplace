package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

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
