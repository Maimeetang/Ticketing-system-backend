package http

import (
	"strconv"
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

// dto
type CreateUserRequest struct {
	Username           string `json:"username" validate:"required,min=4"`
	Password           string `json:"password" validate:"required,min=6"`
	Role               string `json:"role" validate:"required,oneof=CASHIER SCANNER"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

type UpdateUserRequest struct {
	Username           string `json:"username" validate:"required,min=4"`
	Role               string `json:"role" validate:"required,oneof=CASHIER SCANNER"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

// POST /users
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลผู้ใช้งานไม่ถูกต้อง")
	}
	
	if err := validateStruct(&req); err != nil {
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
	}

	err := h.service.Register(&user)
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

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลผู้ใช้งานไม่ถูกต้อง")
	}

	if err := validateStruct(&req); err != nil {
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

// DeleteUser
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		return err
	}

	if user == nil {
		return apperror.NewNotFound("ไม่พบผู้ใช้งาน")
	}

	res := newUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return err
	}

	resList := make([]userResponse, len(users))
	for i := range users {
		resList[i] = newUserResponse(&users[i])
	}
	

	return c.Status(fiber.StatusOK).JSON(resList)
}