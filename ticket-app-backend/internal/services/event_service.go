package services

import (
	"Ticket-app/internal/domain"
	"context"
)

func (service *Service) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	req := `SELECT id, name, room, price, start_time FROM events`
	query, err := service.db.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	var events []domain.Event
	defer query.Close()
	for query.Next() {
		var event domain.Event
		if err = query.Scan(&event.Id, &event.Name, &event.Room, &event.Price, &event.StartTime); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = query.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (service *Service) GetSeatsByEventID(ctx context.Context, id int) ([]domain.Seat, error) {
	req := `SELECT id, row_number, line_number, status FROM seats WHERE event_id = $1`
	query, err := service.db.Query(ctx, req, id)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	var seats []domain.Seat
	for query.Next() {
		var seat domain.Seat
		if err = query.Scan(&seat.Id, &seat.RowNumber, &seat.LineNumber, &seat.Status); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	if err = query.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}
