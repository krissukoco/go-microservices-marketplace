package handler

import (
	"os"

	"github.com/krissukoco/go-microservices-marketplace/cmd/product/handler/client"
	"github.com/krissukoco/go-microservices-marketplace/cmd/product/model"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"
	productPb "github.com/krissukoco/go-microservices-marketplace/proto/product"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	productPb.UnimplementedProductServiceServer

	Pg         *gorm.DB
	AuthClient *client.AuthClient
}

var _ productPb.ProductServiceServer = (*Server)(nil)

func NewServer() (*Server, error) {
	db, err := database.NewPostgresGorm()
	if err != nil {
		return nil, err
	}
	database.AutoMigrate(db, &model.Store{}, &model.Product{},
		&model.ProductReview{}, &model.ProductVariant{},
		&model.ProductReviewMedia{},
	)
	// Clients
	authServerUrl, ok := os.LookupEnv("AUTH_SERVER_URL")
	if !ok {
		authServerUrl = "localhost:50051"
	}
	authClient := &client.AuthClient{
		Url: authServerUrl,
		// TODO: use TLS/SSL certificates in production
		DialOptions: []grpc.DialOption{grpc.WithInsecure()},
	}

	return &Server{Pg: db, AuthClient: authClient}, nil
}
