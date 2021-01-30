# ShareSecret

[![Test](https://github.com/bernardosecades/sharesecret/workflows/Test/badge.svg)](https://github.com/bernardosecades/sharesecret/actions)
[![Super-Linter](https://github.com/bernardosecades/sharesecret/workflows/Super-Linter/badge.svg)](https://github.com/bernardosecades/sharesecret/actions)

ShareSecret is a service to share sensitive information that's both simple and secure.

If you share some text will be display it once and then delete it. After that it's gone forever.

We keep secrets for up to 5 days.

## Why should I trust you?

General we can't do anything with your information even if we wanted to (which we don't). If it's a password for example, we don't know the username or even the application that the credentials are for.

If you include a password, we use it to encrypt the secret. We don't store the password (only a crypted hash) so we can never know what the secret is because we can't decrypt it.

## Create new secret

POST: `http://localhost:8080/secret`

### Header (optional)

```
X-Password: "MyPassword"
```

### Payload

```json
{
    "content": "This is my secret"
}
```

### Example

Without password:

`curl -X POST http://localhost:8080/secret -d "{\"content\":\"This is my secret\"}"`

With password:

`curl -X POST http://localhost:8080/secret -H "X-Password: myPassword" -d "{\"content\":\"This is my secret\"}"`

## See secret

### Header (optional)

```
X-Password: "MyPassword"
```

### Request

GET: `http://localhost:8080/secret/{id}`

### Example

Without password:

`curl http://127.0.0.1:8080/secret/b3eb17a5-bda5-4e83-9690-56967857d03e`

With password:

`curl -H "X-Password: myPassword" http://127.0.0.1:8080/secret/19d38f65-18c3-4d06-9685-9b705ee9d734`

# Use case (Generate secret without password)

UserA create a new secret:

Request:

`curl -X POST http://localhost:8080/secret -d "{\"content\":\"This is my secret for userB\"}"`

Response:

```json
{
    "url": "http://127.0.0.1:8080/secret/90055dba-36aa-4572-8bb0-e4a1f8ecdf54"
}
```

And share the link to UserB to open the link and see the content of the secret:

Request:

`curl http://127.0.0.1:8080/secret/90055dba-36aa-4572-8bb0-e4a1f8ecdf54`

Response:

```json
{
    "content": "This is my secret for userB"
}
```

If UserB or whoever try to access again to the secret:

`curl http://127.0.0.1:8080/secret/90055dba-36aa-4572-8bb0-e4a1f8ecdf54`

He will receive the response not found because already was viewed:

```json
{
    "StatusCode": 404,
    "Error": "sql: no rows in result set"
}
```

# Use case (Generate secret with password)

WIP

# Builds

## Server

Up the server with REST API.

`cd cmd/server && go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"`

## Purge

Command to delete in database secrets already expided.

`cd cmd/purge && go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"`

# Docker

Up the database:

`docker-compose up --build`

Build with version:

`go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"`

Run service:

`./sharesecret`


Travis
https://dave.cheney.net/2018/07/16/using-go-modules-with-travis-ci