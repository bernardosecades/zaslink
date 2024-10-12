ssh root@bernardosecades.com "docker pull bernardosecades/zaslink-api && docker rm -f zaslink-api 2>/dev/null || true && docker run --name zaslink-api \
  --restart unless-stopped \
  -p 8080:8080 \
  -e SECRET_KEY=$SECRET_KEY \
  -e DEFAULT_PASSWORD=$DEFAULT_PASSWORD \
  -e MONGODB_URI=$MONGODB_URI \
  -e MONGODB_NAME=$MONGODB_NAME \
  -e NATS_URL=$NATS_URL \
  bernardosecades/zaslink-api:latest" || echo "API Deployment failed" &
