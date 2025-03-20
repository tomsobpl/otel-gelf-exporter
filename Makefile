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

.PHONY: rebuild_otel_collector rebuild_data_producer restart_otel_collector restart_data_producer
rebuild_otel_collector: # Rebuild otel_collector service
	${dc} down otel_collector data_otelgen
	${dc} build --build-arg USER_ID=$(shell id --user) --build-arg GROUP_ID=$(shell id --group) otel_collector
	${dc} up --detach --remove-orphans otel_collector data_otelgen

restart_otel_collector: # Restart otel_collector service
	${dc} restart otel_collector data_otelgen
