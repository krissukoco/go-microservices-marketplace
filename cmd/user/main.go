package main

import (
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/handler"
	"google.golang.org/grpc"

	authPb "github.com/krissukoco/go-microservices-marketplace/proto/auth"
)

var Server *handler.Server

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
	l, err := net.Listen("tcp", ":11000")
	if err != nil {
		log.Fatal("ERROR listening to port 11000: ", err)
	}
	server := grpc.NewServer()
	registerServices(server)

	log.Fatal(server.Serve(l))

}
