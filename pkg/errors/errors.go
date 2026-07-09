package errors

import "errors"

var ErrFromValidate = errors.New("error from validate")
var ErrUserNotFound = errors.New("user not found")
var ErrGeneratePassword = errors.New("error from generate password")
var ErrBadRequest = errors.New("bad request")
