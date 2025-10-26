# DriftLock API Documentation

This document provides detailed information about the DriftLock API endpoints, authentication, and usage examples.

## Authentication

All API endpoints require authentication using JWT tokens. To authenticate:

1. Register a new account or log in to get an authentication token
2. Include the token in the `Authorization` header as a Bearer token

```http
Authorization: Bearer <your-jwt-token>
```

## API Base URL

```
https://api.driftlock.com/api/v1
```

## Endpoints

### Authentication

#### POST /auth/register
Register a new user account

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

#### POST /auth/login
Authenticate a user and get a token

**Request:**
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### User Management

#### GET /user
Get current user profile

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2023-01-01T00:00:00Z"
}
```

#### PUT /user
Update current user profile

**Request:**
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "John Smith",
  "email": "johnsmith@example.com"
}
```

### Anomalies

#### GET /anomalies
Get list of anomalies

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `status` (optional): "active", "resolved", "all" (default: "active")
- `severity` (optional): "low", "medium", "high", "critical"

#### GET /anomalies/:id
Get a specific anomaly by ID

#### PUT /anomalies/:id/resolve
Mark an anomaly as resolved

**Request:**
```json
{
  "resolution": "Anomaly resolved - was caused by scheduled maintenance"
}
```

#### DELETE /anomalies/:id
Delete an anomaly (admin only)

### Events

#### GET /events
Get list of events

**Query Parameters:**
- `page` (optional)
- `limit` (optional)
- `from` (optional): Start timestamp
- `to` (optional): End timestamp

#### POST /events/ingest
Ingest a new event

**Request:**
```json
{
  "timestamp": "2023-01-01T10:00:00Z",
  "type": "log",
  "source": "api-server-1",
  "data": {
    "level": "error",
    "message": "Database connection failed",
    "details": {
      "host": "db.example.com",
      "error": "timeout"
    }
  }
}
```

### Dashboard

#### GET /dashboard/stats
Get dashboard statistics

**Response:**
```json
{
  "total_anomalies": 42,
  "active_anomalies": 5,
  "events_processed": 12500,
  "avg_severity": "medium"
}
```

#### GET /dashboard/recent
Get recent anomalies

### Billing

#### GET /billing/plans
Get available subscription plans

**Response:**
```json
[
  {
    "id": "price_free",
    "name": "Free Plan",
    "description": "Basic plan with limited features",
    "price": 0,
    "currency": "usd",
    "interval": "month",
    "features": [...]
  }
]
```

#### POST /billing/checkout
Create a checkout session for a plan

**Request:**
```json
{
  "plan_id": "price_pro"
}
```

**Response:**
```json
{
  "session_id": "cs_test_...",
  "url": "https://checkout.stripe.com/pay/..."
}
```

#### GET /billing/subscription
Get current user's subscription

#### DELETE /billing/subscription
Cancel current subscription

#### GET /billing/usage
Get current usage

### Email

#### POST /email/test
Send a test email

**Request:**
```json
{
  "to": "test@example.com",
  "subject": "Test Email",
  "body": "This is a test email"
}
```

### Onboarding

#### GET /onboarding/progress
Get current onboarding progress

#### POST /onboarding/step/complete
Mark an onboarding step as complete

**Request:**
```json
{
  "step_name": "connect_data_source"
}
```

## Error Handling

All API errors return a JSON response with the following structure:

```json
{
  "error": "Error message",
  "code": "error_code"
}
```

## Rate Limiting

The API implements rate limiting:
- 100 requests per minute per IP for unauthenticated endpoints
- 1000 requests per minute per user for authenticated endpoints