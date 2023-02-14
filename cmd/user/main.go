package main

import (
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/handler"
	"google.golang.org/grpc"

	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
)

func init() {
	config.InitializeConfigs()
	database.InitializePostgres()
}

func registerServices(server *grpc.Server) {
	auth.RegisterAuthServiceServer(server, &handler.Server{})
}

func main() {
	l, err := net.Listen("tcp", ":11000")
	if err != nil {
		log.Fatal("ERROR listening to port 11000: ", err)
	}
	server := grpc.NewServer()
	registerServices(server)

	log.Fatal(server.Serve(l))

}
