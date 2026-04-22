# ksk8s - kubeasz Web Management Platform

A web-based K8s cluster installation and lifecycle management platform built on top of [kubeasz](https://github.com/easzlab/kubeasz).

## Architecture

- **Frontend**: Vue 3 + Element Plus + xterm.js
- **Backend**: Go + Gin + GORM
- **Database**: MySQL 8
- **Deployment**: Docker Compose on ignition machine

## Features

- Visual cluster creation wizard with template system
- 7-step installation pipeline (01-07) with real-time log streaming
- Admin approval between each installation step
- Dual-mode config editor: structured form + text (INI/YAML)
- Lifecycle operations: start, stop, upgrade, backup, restore, destroy
- Node management: add/del etcd, master, worker nodes
- Security ops: CA certificate renewal
- Embedded WebSSH terminal
- Full audit trail and persistent log storage

## Quick Start

### Prerequisites

- Docker + Docker Compose
- kubeasz installed at `/etc/kubeasz`
- SSH key for node access

### Development

```bash
# Start all services
docker-compose up -d

# Backend only
cd ksk8s-backend
go run cmd/server/main.go

# Frontend only
cd ksk8s-frontend
npm run dev
```

### Production

```bash
# Copy and edit environment variables
cp .env.example .env
# Edit .env with your DB password and secrets

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MYSQL_ROOT_PASSWORD` | MySQL root password | - |
| `MYSQL_PASSWORD` | MySQL app password | - |
| `JWT_SECRET` | HS256 JWT secret (auto-generated if empty) | - |
| `KSK8S_SSH_KEY` | SSH private key path | `/root/.ssh/id_rsa` |
| `KSK8S_MOCK_EZCTL` | Use mock ezctl for testing | `0` |

## API Endpoints

- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register
- `GET /api/clusters` - List clusters
- `POST /api/clusters` - Create cluster
- `GET /api/clusters/:id/tasks` - List tasks
- `POST /api/clusters/:id/steps/:step/run` - Run step
- `POST /api/clusters/:id/tasks/:taskId/approve` - Approve task
- `POST /api/clusters/:id/tasks/:taskId/abort` - Abort task
- `WS /ws/tasks/:id/logs` - Stream logs via WebSocket
- `WS /ws/ssh` - WebSSH terminal

## Testing

```bash
cd ksk8s-backend
go test ./...
```

## Project Structure

```
ksk8s-backend/
  cmd/server/          # Main entrypoint
  internal/
    handler/           # HTTP handlers
    middleware/        # Auth middleware
    model/             # GORM models
    repository/        # Data access layer
    service/           # Business logic
    websocket/         # WebSocket hub + ring buffer
    ssh/               # SSH client
    tls/               # Self-signed cert generation
    config/            # Boot-time config
  migrations/          # SQL migrations

ksk8s-frontend/
  src/
    api/               # API client
    components/        # Reusable components
    router/            # Vue Router
    stores/            # Pinia stores
    views/             # Page components
```

## License

MIT
