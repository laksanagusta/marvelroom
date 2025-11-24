# HTTP Status Response Changes

## Overview
All PUT/PATCH/POST operations now return HTTP status codes only, without JSON response data, following REST API best practices.

## Response Changes

### Before (with JSON data)
```json
{
  "message": "Business trip created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "2024-12-01",
    // ... other fields
  }
}
```

### After (status only)
- **POST**: Returns `201 Created` with empty body
- **PUT/PATCH**: Returns `200 OK` with empty body

## Updated Endpoints

### Business Trip Operations
| Method | Endpoint | Status Code | Response |
|--------|----------|-------------|----------|
| POST | `/api/v1/business-trips` | `201 Created` | Empty body |
| PUT | `/api/v1/business-trips/{id}` | `200 OK` | Empty body |
| PUT | `/api/v1/business-trips/{id}/with-assignees` | `200 OK` | Empty body |

### Assignee Operations
| Method | Endpoint | Status Code | Response |
|--------|----------|-------------|----------|
| POST | `/api/v1/business-trips/{tripId}/assignees` | `201 Created` | Empty body |
| PUT | `/api/v1/business-trips/{tripId}/assignees/{assigneeId}` | `200 OK` | Empty body |

### Transaction Operations
| Method | Endpoint | Status Code | Response |
|--------|----------|-------------|----------|
| POST | `/api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions` | `201 Created` | Empty body |
| PUT | `/api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions/{transactionId}` | `200 OK` | Empty body |

### Legacy Endpoints (unchanged)
| Method | Endpoint | Status Code | Response |
|--------|----------|-------------|----------|
| POST | `/api/business-trips/{businessTripId}/assignees` | `201 Created` | Empty body |
| POST | `/api/assignees/{assigneeId}/transactions` | `201 Created` | Empty body |

## Error Responses
Error responses still return JSON with error details:

```json
{
  "error": "Validation failed",
  "details": "field validation message"
}
```

## Benefits

1. **REST Compliance**: Follows REST API standards
2. **Bandwidth Efficiency**: Reduces response size
3. **Consistency**: Uniform response pattern for mutations
4. **Simplified Client Logic**: No need to parse response data for successful operations
5. **HTTP Semantics**: Uses HTTP status codes to indicate success/failure

## Client Usage Examples

### Create Business Trip
```bash
# Request
curl -X POST http://localhost:3000/api/v1/business-trips \
  -H "Content-Type: application/json" \
  -d @payload.json

# Response: 201 Created (empty body)
# Success if HTTP status is 201, check response body for errors
```

### Update Business Trip
```bash
# Request
curl -X PUT http://localhost:3000/api/v1/business-trips/{id}/with-assignees \
  -H "Content-Type: application/json" \
  -d @update_payload.json

# Response: 200 OK (empty body)
# Success if HTTP status is 200, check response body for errors
```

### Error Handling
```bash
# If validation fails:
# Response: 400 Bad Request
# Body: {"error": "Validation failed", "details": "..."}
```

## Read Operations (unchanged)
GET operations still return JSON data:

```bash
# Get business trip
GET /api/v1/business-trips/{id}
# Returns: 200 OK with full business trip data

# List business trips
GET /api/v1/business-trips
# Returns: 200 OK with paginated list
```

## Migration Guide for Client Applications

### Before
```javascript
const response = await fetch('/api/v1/business-trips', {
  method: 'POST',
  body: JSON.stringify(payload)
});
const data = await response.json();
if (response.ok) {
  console.log('Created:', data.data);
}
```

### After
```javascript
const response = await fetch('/api/v1/business-trips', {
  method: 'POST',
  body: JSON.stringify(payload)
});

if (response.ok) {
  console.log('Created successfully');
  // Get the created resource with a separate GET request if needed
  const createdData = await fetch(response.headers.get('Location') || `/api/v1/business-trips/${id}`);
}
```

This change makes the API more RESTful and consistent with modern API design principles.