package handlers

import (
	"Ticket-app/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *services.Service
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
