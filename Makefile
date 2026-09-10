.PHONY: test build up down

test:
	@bash ./test.sh

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down -v