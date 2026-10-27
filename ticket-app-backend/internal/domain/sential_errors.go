package domain

import "errors"

var (
	ErrTicketBooked     = errors.New("Ticket is already booked")
	ErrAlreadyCancelled = errors.New("Ticket is already cancelled")
)
