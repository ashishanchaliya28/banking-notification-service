package handler

import (
	"strconv"
	"github.com/banking-superapp/notification-service/model"
	"github.com/banking-superapp/notification-service/service"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct{ svc service.NotificationService }

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) Send(c *fiber.Ctx) error {
	var req model.SendNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return respond(c, fiber.StatusBadRequest, nil, "invalid request body")
	}
	notif, err := h.svc.Send(c.Context(), &req)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusCreated, notif, "")
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "20"), 10, 64)
	notifications, total, err := h.svc.GetNotifications(c.Context(), userID, page, limit)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, fiber.Map{"notifications": notifications, "total": total}, "")
}

func (h *NotificationHandler) RegisterDevice(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	var req model.RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return respond(c, fiber.StatusBadRequest, nil, "invalid request body")
	}
	if err := h.svc.RegisterDevice(c.Context(), userID, &req); err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, fiber.Map{"registered": true}, "")
}

func (h *NotificationHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	var req model.UpdatePreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return respond(c, fiber.StatusBadRequest, nil, "invalid request body")
	}
	if err := h.svc.UpdatePreferences(c.Context(), userID, &req); err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, fiber.Map{"updated": true}, "")
}

func respond(c *fiber.Ctx, status int, data interface{}, errMsg string) error {
	if errMsg != "" {
		return c.Status(status).JSON(fiber.Map{"success": false, "error": errMsg})
	}
	return c.Status(status).JSON(fiber.Map{"success": true, "data": data})
}
