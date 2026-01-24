# Book Library Application

A full-stack web application for managing and browsing a digital book library with user authentication, category management, and admin features.

## 🚀 Features

- **Book Management**: Browse, search, and view books with cover images
- **Category System**: Organize books into categories with full CRUD operations
- **User Authentication**: 
  - Email/password registration and login
  - Google OAuth integration
  - JWT-based authentication with refresh tokens
- **Admin Panel**: Administrative interface for managing books and categories
- **Content Sections**:
  - Home page
  - Book catalog with search functionality
  - YouTube playlist integration
  - Leadership course content
  - Academy teachings section
- **Modern UI**: Built with Material-UI and responsive design

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.24
- **Database**: PostgreSQL 16
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **OAuth**: Google OAuth 2.0
- **Migrations**: golang-migrate/migrate
- **Configuration**: Viper (YAML/JSON/ENV support)
- **Database Driver**: pgx/v5

### Frontend
- **Framework**: React 19 with TypeScript
- **UI Library**: Material-UI (MUI) v7
- **Routing**: React Router v7
- **State Management**: Zustand
- **HTTP Client**: Axios
- **Build Tool**: Create React App
- **Styling**: Emotion (CSS-in-JS)

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Web Server**: Nginx (for frontend production build)
- **Database**: PostgreSQL with persistent volumes

## 📋 Prerequisites

- **Docker Desktop** (for containerized setup)
- **Go 1.24+** (for local backend development)
- **Node.js 16+** (for local frontend development)
- **PostgreSQL 16** (for local database, if not using Docker)

## 🚀 Quick Start with Docker

The easiest way to get started is using Docker Compose:

1. **Clone the repository** (if applicable)

2. **Start all services**:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api
   - PostgreSQL: localhost:5432

4. **Stop services**:
   ```bash
   docker-compose down
   ```

For detailed Docker setup instructions, see [DOCKER_README.md](./DOCKER_README.md).

## 🔧 Local Development Setup

### Backend Setup

1. **Install Go dependencies**:
   ```bash
   cd backend
   go mod download
   ```

2. **Configure database**:
   - Update `config.yaml` or set environment variables:
     ```yaml
     database:
       dsn: 'postgres://postgres:password@localhost:5432/booksite?sslmode=disable'
     ```

3. **Run database migrations**:
   - Migrations are typically run automatically on startup
   - Manual migration files are in `backend/internal/migrations/sql/`

4. **Set environment variables** (optional, can use config.yaml):
   ```bash
   export DATABASE_DSN="postgres://postgres:password@localhost:5432/booksite?sslmode=disable"
   export HTTPSERVER_PORT=8080
   export HTTPSERVER_FRONTEND_URL=http://localhost:3000
   export JWT_SECRET=your-secret-key
   export GOOGLE_CLIENT_ID=your-google-client-id
   export GOOGLE_CLIENT_SECRET=your-google-client-secret
   ```

5. **Run the backend**:
   ```bash
   cd backend/cmd
   go run main.go
   ```

### Frontend Setup

1. **Install dependencies**:
   ```bash
   cd frontend/book-library
   npm install
   ```

2. **Configure API endpoint**:
   - Update `src/api.ts` if your backend runs on a different port:
     ```typescript
     const API_URL = 'http://localhost:8080/api'
     ```

3. **Start development server**:
   ```bash
   npm start
   ```

4. **Build for production**:
   ```bash
   npm run build
   ```

## 📁 Project Structure

```
booksite/
├── backend/                 # Go backend application
│   ├── cmd/
│   │   └── main.go        # Application entry point
│   ├── internal/
│   │   ├── config/        # Configuration management
│   │   ├── httpserver/    # HTTP handlers and routing
│   │   │   ├── middleware/ # Authentication middleware
│   │   │   └── *.go       # Handler files
│   │   ├── jwt/           # JWT token utilities
│   │   ├── logger/        # Logging configuration
│   │   ├── migrations/    # Database migrations
│   │   ├── oauth/         # OAuth providers (Google)
│   │   ├── repo/          # Database repository layer
│   │   └── service/       # Business logic layer
│   └── go.mod
├── frontend/
│   └── book-library/      # React frontend application
│       ├── src/
│       │   ├── components/ # React components
│       │   ├── store/      # State management
│       │   ├── api.ts      # API client
│       │   └── App.tsx     # Main app component
│       ├── public/         # Static assets
│       └── package.json
├── migrations/             # Root-level migrations (if any)
├── uploads/               # Uploaded book covers
├── config.yaml            # Configuration file
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile.backend     # Backend Docker image
└── README.md             # This file
```

## 🔐 Authentication

The application supports two authentication methods:

1. **Email/Password**:
   - Register: `POST /api/register`
   - Login: `POST /api/login`

2. **Google OAuth**:
   - Configure Google OAuth credentials in `config.yaml` or environment variables
   - OAuth flow is handled automatically by the backend

JWT tokens are stored in localStorage on the frontend and automatically included in API requests.

## 📚 API Endpoints

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - Login user
- `GET /api/auth/google` - Initiate Google OAuth
- `GET /api/auth/google/callback` - Google OAuth callback

### Books
- `GET /api/books` - Get all books
- `GET /api/books/uncategorized` - Get uncategorized books
- `POST /api/books/create` - Create new book (multipart/form-data)
- `DELETE /api/books/delete?id={id}` - Delete book
- `POST /api/books/assign` - Assign book to category

### Categories
- `GET /api/categories/with-books` - Get categories with their books
- `POST /api/categories/create` - Create category (requires auth)
- `PUT /api/categories/rename?id={id}` - Rename category (requires auth)
- `DELETE /api/categories/delete?id={id}` - Delete category (requires auth)

## ⚙️ Configuration

Configuration can be provided via:
1. `config.yaml` file (see `config.yaml` for structure)
2. Environment variables (take precedence)
3. `config.json` file

Key configuration options:
- Database connection string
- JWT secrets and TTL
- Google OAuth credentials
- Server port and URLs
- Logger settings

## 🗄️ Database

The application uses PostgreSQL with the following main entities:
- **users**: User accounts
- **books**: Book information with cover images
- **categories**: Book categories
- **book_category_relations**: Many-to-many relationship between books and categories

Migrations are located in `backend/internal/migrations/sql/` and are automatically applied on startup.

## 🐳 Docker Services

- **postgres**: PostgreSQL 16 database
- **backend**: Go backend API server
- **frontend**: React frontend served via Nginx

All services are connected via a Docker network and configured with health checks.

## 📝 Development Notes

- The frontend uses TypeScript for type safety
- Backend follows a layered architecture (handlers → services → repository)
- File uploads (book covers) are stored in the `uploads/covers/` directory
- JWT tokens have configurable TTL (default: 15 minutes access, 7 days refresh)
- CORS is configured to allow requests from the frontend URL

## 🔍 Troubleshooting

### Common Issues

1. **Database connection errors**:
   - Ensure PostgreSQL is running
   - Check database credentials in `config.yaml`
   - Verify database exists: `createdb booksite`

2. **Port conflicts**:
   - Change ports in `docker-compose.yml` or `config.yaml`
   - Frontend: 3000, Backend: 8080, PostgreSQL: 5432

3. **Authentication issues**:
   - Verify JWT secrets are set correctly
   - Check token expiration settings
   - Clear browser localStorage if tokens are corrupted

4. **File upload issues**:
   - Ensure `uploads/` directory exists and has write permissions
   - Check file size limits in backend configuration

## 📄 License

This project is proprietary and not open for public use or redistribution.

## 🙏 Acknowledgments

Built with:
- [React](https://react.dev/)
- [Material-UI](https://mui.com/)
- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
