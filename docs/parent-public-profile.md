# Parent Public Profile — Babysitter Access

## Endpoint

```
GET /api/v1/parents/:id
```

**Auth required:** Yes (Bearer token)  
**Allowed roles:** `babysitter`

---

## Path Parameters

| Parameter | Type   | Description                          |
|-----------|--------|--------------------------------------|
| `id`      | string | UUID of the parent user to look up   |

---

## Response — `200 OK`

```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "full_name": "Jane Doe",
  "location": "Nairobi",
  "primary_location": "Westlands",
  "occupation": "Teacher",
  "preferred_hours": "Morning",
  "profile_picture_url": "https://res.cloudinary.com/example/image/upload/v1234/parents/jane.jpg"
}
```

`profile_picture_url` is an empty string `""` when the parent has not uploaded a photo.

---

## Error Responses

| Status | Body                                      | Reason                                  |
|--------|-------------------------------------------|-----------------------------------------|
| 400    | `{"error": "invalid user id"}`            | `:id` is not a valid UUID               |
| 401    | `{"error": "unauthorised"}`               | Token missing or expired                |
| 403    | `{"error": "forbidden"}`                  | Caller is not a babysitter              |
| 404    | `{"error": "parent not found"}`           | No parent user with that ID             |
| 404    | `{"error": "parent profile not found"}`   | User exists but has no profile yet      |
| 500    | `{"error": "internal server error"}`      | Unexpected server error                 |

---

## Usage Notes

- This endpoint is intended for babysitters who need to see the profile (including photo) of a parent they are in a conversation with. Retrieve the parent's user ID from the conversation object and pass it here.
- Sensitive fields (email, phone, status) are intentionally excluded from this public view.
- Parents access their **own** profile via `GET /api/v1/parents/profile` (parent role only).
