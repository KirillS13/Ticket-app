package domain

import "time"

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Event struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Room      string    `json:"room"`
	Price     int       `json:"price"`
	StartTime time.Time `json:"start_time"`
}

type Seat struct {
	Id         int64  `json:"id"`
	EventId    int64  `json:"event_id"`
	RowNumber  int    `json:"row_number"`
	LineNumber int    `json:"line_number"`
	Status     string `json:"status"`
}

type Booking struct {
	Id        int64     `json:"id"`
	EventId   int64     `json:"event_id"`
	UserId    int64     `json:"user_id"`
	SeatId    int64     `json:"seat_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"`
}
