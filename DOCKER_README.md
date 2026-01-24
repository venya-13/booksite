# Docker Setup Instructions

This project uses Docker Compose to run PostgreSQL, Backend (Go), and Frontend (React) services together.

## Prerequisites

- Docker Desktop installed and running
- Docker Compose v3.8 or higher

## Quick Start

1. **Build and start all services:**
   ```bash
   docker-compose up --build
   ```

2. **Run in detached mode (background):**
   ```bash
   docker-compose up -d --build
   ```

3. **Stop all services:**
   ```bash
   docker-compose down
   ```

4. **Stop and remove volumes (clears database):**
   ```bash
   docker-compose down -v
   ```

## Services

- **PostgreSQL**: Available at `localhost:5432`
  - Database: `booksite`
  - User: `postgres`
  - Password: `postgres`
  - Persistent volume: `postgres_data`

- **Backend (Go)**: Available at `http://localhost:8080`
  - API endpoints: `http://localhost:8080/api/*`
  - Uploads directory: `./uploads` (mounted as volume)

- **Frontend (React)**: Available at `http://localhost:3000`
  - Proxies `/api` requests to backend
  - Built with nginx in production mode

## Environment Variables

You can override environment variables by creating a `.env` file or setting them in the shell:

```bash
export GOOGLE_CLIENT_ID=your-client-id
export GOOGLE_CLIENT_SECRET=your-client-secret
docker-compose up
```

Or modify `docker-compose.yml` directly.

## Database Migrations

The backend should run migrations automatically on startup. If you need to run migrations manually, you can exec into the backend container:

```bash
docker-compose exec backend ./main
```

## Development

For development, you may want to:

1. Mount source code as volumes for hot-reloading (not included in default setup)
2. Use separate docker-compose files for dev/prod
3. Override environment variables for local development

## Troubleshooting

1. **Port already in use**: Change port mappings in `docker-compose.yml`
2. **Database connection errors**: Wait for postgres healthcheck to pass (10-15 seconds)
3. **Build failures**: Check Docker logs: `docker-compose logs [service-name]`
4. **Volume permission issues**: Ensure uploads directory has correct permissions

## Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

## Rebuilding After Changes

```bash
# Rebuild specific service
docker-compose build backend
docker-compose up -d backend

# Rebuild all
docker-compose up --build
```
