# LDAP Service Backend

Go application for working with LDAP and logging to Kafka.

## Requirements

- Docker
- Docker Compose

## Configuration

Before starting, ensure that the `.env` file is downloaded from group

## Launch

1. Build and start the containers (the launch will take approximately 2 minutes):
```bash
docker compose up --build -d
```

2. Check container status:
```bash
docker compose ps
```

## Available Services

- Go application: `http://localhost:8081`
- Kafka UI: `http://localhost:8082`

## Project Structure

- `cmd/` - application entry point
- `config/` - configuration
- `auth/` - authentication
- `middleware/` - middleware components

## Logging

Logs are sent to the Kafka topic `logs`. To view logs, use Kafka UI `http://localhost:8082` or the console consumer:

```bash
docker compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic logs --from-beginning
```

## Stop

```bash
docker compose down -v --remove-orphans
```
