package handlers

import (
	"Ticket-app/internal/domain"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Registration(c echo.Context) error {
	payload := new(domain.RegisterRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err := h.Service.CreateUser(payload.Email, payload.Username, payload.Password, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "user created",
	})
}

func (h *Handler) Login(c echo.Context) error {
	payload := new(domain.LoginRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")

	jwt, err := h.Service.Login(payload.Email, payload.Password, secretKey, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": jwt, "message": "user logged in"})
}
