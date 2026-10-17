package main

import (
	"Ticket-app/internal/middlewares"

	"github.com/labstack/echo/v4/middleware"
)

func (app Application) route(jwtSecret string) {

	app.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.server.GET("/", app.handler.HealthCheck)
	app.server.POST("/register", app.handler.Registration)
	app.server.POST("/login", app.handler.Login)

	protected := app.server.Group("")
	protected.Use(middlewares.JwtAuthMiddleware(jwtSecret))
	protected.POST("/bookings", app.handler.BookingHandler)
	protected.GET("/booking/me", app.handler.CheckBookings)
	protected.DELETE("/bookings/cancel", app.handler.CancelBooking)
}
