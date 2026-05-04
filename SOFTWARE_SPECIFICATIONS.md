# BabyCare Platform — Software Specifications

| | |
|---|---|
| **Document Version** | 1.0 |
| **Date** | May 2026 |
| **Platform Version** | API v1.1.0 |

---

## Table of Contents

1. [Overview](#1-overview)
2. [System Architecture](#2-system-architecture)
3. [User Roles](#3-user-roles)
4. [Mobile Application (Flutter)](#4-mobile-application-flutter)
5. [Admin Panel (Web)](#5-admin-panel-web)
6. [Backend API](#6-backend-api)
7. [Data Models](#7-data-models)
8. [API Endpoints](#8-api-endpoints)
9. [Third-Party Services](#9-third-party-services)
10. [Infrastructure & Deployment](#10-infrastructure--deployment)
11. [Security](#11-security)
12. [Non-Functional Requirements](#12-non-functional-requirements)

---

## 1. Overview

BabyCare is a three-tier platform that connects parents with trusted, verified babysitters. The system consists of:

- A **Flutter mobile application** used by parents and babysitters on iOS and Android.
- A **React web admin panel** used by platform administrators to manage users, approve babysitters, and monitor activity.
- A **Go REST API** that powers both frontends and enforces all business logic, authentication, and data access.

The platform operates a two-sided marketplace model: parents discover and contact babysitters; babysitters manage their own profiles and availability; administrators review and approve babysitter accounts before they are visible to parents.

---

## 2. System Architecture

The platform follows a client–server architecture with a single versioned REST API (`/api/v1`) serving both the mobile app and the admin panel.

```
┌──────────────────────┐     ┌──────────────────────┐
│  Flutter Mobile App  │     │  React Admin Panel   │
│  (iOS & Android)     │     │  (Web – Firebase)    │
└────────┬─────────────┘     └──────────┬───────────┘
         │  HTTPS / REST                │  HTTPS / REST
         └──────────────┬───────────────┘
                        │
              ┌─────────▼──────────┐
              │  BabyCare REST API │
              │  (Go + Gin)        │
              │  Render (cloud)    │
              └────┬───────┬───────┘
                   │       │
          ┌────────▼─┐  ┌──▼──────┐
          │PostgreSQL│  │  Redis  │
          │  (Neon)  │  │ (cache) │
          └──────────┘  └─────────┘
```

**External Services:**

| Service | Purpose |
|---|---|
| Clerk | JWT token generation and verification |
| GetStream | In-app real-time messaging channels |
| Cloudinary | Profile picture and document file storage |
| SendGrid | Transactional email (approval notifications) |

---

## 3. User Roles

| Role | Description |
|---|---|
| **admin** | Platform administrator. Has full access to all admin endpoints. Cannot be seen by other roles. |
| **parent** | A user looking for a babysitter. Can browse, save babysitters, manage their own profile, and start conversations. |
| **babysitter** | A user offering babysitting services. Must be approved by an admin before becoming visible. Can manage their profile, set availability, and respond to conversations. |

---

## 4. Mobile Application (Flutter)

### 4.1 Technology Stack

| Component | Technology |
|---|---|
| Framework | Flutter (Dart) |
| State Management | Provider pattern |
| HTTP Client | Custom `ApiClient` wrapper around the `http` package |
| Secure Storage | `flutter_secure_storage` — stores JWT tokens |
| File Handling | `file_picker` — document and image uploads |
| Fonts | Google Fonts |
| URL Handling | `url_launcher` |

### 4.2 Screens

#### Gateway
- **Gateway Screen** — Entry point. Detects whether the user is logged in and routes them to the appropriate experience (parent or babysitter).

#### Parent Screens
| Screen | Description |
|---|---|
| Parent Login | Email/password login for existing parent accounts |
| Parent Account Creation | Registration form for new parent accounts |
| Parent Discover | Browse and search available, approved babysitters |
| Parent Profile – Sitter View | View a babysitter's full public profile |
| Parent Messages | Conversations inbox and messaging interface |
| Parent Account | Manage own profile, settings |

#### Babysitter Screens
| Screen | Description |
|---|---|
| Sitter Login | Email/password login for existing babysitter accounts |
| Sitter Registration Step 1 | Personal details (name, contact, location, gender) |
| Sitter Registration Step 2 | Professional details (rate, availability, languages, payment) |
| Sitter Registration Step 3 | Document uploads (National ID, LCI letter, CV, profile picture) |
| Sitter Dashboard | Overview of profile status, profile view stats, availability toggle |
| Sitter Profile – Parent View | Preview of how a parent sees the babysitter's profile |
| Sitter Messages | Conversations inbox and messaging interface |
| Sitter Account | Manage own profile and settings |

### 4.3 Providers

| Provider | Responsibilities |
|---|---|
| `AuthProvider` | Login, logout, token persistence, session restoration |
| `BabysitterDashboardProvider` | Profile data, profile view analytics, work status toggle |
| `ConversationsProvider` | Conversation list, message threads, sending messages |
| `ParentProvider` | Babysitter discovery list, saved babysitters, parent profile |

### 4.4 Services

| Service | Responsibilities |
|---|---|
| `ApiClient` | Base HTTP client; attaches Bearer token headers; handles errors |
| `AuthService` | Login and registration API calls |
| `BabysitterService` | Babysitter profile CRUD, work status, profile views |
| `ParentService` | Parent profile management, saved babysitters |
| `ConversationService` | Start conversation, list conversations, send and list messages |
| `SecureStorageService` | Read/write JWT token from device secure storage |

### 4.5 Business Rules

- On first launch, the gateway screen checks secure storage for a saved token and auto-navigates if valid.
- Babysitters must complete a three-step registration including document uploads before their account can be reviewed by an admin.
- Babysitters cannot log in until their account has been approved by an admin.
- Parents can bookmark babysitters using the save feature.
- Only parents can initiate conversations with babysitters.

---

## 5. Admin Panel (Web)

### 5.1 Technology Stack

| Component | Technology |
|---|---|
| Framework | React 18 with TypeScript |
| Build Tool | Vite |
| Styling | Tailwind CSS v3 (custom pink theme) |
| Font | Andika (Google Fonts) |
| Client State | Redux Toolkit |
| Server State | TanStack Query (React Query v5) |
| HTTP Client | Axios with request/response interceptors |
| Routing | React Router v6 |
| Notifications | React Hot Toast |
| Icons | Lucide React |
| Deployment | Firebase Hosting |

### 5.2 Pages

| Route | Page | Description |
|---|---|---|
| `/login` | Login | Email/password login restricted to admin accounts |
| `/` | Dashboard | Stat cards (total users, babysitters, parents, pending approvals) and recent activity table |
| `/users` | Users | Searchable, filterable table of all users with infinite scroll; suspend and delete actions |
| `/users/:id` | User Detail | Full user profile with document previews; approve, suspend, and delete actions |
| `/approvals` | Approvals | Card grid of all babysitter accounts awaiting review; Review Profile, Suspend, and Delete actions |
| `/activity` | Activity | Message activity report with stat cards, activity-level filter, and infinite scroll table |
| `/create-admin` | Create Admin | Form to provision a new admin account with client-side validation |

### 5.3 State Management

- **Redux slices:**
  - `authSlice` — stores the admin's JWT token, user object, and authentication status.
  - `uiSlice` — manages UI state including sidebar collapsed state and session-expired modal visibility.
- **TanStack Query** handles all server-side data fetching, caching, and invalidation for user lists, activity reports, and approvals.

### 5.4 Access Control

- All routes except `/login` are protected. Unauthenticated users are redirected to `/login`.
- Only tokens issued to users with the `admin` role are accepted; non-admin tokens receive `403 Forbidden` from the API.
- A `SessionExpiredModal` is displayed when the API returns `401`, prompting the admin to log in again.

### 5.5 Key Workflows

**Babysitter Approval Flow:**
1. Admin navigates to `/approvals`.
2. Admin clicks "Review Profile" to open the full profile with documents.
3. Admin approves or suspends the babysitter account.
4. On approval, the API sends a notification email to the babysitter via SendGrid.

**User Management:**
- Admins can view all non-admin users (parents and babysitters).
- Accounts can be suspended (reversible) or deleted (soft-delete; data is retained but the account is deactivated).

---

## 6. Backend API

### 6.1 Technology Stack

| Component | Technology |
|---|---|
| Language | Go 1.25 |
| Web Framework | Gin |
| Database | PostgreSQL (Neon managed cloud) |
| Database Migrations | Goose |
| Type-Safe SQL | sqlc v1.27.0 |
| Caching | Redis 7 |
| Containerisation | Docker + Docker Compose |
| Hot Reload | Air |

### 6.2 Project Structure

```
cmd/server/         Application entry point
internal/
  config/           Environment variable loading
  database/         Database connection, migration runner, seed data
  db/               sqlc-generated query functions and models
  handlers/         HTTP handler packages (admin, auth, babysitter, messaging, parent)
  middleware/        Auth and role enforcement middleware
  models/           Shared request/response structs
  router/           Gin router setup and route registration
  services/         External service clients (Clerk, Cloudinary, SendGrid, GetStream, Redis)
db/
  migrations/       Goose SQL migration files
  queries/          Raw SQL query files consumed by sqlc
sqlc/               sqlc configuration
```

### 6.3 Middleware

| Middleware | Description |
|---|---|
| `RequireAuth` | Validates the Bearer JWT using Clerk; sets `clerk_user_id` in context |
| `RequireRole` | Verifies the authenticated user holds one of the required roles; returns `403` otherwise |
| CORS | Configured to allow the admin panel and mobile app origins |
| `gin.Recovery` | Catches panics and returns HTTP 500 |
| `gin.Logger` | Logs all requests and responses |

### 6.4 Authentication Flow

1. Client sends email and password to `POST /api/v1/auth/login`.
2. API looks up the user by email, verifies the bcrypt password hash.
3. For babysitter accounts, the API checks that `is_approved = true` before allowing login.
4. For suspended accounts, login is rejected with `403`.
5. On success, Clerk generates a JWT valid for **90 days** and the full user object is returned.
6. All subsequent requests include the token as `Authorization: Bearer <token>`.
7. Logout is stateless — the client discards the token.

---

## 7. Data Models

### Users

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `full_name` | VARCHAR(255) | Required |
| `email` | VARCHAR(255) | Unique, required |
| `phone` | VARCHAR(20) | Optional |
| `role` | ENUM | `admin`, `babysitter`, `parent` |
| `status` | ENUM | `active`, `suspended`, `deleted`; default `active` |
| `clerk_user_id` | VARCHAR(255) | Unique; linked to Clerk identity |
| `password_hash` | VARCHAR(255) | bcrypt hash |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |
| `deleted_at` | TIMESTAMPTZ | Soft-delete timestamp |

### Babysitter Profiles

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `user_id` | UUID | Foreign key → users |
| `location` | TEXT | |
| `national_id_url` | TEXT | Cloudinary URL |
| `lci_letter_url` | TEXT | Cloudinary URL |
| `cv_url` | TEXT | Cloudinary URL |
| `profile_picture_url` | TEXT | Cloudinary URL |
| `languages` | TEXT[] | Array of language strings |
| `days_per_week` | INTEGER | |
| `hours_per_day` | INTEGER | |
| `rate_type` | ENUM | `hourly`, `daily`, `weekly`, `monthly` |
| `rate_amount` | TEXT | |
| `payment_method` | TEXT | e.g. Mobile Money |
| `is_approved` | BOOLEAN | Default `false`; set by admin |
| `gender` | VARCHAR(10) | |
| `availability` | TEXT[] | e.g. weekdays, weekends |
| `currency` | VARCHAR(10) | Default `UGX` |
| `is_available` | BOOLEAN | Default `true`; toggled by babysitter |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### Parent Profiles

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `user_id` | UUID | Foreign key → users |
| `location` | TEXT | |
| `primary_location` | TEXT | |
| `occupation` | TEXT | |
| `preferred_hours` | TEXT | |
| `profile_picture_url` | TEXT | Cloudinary URL |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### Conversations

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `parent_id` | UUID | Foreign key → users |
| `babysitter_id` | UUID | Foreign key → users |
| `stream_channel_id` | TEXT | GetStream channel identifier |
| `is_locked` | BOOLEAN | Prevents further messages when true |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### Messages

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `conversation_id` | UUID | Foreign key → conversations |
| `sender_id` | UUID | Foreign key → users |
| `content` | TEXT | |
| `is_read` | BOOLEAN | Default `false` |
| `sent_at` | TIMESTAMPTZ | |

### Profile Views

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `babysitter_id` | UUID | Foreign key → users |
| `parent_id` | UUID | Foreign key → users |
| `viewed_at` | TIMESTAMPTZ | |

### Saved Babysitters

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key |
| `parent_id` | UUID | Foreign key → users |
| `babysitter_id` | UUID | Foreign key → users |
| `saved_at` | TIMESTAMPTZ | |

---

## 8. API Endpoints

**Base URL:** `https://babycare-api-88gn.onrender.com/api/v1`

### Health

| Method | Endpoint | Auth | Roles |
|---|---|---|---|
| GET | `/health` | No | Any |

### Authentication

| Method | Endpoint | Auth | Roles | Description |
|---|---|---|---|---|
| POST | `/api/v1/auth/register/parent` | No | — | Register a new parent account |
| POST | `/api/v1/auth/register/babysitter` | No | — | Register a new babysitter account with document uploads |
| POST | `/api/v1/auth/login` | No | — | Authenticate and receive a JWT |
| POST | `/api/v1/auth/logout` | Yes | Any | Stateless logout (client discards token) |

### Admin

| Method | Endpoint | Auth | Roles | Description |
|---|---|---|---|---|
| GET | `/api/v1/admin/users` | Yes | admin | List all non-admin, non-deleted users |
| GET | `/api/v1/admin/users/:id` | Yes | admin | Get a user with their role-specific profile |
| PUT | `/api/v1/admin/babysitters/:id/approve` | Yes | admin | Approve a babysitter account; triggers approval email |
| PUT | `/api/v1/admin/users/:id/suspend` | Yes | admin | Suspend a user account |
| DELETE | `/api/v1/admin/users/:id` | Yes | admin | Soft-delete a user account |
| POST | `/api/v1/admin/create` | Yes | admin | Create a new admin account |
| GET | `/api/v1/admin/activity` | Yes | admin | Get all users with a 30-day message activity label |

### Babysitters

| Method | Endpoint | Auth | Roles | Description |
|---|---|---|---|---|
| GET | `/api/v1/babysitters` | No | — | List approved, active, available babysitters |
| GET | `/api/v1/babysitters/:id` | Yes | parent, babysitter | View a babysitter's full public profile |
| GET | `/api/v1/babysitters/profile` | Yes | babysitter | Get own profile |
| PUT | `/api/v1/babysitters/profile` | Yes | babysitter | Update own profile and profile picture |
| GET | `/api/v1/babysitters/profile/views` | Yes | babysitter | Get total profile view count |
| GET | `/api/v1/babysitters/profile/weekly-views` | Yes | babysitter | Get profile views grouped by week |
| PUT | `/api/v1/babysitters/work-status` | Yes | babysitter | Toggle availability (available / unavailable) |

### Parents

| Method | Endpoint | Auth | Roles | Description |
|---|---|---|---|---|
| GET | `/api/v1/parents/profile` | Yes | parent | Get own profile |
| PUT | `/api/v1/parents/profile` | Yes | parent | Update own profile including profile picture |
| POST | `/api/v1/parents/saved-babysitters` | Yes | parent | Save a babysitter to bookmarks |
| DELETE | `/api/v1/parents/saved-babysitters/:babysitter_id` | Yes | parent | Remove a babysitter from bookmarks |
| GET | `/api/v1/parents/saved-babysitters` | Yes | parent | List all saved babysitters |
| GET | `/api/v1/parents/:id` | Yes | babysitter | View a parent's public profile |

### Messaging

| Method | Endpoint | Auth | Roles | Description |
|---|---|---|---|---|
| POST | `/api/v1/conversations` | Yes | parent | Start a new conversation with a babysitter |
| GET | `/api/v1/conversations` | Yes | parent, babysitter | List all conversations with preview text |
| POST | `/api/v1/conversations/:id/messages` | Yes | parent, babysitter | Send a message |
| GET | `/api/v1/conversations/:id/messages` | Yes | parent, babysitter | List all messages in a conversation |

---

## 9. Third-Party Services

### Clerk — Authentication

Clerk is used for token issuance and verification. On login, the API calls Clerk to generate a signed JWT. On every protected request, `RequireAuth` middleware calls Clerk to verify the token signature and extract the subject. Tokens are valid for 90 days.

### GetStream — Messaging

GetStream provides the real-time channel infrastructure for in-app messaging. When a parent starts a conversation, the API creates or retrieves a GetStream channel and stores the `stream_channel_id` in the `conversations` table. The Flutter app connects directly to GetStream for real-time message delivery.

### Cloudinary — Media Storage

Profile pictures and babysitter verification documents (National ID, LCI letter, CV) are uploaded via the API and stored in Cloudinary. The API returns the Cloudinary URL, which is then stored in the database and served to clients.

### SendGrid — Transactional Email

SendGrid is used to send transactional emails. Currently, approval notification emails are dispatched to babysitters when an admin approves their account.

### Neon — Managed PostgreSQL

Neon provides the managed cloud PostgreSQL instance. The API connects using a `DATABASE_URL` connection string. Schema is version-controlled using Goose migrations.

### Redis — Caching

Redis is used for response caching. Babysitter list responses are cached to reduce database load. The cache is invalidated when a babysitter's profile or work status is updated. Redis is optional at runtime; if unavailable, the API continues to function without caching.

---

## 10. Infrastructure & Deployment

### Backend API

- **Platform:** Render (managed cloud hosting)
- **Production URL:** `https://babycare-api-88gn.onrender.com`
- **Containerised:** Dockerfile provided; Docker Compose used for local development
- **Migrations:** Run automatically on startup via the `migrations.go` runner
- **Hot Reload (dev):** Air watches the source tree and rebuilds on file changes

### Admin Panel

- **Platform:** Firebase Hosting
- **Build Tool:** Vite; output to `dist/`
- **Environment:** `VITE_API_BASE_URL` points to the backend API URL

### Mobile App

- **Framework:** Flutter
- **Targets:** iOS and Android
- **Backend URL:** `https://babycare-api-0prm.onrender.com` (configured in `main.dart`)

### Database

- **Provider:** Neon (cloud PostgreSQL)
- **Migrations:** 8 migration files managed by Goose (`db/migrations/`)

### Local Development (Backend)

```
make dev     → builds Docker image, starts API (port 8080) + Redis (port 6379) with hot reload
make down    → stops all services
make logs    → tail API logs
make clean   → stop services, remove volumes
```

---

## 11. Security

### Authentication & Authorization

- All non-public endpoints require a valid Bearer JWT issued by Clerk.
- Role enforcement is applied at the router level via the `RequireRole` middleware. Requests from users without the required role receive `403 Forbidden`.
- Admin endpoints are entirely inaccessible to parent and babysitter tokens.

### Password Storage

- Passwords are hashed using bcrypt before being stored. Plain-text passwords are never persisted.

### Soft Deletes

- User account deletion is a soft-delete: the `deleted_at` timestamp is set and the account is excluded from all listings, but data is retained for audit purposes.

### Input Validation

- All request bodies are validated on binding. Invalid or malformed requests return `400 Bad Request` before reaching handler logic.

### Document Access

- Verification documents (National ID, LCI letter, CV) are stored in Cloudinary and are only surfaced to admin users through the admin panel.

### Transport Security

- All production traffic runs over HTTPS. The API is served via Render's TLS termination.

### CORS

- CORS is configured on the Gin router. Only allowed origins, methods, and headers are accepted.

---

## 12. Non-Functional Requirements

### Performance

- Redis caching reduces repeated database reads for the babysitter listing, which is the highest-traffic public endpoint.
- The API includes `gin.Recovery()` middleware to prevent single panics from taking down the server.

### Scalability

- The stateless API design allows horizontal scaling on Render.
- PostgreSQL (Neon) and Redis are managed externally, separating compute and data concerns.

### Reliability

- The migration runner on startup ensures the database schema is always up to date with the deployed code.
- Redis is treated as optional; the API degrades gracefully if the cache is unavailable.

### Maintainability

- SQL queries are managed via sqlc, providing type-safe, auto-generated Go query functions from plain SQL files.
- The project is structured by domain (handlers per role, services per integration), making individual concerns easy to locate and modify.
- Database schema changes are versioned through Goose migration files.

### Observability

- All requests and responses are logged by `gin.Logger()`.
- Handler-level errors are logged using Go's standard `log` package with contextual prefixes.

---

*End of Document*
