# Growth Technical Challenge

RESTful API developed in Go using best practices and modern technologies with layered architecture.

## 🚀 Tech Stack

- **Language**: Go 1.23+
- **HTTP Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL 16
- **Documentation**: TypeSpec + OpenAPI 3.1 (with ReDoc, Swagger UI, Scalar)
- **Dependency Injection**: Uber FX
- **Containerization**: Docker + Docker Compose

## 📐 Architecture

The project follows a clean layered architecture:

- **Models**: Domain entities and database schemas
- **Repositories**: Data access layer with interfaces
- **Services**: Business logic and validation rules
- **Handlers**: HTTP request/response handling
- **Routes**: API endpoint definitions

## 📁 Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── database.go          # Database configurations
├── internal/
│   ├── database/
│   │   └── database.go      # Connection and migrations
│   ├── handlers/
│   │   ├── user_handler.go
│   │   ├── employee_handler.go
│   │   ├── department_handler.go
│   │   └── manager_handler.go
│   ├── models/
│   │   ├── user.go
│   │   ├── employee.go
│   │   └── department.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── employee_repository.go
│   │   └── department_repository.go
│   ├── services/
│   │   ├── user_service.go
│   │   ├── employee_service.go
│   │   └── department_service.go
│   ├── routes/
│   │   └── routes.go
│   └── utils/
│       └── validators.go     # CPF and RG validators
├── docs/
│   └── docs.go              # Swagger documentation
├── docker-compose.yml       # Container orchestration
├── Dockerfile               # Application Docker image
├── Makefile                 # Facilitated commands
└── go.mod                   # Project dependencies
```

## 🔧 Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)
- Make (optional, for facilitated commands)

## 🐳 Running with Docker (Recommended)

### 1. Clone the repository

```bash
git clone https://github.com/lkgiovani/growth_technical_challenge.git
cd growth_technical_challenge
```

### 2. Configure environment variables

```bash
cp env.example .env
```

Edit the `.env` file as needed.

### 3. Start containers

```bash
docker-compose up -d
```

Or using Make:

```bash
make up
```

### 4. Check logs

```bash
docker-compose logs -f app
```

Or:

```bash
make logs
```

The API will be available at:

- **API**: http://localhost:8080
- **Documentation**:
  - ReDoc: http://localhost:8080/docs/redoc
  - Swagger UI: http://localhost:8080/docs/swagger
  - Scalar: http://localhost:8080/docs/scalar
  - OpenAPI Spec: http://localhost:8080/docs/openapi.yaml
- **PostgreSQL**: localhost:5432

## 💻 Local Development

### 1. Install dependencies

```bash
go mod download
```

### 2. Configure database

Make sure you have PostgreSQL running and configure environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=growth_db
export DB_SSLMODE=disable
export APP_PORT=8080
```

### 3. Run the application

```bash
go run cmd/main.go
```

## 📊 Domain Model

### Employee (Colaborador)

- **id** (UUIDv7) - Primary key
- **name** (required) - Employee name
- **cpf** (required, unique, validated) - Brazilian CPF
- **rg** (optional, unique if provided) - Brazilian RG
- **department_id** (FK to Department, required) - Department reference

### Department (Departamento)

- **id** (UUIDv7) - Primary key
- **name** (required) - Department name
- **manager_id** (FK to Employee, required) - Manager reference
- **parent_department_id** (FK optional to Department) - Parent department for hierarchy

## 📚 Business Rules

1. **CPF must be unique and valid**
2. **RG, if provided, must also be unique**
3. **The manager must be an existing Employee linked to the same department**
4. **The Parent Department is optional, but cannot create cycles in the hierarchy**

## 🔌 API Endpoints

### Health Check

```bash
GET /health
```

### Employees

#### Create Employee

```bash
POST /api/v1/employees
Content-Type: application/json

{
  "name": "John Doe",
  "cpf": "12345678901",
  "rg": "123456789",
  "department_id": "uuid-here"
}
```

#### Get Employee by ID

```bash
GET /api/v1/employees/:id
```

Returns the employee and their department manager's name.

#### Update Employee

```bash
PUT /api/v1/employees/:id
Content-Type: application/json

{
  "name": "John Doe Updated",
  "cpf": "12345678901",
  "rg": "123456789",
  "department_id": "uuid-here"
}
```

#### Delete Employee

```bash
DELETE /api/v1/employees/:id
```

Soft delete.

#### List Employees with Filters

```bash
POST /api/v1/employees/list
Content-Type: application/json

{
  "filters": {
    "name": "John",
    "cpf": "12345678901",
    "rg": "123456789",
    "department_id": "uuid-here"
  },
  "page": 1,
  "limit": 10
}
```

### Departments

#### Create Department

```bash
POST /api/v1/departments
Content-Type: application/json

{
  "name": "IT Department",
  "manager_id": "uuid-here",
  "parent_department_id": "uuid-here" // optional
}
```

Validates manager_id and prevents cycles.

#### Get Department by ID

```bash
GET /api/v1/departments/:id
```

Returns department, manager, and complete hierarchical tree of subdepartments.

#### Update Department

```bash
PUT /api/v1/departments/:id
Content-Type: application/json

{
  "name": "IT Department Updated",
  "manager_id": "uuid-here",
  "parent_department_id": "uuid-here"
}
```

Prevents cycles in hierarchy.

#### Delete Department

```bash
DELETE /api/v1/departments/:id
```

#### List Departments with Filters

```bash
POST /api/v1/departments/list
Content-Type: application/json

{
  "filters": {
    "name": "IT",
    "manager_name": "John",
    "parent_department_id": "uuid-here"
  },
  "page": 1,
  "limit": 10
}
```

### Managers

#### Get Manager's Employees

```bash
GET /api/v1/managers/:id/employees
```

Returns all employees from the manager's subordinate departments, recursively.

## ⚖️ HTTP Status Codes

The API uses consistent error responses:

- **200 OK** - Successful request
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid filters or malformed request
- **404 Not Found** - Resource not found
- **409 Conflict** - Uniqueness constraint violation (CPF, RG)
- **422 Unprocessable Entity** - Domain validation error (invalid CPF, business rules)
- **500 Internal Server Error** - Server error

## 🧪 Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Create employee
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "cpf": "12345678901",
    "rg": "123456789",
    "department_id": "uuid-here"
  }'

# List employees with filters
curl -X POST http://localhost:8080/api/v1/employees/list \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {"name": "John"},
    "page": 1,
    "limit": 10
  }'
```

## 🛠️ Make Commands

```bash
make help      # Shows all available commands
make build     # Build Docker images
make up        # Start containers
make down      # Stop containers
make logs      # Show logs
make restart   # Restart containers
make clean     # Remove containers and volumes
make test      # Run tests
```

## 🗄️ Database

Migrations are automatically executed by GORM when starting the application. Tables are created based on models defined in `internal/models/`.

### Connect to PostgreSQL

```bash
docker exec -it growth_postgres psql -U postgres -d growth_db
```

## 🔐 Environment Variables

| Variable            | Description              | Default     |
| ------------------- | ------------------------ | ----------- |
| `POSTGRES_USER`     | PostgreSQL user          | `postgres`  |
| `POSTGRES_PASSWORD` | PostgreSQL password      | `postgres`  |
| `POSTGRES_DB`       | Database name            | `growth_db` |
| `POSTGRES_PORT`     | PostgreSQL port          | `5432`      |
| `DB_HOST`           | Database host            | `localhost` |
| `DB_PORT`           | Database port            | `5432`      |
| `DB_SSLMODE`        | PostgreSQL SSL mode      | `disable`   |
| `APP_PORT`          | Application port         | `8080`      |
| `GIN_MODE`          | Gin mode (debug/release) | `debug`     |

## 📖 API Documentation

This project uses **TypeSpec** for API documentation, providing multiple viewing options:

### Documentation Viewers

- **ReDoc**: http://localhost:8080/docs/redoc (Modern, clean interface)
- **Swagger UI**: http://localhost:8080/docs/swagger (Interactive API explorer)
- **Scalar**: http://localhost:8080/docs/scalar (Beautiful, modern UI)
- **OpenAPI Spec**: http://localhost:8080/docs/openapi.yaml (Raw YAML file)

### Generating TypeSpec Documentation

To regenerate the OpenAPI specification from TypeSpec:

```bash
cd docs
npm install
npm run compile
```

Or using Make:

```bash
cd docs
make install
make compile
```

See [docs/README.md](docs/README.md) for more information about TypeSpec.

## 🚀 Production Build

```bash
# Build application
go build -o app cmd/main.go

# Run
./app
```

Or with Docker:

```bash
docker build -t growth-app .
docker run -p 8080:8080 growth-app
```

## ✅ Features

- ✅ Clean layered architecture (Models, Repositories, Services, Handlers)
- ✅ UUID v7 for primary keys
- ✅ CPF and RG validation
- ✅ Prevention of cycles in department hierarchy
- ✅ Soft delete on records
- ✅ Pagination and filtering
- ✅ Consistent error responses
- ✅ Auto migrations with GORM
- ✅ Swagger documentation
- ✅ Docker multi-stage build
- ✅ Comprehensive README

## 📝 License

This project was developed as a technical challenge.

## 👤 Author

**lkgiovani**

- GitHub: [@lkgiovani](https://github.com/lkgiovani)
