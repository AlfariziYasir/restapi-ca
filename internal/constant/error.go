package constant

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrServer = errors.New("something went wrong")

	ErrUrlPathParameter  = errors.New("invalid url path parameter")
	ErrUrlQueryParameter = errors.New("invalid url query parameter")
	ErrRequestBody       = errors.New("invalid request body")
	ErrUnauthorized      = errors.New("you are not authorized to perform this action")
	ErrFieldValidation   = errors.New("field is not valid")

	ErrUserNotFound          = errors.New("user not found")
	ErrEmailRegistered       = errors.New("email already in use")
	ErrEmailNotRegistered    = errors.New("email not registered")
	ErrUserNameNotRegistered = errors.New("username not registered")
	ErrWrongPassword         = errors.New("password incorrect")

	ErrRecordNotFound = errors.New("record not found")
)

func NewErrFieldValidation(err validator.FieldError) error {
	return fmt.Errorf("%s: %w; format must be (%s=%s)", err.Field(), ErrFieldValidation, err.ActualTag(), err.Param())
}
