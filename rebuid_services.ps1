docker compose -f ./deploy/docker-compose.yml stop order-service loyalty-service sso-service
docker compose -f ./deploy/docker-compose.yml rm -f order-service loyalty-service sso-service
docker compose -f ./deploy/docker-compose.yml up --build -d order-service loyalty-service sso-service