# API Usage Guide

This guide provides practical examples and best practices for using the Jeki Backend API.

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
Authorization: Bearer your_jwt_access_token
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
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "abc123def456ghi789...",
  "expires_at": "2024-01-01T12:00:00Z"
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
curl -X POST "http://localhost:8080/v1/auth/refresh?refresh_token=your_refresh_token_here"
```

**Response:**
```json
{
  "access_token": "new_access_token_here...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "expires_at": "2024-01-01T12:00:00Z"
}
```

### 5. Logout

**Request:**
```bash
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Authorization: Bearer your_access_token_here"
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
curl -X GET "http://localhost:8080/v1/users?page=1&limit=10"
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### 2. Get User by ID

**Request:**
```bash
curl -X GET http://localhost:8080/v1/users/123e4567-e89b-12d3-a456-426614174000
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### 3. Update User

**Request:**
```bash
curl -X PUT http://localhost:8080/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com"
  }'
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.smith@example.com",
  "name": "John Smith",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T11:00:00Z"
}
```

### 4. Update User Password

**Request:**
```bash
curl -X PUT http://localhost:8080/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "password": "newsecurepassword456"
  }'
```

### 5. Delete User

**Request:**
```bash
curl -X DELETE http://localhost:8080/v1/users/123e4567-e89b-12d3-a456-426614174000
```

**Response:** `204 No Content`

## JavaScript/Node.js Examples

### Using Fetch API

```javascript
// User Registration
async function registerUser(email, name, password) {
  const response = await fetch('http://localhost:8080/v1/auth/signup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email,
      name,
      password
    })
  });
  
  if (!response.ok) {
    throw new Error('Registration failed');
  }
  
  return await response.json();
}

// User Login
async function loginUser(email, password) {
  const response = await fetch('http://localhost:8080/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email,
      password
    })
  });
  
  if (!response.ok) {
    throw new Error('Login failed');
  }
  
  return await response.json();
}

// Get User Profile (Protected Route)
async function getUserProfile(userId, accessToken) {
  const response = await fetch(`http://localhost:8080/v1/users/${userId}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to get user profile');
  }
  
  return await response.json();
}

// Logout
async function logoutUser(accessToken) {
  const response = await fetch('http://localhost:8080/v1/auth/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  if (!response.ok) {
    throw new Error('Logout failed');
  }
  
  return await response.json();
}
```

### Using Axios

```javascript
import axios from 'axios';

// Configure base URL
const api = axios.create({
  baseURL: 'http://localhost:8080'
});

// Add auth interceptor
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// User registration
const registerUser = async (userData) => {
  try {
    const response = await api.post('/v1/auth/signup', userData);
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Registration failed');
  }
};

// User login
const loginUser = async (credentials) => {
  try {
    const response = await api.post('/v1/auth/login', credentials);
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Login failed');
  }
};

// Get all users
const getUsers = async (page = 1, limit = 10) => {
  try {
    const response = await api.get(`/v1/users?page=${page}&limit=${limit}`);
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Failed to get users');
  }
};
```

## Error Handling

### Common Error Responses

```json
// 400 Bad Request
{
  "error": "Invalid request body: email is required"
}

// 401 Unauthorized
{
  "error": "Invalid credentials"
}

// 404 Not Found
{
  "error": "User not found"
}

// 409 Conflict
{
  "error": "user with this email already exists"
}

// 500 Internal Server Error
{
  "error": "Internal server error"
}
```

### Error Handling Best Practices

1. **Always check response status codes**
2. **Parse error messages from response body**
3. **Handle token expiration gracefully**
4. **Implement retry logic for network errors**
5. **Store tokens securely (localStorage, secure cookies, etc.)**

## Security Best Practices

1. **Token Storage**
   - Store access tokens in memory when possible
   - Use secure cookies for refresh tokens
   - Never store tokens in localStorage for production apps

2. **Token Refresh**
   - Implement automatic token refresh before expiration
   - Handle refresh token expiration gracefully
   - Redirect to login when refresh fails

3. **Input Validation**
   - Validate all user inputs on client side
   - Sanitize data before sending to API
   - Use HTTPS in production

4. **Error Handling**
   - Don't expose sensitive information in error messages
   - Log errors appropriately
   - Implement proper error boundaries

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Authentication endpoints**: 5 requests per minute per IP
- **User management endpoints**: 100 requests per minute per user
- **General endpoints**: 1000 requests per minute per IP

## Testing

### Test Credentials

For testing purposes, you can use these credentials:
- **Email**: `test@example.com`
- **Password**: `testpassword123`

### Test Environment

The API includes a test environment with:
- In-memory database for testing
- Mock authentication for development
- Sample data generation
- Health check endpoint at `/ping`

## Support

For API support and questions:
- Check the interactive documentation at `/swagger/index.html`
- Review the authentication flow documentation
- Create an issue in the repository 