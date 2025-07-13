# API Usage Guide

This guide provides practical examples and best practices for using the maxwash Backend API.

## Getting Started

### Base URL
```
http://localhost:8080
```

### Content Types
- **JSON**: `application/json` for most endpoints
- **Form Data**: `application/x-www-form-urlencoded` for Google OAuth

### Authentication
Protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer your_jwt_accessToken
```

## Authentication Examples

### 1. User Registration

**Request:**
```bash
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "John Doe",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenType": "Bearer",
  "expiresIn": 3600,
  "refreshToken": "abc123def456ghi789...",
  "expiresAt": "2024-01-01T12:00:00Z"
}
```

### 2. User Login

**Request:**
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

**Response:** Same as registration response

### 3. Google OAuth Login

**Request:**
```bash
curl -X POST http://localhost:8080/v1/auth/google \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id_token=your_google_id_token_here"
```

**Response:** Same as registration response

### 4. Token Refresh

**Request:**
```bash
curl -X POST "http://localhost:8080/v1/auth/refresh?refreshToken=your_refreshToken_here"
```

**Response:**
```json
{
  "accessToken": "new_accessToken_here...",
  "tokenType": "Bearer",
  "expiresIn": 3600,
  "expiresAt": "2024-01-01T12:00:00Z"
}
```

### 5. Logout

**Request:**
```bash
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Authorization: Bearer your_accessToken_here"
```

**Response:**
```json
{
  "message": "Successfully logged out"
}
```

## User Management Examples

### 1. Get All Users

**Request:**
```bash
curl -X GET "http://localhost:8080/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer your_accessToken_here"
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "name": "John Doe",
    "lastLoginAt": "2024-01-01T10:00:00Z",
    "lastLogoutAt": "2024-01-01T09:00:00Z",
    "createdAt": "2024-01-01T08:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
]
```

### 2. Get User by ID

**Request:**
```bash
curl -X GET "http://localhost:8080/v1/users/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer your_accessToken_here"
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "lastLoginAt": "2024-01-01T10:00:00Z",
  "lastLogoutAt": "2024-01-01T09:00:00Z",
  "createdAt": "2024-01-01T08:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

### 3. Update User

**Request:**
```bash
curl -X PUT "http://localhost:8080/v1/users/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer your_accessToken_here" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com"
  }'
```

**Response:** Updated user object

### 4. Delete User

**Request:**
```bash
curl -X DELETE "http://localhost:8080/v1/users/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer your_accessToken_here"
```

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

## Error Handling

### Common Error Responses

**400 Bad Request:**
```json
{
  "error": "Invalid request body: email is required"
}
```

**401 Unauthorized:**
```json
{
  "error": "Invalid credentials"
}
```

**404 Not Found:**
```json
{
  "error": "User not found"
}
```

**409 Conflict:**
```json
{
  "error": "User with this email already exists"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

## Best Practices

### 1. Authentication
- Always include the Bearer token in the Authorization header for protected endpoints
- Store tokens securely and refresh them before expiration
- Handle 401 responses by redirecting to login

### 2. Error Handling
- Always check the HTTP status code
- Parse error messages from the response body
- Implement retry logic for transient errors

### 3. Rate Limiting
- Respect rate limits (if implemented)
- Implement exponential backoff for retries
- Cache responses when appropriate

### 4. Data Validation
- Validate input data on the client side
- Handle validation errors gracefully
- Provide clear error messages to users

### 5. Security
- Never store passwords in plain text
- Use HTTPS in production
- Validate tokens on the server side
- Implement proper session management

## SDK Examples

### JavaScript/TypeScript

```typescript
class MaxwashAPI {
  private baseURL: string;
  private accessToken: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  async login(email: string, password: string) {
    const response = await fetch(`${this.baseURL}/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data = await response.json();
    this.accessToken = data.accessToken;
    return data;
  }

  async getUsers(page = 1, limit = 10) {
    if (!this.accessToken) {
      throw new Error('Not authenticated');
    }

    const response = await fetch(
      `${this.baseURL}/v1/users?page=${page}&limit=${limit}`,
      {
        headers: {
          'Authorization': `Bearer ${this.accessToken}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch users');
    }

    return response.json();
  }
}
```

### Python

```python
import requests
from typing import Optional, Dict, Any

class MaxwashAPI:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.accessToken: Optional[str] = None

    def login(self, email: str, password: str) -> Dict[str, Any]:
        response = requests.post(
            f"{self.base_url}/v1/auth/login",
            json={"email": email, "password": password}
        )
        response.raise_for_status()
        
        data = response.json()
        self.accessToken = data["accessToken"]
        return data

    def get_users(self, page: int = 1, limit: int = 10) -> Dict[str, Any]:
        if not self.accessToken:
            raise ValueError("Not authenticated")

        response = requests.get(
            f"{self.base_url}/v1/users",
            params={"page": page, "limit": limit},
            headers={"Authorization": f"Bearer {self.accessToken}"}
        )
        response.raise_for_status()
        
        return response.json()
```

## Testing

### Using curl for Testing

```bash
# Test health endpoint
curl http://localhost:8080/ping

# Test authentication
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","password":"password123"}'

# Test protected endpoint
curl -X GET http://localhost:8080/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Using Postman

1. Import the API collection
2. Set the base URL variable
3. Use the authentication flow to get tokens
4. Test protected endpoints with the token

## Troubleshooting

### Common Issues

1. **401 Unauthorized**: Check if the token is valid and not expired
2. **400 Bad Request**: Validate request body format and required fields
3. **500 Internal Server Error**: Check server logs for details
4. **CORS Issues**: Ensure proper CORS configuration for web clients

### Debug Tips

1. Use verbose curl output: `curl -v`
2. Check response headers for additional information
3. Log request/response data for debugging
4. Use browser developer tools for web applications 