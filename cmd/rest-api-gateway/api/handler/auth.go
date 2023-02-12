package handler

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/schema"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
	"google.golang.org/grpc"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Location  string `json:"location"`
}

// HealthCheck godoc
// @Summary Login
// @Description Login with Email and Password.
// @Tags Auth
// @Accept */*
// @Produce json
// @Success 200 {object} schema.APIResponseOK "User and Token"
// @Security ApiKey
// @Router /auth/login [get]
func Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(422).JSON(schema.UnparsableBodyError())
	}
	// Call User microservice to verify user
	conn, err := grpc.Dial("localhost:11000", grpc.WithInsecure())
	if err != nil {
		log.Println("ERROR connecting to User microservice: ", err)
		return APIErrorFromCode(c, statuscode.ServiceUnavailable)
	}
	defer conn.Close()
	client := auth.NewAuthServiceClient(conn)
	res, err := client.Login(context.Background(), &auth.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		log.Println("ERROR: ", err)
		return APIErrorFromCode(c, statuscode.ServerUnknown)
	}
	if !res.Success {
		return APIErrorFromCode(c, statuscode.EmailOrPasswordInvalid)
	}
	return APIOkWithData(c, map[string]interface{}{
		"token":      res.Token,
		"email":      res.Email,
		"first_name": res.FirstName,
		"last_name":  res.LastName,
	})
}

func AuthRefresh(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if token == "" {
		return APIErrorFromCode(c, statuscode.TokenMissing)
	}
	split := strings.Split(token, " ")
	if len(split) != 2 {
		return APIErrorFromCode(c, statuscode.TokenMalformed)
	}
	token = split[1]
	// Call User microservice to verify user
	conn, err := grpc.Dial("localhost:11000", grpc.WithInsecure())
	if err != nil {
		return c.Status(503).JSON(schema.NewErrorResponse(
			schema.ErrorServiceUnavailable,
			"Internal Server Error",
		))
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), &auth.RefreshRequest{Token: token})
	if err != nil {
		log.Println("ERROR: ", err)
		return c.Status(500).JSON(schema.NewErrorResponse(
			schema.ErrorInternal,
			"Internal Server Error",
		))
	}
	if !res.Success {
		return c.Status(401).JSON(schema.NewErrorResponse(
			schema.ErrorTokenInvalid,
			"Token expired or invalid",
		))
	}
	return c.Status(200).JSON(schema.NewSuccessResponse(map[string]interface{}{
		"email":      res.Email,
		"first_name": res.FirstName,
		"last_name":  res.LastName,
	}))
}

func Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(422).JSON(schema.UnparsableBodyError())
	}
	// TODO: Call User microservice to register user
	return c.Status(200).JSON(schema.NewSuccessResponse(&body))
}
