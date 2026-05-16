package handler

import (
	"strconv"
	"ticketing-system/internal/core/model"
	"ticketing-system/internal/core/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AddUser(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.service.AddUser(user); err != nil {
		return handleError(c, err)
	}
	return c.Status(201).JSON(user)
}

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.service.EditUser(user); err != nil {
		return handleError(c, err)
	}
	return c.Status(204).JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.DeleteUser(uint(id)); err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}

func (h *UserHandler) FindUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	user, err := h.service.FindUserByID(uint(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(200).JSON(user)
}

func (h *UserHandler) FindUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username required"})
	}
	user, err := h.service.FindUserByUsername(username)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(200).JSON(user)
}
