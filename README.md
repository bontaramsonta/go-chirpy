# Go Chirpy

(an example project for learning http server development in go lang)

A microblogging application built with Go that allows users to create and share short messages (chirps).

## Features

- User management with email-based registration
- Create and share chirps (short messages)
- View all chirps or specific chirps by ID
- Content filtering for inappropriate words
- Admin dashboard with metrics
- RESTful API endpoints
- PostgreSQL database integration

## Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Environment variables (see Configuration section)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/bontaramsonta/go-chirpy.git
cd go-chirpy
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
   - Create a PostgreSQL database
   - Run the database migrations (using goose)

4. Configure environment variables:
   - Create a `.env` file in the root directory
   - Add the following variables:
     ```
     DB_URL=postgres://username:password@localhost:5432/database_name
     PLATFORM=dev
     ```

## Running the Application

Start the server:
```bash
go run main.go
```

The server will start on port 8080 by default.

## API Endpoints

### Users
- `POST /api/users` - Create a new user
  - Request body: `{ "email": "user@example.com" }`

### Chirps
- `POST /api/chirps` - Create a new chirp
  - Request body: `{ "user_id": "uuid", "body": "message" }`
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get a specific chirp by ID

### Admin
- `GET /admin/metrics` - View application metrics
- `POST /admin/reset` - Reset metrics and database (dev environment only)

### Health Check
- `GET /api/healthz` - Check API health status

## Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   └── database/          # Database models and queries
├── sql/
│   ├── schema/           # Database migrations
│   └── queries/          # SQL queries
├── apiDocs/              # API documentation
└── .env                  # Environment variables
```

## Development

The project uses:
- `sqlc` for type-safe database queries
- `goose` for database migrations
- `godotenv` for environment variable management
- `lib/pq` for PostgreSQL driver
