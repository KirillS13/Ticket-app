package services

import (
	"Ticket-app/internal/domain"
	"context"
	"time"
)

func (service *Service) CreateEvent(room string, price int, name string, ctx context.Context) error {
	tx, err := service.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var event = domain.Event{
		Name:      name,
		Room:      room,
		Price:     price,
		StartTime: time.Now(),
	}
	var EventId int64
	req := `INSERT INTO events (name, room, price, start_time) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(ctx, req, event.Name, event.Room, event.Price, event.StartTime).Scan(&EventId)
	if err != nil {
		return err
	}
	req = `INSERT INTO seats (event_id, row_number, line_number, status)
			SELECT $1, r, l, 'available'
			FROM generate_series(1, 15) AS r
			CROSS JOIN generate_series(1, 15) AS l`
	_, err = tx.Exec(ctx, req, EventId)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
