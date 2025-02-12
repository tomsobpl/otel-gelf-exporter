dc = docker compose

.PHONY: build up down logs ps rebuild restart
build: # Build docker stack
	${dc} build --build-arg USER_ID=$(shell id --user) --build-arg GROUP_ID=$(shell id --group)

up: # Start docker stack
	${dc} up --detach --remove-orphans

down: # Stop docker stack
	${dc} down

logs: # Show logs
	${dc} logs --follow

ps: # Show docker stack status
	${dc} ps

rebuild: down build up # Rebuild docker stack
restart: down up # Restart docker stack
