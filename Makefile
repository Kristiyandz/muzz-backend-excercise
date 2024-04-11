.PHONY: start-docker
start-docker:
	@open --background -a Docker

.PHONY: build
build: start-docker
	@docker compose build

.PHONY: up
up: start-docker
	@docker compose up

.PHONY: down
down: start-docker
	@docker compose down

.PHONY: nuke
nuke: start-docker
	@docker compose down -v
	@docker system prune -a
