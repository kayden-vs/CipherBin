# CipherBin

A lightweight, self-hosted web application for storing and sharing code snippets.  Built with Go and designed for simplicity and ease of deployment.

## Overview

CipherBin provides a clean, minimal interface for creating, viewing, and managing code snippets with automatic expiration.  It features user authentication, session management, and a responsive design that works on all devices.

**Live Demo:** [cipherbin.onrender.com](https://cipherbin.onrender.com)

## Features

- **Snippet Management** - Create, view, and list code snippets with customizable expiration (1 day, 1 week, or 1 year)
- **Anonymous Snippet Creation** - Create snippets without requiring authentication
- **Author Attribution** - Snippets display author names for authenticated users
- **Auto-Linking** - URLs in snippet content are automatically converted to clickable links
- **User Authentication** - Secure signup and login with bcrypt password hashing
- **Account Management** - User account section with password change functionality
- **Session Management** - Persistent sessions stored in PostgreSQL with automatic cleanup
- **CSRF Protection** - Built-in cross-site request forgery protection
- **TLS Support** - HTTPS enabled for secure connections
- **Docker Ready** - Full Docker and Docker Compose support for easy deployment

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.24+ |
| Router | chi |
| Templates | templ |
| Database | PostgreSQL  |
| Sessions | SCS with Postgresql store |
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
└── docker-compose.yaml
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Postgres 8.0
- [templ](https://templ.guide/) CLI (for template generation)

### Using Docker

The easiest way to run CipherBin is with Docker Compose:

```bash
docker-compose up -d
```

This starts both the Postgres database and the web application.  Access the app at `http://localhost:8080`.

To stop the services:

```bash
docker-compose down
```

### Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:4000` | HTTP server address |
| `-dsn` | - | Postgres data source name |

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