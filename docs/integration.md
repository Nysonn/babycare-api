# BabyCare API — Developer Documentation

| | |
|---|---|
| **Base URL** | https://babycare-api-88gn.onrender.com |
| **Version** | 1.1.0 |
| **Prepared for** | Frontend Development Team |
| **Date** | March 2026 |

---

## Overview

The BabyCare API is a RESTful backend service that powers the BabyCare mobile and web applications. It connects parents with verified babysitters through a secure, role-based platform. All endpoints return JSON. Authentication uses Bearer JWT tokens issued on login.

---

## Authentication

Protected endpoints require an `Authorization` header with a Bearer token obtained from the login endpoint.

Tokens are valid for **90 days**. Include the header as:

```
Authorization: Bearer <your_token>
```

---

## User Roles

| Role | Description | Access |
|---|---|---|
| `admin` | Platform administrator | Full access to all admin endpoints |
| `parent` | Parent looking for a babysitter | Browse babysitters, save, message, manage profile |
| `babysitter` | Babysitter offering services | Manage profile, set work status, respond to messages |

---

## Standard Error Response

```json
{ "error": "descriptive error message" }
```

---

## All Endpoints at a Glance

| Method | Endpoint | Auth | Role |
|---|---|---|---|
| GET | /health | No | Any |
| POST | /api/v1/auth/register/parent | No | — |
| POST | /api/v1/auth/register/babysitter | No | — |
| POST | /api/v1/auth/login | No | — |
| POST | /api/v1/auth/logout | Yes | Any |
| GET | /api/v1/admin/users | Yes | Admin |
| GET | /api/v1/admin/users/:id | Yes | Admin |
| PUT | /api/v1/admin/babysitters/:id/approve | Yes | Admin |
| PUT | /api/v1/admin/users/:id/suspend | Yes | Admin |
| DELETE | /api/v1/admin/users/:id | Yes | Admin |
| POST | /api/v1/admin/create | Yes | Admin |
| GET | /api/v1/admin/activity | Yes | Admin |
| GET | /api/v1/babysitters | No | — |
| GET | /api/v1/babysitters/:id | Yes | Parent, Babysitter |
| PUT | /api/v1/babysitters/profile | Yes | Babysitter |
| GET | /api/v1/babysitters/profile/views | Yes | Babysitter |
| GET | /api/v1/babysitters/profile/weekly-views | Yes | Babysitter |
| PUT | /api/v1/babysitters/work-status | Yes | Babysitter |
| GET | /api/v1/parents/profile | Yes | Parent |
| PUT | /api/v1/parents/profile | Yes | Parent |
| POST | /api/v1/parents/saved-babysitters | Yes | Parent |
| DELETE | /api/v1/parents/saved-babysitters/:babysitter_id | Yes | Parent |
| GET | /api/v1/parents/saved-babysitters | Yes | Parent |
| POST | /api/v1/conversations | Yes | Parent |
| GET | /api/v1/conversations | Yes | Parent, Babysitter |
| POST | /api/v1/conversations/:id/messages | Yes | Parent, Babysitter |
| GET | /api/v1/conversations/:id/messages | Yes | Parent, Babysitter |

---

## 1. Health Check

### `GET /health`

Verify the API server is running. No authentication required.

**Response**
```json
{
  "service": "babycare-api",
  "status": "ok"
}
```

---

## 2. Authentication

### `POST /api/v1/auth/register/parent`

Register a new parent account.

**Request Body**
```json
{
  "full_name": "Ochieng Samuel",
  "email": "samuel.ochieng@gmail.com",
  "phone": "+256772445566",
  "location": "Jinja, Uganda",
  "primary_location": "Ntinda, Kampala",
  "occupation": "Project Manager",
  "preferred_hours": "Flexible",
  "password": "OchiengSecure2026!"
}
```

| Field | Required | Notes |
|---|---|---|
| `full_name` | Yes | |
| `email` | Yes | Must be unique |
| `phone` | No | |
| `location` | Yes | General area |
| `primary_location` | No | Specific home/work location |
| `occupation` | Yes | |
| `preferred_hours` | Yes | |
| `password` | Yes | Minimum 8 characters |

**Response** `201 Created`
```json
{
  "id": "6590a01c-aaf6-47f1-b28b-cf75de37e263",
  "full_name": "Ochieng Samuel",
  "email": "samuel.ochieng@gmail.com",
  "phone": "+256772445566",
  "role": "parent",
  "status": "active",
  "created_at": "2026-03-19T18:30:08.210148Z"
}
```

---

### `POST /api/v1/auth/register/babysitter`

Register a new babysitter account. Uses `multipart/form-data`.

**Request Body** — `Content-Type: multipart/form-data`

| Field | Type | Required | Notes |
|---|---|---|---|
| `full_name` | text | Yes | |
| `email` | text | Yes | Must be unique |
| `phone` | text | No | |
| `location` | text | Yes | |
| `languages` | text | Yes | Comma-separated e.g. `English,Luganda` |
| `password` | text | Yes | Minimum 8 characters |
| `gender` | text | Yes | `male` or `female` |
| `availability` | text | No | Comma-separated days e.g. `Mon,Tue,Wed,Fri` |
| `rate_type` | text | No | `hourly`, `daily`, `weekly`, or `monthly` |
| `rate_amount` | text | No | Numeric string e.g. `25000` |
| `currency` | text | No | Defaults to `UGX` |
| `payment_method` | text | No | `Mobile Money`, `Cash`, or `Bank/Visa Card` |
| `national_id` | file | Yes | Image |
| `lci_letter` | file | Yes | PDF |
| `cv` | file | Yes | PDF |
| `profile_picture` | file | Yes | Image |

**Response** `201 Created`
```json
{
  "id": "9c7f5648-0059-46ce-91a1-d7826d5fc937",
  "full_name": "Mary Nakato",
  "email": "marynakato@example.com",
  "phone": "+256700000002",
  "role": "babysitter",
  "status": "active",
  "created_at": "2026-03-19T18:37:38.172614Z"
}
```

> **Note:** Account is inactive until an admin approves it. Login returns `403` until approval.

---

### `POST /api/v1/auth/login`

Login for all user types (admin, parent, babysitter).

**Request Body**
```json
{
  "email": "kato.emma@outlook.com",
  "password": "Agumya2022!"
}
```

**Response**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-06-17T18:40:19.090884234Z",
  "user": {
    "id": "11b14ba7-d5bc-43cd-812b-a8fa8aecdb74",
    "full_name": "Kato Emmanuel",
    "email": "kato.emma@outlook.com",
    "phone": "+256702112233",
    "role": "parent",
    "status": "active",
    "created_at": "2026-03-19T15:49:14.301273Z"
  }
}
```

> **Note:** Store the token securely (Flutter Secure Storage). Token expires in 90 days. Babysitters must be approved before login succeeds.

---

### `POST /api/v1/auth/logout`

Logout the currently authenticated user.

**Auth:** Bearer token required

**Response**
```json
{ "message": "logged out successfully" }
```

> **Note:** Discard the token on the client side after calling this endpoint.

---

## 3. Admin

All admin endpoints require a valid admin Bearer token. Non-admin tokens receive `403 Forbidden`.

### `GET /api/v1/admin/users`

List all parent and babysitter accounts (admin accounts excluded).

**Response**
```json
[
  {
    "id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
    "full_name": "Erina Ahabwe",
    "email": "edrina@gmail.com",
    "phone": "+1987654321",
    "role": "babysitter",
    "status": "active",
    "created_at": "2026-03-19T15:50:26.796082Z"
  }
]
```

---

### `GET /api/v1/admin/users/:id`

Get full details of a single user including their profile.

**Path Parameters:** `:id` — UUID of the user

**Response**
```json
{
  "id": "11b14ba7-d5bc-43cd-812b-a8fa8aecdb74",
  "full_name": "Kato Emmanuel",
  "email": "kato.emma@outlook.com",
  "phone": "+256702112233",
  "role": "parent",
  "status": "active",
  "created_at": "2026-03-19T15:49:14.301273Z",
  "location": "Queens, NY",
  "occupation": "Teacher",
  "preferred_hours": "Weekday mornings"
}
```

---

### `PUT /api/v1/admin/babysitters/:id/approve`

Approve a babysitter account so they can log in and receive messages.

**Path Parameters:** `:id` — UUID of the babysitter user

**Response**
```json
{ "message": "babysitter approved successfully" }
```

> **Note:** An approval email is automatically sent to the babysitter via SendGrid.

---

### `PUT /api/v1/admin/users/:id/suspend`

Suspend a user account and lock all their conversations.

**Path Parameters:** `:id` — UUID of the user to suspend

**Response**
```json
{ "message": "user suspended successfully" }
```

> **Note:** Suspended users cannot send new messages. A suspension email is sent.

---

### `DELETE /api/v1/admin/users/:id`

Soft-delete a user account (data retained, account deactivated).

**Path Parameters:** `:id` — UUID of the user to delete

**Response**
```json
{ "message": "user deleted successfully" }
```

---

### `POST /api/v1/admin/create`

Create a new admin user account.

**Request Body**
```json
{
  "full_name": "Timo Mugumya",
  "email": "timo.mugumya@gmail.com",
  "password": "admin123"
}
```

**Response**
```json
{
  "id": "c83b7a72-4bde-4e73-b344-968d0e317470",
  "full_name": "Timo Mugumya",
  "email": "timo.mugumya@gmail.com",
  "phone": "",
  "role": "admin",
  "status": "active",
  "created_at": "2026-03-19T18:53:34.485535Z"
}
```

---

### `GET /api/v1/admin/activity`

Get activity report for all users based on message volume (last 30 days).

**Response**
```json
[
  {
    "user_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
    "full_name": "Erina Ahabwe",
    "role": "babysitter",
    "activity_label": "Low",
    "message_count": 1
  }
]
```

> **Activity labels:** Low = 0–10 messages, Medium = 11–50 messages, High = 51+ messages.

---

## 4. Babysitters

### `GET /api/v1/babysitters`

List all approved, active, and **available** babysitter profiles.

**Auth:** No — public endpoint

> **Note:** Babysitters who have set their work status to unavailable are excluded from this list.

**Response**
```json
[
  {
    "user_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
    "full_name": "Erina Ahabwe",
    "email": "edrina@gmail.com",
    "phone": "+1987654321",
    "location": "Manhattan, NY",
    "profile_picture_url": "https://res.cloudinary.com/...",
    "languages": ["English", "Spanish", "French"],
    "days_per_week": 5,
    "hours_per_day": 8,
    "rate_type": "hourly",
    "rate_amount": 25000,
    "payment_method": "Mobile Money",
    "is_approved": true,
    "gender": "female",
    "availability": ["Mon", "Tue", "Wed", "Thu", "Fri"],
    "currency": "UGX",
    "is_available": true
  }
]
```

> Results are cached in Redis for 5 minutes for performance.

---

### `GET /api/v1/babysitters/:id`

Get the full profile of a single babysitter. Also records a profile view automatically.

**Auth:** Yes — Parent or Babysitter token

**Path Parameters:** `:id` — UUID of the babysitter

**Response**
```json
{
  "user_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
  "full_name": "Erina Ahabwe",
  "email": "edrina@gmail.com",
  "phone": "+1987654321",
  "location": "Manhattan, NY",
  "profile_picture_url": "https://res.cloudinary.com/...",
  "languages": ["English", "Spanish", "French"],
  "days_per_week": 5,
  "hours_per_day": 8,
  "rate_type": "hourly",
  "rate_amount": 25000,
  "payment_method": "Mobile Money",
  "is_approved": true,
  "gender": "female",
  "availability": ["Mon", "Tue", "Wed", "Thu", "Fri"],
  "currency": "UGX",
  "is_available": true
}
```

> **Profile view recording:** When a **parent** calls this endpoint, their view is **automatically recorded** in the background. Flutter does not need to call a separate endpoint — just calling `GET /babysitters/:id` is enough. The view will appear in the babysitter's profile views list.

> The `is_available` field shows whether the babysitter is currently accepting work. Display this clearly on the profile screen.

---

### `PUT /api/v1/babysitters/profile`

Update the logged-in babysitter's own profile. Uses `multipart/form-data`. All fields are optional — only send what you want to change.

**Auth:** Yes — Babysitter token

**Request Body** — `Content-Type: multipart/form-data`

| Field | Type | Notes |
|---|---|---|
| `location` | text | |
| `languages` | text | Comma-separated |
| `days_per_week` | text | Integer |
| `hours_per_day` | text | Integer |
| `rate_type` | text | `hourly`, `daily`, `weekly`, or `monthly` |
| `rate_amount` | text | Numeric |
| `currency` | text | e.g. `UGX` |
| `payment_method` | text | `Mobile Money`, `Cash`, or `Bank/Visa Card` |
| `gender` | text | `male` or `female` |
| `availability` | text | Comma-separated days e.g. `Mon,Wed,Fri` |
| `profile_picture` | file | Image (optional) |

**Response**
```json
{
  "user_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
  "full_name": "Erina Ahabwe",
  "location": "Manhattan, NY",
  "languages": ["English", "Spanish"],
  "days_per_week": 5,
  "hours_per_day": 8,
  "rate_type": "hourly",
  "rate_amount": 25000,
  "payment_method": "Mobile Money",
  "is_approved": true,
  "gender": "female",
  "availability": ["Mon", "Wed", "Fri"],
  "currency": "UGX",
  "is_available": true
}
```

---

### `PUT /api/v1/babysitters/work-status`

Set the babysitter's work availability status. Unavailable babysitters do not appear in `GET /babysitters`.

**Auth:** Yes — Babysitter token

**Request Body**
```json
{ "is_available": false }
```

| Field | Required | Notes |
|---|---|---|
| `is_available` | Yes | `true` = accepting work, `false` = not available |

**Response**
```json
{ "message": "work status set to unavailable" }
```

> **Flutter integration:** Show an availability toggle on the babysitter's own profile screen. Call this endpoint whenever the babysitter flips the toggle. Parents browsing `GET /babysitters` will immediately (within 5 min cache refresh) stop seeing the unavailable babysitter.

---

### `GET /api/v1/babysitters/profile/views`

Get the list of parents who viewed the logged-in babysitter's profile.

**Auth:** Yes — Babysitter token

**Restricted fields:** `email`, `phone`, `primary_location`, and `preferred_hours` are only returned when `has_messaged` is `true` (the parent has sent at least one message to this babysitter). When `has_messaged` is `false`, display a default avatar — do not show the parent's profile picture.

**Response**
```json
[
  {
    "id": "0305a2b5-286a-4d47-9935-d8551176af6e",
    "parent_id": "6590a01c-aaf6-47f1-b28b-cf75de37e263",
    "parent_name": "Ochieng Samuel",
    "occupation": "Project Manager",
    "viewed_at": "2026-03-19T19:00:15.000394Z",
    "has_messaged": false
  },
  {
    "id": "1a2b3c4d-...",
    "parent_id": "abc123...",
    "parent_name": "Nakato Grace",
    "occupation": "Doctor",
    "viewed_at": "2026-03-20T10:30:00Z",
    "has_messaged": true,
    "email": "grace.nakato@gmail.com",
    "phone": "+256701234567",
    "primary_location": "Kololo, Kampala",
    "preferred_hours": "Weekday mornings"
  }
]
```

> **Babysitters cannot initiate messaging.** Only parents can start a conversation. If a babysitter wants to respond to a viewer who hasn't messaged yet, they must wait for the parent to send the first message.

---

### `GET /api/v1/babysitters/profile/weekly-views`

Get the total number of profile views in the last 7 days for the authenticated babysitter.

**Auth:** Yes — Babysitter token

**Response**
```json
{
  "babysitter_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
  "view_count": 14,
  "period_days": 7
}
```

---

## 5. Parents

### `GET /api/v1/parents/profile`

Get the logged-in parent's profile.

**Auth:** Yes — Parent token

**Response**
```json
{
  "id": "6590a01c-aaf6-47f1-b28b-cf75de37e263",
  "full_name": "Ochieng Samuel",
  "email": "samuel.ochieng@gmail.com",
  "phone": "+256772445566",
  "status": "active",
  "location": "Jinja, Uganda",
  "primary_location": "Ntinda, Kampala",
  "occupation": "Project Manager",
  "preferred_hours": "Flexible",
  "created_at": "2026-03-19T18:30:08Z"
}
```

---

### `PUT /api/v1/parents/profile`

Update the logged-in parent's profile.

**Auth:** Yes — Parent token

**Request Body**
```json
{
  "location": "Ntinda, Kampala",
  "primary_location": "Kololo, Kampala",
  "occupation": "Doctor",
  "preferred_hours": "7am - 4pm weekdays"
}
```

**Response** — same shape as `GET /api/v1/parents/profile`

---

### `POST /api/v1/parents/saved-babysitters`

Save a babysitter to the parent's saved list. If already saved, the request is silently ignored.

**Auth:** Yes — Parent token

**Request Body**
```json
{ "babysitter_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee" }
```

**Response**
```json
{ "message": "babysitter saved" }
```

---

### `DELETE /api/v1/parents/saved-babysitters/:babysitter_id`

Remove a babysitter from the parent's saved list.

**Auth:** Yes — Parent token

**Path Parameters:** `:babysitter_id` — UUID of the babysitter to remove

**Response**
```json
{ "message": "babysitter removed from saved list" }
```

---

### `GET /api/v1/parents/saved-babysitters`

List all babysitters saved by the authenticated parent. Returns the same fields as `GET /babysitters`.

**Auth:** Yes — Parent token

**Response**
```json
[
  {
    "user_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee",
    "full_name": "Erina Ahabwe",
    "email": "edrina@gmail.com",
    "phone": "+1987654321",
    "location": "Manhattan, NY",
    "profile_picture_url": "https://res.cloudinary.com/...",
    "languages": ["English", "Spanish"],
    "days_per_week": 5,
    "hours_per_day": 8,
    "rate_type": "hourly",
    "rate_amount": 25000,
    "payment_method": "Mobile Money",
    "is_approved": true,
    "gender": "female",
    "availability": ["Mon", "Tue", "Wed"],
    "currency": "UGX",
    "is_available": true
  }
]
```

---

## 6. Messaging

Messaging is powered by Stream Chat. Only parents can initiate a conversation. Email notifications are sent via SendGrid on first unread message.

### `POST /api/v1/conversations`

Start a new conversation between a parent and a babysitter.

**Auth:** Yes — Parent token only

**Request Body**
```json
{ "babysitter_id": "b80b2ffd-c78f-43af-a67b-924b2b6746ee" }
```

**Response**
```json
{
  "id": "f23166e3-3558-473c-ac7d-50be6774d4a3",
  "other_user_name": "Erina Ahabwe",
  "is_locked": false,
  "created_at": "2026-03-19T19:07:52.373207Z"
}
```

> If a conversation between the same parent and babysitter already exists, the existing conversation is returned rather than creating a duplicate.

---

### `GET /api/v1/conversations`

List all conversations for the logged-in user.

**Auth:** Yes — Parent or Babysitter token

**Response**
```json
[
  {
    "id": "f23166e3-3558-473c-ac7d-50be6774d4a3",
    "other_user_name": "Erina Ahabwe",
    "is_locked": false,
    "created_at": "2026-03-19T19:07:52.373207Z"
  }
]
```

---

### `POST /api/v1/conversations/:id/messages`

Send a message in a conversation.

**Auth:** Yes — Parent or Babysitter token

**Path Parameters:** `:id` — UUID of the conversation

**Request Body**
```json
{ "content": "Hi! Are you available this Saturday evening?" }
```

**Response**
```json
{
  "id": "ab8542da-8fec-4ea2-8db3-8f7fae4d1d21",
  "conversation_id": "f23166e3-3558-473c-ac7d-50be6774d4a3",
  "sender_id": "6590a01c-aaf6-47f1-b28b-cf75de37e263",
  "content": "Hi! Are you available this Saturday evening?",
  "is_read": false,
  "sent_at": "2026-03-19T19:09:49.207808Z"
}
```

> Returns `403` if the conversation is locked (user was suspended). An email notification is sent to the recipient only on their first unread message.

---

### `GET /api/v1/conversations/:id/messages`

List all messages in a conversation. Also marks messages as read.

**Auth:** Yes — Parent or Babysitter token

**Path Parameters:** `:id` — UUID of the conversation

**Response**
```json
[
  {
    "id": "ab8542da-8fec-4ea2-8db3-8f7fae4d1d21",
    "conversation_id": "f23166e3-3558-473c-ac7d-50be6774d4a3",
    "sender_id": "6590a01c-aaf6-47f1-b28b-cf75de37e263",
    "content": "Hi! Are you available this Saturday evening?",
    "is_read": false,
    "sent_at": "2026-03-19T19:09:49.207808Z"
  }
]
```

> Calling this endpoint automatically marks all messages from the other participant as read.

---

## 7. Error Code Reference

| HTTP Code | Meaning | Common Cause |
|---|---|---|
| 400 | Bad Request | Missing or invalid request body fields |
| 401 | Unauthorised | Missing, invalid or expired Bearer token |
| 403 | Forbidden | Valid token but insufficient role, account suspended, or babysitter not yet approved |
| 404 | Not Found | Resource does not exist or has been soft-deleted |
| 409 | Conflict | Email address already registered |
| 500 | Server Error | Unexpected internal error — check server logs |

---

## 8. Integration Notes for Flutter Team

### Token Storage
Store the JWT token in Flutter Secure Storage immediately after login. Include it in every protected request as `Authorization: Bearer <token>`. Tokens expire after 90 days.

### Babysitter Registration
Use `multipart/form-data`. Required file fields: `national_id`, `lci_letter`, `cv`, `profile_picture`. Required text fields include `gender` (`male` or `female`). The account will be inactive until admin approval.

### Babysitter Login
If a babysitter tries to login before admin approval, the API returns `403`. Display: *"Your account is pending admin approval."*

### Profile View Recording (How it works)
Profile views are **automatic**. Flutter does not need to call a separate endpoint. Simply calling `GET /api/v1/babysitters/:id` while authenticated as a parent is enough — the backend records the view in the background. No POST call is needed.

### Babysitter Availability (Work Status)
- `GET /babysitters` only returns babysitters where `is_available = true`.
- `GET /babysitters/:id` returns the profile regardless of status, and includes `is_available` so you can show a badge ("Not Available" / "Available") on the profile screen.
- The babysitter updates their status via `PUT /babysitters/work-status`.

### Profile Views — Restricted Parent Info
When a babysitter views their profile views list (`GET /babysitters/profile/views`), each entry includes a `has_messaged` boolean:
- `has_messaged: false` → show only name and occupation. Use a **default avatar** (do not attempt to load a profile picture).
- `has_messaged: true` → full parent details are included (email, phone, primary_location, preferred_hours).

Babysitters cannot message parents. Only parents can initiate conversations.

### Saved Babysitters Flow
1. Parent browses or views a babysitter profile.
2. Parent taps "Save" → `POST /parents/saved-babysitters` with `{ "babysitter_id": "..." }`.
3. Parent can view their saved list via `GET /parents/saved-babysitters`.
4. Parent unsaves via `DELETE /parents/saved-babysitters/:babysitter_id`.

### Messaging Flow
Only parents can start conversations. After starting a conversation, use the returned `id` to send and fetch messages. Poll `GET /conversations/:id/messages` to refresh the message list.

### Profile Picture URLs
Profile picture URLs are hosted on Cloudinary CDN and are publicly accessible. Use them directly in Image widgets.

### Caching
The babysitter list is cached for 5 minutes. Individual profiles are cached for 10 minutes. Caches are invalidated automatically after a profile update or work status change.

### Phone Number Format
Phone numbers should include the country code. Uganda numbers follow the format: `+256XXXXXXXXX`.

### Availability Days Format
Days are returned as an array of short day strings: `["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]`. When sending, pass as a comma-separated string: `Mon,Tue,Wed`.

### Currency
Currency defaults to `UGX` if not specified at registration. Always display the `currency` field alongside `rate_amount`.
