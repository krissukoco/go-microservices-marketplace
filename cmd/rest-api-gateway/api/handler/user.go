package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/response"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
	"github.com/krissukoco/go-microservices-marketplace/internal/pb/auth"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"google.golang.org/grpc"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePassword(c *fiber.Ctx) error {
	var body ChangePasswordRequest
	if err := c.BodyParser(&body); err != nil {
		return response.APIErrorFromCode(c, statuscode.UnparsableBody)
	}
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return response.APIErrorFromCode(c, statuscode.Unauthorized)
	}
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		return response.APIErrorFromCode(c, statuscode.ServiceUnavailable)
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	res, err := client.ChangePassword(context.Background(), &auth.ChangePasswordRequest{
		Id:          userId,
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	})
	if err != nil {
		return response.APIErrorFromCode(c, statuscode.ServerError)
	}
	if res.Status != statuscode.OK {
		return response.APIErrorFromCode(c, res.Status)
	}
	return response.APIOkWithData(c, map[string]string{
		"message": "Password changed successfully",
	})
}
