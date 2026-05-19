package http

import (
	"strconv"
	"ticketing-system/internal/adapters/handler/dto"
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

func (h *UserHandler) AddUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userDomain := domain.User{
		Username:    req.Username,
		Password:    req.Password,
		Role:        domain.UserRole(req.Role),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		ReservePhoneNumber: req.ReservePhoneNumber,
	}

	err := h.service.AddUser(userDomain)
	if err != nil {
		return err 
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user has been added.",
	})
}

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userDomain := domain.User{
		ID:                 uint(userID),
		Role:               domain.UserRole(req.Role),
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		PhoneNumber:        req.PhoneNumber,
		ReservePhoneNumber: req.ReservePhoneNumber,
	}

	if err := h.service.EditUser(userDomain); err != nil {
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

func (h *UserHandler) FindUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	user, err := h.service.FindUserByID(uint(id))
	if err != nil {
		return err
	}

	res := dto.UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Role:               string(user.Role),
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		PhoneNumber:        user.PhoneNumber,
		ReservePhoneNumber: user.ReservePhoneNumber,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *UserHandler) FindUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username parameter required"})
	}

	user, err := h.service.FindUserByUsername(username)
	if err != nil {
		return err
	}

	res := dto.UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Role:               string(user.Role),
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		PhoneNumber:        user.PhoneNumber,
		ReservePhoneNumber: user.ReservePhoneNumber,
	}

	return c.Status(200).JSON(res)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return err
	}

	// filter password out
	resList := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resList[i] = dto.UserResponse{
			ID:                 user.ID,
			Username:           user.Username,
			Role:               string(user.Role),
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			PhoneNumber:        user.PhoneNumber,
			ReservePhoneNumber: user.ReservePhoneNumber,
		}
	}

	return c.Status(fiber.StatusOK).JSON(resList)
}