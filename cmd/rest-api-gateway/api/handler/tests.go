package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
	chatPb "github.com/krissukoco/go-microservices-marketplace/internal/pb/chat"
	"google.golang.org/grpc"
)

func TestGRPC(c *fiber.Ctx) error {
	conn, err := grpc.Dial(config.Api.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		return c.SendStatus(500)
	}
	defer conn.Close()

	client := chatPb.NewChatServiceClient(conn)
	res, err := client.SendChat(context.Background(), &chatPb.Chat{
		Message:   "Hello from the client!",
		Username:  "johndoe",
		Timestamp: 123456789,
	})
	if err != nil {
		log.Println("ERROR in sending chat: ", err)
	}
	log.Println("Response from server: ", res)

	return c.Status(200).JSON(res)
}
