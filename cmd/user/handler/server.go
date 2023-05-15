package handler

import (
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"
	"github.com/krissukoco/go-microservices-marketplace/proto/auth"
	"gorm.io/gorm"
)

// Generate Server
// Server is the gRPC server
type Server struct {
	auth.UnimplementedAuthServiceServer
	Pg        *gorm.DB
	JwtSecret string
}

var _ auth.AuthServiceServer = (*Server)(nil)

func NewServer() (*Server, error) {
	db, err := database.NewPostgresGorm()
	if err != nil {
		return nil, err
	}
	return &Server{Pg: db, JwtSecret: config.Cfg.JWTSecret}, nil
}
