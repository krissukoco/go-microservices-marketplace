package handler

import (
	"context"
	"log"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
	"github.com/krissukoco/go-microservices-marketplace/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	var u model.User
	if err := u.FindByEmail(s.Pg, in.Email); err != nil {
		if err == model.ErrUserNotFound || u.Email != in.Email {
			return nil, status.Error(codes.NotFound, "Email or password is invalid")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if err := u.ComparePassword(in.Password); err != nil {
		return nil, status.Error(codes.NotFound, "Email or password is invalid")
	}
	// Generate JWT Token
	token, err := u.GenerateJWT(s.JwtSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &auth.LoginResponse{
		Token:     token,
		Id:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, r *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	var u model.User
	err := u.FromJWT(s.Pg, r.Token, s.JwtSecret)
	if err != nil {
		log.Println("ERROR getting user from JWT: ", err)
		return nil, status.Error(codes.Unauthenticated, "Token is invalid")
	}
	return &auth.RefreshResponse{
		Id:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}
