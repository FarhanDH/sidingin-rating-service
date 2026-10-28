# Rating Service API Contract

Base URL: `https://api.sidingin.my.id/api/v1`

All responses are JSON. Successful responses use HTTP 2xx; errors use 4xx/5xx with a standard error body.

---

## 1. Get Rating Direction Status

**`GET /ratings/tickets/:order_id/status`**

Fetches the rating status and existing ratings for all directions of a given order. The direction is dynamic based on role authorization.

### Path Parameters

| Parameter  | Type   | Required | Description                  |
| ---------- | ------ | -------- | ---------------------------- |
| `order_id` | string | Yes      | The order identifier (UUID). |

### Response (200 OK)

| Field                          | Type              | Description                                             |
| ------------------------------ | ----------------- | ------------------------------------------------------- |
| `message`                      | string            | Human-readable status message.                          |
| `data.order_id`                | string            | The order identifier.                                   |
| `data.expires_at`              | string (ISO-8601) | Deadline after which rating is no longer allowed.       |
| `data.order_completed_at`      | string (ISO-8601) | Timestamp when the order was completed.                 |
| `data.directions`              | object            | Map of direction key to rating details.                 |
| `data.directions.<key>.status` | string            | Rating status: `"rated"`, `"locked"`, `"pending"`, etc. |
| `data.directions.<key>.stars`  | number            | Stars given (0–5). `0` when not yet rated.              |

### Direction Keys

| Key                      | Description                         |
| ------------------------ | ----------------------------------- |
| `customer_to_technician` | Customer rates the technician.      |
| `customer_to_app`        | Customer rates the app/platform.    |
| `customer_to_service`    | Customer rates the overall service. |
| `technician_to_customer` | Technician rates the customer.      |
| `technician_to_app`      | Technician rates the app/platform.  |

### Example Request

```http
GET /api/v1/ratings/tickets/S-sdfkj/status
```

### Example Response

```json
{
  "message": "Ticket status fetched successfully",
  "data": {
    "order_id": "S-sdfkj",
    "expires_at": "2026-11-10T03:12:00Z",
    "order_completed_at": "2026-11-10T03:12:00Z",
    "directions": {
      "customer_to_technician": {
        "status": "rated",
        "stars": 5
      },
      "customer_to_app": {
        "status": "open",
        "stars": 0
      },
      "customer_to_service": {
        "status": "locked",
        "stars": 0
      }
    }
  }
}
```

---

## 2. Get Rating Tags by Direction

**`GET /ratings/forms`**

Returns the available rating tags (codes and labels) for a specific direction. FrontEnd uses these to render the tag selection UI.

### Query Parameters

| Parameter   | Type   | Required | Description                                                                                                                |
| ----------- | ------ | -------- | -------------------------------------------------------------------------------------------------------------------------- |
| `direction` | string | Yes      | One of: `customer_to_technician`, `customer_to_app`, `customer_to_service`, `technician_to_customer`, `technician_to_app`. |

### Response (200 OK)

| Field              | Type   | Description                                      |
| ------------------ | ------ | ------------------------------------------------ |
| `message`          | string | Human-readable status message.                   |
| `data`             | array  | List of tag objects.                             |
| `data[].code`      | string | Machine-readable tag code (send this in submit). |
| `data[].label`     | string | Human-readable tag label (display this in UI).   |
| `data[].min_stars` | number | Minimum star rating for this tag to apply.       |
| `data[].max_stars` | number | Maximum star rating for this tag to apply.       |

### Example Request

```http
GET /api/v1/ratings/forms?direction=customer_to_technician
```

### Example Response

```json
{
  "message": "Form tags fetched successfully",
  "data": [
    {
      "code": "tech_on_time",
      "label": "Tepat waktu",
      "min_stars": 4,
      "max_stars": 5
    },
    {
      "code": "tech_late",
      "label": "Teknisi ngaret",
      "min_stars": 1,
      "max_stars": 3
    }
  ]
}
```

---

## 3. Submit Rating

**`POST /ratings`**

Submits a customer rating for a specific order and direction.

### Request Body

| Field       | Type          | Required | Description                                                                 |
| ----------- | ------------- | -------- | --------------------------------------------------------------------------- |
| `order_id`  | string        | Yes      | The order identifier.                                                       |
| `direction` | string        | Yes      | One of: `customer_to_technician`, `customer_to_app`, `customer_to_service`. |
| `stars`     | integer (1–5) | Yes      | Star rating given by the customer.                                          |
| `comment`   | string        | No       | Free-text comment from the customer.                                        |
| `tag_codes` | array[string] | No       | List of tag codes (from endpoint #2) applicable to this rating.             |

### Response (201 Created)

| Field            | Type   | Description                              |
| ---------------- | ------ | ---------------------------------------- |
| `message`        | string | Human-readable status message.           |
| `data.rating_id` | string | Unique identifier of the created rating. |

### Example Request

```http
POST /api/v1/ratings
Content-Type: application/json


{
  "order_id": "S-sdfkj",
  "direction": "customer_to_technician",
  "stars": 5,
  "comment": "Teknisi sat set",
  "is_comment_hidden": false,
  "tag_codes": ["tech_on_time", "tech_friendly"]
}
```

### Example Response

```json
{
  "message": "Rating submitted successfully",
  "data": {
    "rating_id": "sldkfjwaoe8fj2qefjqep[fj"
  }
}
```

---

## 4. Error Responses

All error responses follow this shape:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error description",
    "details": null
  }
}
```

| HTTP Status | Error Code              | Meaning                                         |
| ----------- | ----------------------- | ----------------------------------------------- |
| 400         | `VALIDATION_ERROR`      | Request body or parameters invalid.             |
| 401         | `UNAUTHORIZED`          | Missing or invalid authentication.              |
| 403         | `FORBIDDEN`             | Not allowed to rate this order.                 |
| 404         | `NOT_FOUND`             | Order or rating not found.                      |
| 409         | `ALREADY_RATED`         | Rating already submitted for this direction.    |
| 422         | `RATING_LOCKED`         | Rating window expired or order not completable. |
| 500         | `INTERNAL_SERVER_ERROR` | Unexpected server error.                        |

---

## 5. Quick Reference for FrontEnd

| Goal                            | Endpoint                               | Method |
| ------------------------------- | -------------------------------------- | ------ |
| Check if order can be rated     | `/ratings/tickets/:order_id/status`    | GET    |
| Get tag options for a direction | `/ratings/forms?direction=<direction>` | GET    |
| Submit stars + comment + tags   | `/ratings`                             | POST   |

**Typical flow:**

1. `GET /ratings/tickets/:order_id/status` → check `directions.<key>.status` to see if already `"rated"` or `"locked"`.
2. `GET /ratings/forms?direction=customer_to_technician` → render tag chips using `label`, filter by `min_stars`/`max_stars` based on selected stars.
3. `POST /ratings` → send `order_id`, `direction`, `stars`, optional `comment`, and optional `tag_codes`.
