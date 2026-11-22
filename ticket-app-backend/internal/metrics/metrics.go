package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var BookingTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "booking_total",
	Help: "The total number of bookings",
}, []string{"status"})

var HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration_seconds",
	Help: "Buckets: prometheus.DefBuckets",
}, []string{"method", "path", "status"})
