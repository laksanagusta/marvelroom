# Pagination API Documentation

## Overview
The business trip API now supports advanced pagination, filtering, and sorting using a flexible query builder system.

## API Endpoint
**GET** `/api/v1/business-trips`

## Query Parameters

### Pagination
- `page` - Page number (default: 1, min: 1)
- `limit` - Items per page (default: 20, min: 1, max: 100)

### Sorting
- `sort` - Sort fields in format: "field1 direction1,field2 direction2"
  - `direction` can be `asc` or `desc` (default: `asc`)
  - Example: `sort=created_at desc,start_date asc`

### Filtering
Use field-based filtering with the format: `{field}={operator} {value}`

#### Available Fields
- `activity_purpose` - Search in activity purpose
- `destination_city` - Search in destination city
- `start_date` - Filter by start date
- `end_date` - Filter by end date
- `created_at` - Filter by creation date
- `updated_at` - Filter by update date

#### Available Operators
- `eq` - Equals (=)
- `ne` - Not equals (!=)
- `gt` - Greater than (>)
- `gte` - Greater than or equal (>=)
- `lt` - Less than (<)
- `lte` - Less than or equal (<=)
- `like` - LIKE (case-sensitive)
- `ilike` - ILIKE (case-insensitive)
- `in` - IN (comma-separated values)
- `nin` - NOT IN (comma-separated values)
- `is` - IS NULL
- `is_not` - IS NOT NULL

## Response Format

```json
{
  "message": "Business trips retrieved successfully",
  "data": {
    "businessTrips": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "startDate": "2024-12-01",
        "endDate": "2024-12-05",
        "activityPurpose": "Client Meeting",
        "destinationCity": "Jakarta",
        "totalCost": 2500000.00,
        "assignees": [...],
        "createdAt": "2024-11-17T08:00:00Z",
        "updatedAt": "2024-11-17T08:00:00Z"
      }
    ],
    "total": 150,
    "page": 1,
    "limit": 20,
    "totalPages": 8
  }
}
```

## Examples

### Basic Pagination
```bash
# Get first page with default limit (20)
GET /api/v1/business-trips?page=1

# Get second page with 10 items per page
GET /api/v1/business-trips?page=2&limit=10
```

### Sorting
```bash
# Sort by creation date (newest first)
GET /api/v1/business-trips?sort=created_at desc

# Sort by start date ascending, then destination city descending
GET /api/v1/business-trips?sort=start_date asc,destination_city desc
```

### Filtering
```bash
# Search by activity purpose (case-insensitive)
GET /api/v1/business-trips?activity_purpose=ilike meeting

# Filter by destination city
GET /api/v1/business-trips?destination_city=eq Jakarta

# Filter by date range
GET /api/v1/business-trips?start_date=gte 2024-12-01&end_date=lte 2024-12-31

# Filter by multiple cities
GET /api/v1/business-trips?destination_city=in Jakarta,Surabaya,Bandung
```

### Combined Queries
```bash
# Complex query with pagination, sorting, and multiple filters
GET /api/v1/business-trips?page=2&limit=15&sort=created_at desc&destination_city=ilike jakarta&start_date=gte 2024-12-01
```

## Advanced Features

### Case-Insensitive Search
```bash
# Search in activity purpose and destination city
GET /api/v1/business-trips?activity_purpose=ilike client&destination_city=ilike jakarta
```

### Date Range Filtering
```bash
# Business trips in December 2024
GET /api/v1/business-trips?start_date=gte 2024-12-01&end_date=lte 2024-12-31

# Recent trips (last 30 days)
GET /api/v1/business-trips?created_at=gte 2024-10-18
```

### Multi-value Filtering
```bash
# Trips to specific cities
GET /api/v1/business-trips?destination_city=in Jakarta,Surabaya,Medan

# Exclude specific destinations
GET /api/v1/business-trips?destination_city=nin Jakarta,Bandung
```

## Implementation Details

The pagination system uses:
- **QueryBuilder**: Dynamically constructs SQL queries with proper parameter binding
- **QueryParser**: Parses URL query parameters into structured filter objects
- **Field Validation**: Only allows predefined fields to prevent SQL injection
- **Type Conversion**: Automatically converts string values to appropriate types (dates, numbers, etc.)

## Error Handling

Invalid queries return structured error responses:
```json
{
  "error": "Validation failed",
  "details": "invalid field: invalid_field_name"
}
```

Common errors:
- Invalid field names
- Unsupported operators
- Invalid date formats
- Out-of-range pagination values

## Performance Considerations

- Always applies `deleted_at IS NULL` filter automatically
- Uses parameterized queries to prevent SQL injection
- Supports indexed database fields for optimal performance
- Limits maximum page size to 100 items