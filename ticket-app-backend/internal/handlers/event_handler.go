package handlers

import (
	"Ticket-app/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ShowEvents(c echo.Context) error {
	events, err := h.Service.GetAllEvents(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if events == nil {
		events = make([]domain.Event, 0)
	}
	return c.JSON(http.StatusOK, events)
}

func (h *Handler) ShowSeats(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	seats, err := h.Service.GetSeatsByEventID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var occupiedSeats []string
	for _, seat := range seats {
		if seat.Status == "booked" || seat.Status == "occupied" {
			occupiedSeats = append(occupiedSeats, fmt.Sprintf("%d-%d", seat.RowNumber, seat.LineNumber))
		}
	}
	if occupiedSeats == nil {
		occupiedSeats = make([]string, 0)
	}
	return c.JSON(http.StatusOK, occupiedSeats)
}
