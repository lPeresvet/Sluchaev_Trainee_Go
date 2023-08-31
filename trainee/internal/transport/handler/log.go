package handler

import (
	"context"
	"net/http"
	"trainee/internal/core"
	"trainee/internal/core/errors"

	"github.com/gofiber/fiber/v2"
)

type LogService interface {
	GetByUserIdAndMonth(ctx context.Context, request *core.LogRequestDto) (*core.LogResponse, error)
}

type LogHandler struct {
	service LogService
}

func NewLogHandler(service LogService) *LogHandler {
	return &LogHandler{service: service}
}

func (handler *LogHandler) InitRoutes(app *fiber.App) {
	app.Get("/log", handler.GetByIdAndMonth)
	app.Static("/files", "./csv")
}

// Get Log
// @Summary Get log
// @Tag log
// @Description Get log of user's segments operations
// @Accept json
// @Produce json
// @Param segment body core.LogRequestDto true "Log JSON"
// @Success 200 {object} core.LogResponse
// @Failure 405 {object} core.FailureMessage
// @Failure 500 {object} core.FailureMessage
// @Router /log [get]
func (handler *LogHandler) GetByIdAndMonth(ctx *fiber.Ctx) error {
	logRequest := new(core.LogRequestDto)
	err := ctx.BodyParser(logRequest)
	if err != nil {
		return ctx.Status(http.StatusMethodNotAllowed).
			JSON(fiber.Map{"message": "Check JSON in request", "error": err.Error()})
	}

	response, err := handler.service.GetByUserIdAndMonth(ctx.UserContext(), logRequest)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundErrorWithMessage:
			return ctx.Status(http.StatusNotFound).
				JSON(fiber.Map{"message": e.Error()})
		default:
			return ctx.Status(500).
				JSON(fiber.Map{"message": "Unable to search log", "error": err.Error()})
		}
	}
	base := ctx.BaseURL() + "/files"

	response.Link = base + response.Link
	return ctx.Status(http.StatusOK).
		JSON(fiber.Map{"data": response})
}
