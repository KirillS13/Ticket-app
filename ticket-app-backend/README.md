# Ticket Booking Service (Backend)

High-concurrency ticket booking backend built with **Go**, **PostgreSQL**, and **Redis**. Designed to guarantee zero double-bookings under heavy traffic spikes when hundreds of users compete for the exact same seat simultaneously.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL (`pgx/v5`)
- **In-Memory Store / Locking:** Redis
- **Infrastructure:** Docker, Docker Compose
- **Load Testing:** Grafana k6

---

## Concurrency & Double-Booking Prevention

Ticket sales create extreme contention on popular seats. To ensure strict data consistency without sacrificing latency, the booking pipeline uses a **dual-layer concurrency control strategy**:

1. **Fast-Fail Seat Holds (Redis):**
    - When a user selects a seat, a temporary hold key (`hold:event:{eventId}:seat:{row}-{col}`) bound to their `user_id` is stored in Redis with a TTL.
    - During checkout, the service verifies hold existence and ownership in Redis before opening a database transaction, rejecting unauthorized or expired requests (`403` / `404`) at the memory layer.

2. **ACID Transactions & Unique Constraint Handling (PostgreSQL):**
    - Inside a database transaction, existing seat rows are locked via `SELECT ... FOR UPDATE`.
    - For seats not yet materialized in the `seats` table, concurrent `INSERT` races are caught at the database level via a composite unique constraint. PostgreSQL `SQLSTATE 23505` (unique violation) errors are intercepted in the repository/service layer and mapped cleanly to `domain.ErrTicketBooked` (**`409 Conflict`**), preventing `500 Internal Server Error` responses under contention.

---

## Load Testing Benchmarks (k6)

Tested locally in Docker with **100 Virtual Users (VUs)** simultaneously attempting to purchase the **same held seat** (`scripts/load-test.js`):

| Metric | Result |
| :--- | :--- |
| **Concurrent VUs** | 100 (targeting 1 seat) |
| **Successful Bookings (`200/201`)** | **1** (1.00%) |
| **Graceful Rejections (`409 Conflict`)** | **99** (99.00%) |
| **Server Errors (`500`)** | **0** (0.00%) |
| **Throughput** | ~2,008 req/s |
| **Avg Latency** | 25.08 ms |
| **p(95) Latency** | 33.39 ms |

### Running the Load Test

```bash
# 1. Seed a temporary seat hold in Redis for your test user ID
docker exec -it ticket_redis redis-cli SET "hold:event:1:seat:1-1" "<USER_ID>"

# 2. Run the k6 race-condition scenario inside the Compose network
docker run --rm -i -v "${PWD}:/scripts" --network ticket-app-backend_default grafana/k6 run /scripts/load-test.js
```

---

## Project Structure

```text
ticket-app-backend/
├── cmd/
│   └── api/                # Application entrypoint
├── internal/
│   ├── domain/             # Core entities and domain errors (ErrTicketBooked, etc.)
│   ├── services/           # Business logic (hold validation, transactional booking)
│   ├── repository/         # PostgreSQL & Redis data access
│   └── transport/          # HTTP handlers, routing, JWT middleware
├── migrations/             # SQL schema migrations
├── scripts/
│   └── load-test.js        # k6 concurrency stress test
├── Dockerfile
└── docker-compose.yml
```

---

## Getting Started

### Prerequisites
- Docker & Docker Compose

### Run Locally

1. Clone the repository and navigate to the backend directory:
   ```bash
   cd ticket-app-backend
   ```

2. Configure environment variables (copy `.env.example` to `.env` if applicable and set your DB/Redis/JWT secrets):
   ```bash
   cp .env.example .env
   ```

3. Build and start all services (API, PostgreSQL, Redis):
   ```bash
   docker-compose up -d --build
   ```

The API will be available at `http://localhost:8080`.