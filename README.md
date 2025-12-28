# CipherBin

A lightweight, self-hosted web application for storing and sharing code snippets.  Built with Go and designed for simplicity and ease of deployment.

## Overview

CipherBin provides a clean, minimal interface for creating, viewing, and managing code snippets with automatic expiration.  It features user authentication, session management, and a responsive design suitable for developers who need a quick way to share code samples.

**Live Demo:** [cipherbin.onrender.com](https://cipherbin.onrender.com)

## Features

- **Snippet Management** - Create, view, and list code snippets with customizable expiration (1 day, 1 week, or 1 year)
- **User Authentication** - Secure signup and login with bcrypt password hashing
- **Session Management** - Persistent sessions stored in MySQL with automatic cleanup
- **CSRF Protection** - Built-in cross-site request forgery protection
- **TLS Support** - HTTPS enabled for secure connections
- **Docker Ready** - Full Docker and Docker Compose support for easy deployment

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.24+ |
| Router | chi |
| Templates | templ |
| Database | MySQL 8.0 |
| Sessions | SCS with MySQL store |
| Forms | go-playground/form |

## Project Structure

```
CipherBin/
├── cmd/
│   └── web/           # Application entrypoint, handlers, routes, middleware
├── internal/
│   └── models/        # Database models and queries
├── ui/
│   ├── html/          # templ templates
│   └── static/        # CSS, JavaScript, images
├── setup.sql          # Database schema
├── Dockerfile
└── docker-compose. yaml
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- MySQL 8.0
- [templ](https://templ.guide/) CLI (for template generation)

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/kayden-vs/CipherBin.git
   cd CipherBin
   ```

2. **Set up the database**

   ```bash
   mysql -u root -p < setup.sql
   ```

3. **Generate templ files**

   ```bash
   templ generate
   ```

4. **Run the application**

   ```bash
   go run ./cmd/web -addr=:4000 -dsn="user:password@/snippetbox?parseTime=true"
   ```

5. Open your browser at `https://localhost:4000`

### Using Docker

The easiest way to run CipherBin is with Docker Compose:

```bash
docker-compose up -d
```

This starts both the MySQL database and the web application.  Access the app at `http://localhost:8080`.

To stop the services:

```bash
docker-compose down
```

### Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:4000` | HTTP server address |
| `-dsn` | - | MySQL data source name |

## Database Schema

CipherBin uses three tables:

- **snippets** - Stores snippet content with title, creation date, and expiration
- **users** - User accounts with hashed passwords
- **sessions** - Session tokens for authenticated users

See `setup.sql` for the complete schema. 

## Development

### Hot Reload

For development with hot reload, use [Air](https://github.com/cosmtrek/air):

```bash
air
```

Configuration is provided in `.air.toml`.

### Running Tests

```bash
go test ./... 
```

## Contributing

Contributions are welcome.  Please open an issue to discuss proposed changes before submitting a pull request.

## License

This project is open source.  See the repository for license details. 
