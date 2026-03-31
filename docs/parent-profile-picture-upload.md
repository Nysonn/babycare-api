# Parent Profile Picture Upload

The parent profile endpoint now supports profile picture uploads.

## Endpoint

- Method: `PUT`
- Path: `/api/v1/parents/profile`
- Auth: Bearer token for a user with the `parent` role

## Request formats

### `multipart/form-data`

Use this when the user is uploading or replacing a profile picture.

Supported form fields:

- `profile_picture`: image file
- `location`: string
- `occupation`: string
- `preferred_hours`: string
- `primary_location`: string

Behavior:

- If `profile_picture` is included, the backend uploads it to Cloudinary and saves the returned URL.
- If a text field is omitted, the backend keeps the current stored value.
- If a text field is sent as an empty string, the backend clears that field.

Example:

```bash
curl -X PUT http://localhost:8080/api/v1/parents/profile \
  -H "Authorization: Bearer <token>" \
  -F "location=Accra" \
  -F "occupation=Teacher" \
  -F "preferred_hours=Weekdays 8am-5pm" \
  -F "primary_location=East Legon" \
  -F "profile_picture=@/path/to/avatar.jpg"
```

### `application/json`

JSON is still supported for text-only profile updates.

Example:

```json
{
  "location": "Accra",
  "occupation": "Teacher",
  "preferred_hours": "Weekdays 8am-5pm",
  "primary_location": "East Legon"
}
```

## Response change

Both `GET /api/v1/parents/profile` and `PUT /api/v1/parents/profile` now return:

- `profile_picture_url`: string

Example response:

```json
{
  "id": "user-id",
  "full_name": "Jane Doe",
  "email": "jane@example.com",
  "phone": "+233000000000",
  "status": "active",
  "location": "Accra",
  "primary_location": "East Legon",
  "occupation": "Teacher",
  "preferred_hours": "Weekdays 8am-5pm",
  "profile_picture_url": "https://res.cloudinary.com/.../image/upload/...jpg",
  "created_at": "2026-03-31T10:00:00Z"
}
```

## Frontend integration notes

- The file field name must be `profile_picture`.
- Use `multipart/form-data` whenever the user selects an image.
- The backend returns a hosted URL; the frontend does not need to upload directly to Cloudinary.
- Existing text-only profile edits can continue using JSON if preferred.