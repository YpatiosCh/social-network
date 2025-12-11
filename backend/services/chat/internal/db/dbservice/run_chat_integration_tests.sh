#!/usr/bin/env bash
set -euo pipefail

# SCRIPT LOCATION:
# backend/services/chat/internal/db/dbservice/run_chat_integration_tests.sh

# ---------------------------------------------------------
# Resolve directories reliably
# ---------------------------------------------------------

# dbservice/
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# db/  (has migrations folder)
DB_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# repo root: 5 levels up from dbservice/
ROOT_DIR="$(cd "$SCRIPT_DIR/../../../../.." && pwd)"

cd "$ROOT_DIR"

# ---------------------------------------------------------
# Start chat-db
# ---------------------------------------------------------

echo "Starting chat-db via docker-compose..."
docker-compose up -d chat-db

echo "Waiting for chat-db container to become ready..."
for i in {1..30}; do
  if docker exec chat-db pg_isready -U postgres -d social_chat >/dev/null 2>&1; then
    echo "chat-db is ready"
    break
  fi
  echo "Waiting ($i/30)..."
  sleep 2
done

# ---------------------------------------------------------
# Database URL
# ---------------------------------------------------------

DB_URL="postgres://postgres:secret@127.0.0.1:5435/social_chat?sslmode=disable"
export DATABASE_URL="$DB_URL"
export TEST_DATABASE_URL="$DB_URL"

echo "Running migrations (DATABASE_URL=$DATABASE_URL)"

# ---------------------------------------------------------
# IMPORTANT PART:
# Run migrations *from the directory that contains ./migrations*
# ---------------------------------------------------------

cd "$DB_DIR"     # <--- THIS is the key fix

# Now `./migrations` resolves to:
# backend/services/chat/internal/db/migrations   (correct)

go run ../../cmd/migrate/main.go

# Return to repo root
cd "$ROOT_DIR"

# ---------------------------------------------------------
# Run integration tests
# ---------------------------------------------------------

echo "Running chat DB tests..."
go test ./services/chat/internal/db/dbservice -v "$@"

echo "Integration tests complete."
echo "To stop DB: docker-compose down chat-db"
