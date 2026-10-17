package handlers

import (
	"Ticket-app/internal/domain"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) BookingHandler(c echo.Context) error {
	userId, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User Not Authorized")
	}
	payload := new(domain.BookingRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.Service.BookTicket(userId, payload.EventId, payload.SeatId, c.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrTicketBooked) {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Booking Successful",
	})
}

func (h *Handler) CheckBookings(c echo.Context) error {
	userId, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User Not Authorized")
	}

	bookings, err := h.Service.GetBookings(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bookings)
}

func (h *Handler) CancelBooking(c echo.Context) error {
	userId, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User Not Authorized")
	}
	payload := new(domain.CancelBookingRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.Service.DeleteBooking(c.Request().Context(), payload.BookingId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Booking has been deleted successfully",
	})
}
