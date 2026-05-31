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
		var messages []string

		for _, fieldErr := range validationErrors {
			fieldName := strings.ToLower(fieldErr.Field())

			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("ต้องระบุ '%s'", fieldName))

			case "min":
				messages = append(messages, fmt.Sprintf("'%s' ต้องมีความยาวอย่างน้อย %s ตัวอักษร", fieldName, fieldErr.Param()))

			case "oneof":
				if fieldErr.Field() == "Role" {
					messages = append(messages, "บทบาทผู้ใช้งานไม่ถูกต้อง (ต้องเป็น CASHIER หรือ SCANNER)")
				} else {
					messages = append(messages, fmt.Sprintf("'%s' ต้องเป็นค่าใดค่าหนึ่งใน [%s]", fieldName, fieldErr.Param()))
				}

			case "gt":
				messages = append(messages, fmt.Sprintf("'%s' ต้องมีค่ามากกว่า %s", fieldName, fieldErr.Param()))
			}
		}

		return apperror.NewBadRequest(strings.Join(messages, ", "))
	}

	return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
}