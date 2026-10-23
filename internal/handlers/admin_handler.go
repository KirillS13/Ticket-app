package handlers

import (
	"Ticket-app/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateEvent(c echo.Context) error {
	payload := new(domain.CreateEventRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.Service.CreateEvent(payload.Room, payload.Price, payload.Name, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "event created",
	})
}
