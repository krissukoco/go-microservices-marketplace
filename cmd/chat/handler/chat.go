package handler

import (
	"context"
	"log"

	"github.com/krissukoco/go-microservices-marketplace/internal/pb/chat"
)

type Server struct {
	chat.UnimplementedChatServiceServer
}

var _ chat.ChatServiceServer = (*Server)(nil)

func (s *Server) SendChat(ctx context.Context, c *chat.Chat) (*chat.Chat, error) {
	// TODO: Save to DB, send notification, etc.
	log.Println("Received Chat: ", c)
	reply := "Your Message is received: " + c.Message
	return &chat.Chat{Message: reply, Username: c.Username, Timestamp: c.Timestamp}, nil
}
