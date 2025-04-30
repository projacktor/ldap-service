# LDAP sevice
## Innopolis University S25 Distributed and Network Programming course project

Application shows the implementation of remote authentication using OIDC over the ldap protocol using lldap and keycloak with backend on Golang and frontend on NextJS.

## Requirements

- Docker
- Docker Compose
- Pnpm

## Configuration

Before starting, ensure that the `.env` file is downloaded from group

File .env for directory backend_go/
```bash
# .env for backend
HOST=... # Backend host
PORT=8081 # Backend port

KAFKA_BROKERS=... # For example kafka:9092
KAFKA_TOPIC=... # Kafka topic
KAFKA_CONSUMER_GROUP=... # Group for Kafka

PROMETHEUS_PORT=9090 # Port for Prometheus (9090 default)

KEYCLOAK_BASE_URL=... # Address of keycloak service
KEYCLOAK_REALM=... # Your Realm
KEYCLOAK_REST_API_CLIENT_ID=... # Client's ID
KEYCLOAK_REST_API_CLIENT_SECRET=... # Client's secret

LDAP_HOST=... # Address of ldap service. Use ldap://...
LDAP_BASE_DN=... # Configuration for Ldap
LDAP_USER_DN=... # Configuration for Ldap
LDAP_USER_PASSWORD=... # Password for Ldap
```

File .env.local for directory frontend/
```bash
# .env.local for frontend
NEXT_PUBLIC_BACKEND_IP=... # For example, http://localhost:8081 (address for backend)
```

## Launch

### Frontend

1. Using for installing dependencies :
```bash
cd frontend/
pnpm install
```

2. Using for starting frontend on port ":3000":
```bash
pnpm dev -p 3000
``` 

### Backend
1. Build and start the containers (the launch will take approximately 2 minutes):
```bash
cd backend_go
docker compose up --build -d
```

2. Check container status:
```bash
docker compose ps
```

## Available Services
- Frontend: `http://localhost:3000`
- Go application: `http://localhost:8081`
- Kafka UI: `http://localhost:8082`

## Project Structure for backend

- `cmd/` - application entry point
- `config/` - configuration
- `auth/` - authentication
- `middleware/` - middleware components

## Logging

Logs are sent to the Kafka topic `logs`. To view logs, use Kafka UI `http://localhost:8082`

## Usage
Using my "Avaliable Services" I go to the http://localhost:3000 and see box with username and password fields. Fill it and press "Sign in". If Lldap has your username and password and you input valid data you will get page with timer, your name, email, amd roles.