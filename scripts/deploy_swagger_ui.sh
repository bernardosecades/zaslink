ssh root@bernardosecades.com "cd sharesecret && git pull origin master && make run-openapi-ui" || echo "Swagger UI Deployment failed" &
