# note: call scripts from /scripts
build:
	docker-compose -f ./deployments/docker-compose.yml build --no-cache resumepro
run:
	# https://stackoverflow.com/a/2670143/6670698
	-docker rmi deploy_whitelister:latest
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./deployments/docker-compose.yml up --build --remove-orphans resumepro
test:
	docker-compose -f ./deployments/docker-compose.yml up tests

lint:
	docker-compose -f ./deployments/docker-compose.yml up linter

down:
	docker-compose -f ./deployments/docker-compose.yml down