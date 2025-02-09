dc = docker compose
otelcol_builder_image = otel/opentelemetry-collector-builder:0.119.0

.PHONY: docker_build docker_up docker_down docker_logs docker_ps docker_rebuild docker_restart
docker_build: # Build docker stack
	${dc} build --build-arg DOCKER_USER_ID=$(shell id --user) --build-arg DOCKER_GROUP_ID=$(shell id --group)

docker_up: # Start docker stack
	${dc} up --detach --remove-orphans

docker_down: # Stop docker stack
	${dc} down

docker_logs: # Show logs
	${dc} logs --follow

docker_ps: # Show docker stack status
	${dc} ps

docker_rebuild: docker_down docker_build docker_up # Rebuild docker stack
docker_restart: docker_down docker_up # Restart docker stack

.PHONY: otelcol_build otelcol_start
otelcol_build: # Build otelcol image
	docker run --volume "$(shell pwd)/otelcol-builder-config.yaml:/builder-config.yaml" \
	--volume "$(shell pwd)/build/otelcol-dev:/tmp/otelcol-dev" \
	--volume "$(shell pwd):/src" ${otelcol_builder_image} --config /builder-config.yaml

otelcol_start: # Start otelcol
	$(shell pwd)/build/otelcol-dev/otelcol-custom --config $(shell pwd)/otelcol-config.yaml