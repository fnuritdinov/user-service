package models

import (
	"time"

	errs "github.com/fnuritdinov/user-service/pkg/errors"
)

const UserRole = "userRole"
const AdminRole = "adminRole"

type User struct {
	ID       int
	Login    string
	Name     string
	Email    string
	Password string
	Phone    string
	Age      int32
	Role     string
}

func (u *User) Validate() error {
	if len(u.Name) == 0 && len(u.Email) == 0 && len(u.Password) == 0 && len(u.Phone) == 0 && u.Age == 0 {
		return errs.ErrFromValidate
	}
	return nil
}

type UpdateUser struct {
	ID    int
	Name  string
	Phone string
}

func (up *UpdateUser) Validate() error {
	if len(up.Name) == 0 && len(up.Phone) == 0 {
		return errs.ErrFromValidate
	}

	return nil
}

type VerifyRequest struct {
	Email string
	Code  string
}

type LoginRequest struct {
	Email    string
	Password string
}

func (lg *LoginRequest) Validate() error {
	if lg.Email == "" && lg.Password == "" {
		return errs.ErrFromValidate
	}
}

type HashTokenReq struct {
	ID        int
	UserID    int
	Hash      string
	ExpiredAt time.Time
}

type RefreshTokenReq struct {
	RefreshToken string
}

type RefreshAccessTokens struct {
	RefreshToken string
	AccessToken  string
	UserID       int
}
