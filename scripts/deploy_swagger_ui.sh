ssh root@bernardosecades.com "cd zaslink && git pull origin master && make run-openapi-ui" || echo "Swagger UI Deployment failed" &
