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

Password reset must be initiated directly from the Clerk Flutter client SDK. Do not call the BabyCare API for this flow.

Clerk's supported flow is:

```
1. Start a sign-in attempt with strategy: reset_password_email_code
2. Pass the user's email as the identifier
3. Clerk sends the reset code to the user's email
4. Submit the code and the new password with the same strategy
```

### API Status

The old backend endpoint below is deprecated and now returns `410 Gone`:

```
POST /api/v1/auth/forgot-password
```

### Client-Side Clerk Flow

Use the Clerk sign-in resource from your Flutter app to start the reset flow.

#### Step 1: Send reset code

```dart
Future<void> sendPasswordResetCode(String email) async {
  await clerk.client.signIn.create(
    strategy: 'reset_password_email_code',
    identifier: email,
  );
}
```

#### Step 2: Submit code and new password

```dart
Future<void> completePasswordReset({
  required String code,
  required String newPassword,
}) async {
  final result = await clerk.client.signIn.attemptFirstFactor(
    strategy: 'reset_password_email_code',
    code: code,
    password: newPassword,
  );

  if (result.status == 'complete') {
    await clerk.setActive(session: result.createdSessionId);
    return;
  }

  throw Exception('Password reset requires additional Clerk steps.');
}
```

### UI Notes

- Build this as a two-step screen: enter email, then enter reset code and new password.
- Show Clerk's error message when the SDK rejects the request.
- If the result status is `complete`, set the created session as active so the user is signed in immediately after resetting their password.
- If your Clerk Flutter SDK version exposes slightly different method names, keep the same flow and strategy value: `reset_password_email_code`.

### Expected Backend Behavior

If the app still calls `POST /api/v1/auth/forgot-password`, the API will return:

```json
{
  "error": "forgot password must be initiated from the Clerk client SDK using the reset_password_email_code flow"
}
```

with status `410 Gone`.

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
