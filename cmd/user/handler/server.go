package handler

import (
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
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
	database.AutoMigrate(db, &model.User{})
	return &Server{Pg: db, JwtSecret: config.Cfg.JWTSecret}, nil
}
