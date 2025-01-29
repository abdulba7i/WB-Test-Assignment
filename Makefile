DC = docker compose
APP = deployments/docker-compose.yaml
APP_DEV = deployments/docker-compose.yaml
POSTGRES = deployments/postgres.yaml
REDIS =deployments/redis.yaml
KAFKA = deployments/kafka.yaml
CONSUMER = deployments/event-consumer.yaml
APP_SERVICE = app-api
ENV = --env-file .env

.PHONY: app-dev
app-dev:
	${DC} -f ${APP_DEV} ${ENV} up --build -d

.PHONY: app-dev-logs
app-dev-logs:
	${DC} -f ${APP_DEV} logs -f

