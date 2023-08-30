package handler

import (
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

func (h *Handler) CreateSegment(c *fiber.Ctx) error {
	// TODO: Implement segment creation
	panic("Implement segment creation")
}

func (h *Handler) DeleteSegment(c *fiber.Ctx) error {
	// TODO: Implement segment deletion
	panic("Implement segment deletion")
}

func (h *Handler) SegmentUser(c *fiber.Ctx) error {
	// TODO: Implement user segmentation
	panic("Implement user segmentation")
}

func (h *Handler) GetActiveSegments(c *fiber.Ctx) error {
	// TODO: Implement getting active segments
	panic("Implement getting active segments")
}
