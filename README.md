# OTP Authentication Service

A comprehensive backend service in Go that implements OTP-based login and registration, along with basic user management features using Clean Architecture principles.

## Features

- **OTP-based Authentication**: Generate and verify OTP codes for user registration/login
- **Rate Limiting**: Prevent abuse with configurable rate limits
- **JWT Token Authentication**: Secure token-based authentication
- **User Management**: CRUD operations for user management
- **RESTful API**: Complete REST API with Swagger documentation
- **Database Integration**: PostgreSQL with automatic migrations
- **Docker Support**: Full containerization with Docker Compose
- **Clean Architecture**: Well-structured, maintainable codebase

## ✅ Task Requirements Checklist

### 1. OTP Login & Registration ✅
- [x] **User sends phone number** → system generates random OTP
- [x] **OTP printed to console** (no SMS sending)
- [x] **OTP stored temporarily** in database
- [x] **OTP expires after 2 minutes**
- [x] **User submits phone number + OTP**
- [x] **Register new user if not existing**
- [x] **Log in existing user otherwise**
- [x] **Return JWT token upon success**

### 2. Rate Limiting ✅
- [x] **Max 3 requests per phone number within 10 minutes**
- [x] **Database-based rate limiting** (persistent across restarts)
- [x] **Proper HTTP 429 responses**

### 3. User Management ✅
- [x] **REST endpoints for user operations**
- [x] **Retrieve single user details** (`GET /api/v1/users/{id}`)
- [x] **Retrieve list of users with pagination**
- [x] **Search functionality** (by phone number)
- [x] **Store minimum required fields**: phone number, registration date

### 4. Database ✅
- [x] **PostgreSQL database** (chosen for ACID compliance and reliability)
- [x] **Docker Compose setup** with database included
- [x] **Automatic migrations** on startup
- [x] **Proper indexing** for performance

### 5. API Documentation ✅
- [x] **All operations exposed via REST APIs**
- [x] **Swagger/OpenAPI documentation**
- [x] **Interactive API documentation** at `/swagger/index.html`
- [x] **Example requests and responses**

### 6. Architecture & Best Practices ✅
- [x] **Clean Architecture** implementation
- [x] **Clear separation of responsibilities**
- [x] **Maintainable and testable code**
- [x] **Interface-based design**

### 7. Containerization ✅
- [x] **Application Dockerized**
- [x] **Database included in docker-compose**
- [x] **Multi-stage Docker build**
- [x] **Security-hardened containers**

### 8. Deliverables ✅
- [x] **Complete source code**
- [x] **Comprehensive documentation**
- [x] **Local running instructions**
- [x] **Docker running instructions**
- [x] **Example API requests & responses**
- [x] **Database choice justification**

## Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware (auth, CORS)
│   ├── models/         # Domain models and DTOs
│   ├── repository/     # Data access layer
│   └── services/       # Business logic layer
├── docs/              # Generated Swagger documentation
├── Dockerfile         # Application containerization
└── docker-compose.yml # Multi-service orchestration
```

## Technology Stack

- **Language**: Go 1.21
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Architecture**: Clean Architecture

## Database Choice Justification

**PostgreSQL** was chosen for the following reasons:

1. **ACID Compliance**: Ensures data integrity for user authentication
2. **Concurrent Access**: Excellent support for multiple simultaneous connections
3. **JSON Support**: Native JSON data types for flexible schema evolution
4. **Performance**: Optimized for read-heavy workloads with proper indexing
5. **Reliability**: Battle-tested database with excellent community support
6. **Docker Integration**: Easy to containerize and manage

## Quick Start

### Prerequisites

- **For Docker**: Docker and Docker Compose
- **For Local Development**: 
  - Go 1.21 or later
  - PostgreSQL 15 or later

### Option 1: Running with Docker (Recommended)

This is the easiest way to get started as it includes the database setup.

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd otp
   ```

2. **Start all services**:
   ```bash
   # Using newer Docker Compose syntax (recommended)
   docker compose up -d
   
   # Or using legacy syntax
   docker-compose up -d
   ```

3. **Verify services are running**:
   ```bash
   # Using newer Docker Compose syntax
   docker compose ps
   
   # Or using legacy syntax
   docker-compose ps
   ```

4. **Access the application**:
   - **API Base URL**: http://localhost:8080
   - **Swagger Documentation**: http://localhost:8080/swagger/index.html
   - **Health Check**: http://localhost:8080/health

5. **View logs** (optional):
   ```bash
   # Using newer Docker Compose syntax
   docker compose logs -f app
   
   # Or using legacy syntax
   docker-compose logs -f app
   ```

6. **Stop services**:
   ```bash
   # Using newer Docker Compose syntax
   docker compose down
   
   # Or using legacy syntax
   docker-compose down
   ```

### Option 2: Running Locally

For development or if you prefer to run without Docker.

1. **Set up environment variables**:
   ```bash
   cp env.example .env
   # Edit .env with your database configuration
   ```

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**:
   ```bash
   # Install PostgreSQL (Ubuntu/Debian)
   sudo apt-get install postgresql postgresql-contrib
   
   # Or on macOS with Homebrew
   brew install postgresql
   
   # Start PostgreSQL service
   sudo systemctl start postgresql  # Linux
   brew services start postgresql   # macOS
   
   # Create database and user
   sudo -u postgres psql
   CREATE DATABASE otp_db;
   CREATE USER otp_user WITH PASSWORD 'otp_password';
   GRANT ALL PRIVILEGES ON DATABASE otp_db TO otp_user;
   \q
   ```

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the application**:
   - **API Base URL**: http://localhost:8080
   - **Swagger Documentation**: http://localhost:8080/swagger/index.html
   - **Health Check**: http://localhost:8080/health

### Option 3: Using Makefile Commands

The project includes a Makefile with convenient commands:

```bash
# Build the application
make build

# Run locally
make run

# Run tests
make test

# Run with Docker
make docker-run

# Stop Docker services
make docker-stop

# View Docker logs
make docker-logs

# Generate Swagger docs
make swagger

# Health check
make health
```

### First Steps After Starting

1. **Check if the service is running**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **Generate an OTP**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
     -H "Content-Type: application/json" \
     -d '{"phone_number": "+1234567890"}'
   ```

3. **Check the console output** for the OTP code (e.g., "OTP for +1234567890: 123456")

4. **Verify the OTP**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
     -H "Content-Type: application/json" \
     -d '{"phone_number": "+1234567890", "code": "123456"}'
   ```

5. **Explore the API** using Swagger UI at http://localhost:8080/swagger/index.html

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/otp/generate` | Generate OTP for phone number | No |
| POST | `/api/v1/auth/otp/verify` | Verify OTP and authenticate user | No |

### User Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | List users with pagination and search | Yes |
| GET | `/api/v1/users/{id}` | Get user by ID | Yes |
| DELETE | `/api/v1/users/{id}` | Delete user by ID | Yes |

### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |
| GET | `/swagger/*` | Swagger documentation |

## Example API Requests

### 1. Generate OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890"
  }'
```

**Response**:
```json
{
  "message": "OTP sent successfully",
  "expires_in_minutes": 2
}
```

### 2. Verify OTP and Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890",
    "code": "123456"
  }'
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-here",
    "phone_number": "+1234567890",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login_at": "2024-01-01T00:00:00Z"
  },
  "expires_at": "2024-01-02T00:00:00Z"
}
```

### 3. List Users (Authenticated)

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10&search=123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "users": [
    {
      "id": "uuid-here",
      "phone_number": "+1234567890",
      "created_at": "2024-01-01T00:00:00Z",
      "last_login_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### 4. Get User by ID (Authenticated)

```bash
curl -X GET http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `otp_user` | Database user |
| `DB_PASSWORD` | `otp_password` | Database password |
| `DB_NAME` | `otp_db` | Database name |
| `JWT_SECRET` | `your-super-secret-jwt-key-change-in-production` | JWT signing secret |
| `JWT_EXPIRY_HOURS` | `24` | JWT token expiry in hours |
| `OTP_EXPIRY_MINUTES` | `2` | OTP expiry in minutes |
| `OTP_LENGTH` | `6` | OTP code length |
| `RATE_LIMIT_MAX_REQUESTS` | `3` | Max OTP requests per window |
| `RATE_LIMIT_WINDOW_MINUTES` | `10` | Rate limit window in minutes |

## Rate Limiting

The service implements rate limiting for OTP generation:
- **Limit**: 3 requests per phone number
- **Window**: 10 minutes
- **Storage**: Database-based (persistent across restarts)

## Security Features

1. **JWT Authentication**: Secure token-based authentication
2. **Rate Limiting**: Prevents OTP abuse
3. **Input Validation**: Comprehensive request validation
4. **CORS Support**: Configurable cross-origin requests
5. **Non-root Container**: Security-hardened Docker container
6. **Environment-based Configuration**: Secure configuration management

## Development

### Project Structure

```
otp/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database operations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Domain models
│   ├── repository/      # Data access layer
│   └── services/        # Business logic
├── docs/               # Generated documentation
├── Dockerfile          # Container definition
├── docker-compose.yml  # Service orchestration
├── go.mod             # Go module file
└── README.md          # This file
```

### Adding New Features

1. **Domain Models**: Add to `internal/models/`
2. **Repository Layer**: Add to `internal/repository/`
3. **Business Logic**: Add to `internal/services/`
4. **HTTP Handlers**: Add to `internal/handlers/`
5. **Routes**: Update `cmd/server/main.go`

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/services -v
```

## Deployment

### Production Considerations

1. **Environment Variables**: Use proper secrets management
2. **Database**: Use managed PostgreSQL service
3. **SSL/TLS**: Enable HTTPS in production
4. **Monitoring**: Add health checks and logging
5. **Backup**: Implement database backup strategy

### Docker Deployment

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Failed
**Symptoms**: Application fails to start with database connection errors

**Solutions**:
```bash
# For Docker setup
docker compose down -v  # Remove volumes
docker compose up -d    # Restart fresh

# For local setup
# Check if PostgreSQL is running
sudo systemctl status postgresql  # Linux
brew services list | grep postgresql  # macOS

# Verify database exists
psql -h localhost -U otp_user -d otp_db -c "SELECT 1;"
```

#### 2. OTP Not Working
**Symptoms**: OTP generation fails or verification doesn't work

**Solutions**:
```bash
# Check console output for OTP codes
docker compose logs app | grep "OTP for"

# Check rate limiting
# If you get "rate limit exceeded", wait 10 minutes or use a different phone number

# Verify database connectivity
docker compose exec app ./main --health-check
```

#### 3. JWT Token Issues
**Symptoms**: Authentication fails with token errors

**Solutions**:
```bash
# Check JWT secret is set
echo $JWT_SECRET

# Verify token format
# Token should be: "Bearer <your-jwt-token>"

# Check token expiration
# Default expiry is 24 hours
```

#### 4. Port Already in Use
**Symptoms**: Application fails to start with port binding errors

**Solutions**:
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Kill the process or change port in .env
# SERVER_PORT=8081
```

#### 5. Docker Issues
**Symptoms**: Docker containers fail to start

**Solutions**:
```bash
# Clean up Docker resources
docker compose down -v
docker system prune -f

# Rebuild containers
docker compose build --no-cache
docker compose up -d
```

### Logs and Debugging

#### View Application Logs
```bash
# Docker logs
docker compose logs app
docker compose logs -f app  # Follow logs

# Local logs (if using log file)
tail -f logs/app.log
```

#### View Database Logs
```bash
# Docker database logs
docker compose logs postgres
docker compose logs -f postgres

# Local PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log
```

#### Health Checks
```bash
# Application health
curl http://localhost:8080/health

# Database health (Docker)
docker compose exec postgres pg_isready -U otp_user -d otp_db

# Database health (Local)
pg_isready -h localhost -U otp_user -d otp_db
```

### Performance Issues

#### High Memory Usage
```bash
# Check container resource usage
docker stats

# Monitor database connections
docker compose exec postgres psql -U otp_user -d otp_db -c "SELECT count(*) FROM pg_stat_activity;"
```

#### Slow Queries
```bash
# Check database indexes
docker compose exec postgres psql -U otp_user -d otp_db -c "\d+ users"
docker compose exec postgres psql -U otp_user -d otp_db -c "\d+ otps"
```

### Development Tips

#### Reset Everything
```bash
# Complete reset (Docker)
docker compose down -v
docker system prune -f
docker compose up -d

# Complete reset (Local)
dropdb otp_db
createdb otp_db
go run cmd/server/main.go
```

#### Test API Endpoints
```bash
# Health check
curl http://localhost:8080/health

# Generate OTP
curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'

# Check Swagger docs
open http://localhost:8080/swagger/index.html
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Verification Checklist

Before submitting or deploying, verify that all requirements are met:

### ✅ Core Functionality Verification
- [ ] **OTP Generation**: `curl -X POST http://localhost:8080/api/v1/auth/otp/generate -H "Content-Type: application/json" -d '{"phone_number": "+1234567890"}'`
- [ ] **OTP Console Output**: Check application logs for "OTP for +1234567890: XXXXXX"
- [ ] **OTP Verification**: Use the generated OTP to verify and get JWT token
- [ ] **Rate Limiting**: Try generating OTP 4 times within 10 minutes (should get 429 error)
- [ ] **User Management**: List users with JWT token authentication
- [ ] **Database**: Verify tables are created and data is persisted

### ✅ Technical Requirements Verification
- [ ] **Clean Architecture**: Code follows layered architecture (handlers → services → repositories → models)
- [ ] **RESTful API**: All endpoints follow REST conventions
- [ ] **Swagger Documentation**: Accessible at http://localhost:8080/swagger/index.html
- [ ] **Docker**: Application runs successfully with `docker compose up -d`
- [ ] **Database**: PostgreSQL is running and accessible
- [ ] **Tests**: All tests pass with `go test ./...`

### ✅ Production Readiness
- [ ] **Security**: JWT secret is configurable via environment variables
- [ ] **Error Handling**: Proper error responses for all scenarios
- [ ] **Logging**: Application logs are available
- [ ] **Health Checks**: Health endpoint responds correctly
- [ ] **Configuration**: All settings are environment-based

### Quick Verification Commands
```bash
# 1. Start the application
docker compose up -d

# 2. Check health
curl http://localhost:8080/health

# 3. Generate OTP
curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'

# 4. Check logs for OTP
docker compose logs app | grep "OTP for"

# 5. Verify OTP (replace XXXXXX with actual OTP)
curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "code": "XXXXXX"}'

# 6. Test rate limiting (should fail after 3 attempts)
for i in {1..4}; do
  curl -X POST http://localhost:8080/api/v1/auth/otp/generate \
    -H "Content-Type: application/json" \
    -d '{"phone_number": "+1234567891"}'
  echo ""
done

# 7. Check Swagger docs
open http://localhost:8080/swagger/index.html
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.