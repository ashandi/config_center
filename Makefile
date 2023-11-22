SERVICE_IMAGE_NAME=config_center:local
LOCAL_COMPOSE_FILE=deployments/local/docker-compose.yml
LOCAL_ENV_FILE=deployments/local/.env

build:
	docker build -f ./Dockerfile \
		-t $(SERVICE_IMAGE_NAME) .

local-run: build
	docker-compose --env-file $(LOCAL_ENV_FILE) -f $(LOCAL_COMPOSE_FILE) up -d --force-recreate
	docker-compose --env-file $(LOCAL_ENV_FILE) -f $(LOCAL_COMPOSE_FILE) logs -f --tail=50
