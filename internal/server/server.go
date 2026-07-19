package server

import (
	"context"
	"errors"

	user "github.com/fnuritdinov/proto/userpb"
	"github.com/fnuritdinov/user-service/internal/models"
	"github.com/fnuritdinov/user-service/internal/service"
	errs "github.com/fnuritdinov/user-service/pkg/errors"
	"github.com/fnuritdinov/user-service/pkg/logger"
	"go.uber.org/zap"
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

func (s *Server) Register(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error) {

	request := models.User{
		Name:     req.Name,
		Password: req.Password,
		Age:      req.Age,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	err := s.service.Register(ctx, request)
	if err != nil {
		if errors.Is(err, errs.ErrGeneratePassword) {
			return nil, errs.ErrBadRequest
		}
		if errors.Is(err, errs.ErrFromValidate) {
			return nil, errs.ErrBadRequest
		}
		s.lg.Error("error from s.service.Add",
			zap.Error(err))
		return nil, err
	}

	return &user.RegisterUserResponse{
		Message: "successfully registration",
	}, nil
}

func (s *Server) Verify(ctx context.Context, req *user.VerifyUserRequest) (*user.VerifyUserResponse, error) {

	request := models.VerifyRequest{
		Email: req.Email,
		Code:  req.Code,
	}

	id, err := s.service.Verify(ctx, request)
	if err != nil {
		s.lg.Error("error from s.service.Verify", zap.Error(err))
		return nil, err
	}

	return &user.VerifyUserResponse{
		Id: int64(id),
	}, nil
}

func (s *Server) Login(ctx context.Context, req *user.LoginUserRequest) (*user.LoginUserResponse, error) {

	request := models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	accessToken, refreshToken, err := s.service.Login(ctx, request)
	if err != nil {
		s.lg.Error("error from s.service.Login", zap.Error(err))
		return nil, err
	}

	return &user.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *user.RefreshTokenUserRequest) (*user.RefreshTokenUserResponse, error) {

	request := models.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	}

	tokens, err := s.service.RefreshToken(ctx, request)
	if err != nil {
		s.lg.Error("error from s.service.RefreshToken", zap.Error(err))
		return nil, err
	}

	return &user.RefreshTokenUserResponse{
		UserID:       int64(tokens.UserID),
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {

	id := req.Id
	if id < 1 {
		return nil, errors.New("id is empty")
	}

	userFromDB, err := s.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrNotFound
		}
		s.lg.Error("error from s.service.GetByID",
			zap.Error(err))
		return nil, err
	}
	return &user.GetUserResponse{
		User: &user.User{
			Id:    int64(userFromDB.ID),
			Name:  userFromDB.Name,
			Phone: userFromDB.Phone,
			Email: userFromDB.Email,
			Age:   userFromDB.Age,
			Role:  userFromDB.Role,
		},
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {

	request := models.UpdateUser{
		ID:    int(req.Id),
		Name:  req.Name,
		Phone: req.Phone,
	}

	err := s.service.Update(ctx, request)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserNotFound
		}
		if errors.Is(err, errs.ErrFromValidate) {
			return nil, errs.ErrBadRequest
		}
		s.lg.Error("error from s.service.Update",
			zap.Error(err))
		return nil, err
	}

	return &user.UpdateUserResponse{
		Code:    0,
		Message: "user successful updated",
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {

	usersFromDB, err := s.service.GetList(ctx)
	if err != nil {
		s.lg.Error("error from s.service.List",
			zap.Error(err))
		return nil, err
	}

	users := make([]*user.User, 0, len(usersFromDB))

	for _, u := range usersFromDB {
		users = append(users, &user.User{
			Id:    int64(u.ID),
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
			Age:   u.Age,
			Role:  u.Role,
		})
	}

	return &user.ListUsersResponse{
		Users: users,
	}, nil
}

func (s *Server) GetUserMe(ctx context.Context, req *user.GetUserMeRequest) (*user.GetUserMeResponse, error) {

	id := req.Id
	if id < 1 {
		return nil, errors.New("id is empty")
	}

	userFromDB, err := s.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrNotFound
		}
		s.lg.Error("error from s.service.GetByID",
			zap.Error(err))
		return nil, err
	}
	return &user.GetUserMeResponse{
		User: &user.User{
			Id:    int64(userFromDB.ID),
			Name:  userFromDB.Name,
			Phone: userFromDB.Phone,
			Email: userFromDB.Email,
			Age:   userFromDB.Age,
			Role:  userFromDB.Role,
		},
	}, nil
}

func (s *Server) CheckToken(ctx context.Context, req *user.CheckTokenRequest) (*user.CheckTokenResponse, error) {

	request := models.CheckToken{
		AccessToken: req.AccessToken,
	}

	myUser, err := s.service.CheckToken(ctx, request)
	if err != nil {
		s.lg.Error("error from s.service.CheckToken", zap.Error(err))
		return nil, err
	}

	return &user.CheckTokenResponse{
		User: &user.User{
			Id:    int64(myUser.ID),
			Name:  myUser.Name,
			Email: myUser.Email,
			Phone: myUser.Phone,
			Age:   myUser.Age,
			Role:  myUser.Role,
		},
	}, nil
}
