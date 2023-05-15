package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/handler"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/middleware"
)

func UseDefault(app *fiber.App) {
	// @Security ApiKey
	v1 := app.Group("/api/v1")

	// Auth
	auth := v1.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.AuthRefresh)

	// User
	user := v1.Group("/user")
	user.Use(middleware.RequireJWT)
	user.Patch("/password", handler.UpdatePassword)

	// Products
	products := v1.Group("/products")
	products.Get("/", handler.GetAllListing)
}
