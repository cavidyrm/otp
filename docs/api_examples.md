# API Examples

This document provides practical examples of how to use the OTP Authentication API.

## Prerequisites

- The service is running on `http://localhost:8080`
- You have `curl` installed (or use any HTTP client)

## Authentication Flow

### 1. Generate OTP

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890"
  }'
```

**Response:**
```json
{
  "message": "OTP sent successfully",
  "expires_in_minutes": 2
}
```

**Note:** The OTP code will be printed to the console where the server is running.

### 2. Verify OTP and Get JWT Token

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890",
    "code": "123456"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2Nzg5MCIsInBob25lX251bWJlciI6IisxMjM0NTY3ODkwIiwiZXhwIjoxNzA0MDY0MDAwfQ.example",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "phone_number": "+1234567890",
    "created_at": "2024-01-01T12:00:00Z",
    "last_login_at": "2024-01-01T12:00:00Z"
  },
  "expires_at": "2024-01-02T12:00:00Z"
}
```

## User Management

### 3. List Users (Requires Authentication)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "phone_number": "+1234567890",
      "created_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### 4. Search Users

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/users?search=123&page=1&page_size=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Get User by ID

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "phone_number": "+1234567890",
  "created_at": "2024-01-01T12:00:00Z",
  "last_login_at": "2024-01-01T12:00:00Z"
}
```

### 6. Delete User

**Request:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

## System Endpoints

### 7. Health Check

**Request:**
```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Error Responses

### Rate Limit Exceeded

**Response:**
```json
{
  "error": "rate limit exceeded. Please try again later"
}
```

### Invalid OTP

**Response:**
```json
{
  "error": "invalid OTP code"
}
```

### Expired OTP

**Response:**
```json
{
  "error": "OTP has expired"
}
```

### Unauthorized

**Response:**
```json
{
  "error": "Authorization header is required"
}
```

### User Not Found

**Response:**
```json
{
  "error": "User not found"
}
```

## Complete Authentication Flow Example

Here's a complete example of the authentication flow:

```bash
# 1. Generate OTP
curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'

# 2. Check console output for OTP code (e.g., "123456")

# 3. Verify OTP and get token
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "code": "123456"}')

# 4. Extract token (requires jq)
TOKEN=$(echo $RESPONSE | jq -r '.token')

# 5. Use token for authenticated requests
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Authorization: Bearer $TOKEN"
```

## Testing with Different Tools

### Using Postman

1. Import the following collection:
```json
{
  "info": {
    "name": "OTP Authentication API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Generate OTP",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"phone_number\": \"+1234567890\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/auth/otp/generate",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "auth", "otp", "generate"]
        }
      }
    }
  ]
}
```

### Using JavaScript/Fetch

```javascript
// Generate OTP
const generateOTP = async (phoneNumber) => {
  const response = await fetch('http://localhost:8080/api/v1/auth/otp/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ phone_number: phoneNumber }),
  });
  return response.json();
};

// Verify OTP
const verifyOTP = async (phoneNumber, code) => {
  const response = await fetch('http://localhost:8080/api/v1/auth/otp/verify', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ phone_number: phoneNumber, code: code }),
  });
  return response.json();
};

// Get users (authenticated)
const getUsers = async (token) => {
  const response = await fetch('http://localhost:8080/api/v1/users', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return response.json();
};
```

## Rate Limiting

The API implements rate limiting for OTP generation:
- **Limit**: 3 requests per phone number
- **Window**: 10 minutes
- **HTTP Status**: 429 (Too Many Requests)

Example of rate limit response:
```json
{
  "error": "rate limit exceeded. Please try again later"
}
```
