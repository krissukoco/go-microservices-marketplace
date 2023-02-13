package handler

import (
	"context"
	"log"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
)

func (s *Server) Login(ctx context.Context, r *auth.LoginRequest) (*auth.LoginResponse, error) {
	res := &auth.LoginResponse{Status: statuscode.TokenInvalid}
	var u model.User
	if err := u.FindByEmail(r.Email); err != nil {
		if err.Error() == "user not found" {
			return res, nil
		}
		return res, err
	}
	if u.Email != r.Email {
		return res, nil
	}
	if err := u.ComparePassword(r.Password); err != nil {
		return res, nil
	}
	res.Status = statuscode.OK
	// Generate JWT Token
	token, err := u.GenerateJWT()
	if err != nil {
		return res, err
	}
	res.Token = token
	res.Email = u.Email
	res.FirstName = u.FirstName
	res.LastName = u.LastName

	return res, nil
}

func (s *Server) Refresh(ctx context.Context, r *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	res := &auth.RefreshResponse{Status: statuscode.TokenInvalid}
	var u model.User
	err := u.FromJWT(r.Token)
	if err != nil {
		log.Println("ERROR getting user from JWT: ", err)
		return res, nil // Prevent internal error being returned
	}
	res.Email = u.Email
	res.FirstName = u.FirstName
	res.LastName = u.LastName
	res.Status = statuscode.OK
	return res, nil
}
