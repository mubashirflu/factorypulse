# FactoryPulse 🏭

Full-stack industrial machine monitoring and predictive maintenance platform built with **Vue 3**, **Go**, and **PostgreSQL**. Tracks real-time machine telemetry, automates alert generation, and manages maintenance workflows across a factory floor.

## Features

- **Authentication & RBAC** — JWT-based auth with role-based access control (Admin, Engineer, Technician, Production Manager)
- **Machine management** — CRUD operations for factory equipment
- **Real-time telemetry** — Live sensor data (temperature, vibration, pressure) with a simulator for demo purposes
- **Automated alerting** — Rule-based threshold detection that flags machines as `WARNING` or `CRITICAL`
- **Maintenance workflow** — Kanban-style job tracking (`OPEN` → `ASSIGNED` → `IN_PROGRESS` → `COMPLETED`)
- **Analytics** — MTBF/MTTR calculations and alert-frequency reporting via raw SQL aggregation
- **Dockerized** — Full stack (backend, frontend, database) runs with a single `docker compose up`

## Tech stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3, TypeScript, Pinia, Vue Router, Tailwind CSS, Chart.js |
| Backend | Go, Gin, JWT, bcrypt |
| Database | PostgreSQL |
| DevOps | Docker, Docker Compose |

## Architecture


## Screenshots

![Dashboard](screenshots/dashboard.png)
![Machine Detail](screenshots/machine-detail.png)
![Maintenance Board](screenshots/maintenance.png)

## Running locally

```bash
git clone https://github.com/tumhara-username/factorypulse.git
cd factorypulse
docker compose up --build
```

Then visit `http://localhost:5173`.

To seed live sensor data, run the simulator separately:

```bash
cd backend
go run cmd/simulator/main.go
```

## What I learned / engineering decisions

- Designed a normalized relational schema (machines → sensors → readings/alerts/maintenance) instead of a flat structure
- Implemented rule-based anomaly detection as a foundation for future ML-based predictive maintenance
- Used Go's `internal/` package convention to enforce clean separation between HTTP handlers, business logic, and data access
- Containerized with multi-stage Docker builds to keep production images minimal

## Roadmap

- [ ] WebSocket-based live updates (currently polling)
- [ ] Spare parts inventory tracking
- [ ] ML-based predictive maintenance
- [ ] CI/CD pipeline with GitHub Actions