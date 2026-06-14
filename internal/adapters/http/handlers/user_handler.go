package handlers

import (
	"strconv"
	dto "ticketing-system/internal/adapters/http/dto"
	e "ticketing-system/internal/core/error"
	s "ticketing-system/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service s.UserService
}

func NewUserHandler(service s.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("bad request")
	}

	ctx := c.UserContext()

	if err := h.service.RegisterUser(
		ctx,
		req.Username,
		req.Password,
		req.Role,
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
	); err != nil {
		return err 
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new user created",
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("bad request")
	}

	ctx := c.UserContext()

	if err := h.service.UpdateUser(
		ctx,
		uint(userID),
		req.Username,
		req.Role,
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
	); err != nil {
		return err
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user updated",
	})
}

func (h *UserHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	var req dto.UpdateUserStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("bad request")
	}

	if err := validate.Struct(req); err != nil {
		return e.NewBadRequest(err.Error())
	}

	ctx := c.UserContext()

	if err := h.service.UpdateUserStatus(ctx, uint(id), *req.IsActive); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user updated",
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	ctx := c.UserContext()

	user, err := h.service.GetUserByID(ctx, uint(id))
	if err != nil {
		return err
	}

	res := dto.NewUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()

	users, err := h.service.ListUsers(ctx)
	if err != nil {
		return err
	}

	resList := make([]dto.UserResponse, len(users))
	for i := range users {
		resList[i] = dto.NewUserResponse(&users[i])
	}
	

	return c.Status(fiber.StatusOK).JSON(resList)
}