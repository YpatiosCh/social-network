# ---- Build stage ----
FROM social-network/go-base AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build chat service binary
RUN go build -o chat_service ./services/chat

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binary and migrations
COPY --from=build /app/chat_service .
COPY --from=build /app/services/chat/migrations ./migrations

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@chat-db:5432/social_chat?sslmode=disable

CMD ["./chat_service"]
