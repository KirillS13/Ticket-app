package services

import (
	"Ticket-app/internal/domain"
	"context"
	"fmt"
	"time"
)

func (service *Service) BookTicket(userId int64, eventId int64, seats []string, ctx context.Context) error {

	tx, err := service.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, seat := range seats {
		var row, col int
		if _, err = fmt.Sscanf(seat, "%d-%d", &row, &col); err != nil {
			return fmt.Errorf("invalid seat format '%s': %w", seat, err)
		}
		var seatId int64
		var seatStatus string

		req := `SELECT id, status FROM seats WHERE event_id = $1 AND row_number = $2 AND line_number = $3 FOR UPDATE`
		err = tx.QueryRow(ctx, req, eventId, row, col).Scan(&seatId, &seatStatus)

		if err != nil {
			insertSeatReq := `INSERT INTO seats (event_id, row_number, line_number, status) VALUES ($1, $2, $3, 'booked') RETURNING id`
			err = tx.QueryRow(ctx, insertSeatReq, eventId, row, col).Scan(&seatId)
			if err != nil {
				return fmt.Errorf("failed to create seat: %w", err)
			}
		} else {
			if seatStatus == "booked" || seatStatus == "occupied" {
				return domain.ErrTicketBooked
			}
			_, err = tx.Exec(ctx, `UPDATE seats SET status = 'booked' WHERE id = $1`, seatId)
			if err != nil {
				return fmt.Errorf("failed to update seat status: %w", err)
			}
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		insertBookingReq := `
			INSERT INTO bookings (user_id, event_id, seat_id, created_at, expires_at, status) 
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, insertBookingReq, userId, eventId, seatId, time.Now(), expiresAt, "CONFIRMED")
		if err != nil {
			return fmt.Errorf("failed to insert booking: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil

}

func (service *Service) GetBookings(ctx context.Context, userId int64) ([]domain.BookingResponse, error) {
	req := `
		SELECT 
			b.id,
			COALESCE(e.name, 'Event'),
			COALESCE(e.room, 'Main Hall'),
			CONCAT('Row ', s.row_number, ' · Seat ', s.line_number) AS seat_info,
			COALESCE(e.price, 0),
			b.status,
			b.created_at
		FROM bookings b
		JOIN events e ON b.event_id = e.id
		JOIN seats s ON b.seat_id = s.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC`

	rows, err := service.db.Query(ctx, req, userId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get bookings: %w", err)
	}
	defer rows.Close()
	bookings := make([]domain.BookingResponse, 0)
	for rows.Next() {
		var booking domain.BookingResponse
		var createdAt time.Time
		if err = rows.Scan(&booking.ID, &booking.Event, &booking.Room, &booking.Seat, &booking.Price, &booking.Status, &createdAt); err != nil {
			return nil, fmt.Errorf("Failed to get bookings: %w", err)
		}
		booking.Date = createdAt.Format("Jan 02, 2006")
		bookings = append(bookings, booking)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to get bookings: %w", err)
	}
	return bookings, nil
}

func (service *Service) DeleteBooking(ctx context.Context, bookingId int64, userId int64) error {
	tx, err := service.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	req := `SELECT event_id, seat_id, status FROM bookings WHERE id = $1 AND user_id = $2 FOR UPDATE`
	row := tx.QueryRow(ctx, req, bookingId, userId)
	var seatId int64
	var eventId int64
	var currentStatus string
	err = row.Scan(&eventId, &seatId, &currentStatus)
	if err != nil {
		return fmt.Errorf("Failed to delete booking: %w", err)
	}
	if currentStatus == "cancelled" {
		return domain.ErrAlreadyCancelled
	}
	req = `UPDATE bookings SET status = $1 WHERE id = $2 AND user_id = $3 `
	_, err = tx.Exec(ctx, req, "cancelled", bookingId, userId)
	if err != nil {
		return fmt.Errorf("Failed to delete booking: %w", err)
	}
	req = `UPDATE seats SET status = $1 WHERE event_id = $2 AND id = $3 `
	_, err = tx.Exec(ctx, req, "available", eventId, seatId)
	if err != nil {
		return fmt.Errorf("Failed to delete booking: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete booking: %w", err)
	}
	return nil
}
