# ProductionRESTAPI

A REST API written in Go using the standard `net/http` library. No frameworks. Built around a school management domain with teachers, students, and executives as resources. The focus was on middleware design and understanding what a production-grade API setup looks like without reaching for Gin or Chi.

---

## Project Structure

```
ProductionRESTAPI/
├── cmd/
│   └── api/
│       └── server.go                   # Entry point, TLS config, middleware chain
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── root.go                 # GET /
│   │   │   ├── teachers.go             # Full CRUD for /teachers/
│   │   │   ├── students.go             # Stub handlers
│   │   │   └── execs.go                # Stub handlers
│   │   ├── middlewares/
│   │   │   ├── cors.go                 # CORS with origin whitelist
│   │   │   ├── security_headers.go     # OWASP security headers
│   │   │   ├── rate_limiter.go         # IP-based rate limiting
│   │   │   ├── compression.go          # gzip response compression
│   │   │   ├── hpp.go                  # HTTP Parameter Pollution protection
│   │   │   └── response_time.go        # Response time logging
│   │   └── router/
│   │       └── router.go               # Route definitions
│   ├── models/
│   │   ├── teacher.go
│   │   └── student.go
│   └── repository/
│       └── sqlconnect/
│           ├── sqlconfig.go            # DB connection (singleton pool)
│           └── teacher_crud.go         # All DB operations for teachers
└── pkg/
    └── utils/
        ├── error_handler.go            # Centralized error logging
        └── middlewareutils.go          # ApplyMiddlewares helper
```

---

## Requirements

- Go 1.21 or later
- MariaDB or MySQL
- A TLS certificate and key (`cert.pem`, `key.pem`) in the project root — the server runs HTTPS only
- A `.env` file in the project root

---

## Environment Variables

Create a `.env` file in the root of the project:

```env
API_PORT=:8443
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=3306
HOST=localhost
```

---

## Database Setup

The API expects a `teachers` table. Run this against your database before starting the server:

```sql
CREATE TABLE teachers (
    id        INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    email      VARCHAR(150) NOT NULL,
    class      VARCHAR(50)  NOT NULL,
    subject    VARCHAR(100) NOT NULL
);
```

---

## TLS Certificate

The server will not start without a certificate. For local development, generate a self-signed one:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

---

## Installation and Running

```bash
git clone https://github.com/StupidAfCoder/ProductionRESTAPI.git
cd ProductionRESTAPI
go mod tidy
go run cmd/api/server.go
```

---

## API Reference

Only the `/teachers/` resource has a complete implementation. Students and executives routes exist but are stubs that return plain text.

### Teachers

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/teachers/` | Get all teachers. Supports filtering and sorting via query params. |
| `POST` | `/teachers/` | Add one or more teachers. Accepts a JSON array. |
| `PATCH` | `/teachers/` | Partially update multiple teachers. Accepts a JSON array of objects with `id`. |
| `DELETE` | `/teachers/` | Delete multiple teachers. Accepts a JSON array of IDs. |
| `GET` | `/teachers/{id}` | Get a single teacher by ID. |
| `PUT` | `/teachers/{id}` | Replace a teacher record entirely. |
| `PATCH` | `/teachers/{id}` | Partially update a single teacher. |
| `DELETE` | `/teachers/{id}` | Delete a single teacher. |

### Query Parameters (GET /teachers/)

Filter by field value:
```
GET /teachers/?first_name=John&subject=Math
```

Sort by one or more fields:
```
GET /teachers/?sortby=last_name:asc&sortby=subject:desc
```

Valid sort fields: `first_name`, `last_name`, `email`, `class`, `subject`

Valid sort orders: `asc`, `desc`

### Request / Response Format

**Teacher object:**
```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@school.com",
    "class": "10A",
    "subject": "Mathematics"
}
```

**POST /teachers/ — add multiple teachers:**
```bash
curl -k -X POST https://localhost:8443/teachers/ \
  -H "Content-Type: application/json" \
  -d '[{"first_name":"John","last_name":"Doe","email":"j.doe@school.com","class":"10A","subject":"Math"}]'
```

**GET /teachers/ — get all:**
```bash
curl -k https://localhost:8443/teachers/
```

**PATCH /teachers/{id} — partial update:**
```bash
curl -k -X PATCH https://localhost:8443/teachers/1 \
  -H "Content-Type: application/json" \
  -d '{"subject": "Physics"}'
```

**DELETE /teachers/ — bulk delete:**
```bash
curl -k -X DELETE https://localhost:8443/teachers/ \
  -H "Content-Type: application/json" \
  -d '[1, 2, 3]'
```

---

## Middleware

The middleware stack is composable via `ApplyMiddlewares` in `pkg/utils`. Each middleware wraps `http.Handler` — no third-party router required.

| Middleware | What it does |
|---|---|
| `Security_headers` | Injects OWASP-recommended response headers (CSP, HSTS, X-Frame-Options, etc.) |
| `Cors` | Validates `Origin` against a whitelist. Requests with no `Origin` header pass through. |
| `RateLimiter` | IP-based rate limiting with a configurable request limit and reset interval |
| `Compression` | gzip compression for clients that send `Accept-Encoding: gzip` |
| `Hpp` | HTTP Parameter Pollution protection — deduplicates repeated query/body params |
| `ResponseTimeMiddleware` | Captures response status codes for logging |

Currently only `Security_headers` is active in `server.go`. The others are implemented and tested but commented out while endpoints are being built.

---

## Limitations

- **Students and executives are not implemented.** The routes exist and return placeholder strings. Only teachers have actual database operations.
- **No authentication.** There is no token, session, or API key validation anywhere.
- **No input validation beyond type checking.** Fields like email are stored as plain strings with no format validation.
- **`ResponseTimeMiddleware` does not log anything yet.** It captures the status code but the logging side is unfinished.
- **Rate limiter uses `r.RemoteAddr` for IP detection.** This breaks behind a proxy or load balancer — you would need to read `X-Forwarded-For` instead.
- **Requires TLS.** There is no plain HTTP fallback. Running locally requires a self-signed certificate.

---

## Dependencies

- [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) — MySQL/MariaDB driver
- [`github.com/joho/godotenv`](https://github.com/joho/godotenv) — `.env` file loading

---

## License

MIT
