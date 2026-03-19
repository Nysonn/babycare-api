# BabyCare API

A RESTful backend API for the BabyCare platform, connecting parents with trusted babysitters. Built with Go and the Gin framework.

---

## Tech Stack

- **Go 1.25** with the Gin web framework
- **PostgreSQL** via Neon (managed cloud Postgres)
- **Redis** for caching and session management
- **Docker + Docker Compose** for containerised development
- **Air** for hot reload during development
- **sqlc** for type-safe SQL query generation
- **Clerk** for authentication and user management
- **GetStream** for in-app messaging
- **Cloudinary** for media storage
- **SendGrid** for transactional email

---

## Prerequisites

- Go 1.25 or later
- Docker
- Docker Compose

---

## Getting Started

### 1. Configure environment variables

Copy the example file and fill in your own values:

```bash
cp .env.example .env
```

Edit `.env` with your actual credentials for Neon, Cloudinary, SendGrid, Clerk, and GetStream.

### 2. Start the development server

```bash
make dev
```

This builds the Docker image and starts the API with Air hot reload and a Redis instance. The API will be available at `http://localhost:8080`.

---

## Available Make Commands

| Command      | Description                                              |
|--------------|----------------------------------------------------------|
| `make dev`   | Build and start all services with hot reload             |
| `make down`  | Stop all running services                                |
| `make logs`  | Tail live logs from the API service                      |
| `make clean` | Stop services and remove volumes and orphan containers   |

---

## API

**Base URL:** `http://localhost:8080/api/v1`

### Health Check

```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "babycare-api"
}
```

---

## Services (docker-compose)

| Service | Description                                                              |
|---------|--------------------------------------------------------------------------|
| `api`   | Go + Gin application running with Air hot reload on port 8080            |
| `redis` | Redis 7 (Alpine) used for caching and session storage on port 6379       |

> PostgreSQL is not included in docker-compose. The application connects to a Neon cloud Postgres instance defined by `DATABASE_URL` in your `.env` file.

---

## API Overview

**Base URL:** `http://localhost:8080/api/v1`

| Method   | Endpoint                                  | Auth Required | Role                       |
|----------|-------------------------------------------|---------------|----------------------------|
| `GET`    | `/health`                                 | No            | —                          |
| `POST`   | `/api/v1/auth/register/parent`            | No            | —                          |
| `POST`   | `/api/v1/auth/register/babysitter`        | No            | —                          |
| `POST`   | `/api/v1/auth/login`                      | No            | —                          |
| `POST`   | `/api/v1/auth/logout`                     | Yes           | Any                        |
| `GET`    | `/api/v1/admin/users`                     | Yes           | admin                      |
| `GET`    | `/api/v1/admin/users/:id`                 | Yes           | admin                      |
| `PUT`    | `/api/v1/admin/babysitters/:id/approve`   | Yes           | admin                      |
| `PUT`    | `/api/v1/admin/users/:id/suspend`         | Yes           | admin                      |
| `DELETE` | `/api/v1/admin/users/:id`                 | Yes           | admin                      |
| `POST`   | `/api/v1/admin/create`                    | Yes           | admin                      |
| `GET`    | `/api/v1/admin/activity`                  | Yes           | admin                      |
| `GET`    | `/api/v1/babysitters`                     | No            | —                          |
| `GET`    | `/api/v1/babysitters/:id`                 | Yes           | parent, babysitter         |
| `PUT`    | `/api/v1/babysitters/profile`             | Yes           | babysitter                 |
| `GET`    | `/api/v1/babysitters/profile/views`       | Yes           | babysitter                 |
| `GET`    | `/api/v1/parents/profile`                 | Yes           | parent                     |
| `PUT`    | `/api/v1/parents/profile`                 | Yes           | parent                     |
| `POST`   | `/api/v1/conversations`                   | Yes           | parent                     |
| `GET`    | `/api/v1/conversations`                   | Yes           | parent, babysitter         |
| `POST`   | `/api/v1/conversations/:id/messages`      | Yes           | parent, babysitter         |
| `GET`    | `/api/v1/conversations/:id/messages`      | Yes           | parent, babysitter         |

A Postman collection covering all endpoints is available in [`babycare.json`](babycare.json). Import it into Postman and set the `base_url` collection variable. The login requests auto-populate `admin_token`, `parent_token`, and `babysitter_token` via test scripts.

---

## Deployment (Render)

- Create a new Web Service on Render
- Set the environment to Docker
- Set the Dockerfile path to: `Dockerfile`
- Set the Docker target to: `production`
- Add all environment variables from `.env.example` in the Render dashboard under **Environment Variables**. Do not use the `.env` file on Render.
- Before the first deploy, run migrations manually from your local machine against the production Neon database:
  ```bash
  make migrate-up
  ```
- The Neon database is external so migrations run the same way locally and in production.
- Add a Render Redis instance and copy the Redis URL into the `REDIS_URL` environment variable in the Render dashboard.

---

## Third Party Services

| Service       | Purpose                                       | Env Keys                                                  |
|---------------|-----------------------------------------------|-----------------------------------------------------------|
| **Neon**      | Managed PostgreSQL database                   | `DATABASE_URL`                                            |
| **Clerk**     | User authentication and session management    | `CLERK_SECRET_KEY`                                        |
| **GetStream** | Real-time in-app messaging (Stream Chat)      | `STREAM_API_KEY`, `STREAM_API_SECRET`                     |
| **SendGrid**  | Transactional email (approval, notifications) | `SENDGRID_API_KEY`                                        |
| **Cloudinary**| Media storage for profile pictures and documents | `CLOUDINARY_CLOUD_NAME`, `CLOUDINARY_API_KEY`, `CLOUDINARY_API_SECRET` |
| **Redis**     | Response caching (babysitter profiles, messages) | `REDIS_URL`                                            |
