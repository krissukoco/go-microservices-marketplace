package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/response"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
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
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
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
	if res.Status != statuscode.OK {
		return response.APIErrorFromCode(c, res.Status)
	}
	log.Println("User ID: ", res.Id)
	c.Locals("userId", res.Id)

	return c.Next()
}
