package validation

import (
	"fmt"
	"strings"
	"ticketing-system/internal/apperror"

	"github.com/go-playground/validator/v10"
)

func Validate(s any) error {
	var validate = validator.New()

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string

		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()

			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("ต้องระบุ '%s'", fieldName))

			case "min":
				messages = append(messages, fmt.Sprintf("'%s' ต้องมีความยาวอย่างน้อย %s ตัวอักษร", fieldName, fieldErr.Param()))

			case "oneof":
				messages = append(messages, fmt.Sprintf("'%s' ต้องเป็นค่าใดค่าหนึ่งใน [%s]", fieldName, fieldErr.Param()))

			case "gt":
				messages = append(messages, fmt.Sprintf("'%s' ต้องมีค่ามากกว่า %s", fieldName, fieldErr.Param()))
			}
		}

		return apperror.NewBadRequest(strings.Join(messages, ", "))
	}

	return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
}