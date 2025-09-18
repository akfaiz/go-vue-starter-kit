# Go-Vue Starter Kit

A full-stack starter kit built with Go and Vue 3, featuring a modern architecture with clean separation of concerns, authentication, and a responsive UI.

## 🌟 Demo

Check out the live demo: **[https://go-vue-starter-kit.fly.dev](https://go-vue-starter-kit.fly.dev)**

OpenAPI Docs: **[https://go-vue-starter-kit.fly.dev/docs](https://go-vue-starter-kit.fly.dev/docs)**

## 🚀 Features

### Backend (Go)
- **Clean Architecture** with dependency injection using Uber FX
- **RESTful API** with Echo framework
- **JWT Authentication** with access and refresh tokens
- **Database Migrations** with custom migration tool
- **Email Service** with SMTP support
- **Internationalization (i18n)** support
- **Request Validation** with custom validators
- **Structured Logging** with configurable levels
- **Health Check** endpoints
- **OpenAPI Documentation** support

### Frontend (Vue 3)
- **Vue 3** with Composition API
- **TypeScript** for type safety
- **Vite** for fast development and building
- **Vue Router** for navigation
- **Pinia** for state management
- **Vuetify** for UI components
- **Responsive Design** with modern layouts
- **Authentication Flow** with token management
- **Error Handling** with user-friendly messages

### DevOps & Deployment
- **Docker** support with multi-stage builds
- **Docker Compose** for local development
- **Makefile** for common development tasks
- **Hot Reload** for development

## 📋 Prerequisites

- Go 1.25+
- Node.js 22+
- PNPM (for frontend dependencies)
- PostgreSQL
- Docker (optional)

## 🛠️ Installation

### 1. Clone the repository

```bash
git clone https://github.com/akfaiz/go-vue-starter-kit.git
cd go-vue-starter-kit
```

### 2. Environment Configuration

Copy the environment file and configure your settings:

```bash
cp .env.example .env
```

Update the `.env` file with your configuration, especially database and SMTP settings.

### 3. Database Setup

Create a PostgreSQL database and run migrations:

```bash
# Install Go dependencies
go mod tidy

# Run database migrations
go run . migrate up
```

### 4. Frontend Setup

```bash
# Navigate to UI directory
cd ui

# Install dependencies
pnpm install

# Build frontend assets
pnpm build

# Return to root directory
cd ..
```

## 🚀 Development

### Using Makefile (Recommended)

```bash
# Install frontend dependencies
make web-install

# Build frontend
make web-build

# Build backend with embedded frontend
make build-embed

# Run in production mode
make run
```

### Manual Commands

#### Backend Development

```bash
# Run server in development mode
DEV=1 go run . serve

# Run database migrations
go run . migrate up
```

#### Frontend Development

```bash
cd ui

# Start development server
pnpm dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

## 🐳 Docker Development

### Using Docker Compose

```bash
# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Building Docker Image

```bash
# Build the image
docker build -t go-vue-starter-kit .

# Run the container
docker run -p 3000:3000 --env-file .env go-vue-starter-kit
```

## 📁 Project Structure

```
.
├── cmd/                    # CLI commands
│   ├── root.go
│   ├── migrate/           # Database migration command
│   └── serve/             # Server command
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── db/               # Database connection
│   ├── delivery/         # Delivery layer (HTTP)
│   │   └── http/
│   │       ├── handler/  # HTTP handlers
│   │       ├── middleware/ # HTTP middleware
│   │       └── routes/   # Route definitions
│   ├── domain/           # Business domain interfaces
│   ├── hash/             # Hashing utilities
│   ├── lang/             # Internationalization
│   ├── model/            # Data models
│   ├── provider/         # External service providers (e.g., SMTP)
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   └── validator/        # Request validation
├── ui/                   # Vue.js frontend
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── layouts/      # Layout components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API services
│   │   ├── stores/       # Pinia stores
│   │   └── utils/        # Utility functions
│   └── public/           # Static assets
├── web/                  # Embedded web assets
├── db/                   # Database migrations
├── pkg/                  # Public packages
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile           # Docker build configuration
├── Makefile            # Development automation
└── main.go             # Application entry point
```

## 🔧 Available Commands

### Go Commands

```bash
go run . serve              # Start the server
go run . migrate up         # Run database migrations
```

### Make Commands

```bash
make dev                    # Run in development mode
make run                    # Run in production mode
make build                  # Build the application
make build-embed           # Build with embedded frontend
make web-install           # Install frontend dependencies
make web-build             # Build frontend
make clean                 # Clean build artifacts
```

### Frontend Commands (in ui/ directory)

```bash
pnpm dev                   # Start development server
pnpm build                 # Build for production
pnpm preview              # Preview production build
pnpm typecheck            # Run TypeScript type checking
pnpm lint                 # Run ESLint
```

## 🔐 Authentication

The application uses JWT-based authentication with:
- **Access tokens** (short-lived, 15 minutes)
- **Refresh tokens** (long-lived, 7 days)
- **Automatic token refresh** on the frontend
- **Secure token storage** in HTTP-only cookies

## 🧪 Testing

```bash
# Run Go tests
go test ./...

# Run frontend tests
cd ui && pnpm test
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Support

If you have any questions or need help, please open an issue on GitHub.

## 🙏 Acknowledgments

- [Echo](https://echo.labstack.com/) - High performance, minimalist Go web framework
- [Uber FX](https://uber-go.github.io/fx/) - Dependency injection framework
- [Bun](https://bun.uptrace.dev/) - SQL-first Golang ORM
- [Vue 3](https://vuejs.org/) - The progressive JavaScript framework
- [Vuetify](https://vuetifyjs.com/) - Material Design component framework
- [Sneat](https://github.com/themeselection/sneat-vuetify-vuejs-admin-template-free) - Free Vue 3 Admin Template