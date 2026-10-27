package services

import (
	"Ticket-app/internal/domain"
	"context"
	"fmt"
	"log"
	"time"
)

func (service *Service) CancelExpired(ctx context.Context, expirationTime time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-expirationTime)
	tr, err := service.db.Begin(ctx)
	if err != nil {
		return -1, fmt.Errorf("cancel expired transaction: %w", err)
	}
	defer tr.Rollback(ctx)
	req := `SELECT id, user_id, event_id, seat_id FROM bookings WHERE status = 'booked' AND created_at < $1 FOR UPDATE SKIP LOCKED`
	query, err := tr.Query(ctx, req, cutoffTime)
	if err != nil {
		return -1, fmt.Errorf("cancel expired transaction: %w", err)
	}
	defer query.Close()
	var bookings []domain.Booking

	for query.Next() {
		var booking domain.Booking
		if err = query.Scan(&booking.Id, &booking.UserId, &booking.EventId, &booking.SeatId); err != nil {
			return -1, fmt.Errorf("cancel expired transaction: %w", err)
		}
		bookings = append(bookings, booking)
	}
	if err = query.Err(); err != nil {
		return -1, fmt.Errorf("cancel expired transaction: %w", err)
	}
	for _, booking := range bookings {
		req = `UPDATE bookings SET status = 'expired' WHERE id = $1`
		_, err = tr.Exec(ctx, req, booking.Id)
		if err != nil {
			return -1, fmt.Errorf("cancel expired transaction: %w", err)
		}
		req = `UPDATE seats SET status = 'available' WHERE event_id = $1 AND id = $2`
		_, err = tr.Exec(ctx, req, booking.EventId, booking.SeatId)
		if err != nil {
			return -1, fmt.Errorf("cancel expired transaction: %w", err)
		}
	}
	if err = tr.Commit(ctx); err != nil {
		return -1, fmt.Errorf("cancel expired transaction: %w", err)
	}
	return int64(len(bookings)), nil
}

func (service *Service) StartWorker(ctx context.Context, interval time.Duration, expirationTime time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			count, err := service.CancelExpired(ctx, expirationTime)
			if err != nil {
				log.Printf("cancel expired transaction: %v", err)
			} else if count > 0 {
				log.Printf("Successfully canceled and released %d expired bookings", count)
			}
		}

	}
}
