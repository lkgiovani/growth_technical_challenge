# TypeSpec Development Guide

## Overview

TypeSpec is a modern API specification language that provides a better developer experience compared to writing raw OpenAPI/Swagger specs. It's strongly typed, modular, and can generate OpenAPI 3.x specifications.

## Why TypeSpec?

### Advantages over Manual Swagger/OpenAPI

1. **Type Safety**: Catch errors at compile time
2. **DRY Principle**: Reuse models with spreads (`...ModelName`)
3. **Maintainability**: Cleaner syntax, easier to read and update
4. **Consistency**: Enforces patterns across your API
5. **Extensibility**: Easy to add new endpoints and models
6. **Multi-format Output**: Generate OpenAPI 3.0, 3.1, or other formats

### Comparison

**OpenAPI/Swagger (YAML/JSON)**:

```yaml
paths:
  /users/{id}:
    get:
      summary: Get user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/User"
```

**TypeSpec**:

```typespec
@route("/users")
interface UserOperations {
  @get
  @summary("Get user")
  getUser(@path id: uint32): {
    @statusCode statusCode: 200;
    @body body: User;
  };
}
```

## Project Structure

```
docs/
├── main.tsp           # Service configuration & common models
├── models.tsp         # All data models (User, Employee, Department)
├── operations.tsp     # API endpoints organized by resource
├── package.json       # Dependencies
├── tspconfig.yaml     # TypeSpec compiler config
└── schema/            # Generated OpenAPI + HTML viewers
    ├── openapi.yaml
    ├── redoc.html
    ├── swagger.html
    └── scalar.html
```

## Key Concepts

### 1. Models

Models define your data structures:

```typespec
@doc("User model")
model User {
  @doc("User ID")
  id?: uint32;

  @doc("User name")
  name: string;

  @doc("User email")
  @format("email")
  email: string;
}
```

### 2. Operations (Endpoints)

Operations define your API endpoints:

```typespec
@tag("Users")
@route("/users")
interface UserOperations {
  @doc("Get all users")
  @get
  @summary("List users")
  listUsers(): {
    @statusCode statusCode: 200;
    @body body: User[];
  };

  @doc("Create a new user")
  @post
  @summary("Create user")
  createUser(@body user: UserCreateRequest): {
    @statusCode statusCode: 201;
    @body body: User;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}
```

### 3. Model Spreads

Reuse model properties with the spread operator:

```typespec
model Employee {
  id?: string;
  name: string;
  cpf: string;
}

model EmployeeWithManager {
  ...Employee;  // Includes all properties from Employee
  manager_name?: string;
}
```

### 4. Decorators

Decorators add metadata to your models and operations:

- `@doc()` - Add description
- `@summary()` - Short summary
- `@example()` - Example value
- `@pattern()` - Regex validation
- `@minLength()` / `@maxLength()` - String length
- `@format()` - Format hint (email, uuid, etc.)
- `@tag()` - Group operations
- `@route()` - Define path
- `@get/@post/@put/@delete` - HTTP method
- `@statusCode` - HTTP status code
- `@body` - Request/response body
- `@path` - Path parameter
- `@query` - Query parameter

## Development Workflow

### 1. Install Dependencies

```bash
cd docs
npm install
```

### 2. Make Changes

Edit TypeSpec files:

- `models.tsp` - Add/modify data models
- `operations.tsp` - Add/modify endpoints
- `main.tsp` - Update service info or common models

### 3. Compile

```bash
npm run compile
```

This generates `openapi.yaml` in `docs/schema/`.

### 4. Watch Mode

For continuous compilation:

```bash
npm run watch
```

### 5. Format Code

```bash
npm run format
```

## Adding New Endpoints

### Example: Add a new "Projects" resource

1. **Add models** in `models.tsp`:

```typespec
@doc("Project model")
model Project {
  @doc("Project ID")
  @format("uuid")
  id?: string;

  @doc("Project name")
  name: string;

  @doc("Project description")
  description?: string;

  @doc("Owner ID")
  @format("uuid")
  owner_id: string;
}

@doc("Project create request")
model ProjectCreateRequest {
  name: string;
  description?: string;
  owner_id: string;
}
```

2. **Add operations** in `operations.tsp`:

```typespec
@tag("Projects")
@route("/projects")
interface ProjectOperations {
  @doc("List all projects")
  @get
  @summary("List projects")
  listProjects(): {
    @statusCode statusCode: 200;
    @body body: Project[];
  };

  @doc("Create a new project")
  @post
  @summary("Create project")
  createProject(@body project: ProjectCreateRequest): {
    @statusCode statusCode: 201;
    @body body: Project;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };

  @doc("Get project by ID")
  @get
  @summary("Get project")
  getProject(@path id: string): {
    @statusCode statusCode: 200;
    @body body: Project;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };
}
```

3. **Compile**:

```bash
npm run compile
```

4. **Implement in Go**: Create handler, service, repository

5. **Verify**: Check the generated OpenAPI spec and test the viewers

## Common Patterns

### Optional vs Required Fields

```typespec
model Example {
  required_field: string;      // Required
  optional_field?: string;     // Optional
  with_default?: int32 = 10;   // Optional with default
}
```

### Multiple Response Types

```typespec
operation(): {
  @statusCode statusCode: 200;
  @body body: Success;
} | {
  @statusCode statusCode: 404;
  @body body: Error;
} | {
  @statusCode statusCode: 500;
  @body body: Error;
};
```

### Generic Models

```typespec
model ListRequest<T> {
  filters?: T;
  page?: int32 = 1;
  limit?: int32 = 10;
}

model ListResponse<T> {
  data: T[];
  total: int64;
  page: int32;
  limit: int32;
}
```

### Path Parameters

```typespec
@get
getUser(@path id: uint32): Response;

@get
@route("/{managerId}/employees")
getEmployees(@path managerId: string): Response;
```

### Query Parameters

```typespec
@get
searchUsers(
  @query name?: string,
  @query email?: string,
  @query page?: int32
): Response;
```

## Validation

TypeSpec supports various validation decorators:

```typespec
model Employee {
  @pattern("^[0-9]{11}$")
  @example("12345678901")
  cpf: string;

  @minLength(1)
  @maxLength(20)
  rg?: string;

  @format("email")
  email: string;

  @format("uuid")
  id: string;
}
```

## Troubleshooting

### Compilation Errors

If you get compilation errors:

1. Check syntax in `.tsp` files
2. Ensure all imports are correct
3. Verify model names match references
4. Check for missing decorators

### Output Not Generated

1. Verify `tspconfig.yaml` is correct
2. Check `package.json` has all dependencies
3. Run `npm install` again
4. Delete `node_modules` and reinstall

### HTML Viewers Not Working

1. Ensure HTML files are in `docs/schema/`
2. Check that Gin is loading templates correctly
3. Verify the OpenAPI file path in HTML files
4. Check browser console for errors

## Resources

- [TypeSpec Documentation](https://typespec.io)
- [TypeSpec Playground](https://typespec.io/playground)
- [TypeSpec HTTP Library](https://typespec.io/docs/libraries/http)
- [TypeSpec OpenAPI Emitter](https://typespec.io/docs/emitters/openapi3)

## Tips

1. **Use watch mode** during development for instant feedback
2. **Group related models** in separate files if they grow large
3. **Keep operations organized** by resource/domain
4. **Add detailed descriptions** with `@doc()` for better documentation
5. **Use examples** with `@example()` to help API consumers
6. **Validate early** - compile often to catch errors quickly
7. **Leverage model spreads** to avoid duplication
8. **Use generic models** for common patterns (pagination, filtering, etc.)
