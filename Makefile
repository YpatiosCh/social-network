DB_CONTAINER=forum_db
INIT_SQL=database/init.sql
CATEGORIES_SQL=database/categories.sql

.PHONY: run-backend run-frontend run-database db-up db-down db-logs db-psql db-init run-all db-reset mkcert-local

run-backend:
	go run backend/app/main.go


run-frontend:
	go run frontend/app/main.go

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-logs:
	docker logs -f $(DB_CONTAINER)

db-psql:
	docker exec -it $(DB_CONTAINER) psql -U server -d forum_db

db-init:
	docker exec -i $(DB_CONTAINER) psql -U server -d forum_db < $(INIT_SQL)
	docker exec -i $(DB_CONTAINER) psql -U server -d forum_db < $(CATEGORIES_SQL)

db-populate:
	go run backend/internal/populate/mock_db.go
	
db-categories:
	docker exec -i $(DB_CONTAINER) psql -U server -d forum_db < $(CATEGORIES_SQL)

db-reset:
	@echo "Resetting forum_db..."
	docker exec -i $(DB_CONTAINER) psql -U server -d postgres -c "DROP DATABASE IF EXISTS forum_db;"
	docker exec -i $(DB_CONTAINER) psql -U server -d postgres -c "CREATE DATABASE forum_db;"
	make db-init
	@echo "Database reset complete."

run-all:
	@echo "Starting database..."; \
	make db-up; \
	echo "Waiting for Postgres to be healthy..."; \
	until docker compose exec -T postgres pg_isready >/dev/null 2>&1; do \
		sleep 1; \
	done; \
	echo "Postgres is ready!"; \
	make run-backend & \
	BACKEND_PID=$$!; \
	make run-frontend & \
	FRONTEND_PID=$$!; \
	wait $$BACKEND_PID $$FRONTEND_PID

# First install mkcert and nss
# Windows:
# 		choco install mkcert
# 		choco install nss -y 
#MacOS:
# 		brew install mkcert
# 		brew install nss
# Linux:
# 		curl -LO "https://github.com/FiloSottile/mkcert/releases/latest/download/mkcert-v1.4.5-linux-amd64"
# 		chmod +x mkcert-v1.4.5-linux-amd64
# 		sudo mv mkcert-v1.4.5-linux-amd64 /usr/local/bin/mkcert
# 		sudo apt install libnss3-tools

mkcert-local:
	@echo "Installing local CA..."
	mkcert -install
	@echo "Generating localhost certificate..."
	mkcert localhost 127.0.0.1 ::1
	@echo "Done! Certificates generated for localhost."
