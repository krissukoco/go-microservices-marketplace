package main

import (
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/chat/handler"
	chatPb "github.com/krissukoco/go-microservices-marketplace/internal/proto/chat"
	"google.golang.org/grpc"
)

func registerServices(s *grpc.Server) {
	chatPb.RegisterChatServiceServer(s, &handler.Server{})
}

func main() {
	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal("ERROR listening to port 10000: ", err)
	}
	server := grpc.NewServer()
	registerServices(server)

	log.Fatal(server.Serve(l))
}
