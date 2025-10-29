# API Documentation with TypeSpec

This directory contains the TypeSpec specification for the Growth Technical Challenge API.

## What is TypeSpec?

TypeSpec is a language for describing cloud service APIs and generating other API description languages, client and service code, documentation, and other assets. TypeSpec provides highly extensible core language primitives that can describe API shapes common among REST, GraphQL, gRPC, and other protocols.

## Prerequisites

- Node.js 18+ and npm

## Installation

Install dependencies:

```bash
npm install
```

## Usage

### Compile TypeSpec to OpenAPI

Generate OpenAPI 3.0 specification:

```bash
npm run compile
```

This will generate the OpenAPI specification in the `tsp-output/` directory.

### Watch Mode

To automatically recompile when files change:

```bash
npm run watch
```

### Format TypeSpec Files

Format all TypeSpec files:

```bash
npm run format
```

## Project Structure

```
docs/
├── main.tsp           # Main TypeSpec entry point with service configuration
├── models.tsp         # Data models and DTOs
├── operations.tsp     # API operations and endpoints
├── package.json       # Node.js dependencies
├── tspconfig.yaml     # TypeSpec compiler configuration
└── tsp-output/        # Generated OpenAPI specification (gitignored)
```

## TypeSpec Files

### main.tsp

Contains the service metadata and common models:

- Service information (title, version, description)
- Server configuration
- Common response models (ErrorResponse, SuccessResponse)
- Generic list request/response models

### models.tsp

Defines all data models:

- User (model, create, update requests)
- Employee (model, create, update requests, filters)
- Department (model, create, update requests, filters)
- Response models with additional information

### operations.tsp

Defines all API operations organized by resource:

- Health Check operations
- User CRUD operations
- Employee CRUD operations with validation
- Department CRUD operations with hierarchy
- Manager operations

## Features

- ✅ Complete OpenAPI 3.0 specification generation
- ✅ Strongly typed models and operations
- ✅ Detailed documentation with @doc decorators
- ✅ Request/response validation with patterns and constraints
- ✅ HTTP status codes for all responses
- ✅ Organized by domain (Users, Employees, Departments, Managers)
- ✅ Support for complex scenarios (hierarchy, filtering, pagination)

## Benefits over Manual Swagger

1. **Type Safety**: TypeSpec catches errors at compile time
2. **Reusability**: Share models across operations with spreads (...)
3. **Maintainability**: Single source of truth for API contracts
4. **Consistency**: Enforces consistent API patterns
5. **Extensibility**: Easy to add new operations and models
6. **Code Generation**: Can generate client SDKs, mock servers, and more

## Output

The generated OpenAPI specification can be:

- Imported into Swagger UI
- Used to generate client libraries
- Validated against your API implementation
- Shared with API consumers

## Quick Start

```bash
npm install
npm run compile
```

Then start your Go API and visit:

- http://localhost:8080/docs/redoc
- http://localhost:8080/docs/swagger
- http://localhost:8080/docs/scalar

## Development Guide

For detailed information about TypeSpec development, see [TYPESPEC_GUIDE.md](TYPESPEC_GUIDE.md).

## Next Steps

After generating the OpenAPI spec:

1. Start the Go API server
2. Visit any of the documentation viewers
3. Test API endpoints interactively
4. Generate client SDKs if needed
