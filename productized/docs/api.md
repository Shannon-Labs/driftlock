# DriftLock API Documentation

## Overview
DriftLock is an anomaly detection platform that monitors your systems, applications, and infrastructure for unusual behavior patterns. Our platform uses advanced algorithms to detect anomalies in logs, metrics, and traces, alerting you before issues impact your users.

## API Base URL
`https://api.driftlock.com/v1`

## Authentication
DriftLock uses JWT (JSON Web Tokens) for authentication. To access the API, you must include an Authorization header in your requests:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Getting a JWT Token
To obtain an access token, send a POST request to the authentication endpoint:

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "your-email@example.com",
  "password": "your-password"
}
```

Response:
```
{
  "user": {
    "id": 1,
    "email": "your-email@example.com",
    "name": "Your Name",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## API Endpoints

### Users
- `GET /api/v1/user` - Get current user profile
- `PUT /api/v1/user` - Update user profile

### Anomalies
- `GET /api/v1/anomalies` - Get list of anomalies
- `GET /api/v1/anomalies/{id}` - Get specific anomaly
- `PUT /api/v1/anomalies/{id}/resolve` - Mark anomaly as resolved
- `DELETE /api/v1/anomalies/{id}` - Delete anomaly

### Events
- `POST /api/v1/events/ingest` - Ingest events for anomaly detection
- `GET /api/v1/events` - Get list of events

### Dashboard
- `GET /api/v1/dashboard/stats` - Get dashboard statistics
- `GET /api/v1/dashboard/recent` - Get recent anomalies

## Response Format
All API responses follow the same structure:

```
{
  "data": { ... },
  "message": "Success message",
  "timestamp": "2023-07-21T14:30:00Z"
}
```

For errors:
```
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "timestamp": "2023-07-21T14:30:00Z"
}
```

## Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Rate Limiting
The API implements rate limiting to ensure fair usage. Standard accounts are limited to 1000 requests per hour. Enterprise accounts may have higher limits.

## Support
For API support, contact our team at api-support@driftlock.com