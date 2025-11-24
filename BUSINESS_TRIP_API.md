# Business Trip API Documentation

## Overview
This API provides complete business trip management functionality with support for assignees and their financial transactions.

## API Endpoints

### Create Business Trip with Assignees and Transactions
**POST** `/api/v1/business-trips`

Creates a new business trip with multiple assignees, each with their own transactions in a single API call.

#### Request Body
```json
{
  "start_date": "2024-12-01",
  "end_date": "2024-12-05",
  "activity_purpose": "Client Meeting and Project Discussion",
  "destination_city": "Jakarta",
  "spd_date": "2024-11-25",
  "departure_date": "2024-12-01",
  "return_date": "2024-12-05",
  "assignees": [
    {
      "name": "John Doe",
      "spd_number": "SPD-2024-001",
      "employee_id": "EMP-001",
      "employee_name": "John Doe",
      "position": "Senior Software Engineer",
      "rank": "Senior Level",
      "transactions": [
        {
          "name": "Hotel Accommodation",
          "type": "accommodation",
          "subtype": "hotel",
          "amount": 150.00,
          "total_night": 4,
          "description": "Standard room with breakfast",
          "transport_detail": ""
        }
      ]
    }
  ]
}
```

#### Response
```json
{
  "message": "Business trip created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "2024-12-01",
    "end_date": "2024-12-05",
    "activity_purpose": "Client Meeting and Project Discussion",
    "destination_city": "Jakarta",
    "total_cost": 2480.00,
    "assignees": [
      {
        "id": "assignee-id-1",
        "name": "John Doe",
        "spd_number": "SPD-2024-001",
        "employee_id": "EMP-001",
        "employee_name": "John Doe",
        "position": "Senior Software Engineer",
        "rank": "Senior Level",
        "total_cost": 1050.00,
        "transactions": [...]
      }
    ],
    "created_at": "2024-11-17T08:00:00Z",
    "updated_at": "2024-11-17T08:00:00Z"
  }
}
```

### Other Endpoints

#### Business Trip Operations
- `GET /api/v1/business-trips` - List business trips with pagination and filtering
- `GET /api/v1/business-trips/{tripId}` - Get specific business trip
- `PUT /api/v1/business-trips/{tripId}` - Update business trip details
- `DELETE /api/v1/business-trips/{tripId}` - Delete business trip

#### Assignee Operations
- `POST /api/v1/business-trips/{tripId}/assignees` - Add assignee to business trip
- `GET /api/v1/business-trips/{tripId}/assignees` - List assignees
- `GET /api/v1/business-trips/{tripId}/assignees/{assigneeId}` - Get specific assignee
- `PUT /api/v1/business-trips/{tripId}/assignees/{assigneeId}` - Update assignee
- `DELETE /api/v1/business-trips/{tripId}/assignees/{assigneeId}` - Delete assignee

#### Transaction Operations
- `POST /api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions` - Add transaction
- `GET /api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions` - List transactions
- `PUT /api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions/{transactionId}` - Update transaction
- `DELETE /api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions/{transactionId}` - Delete transaction

## Transaction Types

### Types
- `accommodation` - Hotel and lodging expenses
- `transport` - Flight, train, taxi, rental car
- `allowance` - Daily allowances, meal allowances
- `other` - Miscellaneous expenses

### Subtypes
- `hotel` - Hotel accommodation
- `flight` - Flight tickets
- `train` - Train tickets
- `taxi` - Taxi/ride-sharing
- `daily_allowance` - Daily allowance
- `rental_car` - Car rental
- `meal` - Meal expenses
- `other` - Other subtypes

## Query Parameters

For listing endpoints:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `search` - Search in activity purpose
- `destination_city` - Filter by destination city
- `start_date` - Filter by start date (YYYY-MM-DD)
- `end_date` - Filter by end date (YYYY-MM-DD)
- `sort_by` - Sort field (default: created_at)
- `sort_direction` - Sort direction: asc or desc (default: desc)

## Example Usage

### 1. Create a Complete Business Trip
```bash
curl -X POST http://localhost:3000/api/v1/business-trips \
  -H "Content-Type: application/json" \
  -d @test_payload_snake_case.json
```

### 2. List Business Trips
```bash
curl "http://localhost:3000/api/v1/business-trips?page=1&limit=10&search=meeting"
```

### 3. Add Transaction to Existing Assignee
```bash
curl -X POST http://localhost:3000/api/v1/business-trips/{tripId}/assignees/{assigneeId}/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Extra Meal",
    "type": "other",
    "amount": 25.00,
    "description": "Business dinner"
  }'
```

## Validation

All endpoints include comprehensive validation:
- Required fields validation
- Date format validation (YYYY-MM-DD)
- Transaction type/subtype validation
- Business logic validation (e.g., date ranges)
- Unique SPD number validation per business trip

## Error Handling

The API returns consistent error responses:
```json
{
  "error": "Validation failed",
  "details": "start date must be before or equal to end date"
}
```

HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error