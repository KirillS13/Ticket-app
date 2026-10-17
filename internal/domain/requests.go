package domain

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BookingRequest struct {
	EventId int64 `json:"event_id"`
	SeatId  int64 `json:"seat_id"`
}

type CancelBookingRequest struct {
	BookingId int64 `json:"booking_id"`
}
