FROM social-network/go-base AS build

WORKDIR /app

COPY . .

RUN go build -o users_service ./services/users/cmd/server

RUN go build -o migrate ./services/users/cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=build /app/users_service .
COPY --from=build /app/migrate .
COPY --from=build /app/services/users/internal/db/migrations ./migrations
COPY --from=build /app/services/users/internal/db/seeds ./seeds

RUN chmod +x /app/seeds/seed.sh

CMD ["./users_service"]
