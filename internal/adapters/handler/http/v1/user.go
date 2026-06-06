package v1

import (
	"strconv"
	"ticketing-system/internal/adapters/handler/http/v1/dto"
	"ticketing-system/internal/adapters/handler/http/validation"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// POST /users
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลผู้ใช้งานไม่ถูกต้อง")
	}
	
	if err := validation.Validate(&req); err != nil {
		return err
	}

	user := domain.User{
		Username:    req.Username,
		Password:    req.Password,
		Role:        domain.UserRole(req.Role),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		ReservePhoneNumber: req.ReservePhoneNumber,
		IsActive: true,
	}

	err := h.service.RegisterUser(&user)
	if err != nil {
		return err 
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เพิ่มผู้ใช้งานสำเร็จ",
	})
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลผู้ใช้งานไม่ถูกต้อง")
	}

	if err := validation.Validate(&req); err != nil {
		return err
	}

	user := domain.User{
		ID:			 uint(userID),
		Username:    req.Username,
		Role:        domain.UserRole(req.Role),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		ReservePhoneNumber: req.ReservePhoneNumber,
	}

	if err := h.service.UpdateUser(&user); err != nil {
		return err
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัปเดตผู้ใช้งานสำเร็จ",
	})
}

// PATCH /users/:id/disable
func (h *UserHandler) DisableUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	if err := h.service.DisableUser(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

// PATCH /users/:id/enable
func (h *UserHandler) EnableUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	if err := h.service.EnableUser(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	if user == nil {
		return apperror.NewNotFound("ไม่พบผู้ใช้งาน")
	}

	res := dto.NewUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(res)
}

// GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return err
	}

	resList := make([]dto.UserResponse, len(users))
	for i := range users {
		resList[i] = dto.NewUserResponse(&users[i])
	}
	

	return c.Status(fiber.StatusOK).JSON(resList)
}