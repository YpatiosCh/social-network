# retrieveusers package
## Overview

The retrieveusers package provides a batched, cached user lookup mechanism for services that store only user IDs locally but need to return basic user information in API responses.

## Responsibilities

The package is responsible for:

- Retrieving basic user information by ID

- Deduplicating lookup requests

- Using Redis as a read-through cache

- Falling back to the Users service via a batch RPC

- Returning results as a map[int64]models.User

The package does not:

- Mutate caller-owned data structures

- Know about posts, comments, conversations, or domain models

- Perform authorization checks

- Expose transport-specific details to callers

## Data model

The retriever returns the following user model:

```
type User struct {
    UserId   ct.Id
    Username ct.Username
    AvatarId ct.Id
}
```

This represents the minimal user projection required by non-Users services.

## Architecture
```
Caller
  |
  |-- collect user IDs
  |
  |-- UserRetriever.GetUsers(ctx, ids)
          |
          |-- Redis (cache hit)
          |-- Users service (batch RPC, cache miss)
          |-- Redis (write-through)
  |
  |-- caller assigns users into response objects
```

## Public API
```UserRetriever```
```
type UserRetriever struct {
    clients UsersBatchClient
    cache   RedisCache
    ttl     time.Duration
}
```
```Constructor```
```
func NewUserRetriever(
    clients UsersBatchClient,
    cache RedisCache,
    ttl time.Duration,
) *UserRetriever
```
### Parameters
|Name|	Description|
|--|--|
|clients|	Abstraction over the Users service batch RPC
|cache|	Redis-compatible cache implementation
|ttl	|Cache TTL for user entries


```GetUsers```
```
func (r *UserRetriever) GetUsers(
    ctx context.Context,
    ids []int64,
) (map[int64]models.User, error)
```
### Behavior

- Deduplicates input IDs

- Attempts to retrieve users from Redis

- Fetches missing users via a single batch RPC

- Caches retrieved users with the configured TTL

- Returns a map keyed by user ID

### Guarantees

- At most one RPC call per invocation

- Redis is always checked before RPC

- Order of input IDs is preserved by the caller, not the retriever

## Required interfaces
### UsersBatchClient

The caller must supply a client that implements:
```
type UsersBatchClient interface {
    GetBatchBasicUserInfo(
        ctx context.Context,
        userIds []int64,
    ) (*userpb.ListUsers, error)
}
```

This allows each service to use its own Users client implementation.

### RedisCache

The cache must implement:
```
type RedisCache interface {
    GetObj(ctx context.Context, key string, dest any) error
    SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}
```

Any Redis adapter matching this contract may be used.

## Service integration
### Application wiring
```
cache := redis_connector.NewRedisClient("localhost:6379", "", 0)

userRetriever := retrieveusers.NewUserRetriever(
    clients,        // UsersBatchClient
    cache,          // RedisCache
    3*time.Minute,  // TTL
)

app := &Application{
    db:          db,
    clients:     clients,
    userFetcher: userRetriever,
}
```
### Usage patterns
## Hydrating a slice of user IDs
```
userMap, err := s.userFetcher.GetUsers(ctx, ids)
if err != nil {
    return nil, err
}

users := make([]models.User, 0, len(ids))
for _, id := range ids {
    if u, ok := userMap[id]; ok {
        users = append(users, u)
    }
}
```
## Hydrating posts
```
var ids []int64
for _, p := range posts {
    ids = append(ids, p.User.UserId.Int64())
}

userMap, err := s.userFetcher.GetUsers(ctx, ids)
if err != nil {
    return nil, err
}

for i := range posts {
    uid := posts[i].User.UserId.Int64()
    posts[i].User = userMap[uid]
}
```

## Hydrating nested users (e.g. conversations)
```
var ids []int64
for _, c := range conversations {
    for _, m := range c.Members {
        ids = append(ids, m.UserId.Int64())
    }
}

userMap, err := s.userFetcher.GetUsers(ctx, ids)
if err != nil {
    return nil, err
}

for i := range conversations {
    for j := range conversations[i].Members {
        uid := conversations[i].Members[j].UserId.Int64()
        conversations[i].Members[j] = userMap[uid]
    }
}
```
### Error handling semantics

|Condition	|Behavior
|--|--|
|Redis miss	|Fetch from Users service
|Redis failure	|Return error
|Users service failure|	Return error
|Cache write failure|	Logged, ignored
|Missing user ID|	Silently skipped

### Design rationale

- Map-based return avoids implicit mutation

- Batch-only API prevents N+1 RPC calls

- Service-local clients avoid shared dependencies

- Minimal interface contracts enable easy mocking

### Testing guidance

Mock only:

```UsersBatchClient```

```RedisCache```

The retriever contains no domain logic and can be tested independently from application code.

### Summary

retrieveusers provides a reusable, efficient, and decoupled mechanism for resolving user IDs into basic user data across services.

It intentionally avoids domain coupling and object mutation, allowing callers to retain full control over response assembly.