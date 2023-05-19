package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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

func loadTlsCredentials() (credentials.TransportCredentials, error) {
	certFile := "../../certificates/auth-server-cert.pem"
	keyFile := "../../certificates/auth-server-key.pem"
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(cfg), nil
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("ERROR listening to port %d: %v", Port, err)
	}

	tlsCredentials, err := loadTlsCredentials()
	if err != nil {
		log.Fatalf("ERROR loading TLS credentials: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	registerServices(server)

	log.Printf("Starting server on port %d", Port)
	log.Fatal(server.Serve(l))

}
