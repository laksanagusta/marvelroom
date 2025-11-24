# Business Trip API Documentation

This API provides endpoints for managing business trips, assignees, and transactions.

## Base URL
```
http://localhost:5002/api
```

## Database Setup
Before using the API, make sure to run the migration script:
```bash
psql -U your_username -d your_database -f migrations/001_create_business_trip_tables.sql
```

## Authentication
Currently, no authentication is required for these endpoints.

## Response Format
All responses follow this format:
```json
{
  "message": "string",
  "data": "object|array"
}
```

Error responses:
```json
{
  "error": "string",
  "details": "string"
}
```

## Business Trip Endpoints

### Create Business Trip
```http
POST /business-trips
```

**Request Body:**
```json
{
  "startDate": "2024-01-15",
  "endDate": "2024-01-17",
  "activityPurpose": "Meeting with client",
  "destinationCity": "Jakarta",
  "spdDate": "2024-01-10",
  "departureDate": "2024-01-15",
  "returnDate": "2024-01-17",
  "assignees": [
    {
      "name": "John Doe",
      "spd_number": "SPD-2024-001",
      "employee_id": "123456789",
      "position": "Software Engineer",
      "rank": "III/a",
      "transactions": [
        {
          "name": "Hotel Booking",
          "type": "accommodation",
          "subtype": "hotel",
          "amount": 500000,
          "total_night": 2,
          "description": "Hotel di Jakarta Selatan"
        },
        {
          "name": "Flight Ticket",
          "type": "transport",
          "subtype": "flight",
          "amount": 1500000,
          "description": "Garuda Indonesia CGK-DPS"
        },
        {
          "name": "Daily Allowance",
          "type": "allowance",
          "subtype": "daily_allowance",
          "amount": 300000,
          "description": "Uang harian 3 hari"
        }
      ]
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "message": "Business trip created successfully",
  "data": {
    "id": "uuid",
    "startDate": "2024-01-15",
    "endDate": "2024-01-17",
    "activityPurpose": "Meeting with client",
    "destinationCity": "Jakarta",
    "spdDate": "2024-01-10",
    "departureDate": "2024-01-15",
    "returnDate": "2024-01-17",
    "totalCost": 2900000,
    "assignees": [
      {
        "id": "uuid",
        "name": "John Doe",
        "spdNumber": "SPD-2024-001",
        "employeeId": "123456789",
        "position": "Software Engineer",
        "rank": "III/a",
        "totalCost": 2900000,
        "transactions": [
          {
            "id": "uuid",
            "name": "Hotel Booking",
            "type": "accommodation",
            "subtype": "hotel",
            "amount": 500000,
            "totalNight": 2,
            "subtotal": 1000000,
            "description": "Hotel di Jakarta Selatan"
          },
          {
            "id": "uuid",
            "name": "Flight Ticket",
            "type": "transport",
            "subtype": "flight",
            "amount": 1500000,
            "totalNight": null,
            "subtotal": 1500000,
            "description": "Garuda Indonesia CGK-DPS"
          },
          {
            "id": "uuid",
            "name": "Daily Allowance",
            "type": "allowance",
            "subtype": "daily_allowance",
            "amount": 300000,
            "totalNight": null,
            "subtotal": 300000,
            "description": "Uang harian 3 hari"
          }
        ],
        "createdAt": "2024-01-10T10:00:00Z",
        "updatedAt": "2024-01-10T10:00:00Z"
      }
    ],
    "createdAt": "2024-01-10T10:00:00Z",
    "updatedAt": "2024-01-10T10:00:00Z"
  }
}
```

### List Business Trips
```http
GET /business-trips?page=1&limit=10&search=meeting&destination_city=Jakarta&start_date=2024-01-01&end_date=2024-01-31&sort_by=created_at&sort_direction=desc
```

**Query Parameters:**
- `page` (int, default: 1) - Page number
- `limit` (int, default: 10, max: 100) - Items per page
- `search` (string) - Search in activity purpose and destination city
- `destination_city` (string) - Filter by destination city
- `start_date` (string, format: YYYY-MM-DD) - Filter by start date
- `end_date` (string, format: YYYY-MM-DD) - Filter by end date
- `sort_by` (string, default: created_at) - Field to sort by
- `sort_direction` (string, default: desc) - Sort direction (asc/desc)

**Response:**
```json
{
  "message": "Business trips retrieved successfully",
  "data": {
    "businessTrips": [
      {
        "id": "uuid",
        "startDate": "2024-01-15",
        "endDate": "2024-01-17",
        "activityPurpose": "Meeting with client",
        "destinationCity": "Jakarta",
        "totalCost": 2900000,
        "assignees": [...]
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "totalPages": 1
  }
}
```

### Get Business Trip by ID
```http
GET /business-trips/{id}
```

**Response:**
```json
{
  "message": "Business trip retrieved successfully",
  "data": {
    "id": "uuid",
    "startDate": "2024-01-15",
    "endDate": "2024-01-17",
    "activityPurpose": "Meeting with client",
    "destinationCity": "Jakarta",
    "spdDate": "2024-01-10",
    "departureDate": "2024-01-15",
    "returnDate": "2024-01-17",
    "totalCost": 2900000,
    "assignees": [...],
    "createdAt": "2024-01-10T10:00:00Z",
    "updatedAt": "2024-01-10T10:00:00Z"
  }
}
```

### Update Business Trip
```http
PUT /business-trips/{id}
```

**Request Body (all fields optional):**
```json
{
  "startDate": "2024-01-16",
  "endDate": "2024-01-18",
  "activityPurpose": "Updated meeting purpose",
  "destinationCity": "Surabaya",
  "spdDate": "2024-01-11",
  "departureDate": "2024-01-16",
  "returnDate": "2024-01-18"
}
```

### Delete Business Trip
```http
DELETE /business-trips/{id}
```

**Response:**
```json
{
  "message": "Business trip deleted successfully"
}
```

### Get Business Trip Summary
```http
GET /business-trips/{id}/summary
```

**Response:**
```json
{
  "message": "Business trip summary retrieved successfully",
  "data": {
    "business_trip_id": "uuid",
    "total_cost": 2900000,
    "total_assignees": 1,
    "total_transactions": 3,
    "cost_by_type": {
      "accommodation": 1000000,
      "transport": 1500000,
      "allowance": 300000,
      "other": 0
    }
  }
}
```

## Assignee Endpoints

### Add Assignee to Business Trip
```http
POST /business-trips/{businessTripId}/assignees
```

**Request Body:**
```json
{
  "name": "Jane Smith",
  "spd_number": "SPD-2024-002",
  "employee_id": "987654321",
  "position": "Project Manager",
  "rank": "IV/a",
  "transactions": [
    {
      "name": "Taxi from Airport",
      "type": "transport",
      "subtype": "taxi",
      "amount": 150000,
      "description": "Taxi dari CGK ke hotel"
    }
  ]
}
```

## Transaction Endpoints

### Add Transaction to Assignee
```http
POST /assignees/{assigneeId}/transactions
```

**Request Body:**
```json
{
  "name": "Lunch Meeting",
  "type": "other",
  "subtype": "meal",
  "amount": 100000,
  "description": "Makan siang dengan client"
}
```

### Get Assignee Summary
```http
GET /assignees/{id}/summary
```

**Response:**
```json
{
  "message": "Assignee summary retrieved successfully",
  "data": {
    "assignee_id": "uuid",
    "assignee_name": "John Doe",
    "total_cost": 2900000,
    "total_transactions": 3,
    "cost_by_type": {
      "accommodation": 1000000,
      "transport": 1500000,
      "allowance": 300000,
      "other": 0
    }
  }
}
```

## Transaction Types

### Type Values
- `accommodation` - For hotel/room bookings
- `transport` - For transportation (flight, train, taxi, etc.)
- `allowance` - For daily allowances and per diems
- `other` - For miscellaneous expenses

### Subtype Values
- `hotel` - Hotel accommodation
- `flight` - Flight tickets
- `train` - Train tickets
- `taxi` - Taxi/ride-hailing services
- `daily_allowance` - Daily allowance
- `rental_car` - Car rental
- `meal` - Food and meals
- `other` - Other subtypes

## Validation Rules

### Business Trip
- `startDate`, `endDate`, `spdDate`, `departureDate`, `returnDate`: Required, must be valid dates in YYYY-MM-DD format
- `activityPurpose`: Required, 1-255 characters
- `destinationCity`: Required, 1-255 characters
- `assignees`: Required, 1-50 assignees
- Date validation: `startDate <= endDate`, `departureDate <= returnDate`, `spdDate <= departureDate`

### Assignee
- `name`: Required, 1-255 characters
- `spd_number`: Required, 1-100 characters (must be unique within a business trip)
- `employee_number`: Required, 1-50 characters (NIP/employee number from external API)
- `position`: Required, 1-255 characters
- `rank`: Required, 1-100 characters
- `transactions`: Optional, each transaction must pass transaction validation

### Transaction
- `name`: Required, 1-255 characters
- `type`: Required, must be one of: accommodation, transport, other, allowance
- `subtype`: Optional, must be one of: hotel, flight, train, taxi, daily_allowance, rental_car, meal, other
- `amount`: Required, must be non-negative
- `total_night`: Optional, non-negative integer (used for accommodation calculations)
- `description`: Optional, max 1000 characters
- `transport_detail`: Optional, max 1000 characters

### Query Parameters
- `page`: Optional, must be ≥ 1
- `limit`: Optional, must be between 1-100
- `start_date`, `end_date`: Optional, must be valid dates in YYYY-MM-DD format
- `sort_direction`: Optional, must be 'asc' or 'desc'

## Subtotal Calculation
- For `accommodation` type: `subtotal = amount × total_night` (if total_night > 0)
- For other types: `subtotal = amount`

## Example Usage with curl

### Create a business trip
```bash
curl -X POST http://localhost:5002/api/business-trips \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2024-01-15",
    "endDate": "2024-01-17",
    "activityPurpose": "Client Meeting",
    "destinationCity": "Jakarta",
    "spdDate": "2024-01-10",
    "departureDate": "2024-01-15",
    "returnDate": "2024-01-17",
    "assignees": [
      {
        "name": "John Doe",
        "spd_number": "SPD-2024-001",
        "employee_id": "123456789",
        "position": "Software Engineer",
        "rank": "III/a",
        "transactions": [
          {
            "name": "Hotel Booking",
            "type": "accommodation",
            "subtype": "hotel",
            "amount": 500000,
            "total_night": 2,
            "description": "Hotel di Jakarta"
          }
        ]
      }
    ]
  }'
```

### List business trips
```bash
curl "http://localhost:5002/api/business-trips?page=1&limit=10"
```

### Get business trip summary
```bash
curl "http://localhost:5002/api/business-trips/{business-trip-id}/summary"
```

## Error Handling

### Common Error Codes
- `400 Bad Request` - Validation errors, invalid request format
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server-side errors

### Error Response Format
```json
{
  "error": "Validation failed",
  "details": "startDate: must be a valid date"
}
```

## Environment Variables
Set these in your `.env` file:
```env
DATABASE_DSN=postgres://username:password@localhost/business_trip?sslmode=disable
PORT=5002
CORS_ALLOW_ORIGINS=http://localhost:3000
```

## Database Schema
The API uses three main tables:
1. `business_trips` - Stores business trip information
2. `assignees` - Stores employee/assignee information linked to business trips
3. `transactions` - Stores transaction details linked to assignees

Each table has proper foreign key relationships and constraints to ensure data integrity.