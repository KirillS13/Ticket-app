package domain

type BookingResponse struct {
	ID     int64  `json:"id"`
	Event  string `json:"event"`
	Room   string `json:"room"`
	Seat   string `json:"seat"`
	Price  int64  `json:"price"`
	Date   string `json:"date"`
	Status string `json:"status"`
}
