package handler

import (
	"fmt"
	"strconv"

	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	{
		api.Post("/create_segment", h.CreateSegment)
		api.Delete("/delete_segment", h.DeleteSegment)
		api.Post("/segment_user/:user_id", h.SegmentUser)
		api.Get("/active_segment/:user_id", h.GetActiveSegments)
	}

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

	// TODO: Create segment

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "segment created successfully",
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

	// TODO: Delete segment

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
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

func (h *Handler) SegmentUser(c *fiber.Ctx) error {
	_, err := GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// TODO: Segment user

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{})
}

func (h *Handler) GetActiveSegments(c *fiber.Ctx) error {
	_, err := GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// TODO: Get active segments

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{})
}
