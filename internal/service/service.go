package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnuritdinov/user-service/internal/models"
	"github.com/fnuritdinov/user-service/internal/repository"
	"github.com/fnuritdinov/user-service/pkg/cache"
	errs "github.com/fnuritdinov/user-service/pkg/errors"
	"github.com/fnuritdinov/user-service/pkg/jwt"
	"github.com/fnuritdinov/user-service/pkg/password"
)

type Service struct {
	repo    repository.Repository
	myCache cache.ICache
}

func New(repo repository.Repository, myCache cache.ICache) Service {
	return Service{
		repo:    repo,
		myCache: myCache}
}

func (s *Service) Register(ctx context.Context, request models.User) error {
	err := request.Validate()
	if err != nil {
		return errs.ErrFromValidate
	}

	exists, err := s.repo.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user with this email already exists")
	}

	hashPassword, err := password.HashPassword(request.Password)
	if err != nil {
		return errs.ErrGeneratePassword
	}

	s.myCache.Save(ctx, request.Email, cache.CacheMemory{
		Email:    request.Email,
		Age:      request.Age,
		Phone:    request.Phone,
		Name:     request.Name,
		Password: hashPassword,
		Role:     request.Role,
	}, 5*time.Minute)

	return nil
}

type CacheMemory struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Age      int
	Role     string
}

func (s *Service) Verify(ctx context.Context, request models.VerifyRequest) (int, error) {

	var data CacheMemory
	err := s.myCache.Get(ctx, request.Email, &data)
	if err != nil {
		return 0, errors.New("user not found in cache")
	}

	id, err := s.repo.Register(ctx, models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Phone:    data.Phone,
		Age:      int32(data.Age),
		Role:     data.Role,
	})
	if err != nil {
		return 0, fmt.Errorf("error from s.repo.Register %w", err)
	}

	return id, nil
}

func (s *Service) Login(ctx context.Context, request models.LoginRequest) (string, string, error) {

	err := request.Validate()
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByEmail(ctx, request.Email)
	if err != nil {
		return "", "", fmt.Errorf("error from s.repo.GetByEmail %w", err)
	}

	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return "", "", fmt.Errorf("error from password.Compare %w", err)
	}

	accessToken, err := jwt.GenerateToken(user.ID, request.Email, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("error from generate accessToken %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("error from generate refreshToken %w", err)
	}

	hash := jwt.HashRefreshToken(refreshToken)

	err = s.repo.SaveRefreshToken(ctx, models.HashTokenReq{
		UserID:    user.ID,
		Hash:      hash,
		ExpiredAt: time.Now().Add(168 * time.Hour),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshToken(ctx context.Context, request models.RefreshTokenReq) (models.RefreshAccessTokens, error) {
	hash := jwt.HashRefreshToken(request.RefreshToken)

	token, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from s.repo.GetRefreshTokenByHash %w", err)
	}

	user, err := s.repo.GetByID(ctx, int64(token.UserID))
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from s.repo.GetByID %w", err)
	}

	if time.Now().After(token.ExpiredAt) {
		return models.RefreshAccessTokens{}, fmt.Errorf("time is expired %w", err)
	}

	err = s.repo.DeleteRefreshTokenByID(ctx, token.ID)
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from  s.repo.DeleteRefreshTokenByID %w", err)
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from generate refreshToken %w", err)
	}

	newHash := jwt.HashRefreshToken(newRefreshToken)

	err = s.repo.SaveRefreshToken(ctx, models.HashTokenReq{
		UserID:    token.UserID,
		Hash:      newHash,
		ExpiredAt: time.Now().Add(time.Hour * 168),
	})
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from s.repo.SaveRefreshToken")
	}

	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.RefreshAccessTokens{}, fmt.Errorf("error from generate accessToken %w", err)
	}

	return models.RefreshAccessTokens{
		RefreshToken: newHash,
		AccessToken:  accessToken,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("error from s.repo.GetByID %w", err)
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, request models.UpdateUser) error {

	err := request.Validate()
	if err != nil {
		return errs.ErrFromValidate
	}

	if request.ID < 1 {
		return errs.ErrFromValidate
	}

	err = s.repo.Update(ctx, request)
	if err != nil {
		return fmt.Errorf("error from s.repo.Update %w", err)
	}

	return nil
}

func (s *Service) GetList(ctx context.Context) ([]models.User, error) {

	users, err := s.repo.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from s.repo.GetList %w", err)
	}

	return users, nil
}

func (s *Service) CheckToken(ctx context.Context, token models.CheckToken) (models.User, error) {

	claims, err := jwt.ParseToken(token.AccessToken)
	if err != nil {
		return models.User{}, fmt.Errorf("error from jwt.ParseToken: %w", err)
	}

	user, err := s.repo.GetByID(ctx, int64(claims.UserID))
	if err != nil {
		return models.User{}, fmt.Errorf("error from s.repo.GetByID: %w", err)
	}

	return user, nil
}
