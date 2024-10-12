ssh root@bernardosecades.com "docker pull bernardosecades/zaslink-notifier && docker rm -f zaslink-notifier 2>/dev/null || true && docker run --name zaslink-notifier \
  --restart unless-stopped \
  -e NOTIFIER_TELEGRAM_BOT_TOKEN=$NOTIFIER_TELEGRAM_BOT_TOKEN \
  -e NOTIFIER_TELEGRAM_USER_ID=$NOTIFIER_TELEGRAM_USER_ID \
  -e NATS_URL=$NATS_URL \
  bernardosecades/zaslink-notifier:latest" || echo "NOTIFIER Deployment failed" &
  