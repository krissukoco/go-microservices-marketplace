package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
	productPb "github.com/krissukoco/go-microservices-marketplace/proto/product"
	"google.golang.org/grpc"
)

func GetAllProducts(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	conn, err := grpc.Dial(config.Api.ProductServiceUrl, grpc.WithInsecure())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer conn.Close()
	client := productPb.NewProductServiceClient(conn)
	res, err := client.GetByFilters(c.Context(), &productPb.GetByFiltersRequest{
		Page:  int64(page),
		Limit: int64(limit),
	})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(res)
}
