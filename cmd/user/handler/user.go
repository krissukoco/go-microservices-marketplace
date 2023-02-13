package handler

import (
	"context"
	"log"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
	"github.com/krissukoco/go-microservices-marketplace/pkg/pb/auth"
)

func (s *Server) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	res := &auth.ChangePasswordResponse{Success: false}
	var u model.User
	err := u.FindByID(req.Id)
	if err != nil {
		log.Println("ERROR getting user from JWT: ", err)
		return res, nil // Prevent internal error being returned
	}
	if err := u.ComparePassword(req.OldPassword); err != nil {
		return res, nil
	}
	if err := u.ChangePassword(req.NewPassword); err != nil {
		return res, err
	}
	err = u.Save()
	if err != nil {
		return res, err
	}
	res.Success = true
	return res, nil
}
