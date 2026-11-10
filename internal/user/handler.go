package user

import (
	"errors"

	"api-starter/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Handler translates HTTP requests into service calls.
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /users.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req model.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		return toHTTPError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// List handles GET /users.
func (h *Handler) List(c *fiber.Ctx) error {
	res, err := h.service.List(c.UserContext())
	if err != nil {
		return toHTTPError(err)
	}
	return c.JSON(res)
}

// GetByID handles GET /users/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	res, err := h.service.GetByID(c.UserContext(), uint(id))
	if err != nil {
		return toHTTPError(err)
	}
	return c.JSON(res)
}

// toHTTPError maps service errors to HTTP errors. Unknown errors are logged
// and hidden from the client.
func toHTTPError(err error) error {
	var invalid *ValidationError
	switch {
	case errors.As(err, &invalid):
		return fiber.NewError(fiber.StatusBadRequest, invalid.Error())
	case errors.Is(err, ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrEmailTaken):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		log.Error().Err(err).Msg("user handler")
		return fiber.ErrInternalServerError
	}
}
