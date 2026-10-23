package services

import (
	"Ticket-app/internal/domain"
	"context"
	"fmt"
	"time"
)

func (service *Service) BookTicket(userId int64, eventId int64, seatId int64, ctx context.Context) error {

	tx, err := service.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var status string
	req := `SELECT status FROM seats WHERE event_id = $1 AND id = $2 FOR UPDATE`
	query := tx.QueryRow(ctx, req, eventId, seatId)
	err = query.Scan(&status)
	if err != nil {
		return fmt.Errorf("Failed to book ticket: %w", err)
	}
	if status != "available" {
		return domain.ErrTicketBooked
	}

	booking := domain.Booking{
		UserId:    userId,
		EventId:   eventId,
		SeatId:    seatId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Status:    "Booked",
	}
	req = `UPDATE seats SET status = 'booked' WHERE event_id = $1 AND id = $2`
	_, err = tx.Exec(ctx, req, eventId, seatId)
	if err != nil {
		return fmt.Errorf("Failed to book ticket: %w", err)
	}
	req = `INSERT INTO bookings (user_id, event_id, seat_id, created_at, expires_at, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, req, booking.UserId, booking.EventId, booking.SeatId, booking.CreatedAt, booking.ExpiresAt, booking.Status)
	if err != nil {
		return fmt.Errorf("Failed to book ticket: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to book ticket: %w", err)
	}
	return nil
}

func (service *Service) GetBookings(ctx context.Context, userId int64) ([]domain.Booking, error) {
	req := `SELECT id, event_id, seat_id, created_at, expires_at, status FROM bookings WHERE user_id = $1`
	rows, err := service.db.Query(ctx, req, userId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get bookings: %w", err)
	}
	defer rows.Close()
	bookings := make([]domain.Booking, 0)
	for rows.Next() {
		var booking domain.Booking
		booking.UserId = userId
		if err = rows.Scan(&booking.Id, &booking.EventId, &booking.SeatId, &booking.CreatedAt, &booking.ExpiresAt, &booking.Status); err != nil {
			return nil, fmt.Errorf("Failed to get bookings: %w", err)
		}
		bookings = append(bookings, booking)
		rows.Err()
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
