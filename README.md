# go-fiber-service

A production-oriented backend service built with **Golang** and **Fiber**, emphasizing maintainability, clear boundaries, and long-term scalability. This project demonstrates authentication using **JWT**, a structured **Clean Architecture** approach, and API documentation generated automatically using **Swagger (Swag)**.

This repository is intended for **portfolio purposes** and as a reference for building scalable Go services.

---

## ✨ Key Capabilities

- RESTful API built with **Fiber v2**
- **JWT-based authentication** (access token)
- **Dedicated Role-Based Access Control (RBAC)** local middleware (`internal/rbac`)
- Clean Architecture (separation of concerns)
- Request validation using `go-playground/validator`
- PostgreSQL integration using **GORM**
- Centralized and structured logging via a dedicated `logger` module using **Logrus**
- Environment configuration via `.env`
- **Swagger API documentation** (auto-generated)
- Test-ready setup using `stretchr/testify`

---

## 🏗️ Architecture Overview

This project applies **Clean Architecture principles** to enforce clear boundaries between business logic and infrastructure concerns:

- Frameworks and third-party libraries live at the outer layer
- Business rules remain independent and testable
- Infrastructure concerns such as **logging**, database access, and HTTP handling are isolated

Core layers are separated into:

- **Handlers / Controllers** – HTTP request handling and response mapping
- **Use Cases / Services** – business logic orchestration
- **Repositories** – data access abstraction
- **Entities / Domain Models** – core business definitions

Cross-cutting concerns like logging and configuration are centralized (e.g. `logger.go`) to avoid duplication and tight coupling, improving long-term maintainability and operational visibility.

---

## 🛠️ Tech Stack

- **Go** 1.24
- **Fiber** v2
- **GORM** (PostgreSQL)
- **JWT** (`golang-jwt/jwt`)
- **Swagger / OpenAPI** (`swaggo/swag`)
- **Validator** (`go-playground/validator`)
- **Logrus**

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/fatihrizqon/gofiber-microservice.git
cd go-fiber-service
```

### 2. Setup Environment Variables

Create a `.env` file:

```env
APP_PORT=3000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_fiber_service

JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

The server will start on:

```
http://localhost:3000
```

---

## 🔐 Authentication

Authentication is implemented using **JWT**:

- Users authenticate via a login endpoint
- A signed JWT is returned upon successful authentication
- Protected routes require a valid `Authorization: Bearer <token>` header

JWT handling is isolated within the authentication layer to keep business logic clean.

---

## 📚 API Documentation (Swagger)

This project uses **Swag** to generate Swagger documentation automatically from code annotations.

### Install Swag CLI

Make sure you have Go installed, then run:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger Docs

```bash
swag init
```

### Access Swagger UI

Once the application is running:

```
http://localhost:3000/swagger/index.html
```

---

## 🧪 Testing

Tests are structured to support unit and service-level testing.

Run all tests:

```bash
go test ./...
```

---

## 📁 Project Structure (Simplified)

```
├── config/
├── docs/          # Swagger generated files
├── util/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── entity/
│   ├── middleware/
│   └── rbac/      # Dedicated RBAC local middleware
├── logger/
├── middleware/
├── router/
├── test/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

---

## 🎯 Purpose

This repository demonstrates senior-level backend engineering practices, including:

- Designing a maintainable Go service with clear architectural boundaries
- Implementing stateless JWT authentication suitable for distributed systems
- Centralizing logging through a shared logger to improve observability and debuggability
- Maintaining framework-agnostic business logic
- Applying API documentation and testing best practices

It is intended as a **senior backend portfolio project** and a reference for designing maintainable Go services in real-world environments.

---

## 📄 License

This project is open-source and available under the **MIT License**.

---

**Author**  
Fatih Rizqon


## Next Sprint Additions

### New Backend Services
- `StockMovementService` — FIFO & moving-average valuation, stock balance queries
- `RBACMiddleware` — permission-based access control (`RequirePermission`, `RequireAnyPermission`)
- `SeedPostingRules` — seeds default double-entry rules for all transaction types

### New Unit Tests (3 files, ~40 test cases)
- `posting_engine_test.go` — WorkflowTransition, WorkflowActionToState
- `sales_invoice_test.go` — Create validation, FindAll/FindById, full workflow lifecycle
- `vendor_bill_test.go` — Full workflow lifecycle, period close/lock transitions

### Frontend
- `CreateInvoiceModal` — multi-line Sales Invoice creation form with live totals
- `CreateVendorBillModal` — multi-line Vendor Bill creation form with tax calculation
- Both wired into their respective list pages

### Repository Refactors
- `IVendorBillRepository` expanded: `Create`, `UpdateStatus`, `HasOpenPayments`, `GenerateBillNo`
- `VendorBillService.Transition` now uses repo methods (testable without raw DB)
