# 📦 Inventory Management System

A secure, high-performance REST API built with **Go** and **Fiber** for managing products. This API features JWT authentication, efficient PostgreSQL integration, Dockerized deployment, and Swagger documentation.

---
```
I have used AI to help me make README and Swagger Docs .Moreover in API i have mainly used AI in correcting my code and for Bolierpate code which I have picked from previous projects
```

1. **Simplest Way to set Up**
   ```bash
   git clone https://github.com/lokesh2201013/inventory-management.git

   Then

   docker-compose up -d
   ```
You are all set to go

Code snippet

## ✨ Features

- 🔐 **JWT Authentication** – Secure endpoints using JSON Web Tokens  
- 🔄 **RESTful API** – Standardized endpoints for product management  
- 🗃️ **PostgreSQL Database** – Integrated with GORM ORM  
- 📚 **Swagger Documentation** – Interactive API docs available  
- ⚙️ **.env Configuration** – Simple environment-based config  

---

## 🚀 Advanced Highlights

### 1. 📉 Efficient Rate Limiting

- Middleware: `fiber/middleware/limiter`
- Config: `100 requests/min per IP`
- ✅ Burst-tolerant and prevents abuse

### 2. ⚡ Database Connection Pooling

- `SetMaxOpenConns(25)`
- `SetMaxIdleConns(10)`
- `SetConnMaxLifetime(5 * time.Minute)`
- ✅ Reduces latency and improves scalability

### 3. 🧾 Structured Logging with Zap

- Logs in **JSON format** (machine-readable)
- Integrated as Fiber middleware
- Log levels: `info`, `warn`, `error`

### 4. 🐳 Multi-Stage Docker Builds

- **Builder Stage**: Compile Go app
- **Final Stage**: Minimal Alpine image (~15MB)
- ✅ Faster deployments, smaller image size

---

## 📖 API Documentation

Once running, access Swagger UI at:  
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### 🧪 Endpoints

| Method | Endpoint                               | Description                            | Auth Required |
|--------|----------------------------------------|----------------------------------------|---------------|
| POST   | `/register`                            | Register a new user                    | ❌ No          |
| POST   | `/login`                               | Authenticate and get JWT              | ❌ No          |
| POST   | `/products`                            | Create a new product                   | ✅ Yes         |
| GET    | `/products`                            | Get all products (paginated)          | ✅ Yes         |
| GET    | `/products/by-id?product_id=<uuid>`    | Get a product by ID                   | ✅ Yes         |
| GET    | `/products/quantity?most=true`         | Get product with highest quantity     | ✅ Yes         |
| PUT    | `/products/:id/quantity`               | Update quantity of a product          | ✅ Yes         |

---

## ⚙️ Getting Started

### 🧰 Prerequisites

- ✅ Go 1.24.2+
- ✅ PostgreSQL
- ✅ Docker (Recommended)

---

### 1️⃣ Local Setup (Dev)

1. **Clone the repo**
   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```
Create a .env file:
Create a file named .env in the root directory and add your database credentials and JWT secret.

Code snippet

### 2  env.example
```
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=products_db
DB_PORT=5432

JWT_SECRET=your-very-secret-key
Install dependencies:
```
Bash
```
go mod tidy
```

Run the application:

Bash
```
go run main.go
```
The server will start on http://localhost:8080.

## 🐳 Docker Deployment (Using Prebuilt Image)

The easiest way to get started is by using Docker Compose with a prebuilt Docker Hub image.

> 🐙 Image: [`lokesh220/fi-assignment-app:latest`](https://hub.docker.com/r/lokesh220/fi-assignment-app)

### 🔧 docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres
    container_name: postgres_inventory
    restart: always
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: inventory
    ports:
      - "5433:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  app:
    image: lokesh220/fi-assignment-app:latest
    container_name: go_inventory_app
    depends_on:
      - postgres
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: always

volumes:
  pgdata:

```
```

docker-compose up --build
```
The API will be available at http://localhost:8080.

```
📁 Project Structure
.
├── controllers/    # Handlers for API routes
├── database/       # Database connection and migration logic
├── docs/           # Swagger auto-generated files
├── middleware/     # Custom middleware (e.g., auth)
├── models/         # GORM data models (structs)
├── routes/         # Route definitions
├── utils/          # Utility functions (e.g., JWT)
├── .env            # Environment variables (not committed)
├── go.mod          # Go module dependencies
├── main.go         # Application entry point
└── README.md






