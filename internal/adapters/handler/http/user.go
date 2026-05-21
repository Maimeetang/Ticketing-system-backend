package http

import (
	"strconv"
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

type CreateUserRequest struct {
	Username           string `json:"username" validate:"required,min=4"`
	Password           string `json:"password" validate:"required,min=6"`
	Role               string `json:"role" validate:"required,oneof=cashier scanner"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

type UpdateUserRequest struct {
	Role               string `json:"role" validate:"required,oneof=cashier scanner"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
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
		"message": "user has been added.",
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user := domain.User{
		ID:                 uint(userID),
		Role:               domain.UserRole(req.Role),
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		PhoneNumber:        req.PhoneNumber,
		ReservePhoneNumber: req.ReservePhoneNumber,
	}

	if err := h.service.UpdateUser(&user); err != nil {
		return err
	}
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user has been update.",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		return err
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