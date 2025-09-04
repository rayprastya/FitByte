# FitByte API

A comprehensive fitness tracking REST API built with Go, Gin, and GORM, following clean architecture principles and Repository-Service-Handler pattern.

## ✨ Features

- 🚀 **Gin Framework** - High-performance HTTP web framework
- 🗄️ **GORM Integration** - Powerful ORM with SQLite database
- 🏗️ **Clean Architecture** - Repository-Service-Handler pattern
- 📊 **Advanced Filtering** - Complex query parameters with pagination
- 🔄 **Base Entity** - Consistent database schema across all tables
- 📝 **Structured Logging** - Efficient logging with zerolog
- 🔒 **CORS Support** - Cross-origin resource sharing
- 📋 **RESTful API** - Following REST conventions
- 🎯 **Dependency Injection** - Clean separation of concerns
- ⚡ **Auto-Migration** - Database schema automatically managed

## 🏛️ Architecture Overview

### Why Repository-Service-Handler Pattern?

We chose this architecture for several key reasons:

1. **Separation of Concerns**: Each layer has a single responsibility
   - **Repository**: Data access and database queries
   - **Service**: Business logic and validation
   - **Handler**: HTTP request/response handling

2. **Testability**: Easy to mock dependencies and unit test each layer
3. **Maintainability**: Changes in one layer don't affect others
4. **Scalability**: Easy to add new features following the same pattern
5. **Database Independence**: Repository layer abstracts database operations

### Base Entity Design

All database tables inherit from `BaseEntity`:
```go
type BaseEntity struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
}
```

**Benefits:**
- **Consistency**: All tables have the same base fields
- **Soft Delete**: Use `is_active` flag instead of hard deletion
- **Audit Trail**: Automatic timestamps for creation and updates
- **Maintainability**: Common fields managed in one place

## 📁 Project Structure

```
fitbyte/
├── main.go                    # Application entry point
├── cmd/
│   ├── app/main.go           # Alternative main entry
│   └── server/               # Server configuration
│       ├── handlers.go       # Handler initialization
│       └── router.go         # Route definitions
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── database/             # Database connection & migrations
│   │   └── connection.go
│   ├── entities/             # Data models & DTOs
│   │   ├── base.go          # Base entity
│   │   ├── user.go          # User entities
│   │   ├── activity.go      # Activity entities
│   │   └── response.go      # API response models
│   ├── repositories/         # Data access layer
│   │   ├── user.go          # User repository
│   │   └── activity.go      # Activity repository
│   ├── services/            # Business logic layer
│   │   ├── user.go          # User service
│   │   └── activity.go      # Activity service
│   ├── handlers/            # HTTP handlers
│   │   ├── health.go        # Health check handlers
│   │   ├── user.go          # User HTTP handlers
│   │   └── activity.go      # Activity HTTP handlers
│   └── routes/              # Route setup (alternative)
│       └── routes.go
└── pkg/                     # Shared utilities
    ├── logger.go            # Logging middleware
    ├── recovery.go          # Recovery middleware
    └── cors.go              # CORS middleware
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23+ or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd fitbyte
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   # or
   go run cmd/app/main.go
   ```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Health Check
- `GET /api/v1/health/` - Health status
- `GET /api/v1/health/ready` - Readiness check

### Users
- `GET /api/v1/users?limit=5&offset=0&isActive=true` - Get users with filtering
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Soft delete user

### Activities
- `GET /api/v1/activity?limit=5&offset=0&activityType=RUNNING&doneAtFrom=2023-01-01T00:00:00Z&doneAtTo=2023-12-31T23:59:59Z&caloriesBurnedMin=50&caloriesBurnedMax=500` - Get activities with advanced filtering
- `POST /api/v1/activity` - Create new activity
- `GET /api/v1/activity-types` - Get available activity types

## 🔍 Advanced Filtering System

### Activity Filtering Parameters

All parameters are optional and use **AND** logic:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `limit` | number | Results per page (default: 5) | `?limit=10` |
| `offset` | number | Skip results (default: 0) | `?offset=20` |
| `activityType` | string | Filter by activity type enum | `?activityType=RUNNING` |
| `doneAtFrom` | ISO Date | Activities after/equal date | `?doneAtFrom=2023-01-01T00:00:00Z` |
| `doneAtTo` | ISO Date | Activities before/equal date | `?doneAtTo=2023-12-31T23:59:59Z` |
| `caloriesBurnedMin` | number | Minimum calories burned | `?caloriesBurnedMin=50` |
| `caloriesBurnedMax` | number | Maximum calories burned | `?caloriesBurnedMax=500` |

**Note**: Invalid parameter values are ignored, defaults are used for limit/offset.

## 📋 Data Models

### User Response
```json
{
  "id": 1,
  "email": "john@example.com",
  "name": "John Doe",
  "preference": "CARDIO",
  "weightUnit": "KG", 
  "heightUnit": "CM",
  "weight": 75.5,
  "height": 180.0,
  "imageUri": "https://example.com/image.jpg",
  "is_active": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Activity Response
```json
{
  "activityId": "1",
  "activityType": "RUNNING",
  "doneAt": "2023-01-01T08:00:00Z",
  "durationInMinutes": 30,
  "caloriesBurned": 300.0,
  "createdAt": "2023-01-01T08:30:00Z"
}
```

### API Response Format
```json
{
  "success": true,
  "message": "Activities retrieved successfully",
  "data": {
    "activities": [...],
    "meta": {
      "total": 100,
      "limit": 5,
      "offset": 0
    }
  }
}
```

## 🛠️ Development

### Adding New Features

Following the Repository-Service-Handler pattern:

1. **Create Entity** in `internal/entities/`
2. **Create Repository** in `internal/repositories/`
3. **Create Service** in `internal/services/`
4. **Create Handler** in `internal/handlers/`
5. **Add to Handlers struct** in `cmd/server/handlers.go`
6. **Add Routes** in `cmd/server/router.go`

### Example: Adding Products

```go
// internal/entities/product.go
type Product struct {
    BaseEntity
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// cmd/server/handlers.go
type Handlers struct {
    HealthHandler   *handlers.HealthHandler
    UserHandler     *handlers.UserHandler
    ActivityHandler *handlers.ActivityHandler
    ProductHandler  *handlers.ProductHandler  // Add new handler
}

// cmd/server/router.go
products := v1.Group("/products")
{
    products.GET("/", h.ProductHandler.GetProducts)
    products.POST("/", h.ProductHandler.CreateProduct)
}
```

## 🏗️ Building for Production

```bash
# Build the application
go build -o fitbyte main.go

# Run the binary
./fitbyte
```

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | Database connection string | `fitbyte.db` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

## 🏃‍♂️ Database

- **Database**: SQLite (development) / PostgreSQL (production ready)
- **ORM**: GORM v2 with auto-migrations
- **Seeded Data**: Pre-populated activity types with calorie rates
- **Soft Delete**: All entities use `is_active` flag

### Available Activity Types
- RUNNING (10 cal/min)
- WALKING (5 cal/min)
- CYCLING (8 cal/min)
- SWIMMING (12 cal/min)
- WEIGHT_LIFTING (6 cal/min)
- YOGA (3 cal/min)
- CARDIO (9 cal/min)

## 🧪 Architecture Benefits

### Repository Layer
- **Database abstraction**: Easy to switch databases
- **Query optimization**: Centralized query logic
- **Testability**: Easy to mock data layer

### Service Layer
- **Business logic**: All validation and processing
- **Transaction management**: Handle complex operations
- **Error handling**: Consistent error responses

### Handler Layer
- **HTTP concerns**: Request parsing, response formatting
- **Authentication**: JWT validation (ready for implementation)
- **Rate limiting**: Request throttling (ready for implementation)

## 🚀 Next Steps

- [x] GORM database integration
- [x] Repository-Service-Handler pattern
- [x] Base entity with soft delete
- [x] Advanced filtering system
- [x] Auto-migrations and seeding
- [ ] JWT authentication
- [ ] Rate limiting
- [ ] API documentation (Swagger)
- [ ] Unit & integration tests
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] PostgreSQL production setup

## 📄 License

This project is licensed under the MIT License.