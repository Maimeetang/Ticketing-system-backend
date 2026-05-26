package http

import (
	"fmt"
	"strings"
	"ticketing-system/internal/apperror"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fieldName := strings.ToLower(fieldErr.Field())

			switch fieldErr.Tag() {
			case "required":
				return apperror.NewBadRequest(fmt.Sprintf("field '%s' is required", fieldName))
			case "min":
				return apperror.NewBadRequest(fmt.Sprintf("field '%s' must be at least %s characters long", fieldName, fieldErr.Param()))
			case "oneof":
				if fieldErr.Field() == "Role" {
					return apperror.NewBadRequest("invalid user role: must be either 'CASHIER' or 'SCANNER'")
				}
				return apperror.NewBadRequest(fmt.Sprintf("field '%s' must be one of: %s", fieldName, fieldErr.Param()))
			case "gt":
				return apperror.NewBadRequest(fmt.Sprintf("field '%s' must be greater than %s", fieldName, fieldErr.Param()))
			}
		}
	}

	return apperror.NewBadRequest("request validation failed")
}
