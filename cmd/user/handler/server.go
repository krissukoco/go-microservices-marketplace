package handler

import "github.com/krissukoco/go-microservices-marketplace/internal/proto/auth"

// Generate Server
// Server is the gRPC server
type Server struct {
	auth.UnimplementedAuthServiceServer
}

var _ auth.AuthServiceServer = (*Server)(nil)
