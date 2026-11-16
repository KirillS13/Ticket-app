package domain

import "errors"

var (
	ErrTicketBooked     = errors.New("Ticket is already booked")
	ErrAlreadyCancelled = errors.New("Ticket is already cancelled")
	ErrAlreadyHeld      = errors.New("Seat is already held")
	NotFound            = errors.New("Seat not found")
	ErrWrongUser        = errors.New("Wrong user")
)
