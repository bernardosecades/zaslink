ssh root@bernardosecades.com "docker pull bernardosecades/api-share-secret && docker rm -f sharesecret-api 2>/dev/null || true && docker run --name sharesecret-api \
  --restart unless-stopped \
  -p 8080:8080 \
  -e SECRET_KEY=$SECRET_KEY \
  -e DEFAULT_PASSWORD=$DEFAULT_PASSWORD \
  -e MONGODB_URI=$MONGODB_URI \
  -e MONGODB_NAME=$MONGODB_NAME \
  bernardosecades/api-share-secret:latest" || echo "API Deployment failed" &
