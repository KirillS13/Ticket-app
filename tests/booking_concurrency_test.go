package tests

import (
	"Ticket-app/internal/database"
	"Ticket-app/internal/domain"
	"Ticket-app/internal/services"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/joho/godotenv"
)

func TestBookingConcurrency(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("Предупреждение: .env файл не найден: %v", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_TEST_NAME"))
	db, err := database.NewPostgresPool(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	EventId := 1
	SeatId := 2
	db.Exec(context.Background(), "DELETE FROM bookings WHERE seat_id = $1", SeatId)
	db.Exec(context.Background(), "DELETE FROM seats WHERE id = $1", SeatId)
	db.Exec(context.Background(), "DELETE FROM events WHERE id = $1", EventId)
	db.Exec(context.Background(), "DELETE FROM users WHERE id >= 1001 AND id <= 1050")

	sqlBytes, err := os.ReadFile("test_seeders/test_seeder_001.sql")
	if err != nil {
		t.Fatalf("Ошибка чтения SQL-файла: %v", err)
	}

	_, err = db.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		t.Fatalf("Ошибка выполнения SQL-сидера: %v", err)
	}
	t.Cleanup(func() {
		db.Exec(context.Background(), "DELETE FROM bookings WHERE seat_id = $1", SeatId)
		db.Exec(context.Background(), "DELETE FROM seats WHERE id = $1", 2)
		db.Exec(context.Background(), "DELETE FROM events WHERE id = $1", EventId)
		db.Exec(context.Background(), "DELETE FROM users WHERE id >= 1001 AND id <= 1050")
	})
	wg := sync.WaitGroup{}
	var successfulCount atomic.Int32
	var failedCount atomic.Int32
	service := services.NewUserService(db)
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func(userId int64) {
			defer wg.Done()
			err := service.BookTicket(userId, int64(EventId), int64(SeatId), context.Background())
			if err != nil {
				if errors.Is(err, domain.ErrTicketBooked) {
					failedCount.Add(1)
					return
				}
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}
			successfulCount.Add(1)
			return
		}(int64(1001 + i))
	}
	wg.Wait()
	if successfulCount.Load() != 1 {
		t.Errorf("Successfully count : %d", successfulCount.Load())
	}
	if failedCount.Load() != 49 {
		t.Errorf("Failed count : %d", failedCount.Load())
	}
	t.Logf("Successfully count : %d", successfulCount.Load())
	t.Logf("Failed count : %d", failedCount.Load())
	return
}
