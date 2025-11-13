### Run containers 

```
docker compose build
docker compose up
```



### Access databases with
```
docker-compose exec users-db psql -U postgres -d social_users
```

```
docker-compose exec chat-db psql -U postgres -d social_chat
```

```
docker-compose exec forum-db psql -U postgres -d social_forum
```

