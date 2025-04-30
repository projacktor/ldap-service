#######################################
# 1) Builder stage
#######################################
FROM golang:1.24-alpine AS builder

# Install git (needed if you fetch modules from private repos) 
# and ca-certificates (for module downloads over HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy only go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Build a statically‑linked binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o ldap-app ./cmd

#######################################
# 2) Runtime stage
#######################################
FROM alpine:3.21

# Needed for TLS (if you ever enable HTTPS or talk to services over TLS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/ldap-app .

# (Optional) copy a default .env if you want to bake one in:
# COPY .env .

# Expose your HTTP+metrics port
EXPOSE 8081

# Run as non‑root for better security
USER nobody:nogroup

ENTRYPOINT ["./ldap-app"]
