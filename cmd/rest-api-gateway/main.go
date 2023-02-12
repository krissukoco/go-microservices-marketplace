package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	_ "github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/docs"

	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/router"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
)

type Server struct {
	Api      *fiber.App
	Port     int
	Postgres *sql.DB
}

func (s *Server) ListenAndServe() error {
	return s.Api.Listen(fmt.Sprintf(":%d", s.Port))
}

// @title Marketplace Rest API
// @version 0.0.1
// @description Gateway API for the Marketplace microservices
// @contact.name Kris Sukoco
// @contact.email kristianto.sukoco@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @securitydefinitions.apikey  APIKey
// @in                          header
// @name                        x-api-key
// @BasePath /api/v1
// @schemes http
func newServer() (*Server, error) {
	srv := &Server{
		Port: 8000,
	}

	config.InitializeConfigs()
	// database.InitializePostgres()

	// API
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	router.UseDefault(app)
	app.Get("/docs/*", swagger.HandlerDefault)
	srv.Api = app

	return srv, nil
}

func main() {
	srv, err := newServer()
	if err != nil {
		log.Fatal("cannot start server! ", err)
	}
	log.Fatal(srv.ListenAndServe())
}
