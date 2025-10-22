package application

import "github.com/go-playground/validator/v10"

func ValidateDTO(dto interface{}) error {
	validate := validator.New()
	return validate.Struct(dto)
}
