package models

import (
	errs "user-service/pkg/errors"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone    string
	Age      int32
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
