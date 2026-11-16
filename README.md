# Ticket Booking Platform 🎫

A full-stack event ticket booking application built with a high-concurrency **Go** backend and a modern web frontend. Designed with a focus on **Clean Architecture**, fault tolerance, and strict race condition prevention during peak ticket sales.

## 🗂 Repository Structure

```text
Ticket-app/
├── ticket-app-backend/     # REST API service (Go, Echo, PostgreSQL, Redis)
└── ticket-app-frontend/    # Web client application
```

> 📖 Detailed API documentation, database schema, and backend setup instructions are available in [ticket-app-backend/README.md](./ticket-app-backend/README.md).

---

## 🚀 Key Engineering Features

* **Two-Layer Double-Booking Protection:**
    * **Temporary Seat Holds (Redis):** Locks selected seats for a specific user before checkout (`hold:event:<id>:seat:<row>-<col>`).
    * **Database-Level ACID Guarantees (PostgreSQL):** Uses pessimistic locking (`SELECT ... FOR UPDATE`) and unique constraints to eliminate race conditions.
* **Load-Tested Concurrency (k6):** Verified under simultaneous load (100 concurrent virtual users competing for the exact same seat) — guaranteeing strictly 1 successful booking and 0 overbookings.
* **Clean Architecture:** Clear separation of concerns across `domain`, `handlers`, `services`, and `middlewares` layers.
* **JWT Authentication:** Stateless authentication securing user sessions and booking endpoints.

---

## 🛠 Tech Stack

**Backend:**
* **Language & Framework:** Go, Echo
* **Databases & Caching:** PostgreSQL, Redis
* **Infrastructure & Testing:** Docker, Docker Compose, k6 (load testing)

**Frontend:**
* **Package Manager:** pnpm
* **Stack:** *(укажи свой стек из `package.json`, например: React / Next.js / TypeScript / Tailwind CSS)*

---

## ⚡ Quick Start

### 1. Start Backend & Databases
```bash
cd ticket-app-backend
cp .env.example .env
docker-compose up -d --build
```

### 2. Start Frontend
```bash
cd ticket-app-frontend
pnpm install
pnpm dev
```

### 3. Run Concurrency Load Test (Optional)
```bash
cd ticket-app-backend
docker run --rm -i -v "${PWD}:/scripts" --network ticket-app-backend_default grafana/k6 run /scripts/load-test.js
```