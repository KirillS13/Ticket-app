package handlers

import (
	"Ticket-app/internal/domain"
	"errors"
	"net/http"
	"strconv"

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
	if len(payload.Seats) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No seats selected"})
	}
	err := h.Service.BookTicket(c.Request().Context(), userId, payload.EventId, payload.Seats)
	if err != nil {
		if errors.Is(err, domain.ErrTicketBooked) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		if errors.Is(err, domain.NotFound) {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, domain.ErrWrongUser) {
			return c.JSON(http.StatusForbidden, err.Error())
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

func (h *Handler) HoldSeat(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User Not Authorized")
	}
	eventID := c.Param("id")
	eventIdInt, err := strconv.ParseInt(eventID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	payload := new(domain.HoldRequest)
	if err = c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.Service.HoldSeat(c.Request().Context(), payload.Seats, eventIdInt, userID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyHeld) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "seat has been held successfully")
}

func (h *Handler) ReleaseSeats(c echo.Context) error {
	userId, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User Not Authorized")
	}
	eventID := c.Param("id")
	eventIdInt, err := strconv.ParseInt(eventID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	payload := new(domain.ReleaseRequest)
	if err = c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = h.Service.ReleaseSeat(c.Request().Context(), payload.Seats, eventIdInt, userId)
	if err != nil {
		if errors.Is(err, domain.NotFound) {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, domain.ErrWrongUser) {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Seat has been released successfully")
}
