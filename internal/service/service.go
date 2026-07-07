package service

import (
	"context"
	"fmt"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/pkg/errors"
	"user-service/pkg/password"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, request models.User) (int, error) {
	err := request.Validate()
	if err != nil {
		return 0, errors.ErrFromValidate
	}

	hashPassword, err := password.HashPassword(request.Password)
	if err != nil {
		return 0, errors.ErrGeneratePassword
	}

	userID, err := s.repo.Add(ctx, models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
		Phone:    request.Phone,
		Age:      request.Age,
	})
	if err != nil {
		return 0, fmt.Errorf("error from s.repo.Add")
	}

	return userID, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, request models.UpdateUser) error {

	err := request.Validate()
	if err != nil {
		return errors.ErrFromValidate
	}

	if request.ID < 1 {
		return errors.ErrFromValidate
	}

	err = s.repo.Update(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
