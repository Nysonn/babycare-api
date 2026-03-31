# Conversation Previews API

Returns the authenticated user's conversation list enriched with the most recent
message in each thread. This is the primary endpoint for populating a messaging
inbox / conversation list screen.

---

## Endpoint

```
GET /api/v1/conversations/previews
```

### Authentication

Requires a valid session token in the `Authorization` header.  
Accessible to users with role **parent** or **babysitter**.

---

## Response

**200 OK** — array of conversation preview objects (ordered by most recently
updated conversation first). Conversations that have no messages yet are omitted.

```json
[
  {
    "conversation_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "other_user_name": "Jane Smith",
    "other_user_profile_picture_url": "https://res.cloudinary.com/example/image/upload/v1/babysitters/jane.jpg",
    "is_locked": false,
    "last_message_id": "8c9d1e2f-3a4b-5c6d-7e8f-9a0b1c2d3e4f",
    "last_sender_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
    "preview_text": "Hey, are you available this weekend?",
    "is_read": false,
    "last_message_sent": "2026-04-01T10:30:00Z"
  }
]
```

### Response fields

| Field | Type | Description |
|---------------------|----------|-------------|
| `conversation_id` | UUID | Unique ID of the conversation thread. |
| `other_user_name` | string | Full name of the other participant. |
| `other_user_profile_picture_url` | string | Profile picture URL of the other participant. Empty string `""` if no picture has been uploaded. |
| `is_locked` | boolean | `true` when the conversation has been administratively locked (e.g. babysitter suspended). |
| `last_message_id` | UUID | ID of the most recent message in this thread. |
| `last_sender_id` | UUID | ID of the user who sent the last message. |
| `preview_text` | string | Content of the last message — display directly as the preview snippet. |
| `is_read` | boolean | Whether the last message has been read by the recipient. |
| `last_message_sent` | ISO 8601 | Timestamp of the last message. |

---

## Showing unread indicators

A message is **unread for the current user** when **both** of the following are true:

```
is_read == false  AND  last_sender_id != <current user's ID>
```

Use this combination to decide whether to render an unread badge or bold the
conversation row. If `last_sender_id` equals the current user's own ID the
message is outgoing, so no unread indicator should be shown regardless of
`is_read`.

---

## Error responses

| Status | Body | Reason |
|--------|------|--------|
| `401 Unauthorized` | `{"error": "unauthorised"}` | Missing or invalid session token. |
| `403 Forbidden` | `{"error": "forbidden"}` | Authenticated user does not have the required role. |
| `500 Internal Server Error` | `{"error": "internal server error"}` | Database failure. |

---

## Example — Flutter / Dart

```dart
final response = await http.get(
  Uri.parse('$baseUrl/api/v1/conversations/previews'),
  headers: {'Authorization': 'Bearer $token'},
);

if (response.statusCode == 200) {
  final List<dynamic> data = jsonDecode(response.body);
  final previews = data.map((e) => ConversationPreview.fromJson(e)).toList();
}
```

### Suggested Dart model

```dart
class ConversationPreview {
  final String conversationId;
  final String otherUserName;
  final String otherUserProfilePictureUrl;
  final bool isLocked;
  final String lastMessageId;
  final String lastSenderId;
  final String previewText;
  final bool isRead;
  final DateTime lastMessageSent;

  ConversationPreview({
    required this.conversationId,
    required this.otherUserName,
    required this.otherUserProfilePictureUrl,
    required this.isLocked,
    required this.lastMessageId,
    required this.lastSenderId,
    required this.previewText,
    required this.isRead,
    required this.lastMessageSent,
  });

  factory ConversationPreview.fromJson(Map<String, dynamic> json) {
    return ConversationPreview(
      conversationId: json['conversation_id'] as String,
      otherUserName: json['other_user_name'] as String,
      otherUserProfilePictureUrl: json['other_user_profile_picture_url'] as String,
      isLocked: json['is_locked'] as bool,
      lastMessageId: json['last_message_id'] as String,
      lastSenderId: json['last_sender_id'] as String,
      previewText: json['preview_text'] as String,
      isRead: json['is_read'] as bool,
      lastMessageSent: DateTime.parse(json['last_message_sent'] as String),
    );
  }

  /// Returns true when the current user has not yet read this message.
  bool hasUnread(String currentUserId) =>
      !isRead && lastSenderId != currentUserId;
}
```
