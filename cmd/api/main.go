package main

import (
	"Ticket-app/internal/database"
	"Ticket-app/internal/handlers"
	"Ticket-app/internal/services"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Application struct {
	server  *echo.Echo
	handler handlers.Handler
}

func main() {
	e := echo.New()
	godotenv.Load(".env")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := database.NewPostgresPool(context.Background(), dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()
	userService := services.NewUserService(db)

	h := handlers.Handler{
		Service: userService,
	}

	app := Application{
		server:  e,
		handler: h,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go h.Service.StartWorker(ctx, 30*time.Second, 15*time.Minute)

	app.route(secretKey)
	go func() {
		if err = app.server.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = app.server.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
