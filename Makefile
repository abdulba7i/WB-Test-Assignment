DC = docker compose
APP = docker-compose.yaml
APP_DEV = docker-compose.yaml
ENV = --env-file .env

.PHONY: app-dev app-dev-logs

app-dev:
	${DC} -f ${APP_DEV} ${ENV} up --build -d

app-dev-logs:
	${DC} -f ${APP_DEV} logs -f

app-stop:
	${DC} -f ${APP_DEV} stop

app-down:
	${DC} -f ${APP_DEV} down