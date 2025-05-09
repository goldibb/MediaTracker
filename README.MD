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

- **Advanced Sorting:** 
  - Sort books by title, author, publication year, or date added
- **Pagination:** 
  - Automatic pagination of collections for better performance and readability
- **Book Editing:** 
  - Full ability to edit all properties of books in your collection

## Planned Features

- **Tags and Categories:** Add custom tags and categorize your media
- **Statistics:** Track reading progress and generate statistics
- **Export/Import Data:** Easily transfer collections between platforms
- **Additional API Integration:** Support for Google Books API and other sources
- **User System:** Individual accounts and permission management

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
git clone https://github.com/goldibb/MediaTracker.git
cd MediaTracker
```

2. Run the application:
```bash
docker-compose up --build
```

3. Open your browser and go to:
   - Frontend: http://localhost:5173/books
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
PORT=5137
APP_ENV=local
BLUEPRINT_DB_HOST=psql_bp
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=blueprint
BLUEPRINT_DB_USERNAME=melkey
BLUEPRINT_DB_PASSWORD=password1234
BLUEPRINT_DB_SCHEMA=public

```

4. Run the application:
```bash
go run cmd/api/main.go
```
5. Access the application:
    - Backend API: http://localhost:5137
    - Frontend interface: http://localhost:5137/books

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


## License

[MIT](https://choosealicense.com/licenses/mit/)

---

*This project is under active development.*
