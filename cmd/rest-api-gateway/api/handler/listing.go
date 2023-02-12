package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/schema"
)

// Listing listing
// swagger:model Listing
type Listing struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Listing
// @Accept */*
// @Produce json
// @Success 200 {object} schema.APIResponseOK{data=Listing} "Listing"
// @Security ApiKey
// @Router /products [get]
func GetAllListing(c *fiber.Ctx) error {
	return c.Status(200).JSON(schema.NewSuccessResponse(&Listing{
		ID:    1,
		Title: "Listing 1",
	}))
}
