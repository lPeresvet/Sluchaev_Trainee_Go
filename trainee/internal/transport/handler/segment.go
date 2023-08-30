package handler

import (
	"context"
	"net/http"
	"trainee/internal/core"
	"trainee/internal/core/errors"

	"github.com/gofiber/fiber/v2"
)

type SegmentService interface {
	Create(ctx context.Context, segment *core.Segment) (*core.Segment, error)
	Delete(ctx context.Context, segment *core.Segment) error
}

type SegmentHandler struct {
	service SegmentService
}

func NewSegmentHandler(service SegmentService) *SegmentHandler {
	return &SegmentHandler{service: service}
}

func (handler *SegmentHandler) InitRoutes(app *fiber.App) {
	app.Post("/segment", handler.Create)
	app.Delete("/segment", handler.Delete)
}

// Create Segment
// @Summary Add new segment
// @Tag segments
// @Description Create new segment by slug
// @Accept json
// @Produce json
// @Param segment body core.Segment true "Segment JSON"
// @Success 201 {object} core.SegmentResponse
// @Failure 405 {object} core.FailureMessage
// @Failure 500 {object} core.FailureMessage
// @Router /segment [post]
func (handler *SegmentHandler) Create(ctx *fiber.Ctx) error {
	segment := new(core.Segment)
	err := ctx.BodyParser(segment)
	if err != nil {
		return ctx.Status(http.StatusMethodNotAllowed).
			JSON(fiber.Map{"message": "Check JSON in request"})
	}
	segment, err = handler.service.Create(ctx.UserContext(), segment)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"message": "Unable to create segment"})
	}
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"data": segment})
}

// Delete Segment
// @Summary Delete segment
// @Tag segments
// @Description Delete segment by slug
// @Accept json
// @Produce json
// @Param segment body core.Segment true "Segment JSON"
// @Success 202 {object} core.SegmentResponse
// @Failure 405 {object} core.FailureMessage
// @Failure 404 {object} core.FailureMessage
// @Failure 500 {object} core.FailureMessage
// @Router /segment [delete]
func (handler *SegmentHandler) Delete(ctx *fiber.Ctx) error {
	segment := new(core.Segment)
	err := ctx.BodyParser(segment)
	if err != nil {
		return ctx.Status(http.StatusMethodNotAllowed).
			JSON(fiber.Map{"message": "Check JSON in request"})
	}
	err = handler.service.Delete(ctx.UserContext(), segment)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundError:
			return ctx.Status(http.StatusNotFound).
			JSON(fiber.Map{"message": e.Error()})
		default:
			return ctx.Status(http.StatusNotFound).
			JSON(fiber.Map{"message": "Unable to delete segment"})
		}
		
	}
	return ctx.Status(http.StatusAccepted).
	JSON(fiber.Map{"message": "Segment <" + segment.Slug + "> deleted"})
}
