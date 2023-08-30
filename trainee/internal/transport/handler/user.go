package handler

import (
	"context"
	"net/http"
	"trainee/internal/core"
	"trainee/internal/core/errors"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetById(ctx context.Context, userId int64) (*core.User, error)
	ProcessUserData(ctx context.Context, userRequest *core.UserRequest) (*core.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) InitRoutes(app *fiber.App) {
	app.Get("/users", handler.GetById)
	app.Post("/users", handler.ProcessUser)
}

// GetById
// @Summary Get user by id
// @Tag users
// @Description Return user with it's active segments
// @Accept json
// @Produce json
// @Param user body core.UserId true "User JSON"
// @Success 200 {object} core.UserResponse
// @Failure 405 {object} core.FailureMessage
// @Failure 404 {object} core.FailureMessage
// @Router /users [get]
func (handler *UserHandler) GetById(ctx *fiber.Ctx) error {
	userRequest := new(core.UserRequest)
	err := ctx.BodyParser(userRequest)
	if err != nil {
		return ctx.Status(http.StatusMethodNotAllowed).
			JSON(fiber.Map{"message": "Check JSON in request", "error": err.Error()})
	}

	user, err := handler.service.GetById(ctx.UserContext(), userRequest.Id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundErrorWithMessage:
			return ctx.Status(http.StatusNotFound).
				JSON(fiber.Map{"message": e.Error()})
		default:
			return ctx.Status(500).
				JSON(fiber.Map{"message": "Unable to search user", "error": err.Error()})
		}
	}

	return ctx.Status(http.StatusOK).
		JSON(fiber.Map{"data": user})
}

// Process User
// @Summary Edit active segments
// @Tag users
// @Description Add and delete user's active segment
// @Accept json
// @Produce json
// @Param user body core.UserRequest true "User Request JSON"
// @Success 202 {object} core.UserResponse
// @Failure 405 {object} core.FailureMessage
// @Failure 404 {object} core.FailureMessage
// @Failure 500 {object} core.FailureMessage
// @Router /users [post]
func (handler *UserHandler) ProcessUser(ctx *fiber.Ctx) error {
	userRequest := new(core.UserRequest)
	err := ctx.BodyParser(userRequest)
	if err != nil {
		return ctx.Status(http.StatusMethodNotAllowed).
			JSON(fiber.Map{"message": "Check JSON in request", "error": err.Error()})
	}
	userResponse, err := handler.service.ProcessUserData(ctx.UserContext(), userRequest)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundErrorWithMessage:
			return ctx.Status(http.StatusNotFound).
				JSON(fiber.Map{"message": e.Error()})
		default:
			return ctx.Status(500).
				JSON(fiber.Map{"message": "Unable to edit user", "error": err.Error()})
		}
	}

	return ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"data": userResponse})
}
