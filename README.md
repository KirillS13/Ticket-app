# Ticket Booking Platform

A full-stack event ticket booking application built to handle concurrent seat reservations safely and reliably.

## Architecture Overview

- **[Backend (`./ticket-app-backend`)](./ticket-app-backend/README.md):** Go REST API backed by PostgreSQL and Redis. Implements two-tier concurrency protection (Redis TTL seat holds + PostgreSQL transactional locking and unique constraint conflict mapping) verified under k6 load tests.
- **Frontend (`./ticket-app-frontend`):** Interactive seat map and booking flow. *(Update path if your frontend folder has a different name.)*

## Key Features

- **JWT Authentication:** Stateless user registration, login, and protected booking endpoints.
- **Real-Time Seat Reservation:** Temporary seat holds in Redis prevent checkout collisions while users complete their order.
- **Zero Double-Booking Guarantee:** Database-level constraints and row-level locking ensure each seat can only be sold once, even under simultaneous burst traffic.
- **Containerized Environment:** One-command local startup via Docker Compose.

## Quick Start

Start the backend infrastructure (Go API, PostgreSQL, Redis):

```bash
cd ticket-app-backend
docker-compose up -d --build
```

For detailed architecture notes, API setup, and k6 load test benchmarks, see the [Backend Documentation](./ticket-app-backend/README.md).