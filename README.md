# Product Management API

A comprehensive REST API built with Go (Golang) for managing products, featuring authentication, authorization, and following Clean Architecture principles.

## 🚀 Features

### Core Functionality
- **Product Management**: Complete CRUD operations for products
- **User Authentication**: JWT-based authentication system
- **Role-based Authorization**: Admin and user roles
- **RESTful API**: Well-structured endpoints following REST principles

### Technical Features
- **Clean Architecture**: Organized code structure with clear separation of concerns
- **Microservice Ready**: Designed to be easily extended into microservices
- **Database Integration**: PostgreSQL with GORM ORM
- **Middleware**: Logging, CORS, Recovery, and Authentication middleware
- **API Documentation**: Swagger/OpenAPI integration
- **Docker Support**: Full containerization with Docker Compose
- **Testing**: Comprehensive unit and integration tests

### Security
- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt for password security
- **Input Validation**: Request validation and sanitization
- **CORS Configuration**: Configurable cross-origin resource sharing

## 🏗️ Architecture

The project follows **Clean Architecture** principles with the following structure:

```
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domain/           # Domain layer (entities, repositories, services)
│   │   ├── entity/       # Domain entities
│   │   ├── repository/   # Repository interfaces
│   │   └── service/      # Service interfaces
│   ├── infrastructure/   # Infrastructure layer
│   │   ├── database/     # Database connection and setup
│   │   └── repository/   # Repository implementations
│   ├── interfaces/       # Interface layer
│   │   └── http/         # HTTP handlers, middleware, routing
│   └── usecase/         # Use case implementations
├── pkg/                  # Public packages
│   └── jwt/             # JWT utilities
├── docs/                # Generated API documentation
├── scripts/             # Database and deployment scripts
└── test/               # Tests
    ├── mocks/          # Mock implementations
    └── unit/           # Unit tests
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Documentation**: Swagger/OpenAPI
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose
- **Environment**: Godotenv

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Docker and Docker Compose (optional)
- Make (optional, for using Makefile)

## 🚀 Quick Start

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd product-management
   ```

2. **Start all services**:
   ```bash
   make docker-compose-up
   # or
   docker-compose up -d
   ```

3. **Access the application**:
   - API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html
   - pgAdmin: http://localhost:5050 (admin@example.com / admin123)

### Option 2: Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd product-management
   make setup
   ```

2. **Install development tools**:
   ```bash
   make install-tools
   ```

3. **Setup PostgreSQL database**:
   ```bash
   # Create database
   createdb product_management
   # Or use the Makefile
   make db-create
   ```

4. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. **Generate Swagger docs and run**:
   ```bash
   make dev
   # or
   make swagger && make run
   ```

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/v1/auth/register` | Register a new user |
| POST   | `/api/v1/auth/login` | Login user |
| GET    | `/api/v1/auth/profile` | Get user profile |
| PUT    | `/api/v1/auth/profile` | Update user profile |
| POST   | `/api/v1/auth/change-password` | Change password |
| POST   | `/api/v1/auth/refresh` | Refresh token |
| POST   | `/api/v1/auth/logout` | Logout user |

### Product Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET    | `/api/v1/products` | List products | No |
| GET    | `/api/v1/products/{id}` | Get product by ID | No |
| GET    | `/api/v1/products/search` | Search products | No |
| POST   | `/api/v1/products` | Create product | Yes |
| PUT    | `/api/v1/products/{id}` | Update product | Yes |
| DELETE | `/api/v1/products/{id}` | Delete product | Yes |
| PUT    | `/api/v1/products/{id}/stock` | Update stock | Yes |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check |
| GET    | `/ready` | Readiness check |
| GET    | `/live` | Liveness check |

### Example Requests

**Register User:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "user123",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Create Product:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Awesome Product",
    "description": "This is an awesome product",
    "price": 29.99,
    "stock": 100,
    "category": "Electronics",
    "image_url": "https://example.com/image.jpg"
  }'
```

**List Products with Filters:**
```bash
curl "http://localhost:8080/api/v1/products?page=1&page_size=10&category=Electronics&min_price=10&max_price=100&is_active=true"
```

## 🧪 Testing

### Run All Tests
```bash
make test
```

### Run Unit Tests Only
```bash
make test-unit
```

### Run Integration Tests Only
```bash
make test-integration
```

### Generate Test Coverage Report
```bash
make test-coverage
```

## 🚀 Deployment

### Using Docker

1. **Build the image**:
   ```bash
   make docker-build
   ```

2. **Run with environment variables**:
   ```bash
   docker run -p 8080:8080 --env-file .env product-management:latest
   ```

### Using Docker Compose

The `docker-compose.yml` includes:
- **API service**: The main application
- **PostgreSQL**: Database service
- **Redis**: For caching (optional)
- **pgAdmin**: Database management tool

```bash
# Start all services
make docker-compose-up

# Stop all services
make docker-compose-down

# View logs
make docker-compose-logs
```

## ⚙️ Configuration

Configuration is managed through environment variables. See `.env.example` for all available options:

### Database Configuration
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

### JWT Configuration
- `JWT_SECRET`: Secret key for JWT signing
- `JWT_EXPIRES_IN`: Token expiration time

### Server Configuration
- `PORT`: Application port
- `GIN_MODE`: Gin mode (debug, release, test)

### CORS Configuration
- `ALLOWED_ORIGINS`: Allowed origins for CORS
- `ALLOWED_METHODS`: Allowed HTTP methods
- `ALLOWED_HEADERS`: Allowed headers

## 📊 Monitoring and Observability

### Health Checks
- **Health**: `/health` - Overall system health
- **Readiness**: `/ready` - Service readiness
- **Liveness**: `/live` - Service liveness

### Logging
- Structured JSON logging
- Request/response logging middleware
- Error tracking and recovery

### Metrics
The application is designed to easily integrate with monitoring tools like Prometheus and observability platforms.

## 🔧 Development

### Available Make Commands

```bash
make help           # Show all available commands
make setup          # Setup development environment
make run            # Run the application
make test           # Run all tests
make swagger        # Generate Swagger documentation
make docker-build   # Build Docker image
make clean          # Clean build artifacts
```

### Code Quality

- **Linting**: Use `golangci-lint` for code quality checks
- **Formatting**: Use `go fmt` for code formatting
- **Testing**: Comprehensive unit and integration tests
- **Documentation**: Swagger/OpenAPI documentation

### Project Structure Guidelines

- **Domain Layer**: Contains business entities and interfaces
- **Use Case Layer**: Contains business logic implementations
- **Infrastructure Layer**: Contains external dependencies (database, etc.)
- **Interface Layer**: Contains HTTP handlers and middleware

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Create a pull request

### Development Guidelines

- Follow Go coding conventions
- Write comprehensive tests
- Update documentation
- Use meaningful commit messages
- Keep functions small and focused

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For questions or issues, please:

1. Check the existing issues
2. Create a new issue with detailed information
3. Provide steps to reproduce any bugs

## 🚀 Roadmap

- [ ] OAuth2 integration (Google, GitHub)
- [ ] Rate limiting
- [ ] Caching with Redis
- [ ] Event-driven architecture
- [ ] Microservice decomposition
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline setup
- [ ] Performance monitoring
- [ ] Advanced search capabilities
- [ ] File upload for product images

---

**Happy coding! 🎉**