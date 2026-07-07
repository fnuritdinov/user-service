package server

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/service"
	"user-service/pkg/logger"
	user "user-service/userpb/v1"
)

type Server struct {
	user.UnimplementedUserServiceServer
	lg      logger.Logger
	service service.Service
}

func New(lg logger.Logger, service service.Service) *Server {
	return &Server{
		lg:      lg,
		service: service,
	}
}

func (s *Server) Add(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {

	request := models.User{
		Name:     req.Name,
		Password: req.Password,
		Age:      req.Age,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	userID, err := s.service.Add(ctx, request)
	if err != nil {
		return nil, err
	}

	return &user.CreateUserResponse{
		Id: int64(userID),
	}, nil
}

func (s *Server) GetByID(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	id := req.Id
	if id < 1 {
		return &user.GetUserResponse{}, errors.New("invalid id")
	}

	userFromDB, err := s.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user.GetUserResponse{
		Id:    int64(userFromDB.ID),
		Name:  userFromDB.Name,
		Phone: userFromDB.Phone,
		Email: userFromDB.Email,
		Age:   int64(userFromDB.Age),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	request := models.UpdateUser{
		ID:    int(req.Id),
		Name:  req.Name,
		Phone: req.Phone,
	}

	err := s.service.Update(ctx, request)
	if err != nil {
		return nil, err
	}

	return &user.UpdateUserResponse{
		Code:    0,
		Message: "user successful updated",
	}, nil
}
