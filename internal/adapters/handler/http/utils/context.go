package utils

import (
	"ticketing-system/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

func GetUserID(c *fiber.Ctx) (uint, error) {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return 0, apperror.NewUnauthorized("ไม่ได้รับอนุญาต: ไม่พบข้อมูลผู้ใช้งาน")
	}

	if userID, ok := userIDLocal.(uint); ok {
		return userID, nil
	}

	if userIDFloat, ok := userIDLocal.(float64); ok {
		return uint(userIDFloat), nil
	}

	return 0, apperror.NewInternalServerError("ไม่สามารถแปลงประเภทข้อมูลผู้ใช้งานได้")
}