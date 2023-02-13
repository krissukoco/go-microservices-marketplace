package main

import (
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/handler"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
	"golang.org/x/crypto/bcrypt"
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
	var u model.User
	err := u.FindByID("user_47a03ca8-6d56-42dd-901f-2af3dacb0e56")
	if err != nil {
		log.Fatal("ERROR getting user by ID ", err)
	}
	h, err := bcrypt.GenerateFromPassword([]byte("password"), 7)
	if err != nil {
		log.Fatal("ERROR generating hash ", err)
	}
	u.Password = string(h)
	err = u.Save()
	if err != nil {
		log.Fatal("ERROR saving user ", err)
	}

	l, err := net.Listen("tcp", ":11000")
	if err != nil {
		log.Fatal("ERROR listening to port 11000: ", err)
	}
	server := grpc.NewServer()
	registerServices(server)

	log.Fatal(server.Serve(l))

}
