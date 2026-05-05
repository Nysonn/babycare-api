# BabyCare API — New Features Integration Guide

| | |
|---|---|
| **Base URL** | https://babycare-api-88gn.onrender.com |
| **Prepared for** | Flutter Frontend Team |
| **Date** | May 2026 |

This document covers the three new features added to the API. Integrate each section independently.

---

## Table of Contents

1. [Forgot Password](#1-forgot-password)
2. [Report a User](#2-report-a-user)

---

## 1. Forgot Password

Allows a user who has forgotten their password to trigger a reset email. Clerk handles sending the email and the reset link — no extra setup is needed on the Flutter side beyond calling this endpoint.

### Endpoint

```
POST /api/v1/auth/forgot-password
```

**Authentication:** Not required (public endpoint).

### Request Body

```json
{
  "email": "user@example.com"
}
```

| Field   | Type   | Required | Description                    |
|---------|--------|----------|--------------------------------|
| `email` | string | Yes      | The email address of the account. Must be a valid email format. |

### Success Response

**200 OK** — always returned regardless of whether the email exists, to prevent account enumeration.

```json
{
  "message": "If an account with that email exists, a reset link has been sent."
}
```

### Error Responses

| Status | Body | When |
|--------|------|------|
| `400 Bad Request` | `{ "error": "valid email is required" }` | Missing or malformed email field. |
| `500 Internal Server Error` | `{ "error": "internal server error" }` | Unexpected server fault. |

### Flutter Integration Notes

- Always show the same success message to the user regardless of the response — do **not** tell them whether the account exists.
- Show a loading indicator while the request is in flight, then navigate the user to a confirmation screen that says something like *"Check your inbox for a password reset link."*
- The reset link opened from the email is handled outside the app (Clerk's hosted page). You do not need to build a reset password screen in Flutter.

### Example (Dart)

```dart
Future<void> forgotPassword(String email) async {
  final response = await http.post(
    Uri.parse('$baseUrl/api/v1/auth/forgot-password'),
    headers: {'Content-Type': 'application/json'},
    body: jsonEncode({'email': email}),
  );

  if (response.statusCode == 200) {
    // Always show the same neutral confirmation message
    showConfirmationScreen();
  } else {
    final body = jsonDecode(response.body);
    throw Exception(body['error'] ?? 'Something went wrong');
  }
}
```

---

## 2. Report a User

Allows an authenticated parent or babysitter to report another user for spam, harassment, inappropriate behaviour, or other concerns. The report is stored and reviewed by the admin team.

### Endpoint

```
POST /api/v1/reports
```

**Authentication:** Required — Bearer token.  
**Roles:** `parent`, `babysitter`.

### Request Body

```json
{
  "reported_user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "report_type": "spam",
  "description": "This user keeps sending irrelevant promotional messages."
}
```

| Field              | Type   | Required | Description |
|--------------------|--------|----------|-------------|
| `reported_user_id` | UUID   | Yes      | The ID of the user being reported. |
| `report_type`      | string | Yes      | One of: `spam`, `harassment`, `inappropriate`, `other`. |
| `description`      | string | No       | Optional free-text details about the report (recommended). |

#### Valid `report_type` values

| Value           | Use when |
|-----------------|----------|
| `spam`          | Unsolicited or repeated irrelevant messages. |
| `harassment`    | Threatening, abusive, or bullying behaviour. |
| `inappropriate` | Offensive content, profile pictures, or language. |
| `other`         | Anything that does not fit the above categories. |

### Success Response

**201 Created**

```json
{
  "message": "Report submitted successfully",
  "report_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### Error Responses

| Status | Body | When |
|--------|------|------|
| `400 Bad Request` | `{ "error": "reported_user_id is required" }` | Missing required field. |
| `400 Bad Request` | `{ "error": "invalid reported_user_id" }` | The UUID is not valid. |
| `400 Bad Request` | `{ "error": "report_type must be one of: spam, harassment, inappropriate, other" }` | Invalid report type sent. |
| `400 Bad Request` | `{ "error": "you cannot report yourself" }` | Reporter and reported user are the same person. |
| `401 Unauthorized` | `{ "error": "unauthorized" }` | Missing or invalid token. |
| `404 Not Found` | `{ "error": "reported user not found" }` | The `reported_user_id` does not match any active user. |
| `500 Internal Server Error` | `{ "error": "internal server error" }` | Unexpected server fault. |

### Flutter Integration Notes

- Place a **"Report"** option in the profile/overflow menu on any babysitter or parent profile screen.
- Show a bottom sheet or dialog letting the user choose the `report_type` from the four options, plus an optional text field for `description`.
- After a successful `201`, show a brief toast/snackbar: *"Report submitted. Our team will review it shortly."*
- You do **not** need to handle the `report_id` in the response — it is provided for debugging purposes only.
- The reported user is **not** notified. No changes are made to their account until an admin reviews the report.

### Example (Dart)

```dart
Future<void> reportUser({
  required String reportedUserId,
  required String reportType,
  String? description,
}) async {
  final token = await getStoredToken(); // your token retrieval

  final response = await http.post(
    Uri.parse('$baseUrl/api/v1/reports'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode({
      'reported_user_id': reportedUserId,
      'report_type': reportType,
      if (description != null && description.isNotEmpty)
        'description': description,
    }),
  );

  if (response.statusCode == 201) {
    showSuccessToast('Report submitted. Our team will review it shortly.');
  } else {
    final body = jsonDecode(response.body);
    throw Exception(body['error'] ?? 'Failed to submit report');
  }
}
```

### Suggested UI Flow

```
Profile Screen
  └── ⋮ Menu / Report button
        └── Bottom Sheet: "Why are you reporting this user?"
              ├── Spam
              ├── Harassment
              ├── Inappropriate content
              └── Other
                    └── (Optional) Text field: "Add details (optional)"
                          └── [Submit Report] button
                                └── Toast: "Report submitted"
                                      └── Dismiss / Close
```

---

## Headers Reference

| Header          | Value                        | When required |
|-----------------|------------------------------|---------------|
| `Content-Type`  | `application/json`           | All POST/PUT requests |
| `Authorization` | `Bearer <token>`             | All authenticated endpoints |
