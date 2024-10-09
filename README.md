# Zaslink

---

[![Build](https://github.com/bernardosecades/zaslink/actions/workflows/build.yml/badge.svg)](https://github.com/bernardosecades/zaslink/actions/workflows/build.yml)
[![Deploy](https://github.com/bernardosecades/zaslink/actions/workflows/deploy.yml/badge.svg)](https://github.com/bernardosecades/zaslink/actions/workflows/deploy.yml)

## 📋 Table of Contents

1. 🐺 [What is this API ?](#what-is-this-api)
2. ✨ [Production and development links](#production-and-development-links)
3. 🔨 [Installation](#installation)
4. 🐳 [Docker](#docker)
5. 💯 [Tests](#tests)
6. 🌿 [Env variables](#env-variables)
7. ☑️ [Code analysis and consistency](#code-analysis-and-consistency)
8. 🐙 [GitHub Actions](#github-actions)
9. ✨ [Misc commands](#misc-commands)
10. ©️ [License](#license)
11. Observability (WIP)

## <a name="what-is-this-api">🐺 What is this API ?</a>

ZasLink secret is a link that can be accessed only once. It’s a single-use URL.

This api allow share sensitive information that's both simple and secure.

Keep sensitive info out of your email and chat logs.

When you share sensitive information, such as passwords or private links, through email or chat, copies of that data can be stored in various locations. By using a one-time link, the information is accessible only once and can't be viewed again by anyone else. This ensures that your sensitive data is shared securely, with only one recipient able to access it. It's similar to sending a self-destructing message.

## <a name="production-and-development-links">✨ Production and development links</a>

### 🌐 Production

The production version of this API is available at **[docs.zaslink.com](https://docs.zaslink.com)**.

The production server is updated automatically with the latest version of the API when a new release is created.

The Web App is WIP: **[www.zaslink.com](https://www.zaslink.com)**. React App: https://github.com/bernardosecades/zaslink-app

### 🛠️ Development

The development version of this API is available at **[http://localhost:4000](http://localhost:4000)**.

## <a name="installation">🔨 Installation</a>

To install this project, you will need to have on your machine :

- Go 1.23
- Docker

Then, run the following commands :

```bash
# Run zaslink swagger ui: http://localhost:4000
make run-openapi-ui

# Run docker compose
docker compose up mongo up -d

# Run zaslink api: http://localhost:800/healthz
DEFAULT_PASSWORD=@myPassword MONGODB_NAME=share_secret MONGODB_URI=mongodb://root:example@localhost:27017 SECRET_KEY=11111111111111111111111111111111 go run ./cmd/api/main.go
```

## <a name="docker">🐳 Docker</a>

This app is Docker ready !

The Dockerfile is available at the root of the project. It uses a multi-stage build to optimize the image size and distroless image to reduce the attack surface.

Build image with tag version:

`sudo docker build --tag bernardosecades/zaslink-api:latest`

Build container from image (Example):

`docker run --name zaslink-api \
  --restart unless-stopped \
  -p 8080:8080 \
  -e SECRET_KEY=11111111111111111111111111111111 \
  -e DEFAULT_PASSWORD=@myPassword \
  -e MONGODB_URI=mongodb://root:example@bernardosecades.com:27017 \
  -e MONGODB_NAME=share_secret \
  bernardosecades/zaslink-api:latest
`

## <a name="tests">💯 Tests</a>

### 🧪 Unit and Integration tests

To run the tests available in this project thanks to Docker, multiple commands are available :

```bash
# Integration tests
make test-integration

# Unit tests
make test-unit

# All test
make test-all
```

### ▶️ Commands

You can execute util commands, you can see executing `make help`

```bash
Targets:
  Swagger UI
    run-openapi-ui      runs swagger ui: http://localhost:4000/
  Quality
    check-quality       runs code quality checks
    lint                go linting. Update and use specific lint tool and options
    lint-fix            go linting. Update and use specific lint tool and options
    vet                 go vet
    fmt                 runs go formatter
    tidy                run go mod tidy
  Test
    test-all            runs tests and create generates coverage report
    coverage            displays test coverage report in html mode
  All
    all                 quality checks and tests
  Help
    help                Show this help.

```

## <a name="env-variables">🌿 Env variables</a>

Environment variables are :

|         Name         |                   Description                   | Required |                 Default value                  | 
|:--------------------:|:-----------------------------------------------:|:--------:|:----------------------------------------------:|
|     `SECRET_KEY`     |              Used to encrypt data               |    ✅     |       `11111111111111111111111111111111`       | 
|  `DEFAULT_PASSWORD`  | Used to encrypt data (combined with secret key) |    ✅     |                 `@myPassword`                  | 
|    `MONGODB_URI`     |       url mongo db with user and password       |    ✅     |    `mongodb://root:example@localhost:27017`    | 
| `MONGODB_NAME`       |                  mongo db name                  |    ✅     |                 `share_secret`                 | 



## <a name="code-analysis-and-consistency">☑️ Code analysis and consistency</a>

### 🔍 Code linting & formatting

In order to keep the code clean, consistent and free of bad practices.

`make lint` 
`make lint-fix`

Check file: [.golangci.yml](.golangci.yml)


## <a name="versions">📈 Releases</a>

You can see releases in: https://github.com/bernardosecades/zaslink/actions/workflows/deploy.yml

## <a name="github-actions">🐙 GitHub Actions</a>

This project uses **GitHub Actions** to automate some boring tasks.

You can find all the workflows in the **[.github/workflows directory](https://github.com/bernardosecades/zaslink/tree/master/.github/workflows).**

### 🎢 Workflows

|                                                        Name                                                        |                                                                             Description & Status                                                                             |    Triggered on    |    
|:------------------------------------------------------------------------------------------------------------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:------------------:|
|               **[⚙️ Build](https://github.com/bernardosecades/zaslink/actions/workflows/build.yml)**               |                     Execute linter, unit test, integration test and if everything is OK build docker image and push in docker hub                                            | `push` on `master` | 
| **[🚀 Deploy To Production Workflow](https://github.com/bernardosecades/zaslink/actions/workflows/deploy.yml)**    | Connect to machine via ssh and with last docker image run new containers. As well API DOCS (Swagger UI) is deployed (it will use file open api V3: docs/openapi/secret.yaml) |      `manual`      | 

## <a name="license">©️ License</a>

This project is licensed under the [MIT License](LICENSE).
