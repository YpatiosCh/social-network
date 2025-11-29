FROM social-network/go-base AS build

WORKDIR /app

COPY . .

RUN go build -o notifications_service ./services/notifications/cmd/server

RUN go build -o migrate ./services/notifications/cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=build /app/notifications_service .
COPY --from=build /app/migrate .
COPY --from=build /app/services/notifications/internal/db/migrations ./migrations
# COPY --from=build /app/services/notifications/internal/db/seeds ./seeds

# RUN chmod +x /app/seeds/seed.sh


CMD ["./notifications_service"]
