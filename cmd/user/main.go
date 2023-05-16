package main

import (
	"fmt"
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/handler"
	"google.golang.org/grpc"

	authPb "github.com/krissukoco/go-microservices-marketplace/proto/auth"
)

var (
	Server *handler.Server
	Port   int = 11000
)

func init() {
	config.InitializeConfigs()
	// database.InitializePostgres()
	srv, err := handler.NewServer()
	if err != nil {
		log.Fatal("ERROR initializing server: ", err)
	}
	Server = srv
}

func registerServices(server *grpc.Server) {
	authPb.RegisterAuthServiceServer(server, Server)
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("ERROR listening to port %d: %v", Port, err)
	}
	server := grpc.NewServer()
	registerServices(server)

	log.Printf("Starting server on port %d", Port)
	log.Fatal(server.Serve(l))

}
