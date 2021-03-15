# ShareSecret

[![Test](https://github.com/bernardosecades/sharesecret/workflows/Test/badge.svg)](https://github.com/bernardosecades/sharesecret/actions)
[![Super-Linter](https://github.com/bernardosecades/sharesecret/workflows/Super-Linter/badge.svg)](https://github.com/bernardosecades/sharesecret/actions)

ShareSecret is a service to share sensitive information that's both simple and secure.

If you share some text will be display it once and then delete it. After that it's gone forever.

We keep secrets for up to 5 days.

## Why should I trust you?

General we can't do anything with your information even if we wanted to (which we don't). If it's a password for example, we don't know the username or even the application that the credentials are for.

If you include a password, we use it to encrypt the secret. We don't store the password (only a crypted hash) so we can never know what the secret is because we can't decrypt it.

## Demo

![Share Secret Demo](demo/demo.gif)

## Create new secret

POST: `http://localhost:8080/secret`

### Header (optional)

```bash
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

```bash
curl -X POST http://localhost:8080/secret -d "{\"content\":\"This is my secret\"}"
```

With password:

```bash
curl -X POST http://localhost:8080/secret -H "X-Password: myPassword" -d "{\"content\":\"This is my secret\"}"
```

## See secret

### Header (optional)

```bash
X-Password: "MyPassword"
```

### Request

GET
```bash 
http://localhost:8080/secret/{id}
```

### Example

Without password:

```bash 
curl http://127.0.0.1:8080/secret/b3eb17a5-bda5-4e83-9690-56967857d03e
```

With password:

```bash
curl -H "X-Password: myPassword" http://127.0.0.1:8080/secret/19d38f65-18c3-4d06-9685-9b705ee9d734
```

## Docker Compose

Up webserver and MySQL:

```bash
docker-compose up --build
```

Now you will can access to: localhost:8080/secret (POST/GET)

If you execute the command:

```bash
docker-compose ps
```

You will see two containers:

- sharesecret_mysql_1
- sharesecret_web_1

In the container "sharesecret_web_1" we already compile two binary (server and purge) in the Dockerfile and run the server. If
you want to execute the binary "purge" (delete from database all secrets expired):

```bash
docker exec -it sharesecret_web_1 ./cmd/purge/purge
```

Note: You can compile the sever and purge using flags for version:

```bash
cd cmd/server && go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"
```

```bash
cd cmd/purge && go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"
```

```bash
docker build -t sharesecret-server .
docker run --env-file=../.ENV_PROD_SHARESECRET --rm -d -p 8080:8080 sharesecret-server
```

Remove docker images not used by any container:

```bash
docker system prune -a 
```

See ip container, you can check with `docker inspect <container-ID>`


# PROTO

Install: 

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc

See example: https://github.com/TutorialEdge/go-grpc-tutorial

Generate secret.pb.go and secret_grpc.pb.go:

Without gateway (only gRPC):
```bash
protoc -I=proto --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. proto/secret.proto  
```

With gateway (API REST):

```bash 
protoc -I=proto -I /Users/admin/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. proto/secret.proto --grpc-gateway_out=logtostderr=true:./genproto 
```

https://blog.friendsofgo.tech/posts/introduccion-a-grpc/


REFERENCIA:

https://medium.com/@nate510/structuring-go-grpc-microservices-dd176fdf28d0 -> HERE https://medium.com/@nate510/structuring-go-grpc-microservices-dd176fdf28d0

Este sobre todo: 

https://blog.friendsofgo.tech/posts/introduccion-a-grpc/

