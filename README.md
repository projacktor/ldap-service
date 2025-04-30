# 🛡️ LDAP sevice

## Innopolis University S25 Distributed and Network Programming course project

The application shows the implementation of remote authentication using OIDC over the ldap protocol using lldap and
keycloak with backend on Golang and frontend on NextJS.

## 🔧 Tech Stack

[![Go][Go]][Go-url]
[![Docker][Docker]][Docker-url]
[![Docker Compose][Docker-Compose]][Docker-Compose-url]
[![Kubernetes][Kubernetes]][k8s-url]
[![Nginx][Nginx]][nginx-url]
[![Keycloak][Keycloak]][keycloak-url]
[![LLDAP][LLDAP]][lldap-url]
[![Kafka][Kafka]][Kafka-url]

[![React][React]][react-url]
[![NextJS][Nextjs]][Next-url]
[![TypeScript][TypeScript]][ts-url]
[![Tailwind][Tailwind CSS]][Tailwind-url]
[![ShadCN UI][Shadcnui]][shadcn-url]
[![React Hook Form][RHF]][rhf-url]
[![pnpm][pnpm]][pnpm-url]
[![Prettier][prettier]][prettier-url]

## 🚀 Installation requirements

> - Docker
> - Docker Compose
> - pnpm

## ⚙️ Configuration

### 🔙 Backend

Before starting, ensure that the `.env` file is downloaded from group

File .env for directory backend_go/

```bash
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

### 🌐 Frontend

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

Using my "Avaliable Services" I go to the http://localhost:3000 and see box with username and password fields. Fill it
and press "Sign in". If Lldap has your username and password and you input valid data you will get page with timer, your
name, email, amd roles.

## License

We use [SOL License](LICENSE) to protect against the plagiarism by students and [MIT License](LICENSE-MIT) for the other
projects to use


[Go]: https://img.shields.io/badge/Go-000000?style=for-the-badge&logo=go

[Go-url]: https://go.dev/

[RHF]: https://img.shields.io/badge/React_Hook_Form-000000?style=for-the-badge&logo=reacthookform

[rhf-url]: https://react-hook-form.com/

[Kubernetes]: https://img.shields.io/badge/Kubernetes-000000?style=for-the-badge&logo=kubernetes

[k8s-url]: https://kubernetes.io/

[Nginx]: https://img.shields.io/badge/Nginx-000000?style=for-the-badge&logo=nginx

[nginx-url]: https://nginx.org/

[Keycloak]: https://img.shields.io/badge/Keycloak-000000?style=for-the-badge&logo=keycloak

[keycloak-url]: https://www.keycloak.org/

[LLDAP]: https://img.shields.io/badge/lldap-000000?style=for-the-badge&logo=openldap

[lldap-url]: https://github.com/lldap/lldap

[Python]: https://img.shields.io/badge/Python_3.12-000000?style=for-the-badge&logo=python

[Python-url]: https://www.python.org/downloads/

[uv]: https://img.shields.io/badge/uv-000000?style=for-the-badge&logo=python

[uv-url]: https://github.com/astral-sh/uv

[FastAPI]: https://img.shields.io/badge/FastAPI-000000?style=for-the-badge&logo=fastapi

[FastAPI-url]: https://fastapi.tiangolo.com/

[Pydantic]: https://img.shields.io/badge/Pydantic-000000?style=for-the-badge&logo=pydantic

[Pydantic-url]: https://docs.pydantic.dev/latest/

[MWS-GPT-API]: https://img.shields.io/badge/MWS_GPT_API-000000?style=for-the-badge&logo=openai

[MWS-GPT-API-url]: https://api.gpt.mws.ru/

[LangChain]: https://img.shields.io/badge/LangChain-000000?style=for-the-badge&logo=langchain

[LangChain-url]: https://www.langchain.com/

[Ruff]: https://img.shields.io/badge/Ruff-000000?style=for-the-badge&logo=ruff

[Ruff-url]: https://docs.astral.sh/ruff/

[pre-commit]: https://img.shields.io/badge/pre--commit-000000?style=for-the-badge&logo=pre-commit

[pre-commit-url]: https://pre-commit.com/

[Docker]: https://img.shields.io/badge/Docker-000000?style=for-the-badge&logo=docker

[Docker-url]: https://www.docker.com/

[Docker-Compose]: https://img.shields.io/badge/Docker_Compose-000000?style=for-the-badge&logo=docker

[Docker-Compose-url]: https://docs.docker.com/compose/

[NextJS]: https://img.shields.io/badge/Next-black?style=for-the-badge&logo=next.js&logoColor=white

[Next-url]: https://nextjs.org/

[Tailwind CSS]: https://img.shields.io/badge/tailwind-000000?style=for-the-badge&logo=tailwindCSS

[Tailwind-url]: https://tailwindcss.com/

[pnpm]: https://img.shields.io/badge/pnpm-000000.svg?style=for-the-badge&logo=pnpm&logoColor=f69220

[pnpm-url]: https://pnpm.io/

[TypeScript]: https://img.shields.io/badge/typescript-000000.svg?style=for-the-badge&logo=typescript&logoColor=white

[ts-url]: https://www.typescriptlang.org/

[Shadcnui]: https://img.shields.io/badge/shadcn/ui-000000.svg?style=for-the-badge&2F&logo=shadcnui&color=131316

[shadcn-url]: https://ui.shadcn.com/

[json]: https://img.shields.io/badge/json-000000.svg?style=for-the-badge&logo=json&logoColor=white

[json-url]: https://www.json.org/json-en.html

[React]: https://img.shields.io/badge/react-000000.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB

[react-url]: https://react.dev/

[react-query]: https://img.shields.io/badge/React_Query-000000.svg?style=for-the-badge&logo=ReactQuery&logoColor=white

[rq-url]: https://tanstack.com/query/latest/docs/framework/react/overview

[prettier]: https://img.shields.io/badge/prettier-000000.svg?style=for-the-badge&logo=prettier&logoColor=F7BA3E

[prettier-url]: https://prettier.io/

[Kafka]: https://img.shields.io/badge/Kafka-000000?style=for-the-badge&logo=apachekafka&logoColor=white

[Kafka-url]: https://kafka.apache.org/