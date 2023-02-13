package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/response"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
	"google.golang.org/grpc"
)

func RequireJWT(c *fiber.Ctx) error {
	authToken := c.Get("Authorization", "")
	if authToken == "" {
		return response.APIErrorFromCode(c, statuscode.TokenMissing)
	}
	if !strings.HasPrefix(authToken, "Bearer ") {
		return response.APIErrorFromCode(c, statuscode.TokenMalformed)
	}
	bearer := strings.TrimPrefix(authToken, "Bearer ")
	if bearer == "" {
		return response.APIErrorFromCode(c, statuscode.TokenMalformed)
	}
	// Validate JWT to user microservice
	conn, err := grpc.Dial("localhost:11000", grpc.WithInsecure())
	if err != nil {
		return response.APIErrorFromCode(c, statuscode.ServiceUnavailable)
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), &auth.RefreshRequest{
		Token: bearer,
	})
	if err != nil {
		return response.APIErrorFromCode(c, statuscode.ServerError)
	}
	if !res.Success {
		return response.APIErrorFromCode(c, statuscode.TokenInvalid)
	}
	// TODO: Get user ID instead
	c.Locals("userId", res.Email)

	return c.Next()
}
