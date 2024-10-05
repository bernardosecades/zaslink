# ShareSecret

[![Build](https://github.com/bernardosecades/sharesecret/actions/workflows/build.yml/badge.svg)](https://github.com/bernardosecades/sharesecret/actions/workflows/build.yml)

ShareSecret is a service to share sensitive information that's both simple and secure.

If you share some text will be display it once and then delete it. After that it's gone forever.

We keep secrets for up to 24 hours by default.

## Why should I trust you?

General we can't do anything with your information even if we wanted to (which we don't). If it's a password for example, we don't know the username or even the application that the credentials are for.

If you include a password, we use it to encrypt the secret. We don't store the password (only a crypted hash) so we can never know what the secret is because we can't decrypt it.

## Open API v3 - Swagger UI

`make run-openapi-ui`

`http://localhost:4000/`

## Structure folder

Based on: https://github.com/golang-standards/project-layout

## Docker

Build image with tag version:

`sudo docker build --tag bernardosecades/api-share-secret:latest`

Build container from image (Example):

`docker run --name sharesecret-api -p 8080:8080 -e SECRET_KEY=11111111111111111111111111111111 -e DEFAULT_PASSWORD=@myPassword -e MONGODB_URI=mongodb://root:example@bernardosecades.com:27017 -e MONGODB_NAME=share_secret bernardosecades/api-share-secret:latest`

Push to docker hub:

`sudo docker login -u bernardosecades -p <YOUR_PASSWORD>`

`sudo docker push bernardosecades/api-share-secret:latest`