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
	EventId int64    `json:"event_id"`
	Seats   []string `json:"seats"`
}

type CancelBookingRequest struct {
	BookingId int64 `json:"booking_id"`
}

type CreateEventRequest struct {
	Name  string `json:"name"`
	Room  string `json:"room"`
	Price int    `json:"price"`
}
