# Architecture Documentation

## Overview

This OTP Authentication Service is built using **Clean Architecture** principles, ensuring separation of concerns, testability, and maintainability. The application follows a layered architecture pattern with clear boundaries between different components.

## Architecture Layers

### 1. Presentation Layer (Handlers)
**Location**: `internal/handlers/`

The presentation layer handles HTTP requests and responses. It's responsible for:
- Request validation
- Response formatting
- HTTP status codes
- Content-Type headers

**Components**:
- `auth_handler.go`: Handles OTP generation and verification
- `user_handler.go`: Handles user management operations

**Key Features**:
- Swagger/OpenAPI documentation annotations
- Input validation using Gin's binding
- Consistent error response format
- CORS support

### 2. Business Logic Layer (Services)
**Location**: `internal/services/`

The business logic layer contains the core application logic. It's responsible for:
- Business rules implementation
- Data transformation
- Orchestrating operations between repositories
- Authentication and authorization logic

**Components**:
- `auth_service.go`: OTP generation, verification, and JWT token management
- `user_service.go`: User management operations

**Key Features**:
- Rate limiting implementation
- OTP generation and validation
- JWT token creation and validation
- User registration and authentication

### 3. Data Access Layer (Repositories)
**Location**: `internal/repository/`

The data access layer abstracts database operations. It's responsible for:
- Database queries
- Data persistence
- Data retrieval
- Transaction management

**Components**:
- `user_repository.go`: User CRUD operations
- `otp_repository.go`: OTP storage and retrieval

**Key Features**:
- Interface-based design for testability
- PostgreSQL-specific implementations
- Pagination support
- Search functionality

### 4. Domain Layer (Models)
**Location**: `internal/models/`

The domain layer contains the core business entities and rules. It's responsible for:
- Data structures
- Business entity definitions
- Validation rules
- Domain-specific logic

**Components**:
- `user.go`: User entity and related structures
- `otp.go`: OTP entity and related structures
- `auth.go`: Authentication-related structures

**Key Features**:
- JWT claims implementation
- Pagination query structures
- Response DTOs
- Domain validation methods

## Database Design

### Schema Overview

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_login_at TIMESTAMP
);

-- OTPs table
CREATE TABLE otps (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE
);

-- Indexes for performance
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_otps_phone_number ON otps(phone_number);
CREATE INDEX idx_otps_expires_at ON otps(expires_at);
CREATE INDEX idx_otps_created_at ON otps(created_at);
```

### Database Choice: PostgreSQL

**Why PostgreSQL?**

1. **ACID Compliance**: Essential for user authentication data integrity
2. **Concurrent Access**: Excellent support for multiple simultaneous connections
3. **JSON Support**: Native JSON data types for future extensibility
4. **Performance**: Optimized for read-heavy workloads with proper indexing
5. **Reliability**: Battle-tested database with excellent community support
6. **Docker Integration**: Easy to containerize and manage

## Security Implementation

### 1. JWT Authentication
- **Algorithm**: HMAC-SHA256
- **Claims**: User ID, phone number, expiration time
- **Expiry**: Configurable (default: 24 hours)
- **Secret**: Environment variable (change in production)

### 2. Rate Limiting
- **Storage**: Database-based (persistent across restarts)
- **Limit**: 3 requests per phone number
- **Window**: 10 minutes
- **HTTP Status**: 429 (Too Many Requests)

### 3. Input Validation
- **Request Binding**: Gin's built-in validation
- **Phone Number**: Required field validation
- **OTP Code**: Required field validation
- **Pagination**: Min/max value constraints

### 4. CORS Support
- **Origin**: Configurable (default: all origins)
- **Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Headers**: Authorization, Content-Type, etc.

## Configuration Management

### Environment Variables
The application uses environment variables for configuration:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=otp_user
DB_PASSWORD=otp_password
DB_NAME=otp_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24

# OTP Configuration
OTP_EXPIRY_MINUTES=2
OTP_LENGTH=6

# Rate Limiting
RATE_LIMIT_MAX_REQUESTS=3
RATE_LIMIT_WINDOW_MINUTES=10
```

### Configuration Structure
```go
type Config struct {
    Server     ServerConfig
    Database   DatabaseConfig
    JWT        JWTConfig
    OTP        OTPConfig
    RateLimit  RateLimitConfig
}
```

## API Design

### RESTful Endpoints

#### Authentication
- `POST /api/v1/auth/otp/generate` - Generate OTP
- `POST /api/v1/auth/otp/verify` - Verify OTP and authenticate

#### User Management
- `GET /api/v1/users` - List users (with pagination and search)
- `GET /api/v1/users/{id}` - Get user by ID
- `DELETE /api/v1/users/{id}` - Delete user

#### System
- `GET /health` - Health check
- `GET /swagger/*` - API documentation

### Response Format
All API responses follow a consistent format:

**Success Response**:
```json
{
  "data": {...},
  "message": "Success message"
}
```

**Error Response**:
```json
{
  "error": "Error message"
}
```

## Testing Strategy

### Test Structure
- **Unit Tests**: Individual component testing
- **Integration Tests**: Service layer testing with mocked repositories
- **Mock Objects**: Repository interfaces for isolated testing

### Test Coverage
- **Services**: Core business logic testing
- **Handlers**: HTTP request/response testing
- **Models**: Domain logic testing

### Example Test Pattern
```go
func TestAuthService_GenerateOTP(t *testing.T) {
    // Setup
    cfg := &config.Config{...}
    userRepo := &mockUserRepository{...}
    otpRepo := &mockOTPRepository{...}
    authService := NewAuthService(userRepo, otpRepo, cfg)

    // Execute
    response, err := authService.GenerateOTP(ctx, phoneNumber)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response)
    // ... more assertions
}
```

## Deployment Architecture

### Docker Containerization
- **Multi-stage Build**: Optimized for production
- **Non-root User**: Security hardening
- **Alpine Linux**: Minimal attack surface
- **Health Checks**: Built-in health monitoring

### Docker Compose
- **Application**: Go service
- **Database**: PostgreSQL 15
- **Networking**: Isolated network
- **Volumes**: Persistent data storage

### Production Considerations
1. **Environment Variables**: Use proper secrets management
2. **Database**: Use managed PostgreSQL service
3. **SSL/TLS**: Enable HTTPS in production
4. **Monitoring**: Add health checks and logging
5. **Backup**: Implement database backup strategy

## Performance Considerations

### Database Optimization
- **Indexes**: Strategic indexing for common queries
- **Connection Pooling**: Configured connection limits
- **Query Optimization**: Efficient SQL queries

### Application Optimization
- **Goroutines**: Concurrent request handling
- **Memory Management**: Efficient data structures
- **Caching**: Consider Redis for session management

### Scalability
- **Horizontal Scaling**: Stateless application design
- **Load Balancing**: Multiple service instances
- **Database Scaling**: Read replicas for read-heavy workloads

## Error Handling

### Error Types
1. **Validation Errors**: Invalid input data
2. **Authentication Errors**: Invalid credentials
3. **Authorization Errors**: Insufficient permissions
4. **Rate Limit Errors**: Too many requests
5. **System Errors**: Internal server errors

### Error Response Format
```json
{
  "error": "Human-readable error message",
  "code": "ERROR_CODE",
  "details": {...}
}
```

## Monitoring and Logging

### Logging Strategy
- **Structured Logging**: JSON format for easy parsing
- **Log Levels**: DEBUG, INFO, WARN, ERROR
- **Context Information**: Request ID, user ID, etc.

### Health Checks
- **Database Connectivity**: Verify database connection
- **Service Health**: Application status
- **Dependencies**: External service availability

## Future Enhancements

### Potential Improvements
1. **Redis Integration**: Session management and caching
2. **SMS Integration**: Real SMS delivery for OTP
3. **Email Support**: Email-based authentication
4. **Multi-factor Authentication**: Additional security layers
5. **Audit Logging**: Comprehensive activity tracking
6. **API Versioning**: Backward compatibility support

### Scalability Features
1. **Microservices**: Service decomposition
2. **Message Queues**: Asynchronous processing
3. **Event Sourcing**: Audit trail and event replay
4. **CQRS**: Command Query Responsibility Segregation

## Conclusion

This architecture provides a solid foundation for a production-ready OTP authentication service. The clean separation of concerns, comprehensive testing, and security considerations make it suitable for enterprise use. The modular design allows for easy extension and maintenance as requirements evolve.
