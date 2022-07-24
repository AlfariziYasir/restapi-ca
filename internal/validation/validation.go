package validation

import (
	"restapi/internal/constant"

	"github.com/go-playground/validator/v10"
)

func Struct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return constant.NewErrFieldValidation(e)
		}
	}

	return nil
}
