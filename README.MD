# MediaTracker

## Project Description

MediaTracker is an application for tracking media collections, allowing you to manage books, movies, and TV series. The application enables adding, searching, updating, and removing items from your collection.

## Features

- **Book Management:**
  - Add new books manually or through search
  - Display books divided into "Read" and "Unread" categories
  - Update and delete books

- **OpenLibrary API Integration:**
  - Search books by title, author, or ISBN
  - Fetch covers and metadata

- **User Interface:**
  - Modern, responsive interface built with HTMX and Tailwind CSS
  - Dynamic updates without page reloads

## TODO
- TV/Movie tracking
- Ratings
- Book Descriptions 
- Better UI/UX 
- Notes

## Technologies

- **Backend:**
  - Go (Golang)
  - Gin Framework
  - PostgreSQL
  - Docker

- **Frontend:**
  - HTMX
  - Tailwind CSS
  - Vanilla JavaScript

## Running the Project

### Requirements

- Docker and Docker Compose
- Go 1.24 (for local development)

### Using Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/your-username/MediaTracker.git
cd MediaTracker
```

2. Run the application:
```bash
docker compose up --build
```

3. Open your browser and go to:
   - Frontend: http://localhost:5173
   - API: http://localhost:5137/api/books

### Locally (for developers)

1. Clone the repository and navigate to the project directory

2. Run PostgreSQL database:
```bash
docker compose up -d psql_bp
```

3. Set environment variables:
```bash
# Windows PowerShell
$env:BLUEPRINT_DB_HOST="localhost"
$env:PORT="5137"
$env:BLUEPRINT_DB_DATABASE="blueprint"
$env:BLUEPRINT_DB_USERNAME="melkey"
$env:BLUEPRINT_DB_PASSWORD="password1234"
$env:BLUEPRINT_DB_SCHEMA="public"
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## Project Structure

```
MediaTracker/
├── cmd/
│   └── api/
│       └── main.go       # Application entry point
├── internal/
│   ├── database/         # Database access layer
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Data models
│   ├── server/           # Server configuration
│   └── services/         # Business layer
├── frontend/             # HTML, CSS, JS files
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Docker image definition
├── .env                  # Environment variables
└── init.sql              # Database initialization script
```

## API

### Books

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/books` | GET | Get all books |
| `/api/books/grouped` | GET | Get books grouped by status |
| `/api/books/:id` | GET | Get book details |
| `/api/books` | POST | Add a new book |
| `/api/books/:id` | PUT | Update a book |
| `/api/books/:id` | DELETE | Delete a book |
| `/api/books/search` | POST | Search books in OpenLibrary |

## Development

### Installing Dependencies

```bash
go mod download
```

### Running Tests

```bash
go test ./...
```

### Live Reload

```bash
make watch
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

---

*This project is under active development.*