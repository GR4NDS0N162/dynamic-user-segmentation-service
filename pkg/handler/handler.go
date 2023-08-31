package handler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	{
		api.Post("/create_segment", h.CreateSegment)
		api.Delete("/delete_segment", h.DeleteSegment)
		api.Post("/segment_user/:user_id", h.SegmentUser)
		api.Get("/active_segment/:user_id", h.GetActiveSegments)
		api.Get("/get_history", h.GetHistory)
	}

	app.Get("/download/:filename", h.GetFile)

	return app
}

type SegmentInput struct {
	Slug string
}

func (h *Handler) CreateSegment(c *fiber.Ctx) error {
	input := SegmentInput{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if input.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "slug cannot be empty",
		})
	}

	id, affected, err := h.service.CreateSegment(input.Slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if !affected {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": "segment already exists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"id": id,
	})
}

func (h *Handler) DeleteSegment(c *fiber.Ctx) error {
	input := SegmentInput{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if input.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "slug cannot be empty",
		})
	}

	deleted, err := h.service.DeleteSegment(input.Slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	} else if !deleted {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "segment doesn't exists",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(&fiber.Map{
		"message": "segment deleted successfully",
	})
}

func GetUserId(c *fiber.Ctx) (int, error) {
	id := c.Params("user_id", "")
	if id == "" {
		return 0, fmt.Errorf("user_id cannot be empty")
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("user_id must be a number")
	}

	return userId, nil
}

type SegmentationInput struct {
	Add    []string
	Remove []string
}

func (h *Handler) SegmentUser(c *fiber.Ctx) error {
	userId, err := GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	input := SegmentationInput{}
	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err = h.service.AddUserToSegments(userId, input.Add); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err = h.service.RemoveUserFromSegments(userId, input.Remove); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "user successfully segmented",
	})
}

func (h *Handler) GetActiveSegments(c *fiber.Ctx) error {
	userId, err := GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	slugs, err := h.service.GetActiveSegments(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&slugs)
}

type HistoryInput struct {
	Year  int
	Month int
}

func (h *Handler) GetHistory(c *fiber.Ctx) error {
	input := HistoryInput{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if !(1 <= input.Month && input.Month <= 12) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "month is out of range",
		})
	}

	filename, err := h.service.GetFile(input.Year, input.Month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"url": fmt.Sprintf("http://localhost:%v/download/%s", os.Getenv("WEB_PORT"), filename),
	})
}

func (h *Handler) GetFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	return c.Download(filename)
}
