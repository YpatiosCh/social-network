# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build users service binary
RUN go build -o users_service ./services/users

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binary and migrations
COPY --from=build /app/users_service .
COPY --from=build /app/services/users/migrations ./migrations

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@users-db:5432/social_users?sslmode=disable

CMD ["./users_service"]
