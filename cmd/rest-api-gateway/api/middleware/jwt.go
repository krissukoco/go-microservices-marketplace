package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/schema"
	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
	"google.golang.org/grpc"
)

func RequireJWT(c *fiber.Ctx) error {
	authToken := c.Get("Authorization", "")
	if authToken == "" {
		return c.Status(401).JSON(schema.NewErrorResponse(
			schema.ErrorTokenMissing,
			"Token is missing",
		))
	}
	if !strings.HasPrefix(authToken, "Bearer ") {
		return c.Status(401).JSON(schema.NewErrorResponse(
			schema.ErrorTokenMalformed,
			"Token is not of Bearer type",
		))
	}
	bearer := strings.TrimPrefix(authToken, "Bearer ")
	if bearer == "" {
		return c.Status(401).JSON(schema.NewErrorResponse(
			schema.ErrorTokenMalformed,
			"Token is malformed",
		))
	}
	// Validate JWT to user microservice
	conn, err := grpc.Dial("localhost:11000", grpc.WithInsecure())
	if err != nil {
		return c.Status(503).JSON(schema.NewErrorResponse(
			schema.ErrorServiceUnavailable,
			"Internal Server Error",
		))
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), &auth.RefreshRequest{
		Token: bearer,
	})
	if err != nil {
		return c.Status(500).JSON(schema.NewErrorResponse(
			schema.ErrorInternal,
			"Internal Server Error",
		))
	}
	if !res.Success {
		return c.Status(401).JSON(schema.NewErrorResponse(
			schema.ErrorTokenInvalid,
			"Token is invalid",
		))
	}
	// TODO: Get user ID instead
	c.Locals("userId", res.Email)

	return c.Next()
}
