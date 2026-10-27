# Ticket Booking API 🎫

A high-performance and reliable RESTful API service for event ticket booking, built with **Go**, the **Echo** framework, and **PostgreSQL**.

This project is designed with a strong focus on **Clean Architecture**, fault tolerance, and race condition prevention (double-booking) during concurrent seat reservations.

---

## 🚀 Key Features

* **Clean Architecture**: Clear separation of concerns into distinct application layers (`domain`, `handlers`, `services`, `middlewares`).
* **Double-Booking Protection**: Uses pessimistic locking (`SELECT ... FOR UPDATE`) in PostgreSQL to prevent concurrent reservations of the same seat.
* **Transactional Integrity**: All mutating business operations (booking creation, cancellation) are strictly atomic and executed within ACID transactions.
* **Security & Authentication**: JWT-based user authentication (JSON Web Tokens) and password hashing using `bcrypt`.
* **Standardized Error Handling**: Leverages Sentinel Errors to cleanly map business conflicts (400 Bad Request, 409 Conflict) and system errors (500 Internal Server Error).

---

## 🛠 Tech Stack

* **Language**: Go (1.20+)
* **Web Framework**: Echo v4
* **Database**: PostgreSQL
* **Containerization**: Docker & Docker Compose
* **Authentication**: JWT (`golang-jwt`)

---

## 🗺️ Roadmap

- [x] **Phase 1: Core API & Transaction Safety**
  - [x] Clean Architecture setup (`domain`, `handlers`, `services`)
  - [x] JWT Authentication & password hashing (`bcrypt`)
  - [x] Core booking flow (`BookTicket`, `GetBookings`, `CancelBooking`)
  - [x] Pessimistic locking (`SELECT ... FOR UPDATE`) to prevent race conditions
  - [x] Transactional integrity & Sentinel Error mapping

- [ ] **Phase 2: Reliability & Automation** *(In Progress)*
  - [ ] Concurrent load testing (stress test with 50–100 parallel booking requests)
  - [ ] Background worker (`time.Ticker`) for automatic expiration of unpaid/pending bookings

- [ ] **Phase 3: Performance & Scalability**
  - [ ] Redis caching layer for available seats and event catalogs
  - [ ] Rate limiting middleware for API protection

- [ ] **Phase 4: Observability & DevOps**
  - [ ] Structured logging using `slog` or `zap`
  - [ ] Metrics collection with Prometheus and Grafana
  - [ ] CI/CD pipeline via GitHub Actions
