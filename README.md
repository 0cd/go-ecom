# Go E-Commerce API

A simple e-commerce REST API built with Go, PostgreSQL, and [chi](https://github.com/go-chi/chi).

## Features

- User registration, login, JWT authentication
- Product and order management
- Admin endpoints
- PostgreSQL migrations with Goose
- SQLC for type-safe database access

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Setup

1. **Clone the repository:**

   ```sh
   git clone https://github.com/0cd/go-ecom.git
   cd go-ecom
   ```

2. **Copy and edit environment variables:**

   ```sh
   cp .env.example .env
   # Edit .env as needed
   ```

3. **Start PostgreSQL with Docker Compose:**

   ```sh
   docker-compose up -d
   ```

4. **Run database migrations:**

   ```sh
   go install github.com/pressly/goose/v3/cmd/goose@latest
   goose up
   ```

5. **Install dependencies:**

   ```sh
   go mod tidy
   ```

6. **Start the API server:**
   ```sh
   go run cmd/*.go
   ```

The server will run on [http://localhost:1337](http://localhost:1337).
