package handler

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/api/response"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"github.com/krissukoco/go-microservices-marketplace/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Location        string `json:"location"`
}

// HealthCheck godoc
// @Summary Login
// @Description Login with Email and Password.
// @Tags Auth
// @Accept */*
// @Produce json
// @Success 200 {object} schema.APIResponseOK "User and Token"
// @Router /auth/login [get]
func Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return response.APIErrorFromCode(c, statuscode.UnparsableBody)
	}
	// Call User microservice to verify user
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Println("ERROR connecting to User microservice: ", err)
		return response.APIErrorFromCode(c, statuscode.ServiceUnavailable)
	}
	defer conn.Close()
	client := auth.NewAuthServiceClient(conn)
	res, err := client.Login(context.Background(), &auth.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		log.Println("ERROR: ", err)
		return response.APIErrorFromCode(c, statuscode.ServerError)
	}
	return response.APIOkWithData(c, map[string]interface{}{
		"token":      res.Token,
		"email":      res.Email,
		"first_name": res.FirstName,
		"last_name":  res.LastName,
	})
}

func AuthRefresh(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if token == "" {
		return response.APIErrorFromCode(c, statuscode.TokenMissing)
	}
	split := strings.Split(token, " ")
	if len(split) != 2 {
		return response.APIErrorFromCode(c, statuscode.TokenMalformed)
	}
	token = split[1]
	// Call User microservice to verify user
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Println("ERROR connecting to User microservice: ", err)
		return response.APIErrorFromCode(c, statuscode.ServerError)
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), &auth.RefreshRequest{Token: token})
	if err != nil {
		log.Println("ERROR refresh auth service: ", err)
		return response.APIErrorFromCode(c, statuscode.ServerError)
	}
	return response.APIOkWithData(c, map[string]interface{}{
		"id":         res.Id,
		"email":      res.Email,
		"first_name": res.FirstName,
		"last_name":  res.LastName,
	})
}

func Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return response.APIErrorFromCode(c, statuscode.UnparsableBody)
	}
	// TODO: Call User microservice to register user
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Println("ERROR connecting to User microservice: ", err)
		return response.APIErrorFromCode(c, statuscode.ServiceUnavailable)
	}
	defer conn.Close()
	client := auth.NewAuthServiceClient(conn)
	res, err := client.RegisterUser(context.Background(), &auth.RegisterRequest{
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})
	if err != nil {
		log.Println("ERROR: ", err)
		st, ok := status.FromError(err)
		if !ok {
			return response.APIErrorFromCode(c, statuscode.UnknownError)
		}
		code, m := statuscode.ParseGrpcErrMsg(st.Message())
		return response.APIErrorFromCode(c, code, m)
	}
	return response.APIOkWithData(c, res)
}
