package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/response"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"github.com/krissukoco/go-microservices-marketplace/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if ok {
			return response.APIErrorFromCode(c, 400)
		}
		code, msg := statuscode.ParseGrpcErrMsg(st.Message())
		return response.APIErrorFromCode(c, code, msg)
	}
	// log.Println("User ID: ", res.Id)
	c.Locals("userId", res.Id)

	return c.Next()
}
