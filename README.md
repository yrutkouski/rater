# Rater

## Stack

- **Node.js** 22
- **Go** 1.21

## Local Development

### Running the Application

**1. Frontend with mocked backend**
```bash
make frontend
```
- Frontend: `http://localhost:3000`
- Mock API: `http://localhost:8080`

**2. Backend with database**
```bash
make backend
```
- Backend: `http://localhost:8080`
- Database: `localhost:5432` (Docker, data persists)

**3. Full application**
```bash
make application
```
- Frontend: `http://localhost:3000` (dev mode, local)
- Backend: `http://localhost:8080` (Docker)
- Database: `localhost:5432` (Docker, data persists)
- Backend and database in Docker, frontend runs locally

### Commands

```bash
make frontend     # Frontend + mock API
make backend      # Backend + database (Docker DB, local BE)
make application  # Full app in Docker
make stop         # Stop all Docker services
make clean        # Clean build artifacts and Docker volumes
```

### API Endpoints

- `GET /api/health` - Health check
- `GET /api/data` - Data endpoint
