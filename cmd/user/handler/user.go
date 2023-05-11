package handler

import (
	"context"
	"log"

	"github.com/krissukoco/go-microservices-marketplace/cmd/user/model"
	"github.com/krissukoco/go-microservices-marketplace/internal/proto/auth"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
)

func (s *Server) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	res := &auth.ChangePasswordResponse{Status: statuscode.Unauthorized}
	var u model.User
	log.Println("User ID: ", req.Id)
	err := u.FindByID(req.Id)
	if err != nil {
		log.Println("ERROR getting user by ID ", err)
		res.Status = statuscode.TokenInvalid
		return res, nil // Prevent internal error being returned
	}
	log.Println("User: ", u)
	if err := u.ComparePassword(req.OldPassword); err != nil {
		log.Println("ERROR on compare password ", err)
		res.Status = statuscode.EmailOrPasswordInvalid
		return res, nil
	}
	if err := u.ChangePassword(req.NewPassword); err != nil {
		log.Println("ERROR on change password ", err)
		res.Status = statuscode.PasswordMalformed
		return res, err
	}
	err = u.Save()
	if err != nil {
		return res, err
	}
	res.Status = statuscode.OK
	return res, nil
}
