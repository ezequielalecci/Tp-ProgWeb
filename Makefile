DB_CONTAINER ?= gestion_futbol_db
BASE_URL     ?= http://localhost:8080
TIMEOUT      ?= 90

# Version de templ: se lee de go.mod para que el generador y la libreria coincidan.
TEMPL_VERSION := $(shell awk '$$1 == "github.com/a-h/templ" {print $$2}' go.mod)
GO_IMAGE      ?= golang:1.26-alpine

export DB_CONTAINER BASE_URL TIMEOUT

.PHONY: help run test generate sqlc templ deps up wait prepare clean logs

help:
	@echo "make run       Genera el codigo, levanta la base y el servidor y espera a que respondan"
	@echo "make test      Ejecuta todo: generacion, build, contenedores, pruebas y limpieza"
	@echo "make generate  Regenera el codigo de sqlc y de templ"
	@echo "make sqlc      Regenera solo el codigo de sqlc (db/sqlc)"
	@echo "make templ     Regenera solo el codigo de templ (views/*_templ.go)"
	@echo "make deps      Actualiza go.mod y go.sum (go mod tidy)"
	@echo "make prepare   Tareas previas de las pruebas (generar, build, levantar y esperar)"
	@echo "make clean     Borra contenedores y volumenes"
	@echo "make logs      Muestra los logs de los contenedores"
	@echo "Variables: BASE_URL=$(BASE_URL) DB_CONTAINER=$(DB_CONTAINER) TIMEOUT=$(TIMEOUT)"

# Generacion de codigo ------------------------------------------------------
sqlc:
	@echo ">> sqlc generate"
	@if command -v sqlc >/dev/null 2>&1; then sqlc generate; \
	else docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate; fi

templ:
	@echo ">> templ generate ($(TEMPL_VERSION))"
	@if command -v templ >/dev/null 2>&1; then templ generate; \
	elif command -v go >/dev/null 2>&1; then go run github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION) generate; \
	else docker run --rm -u "$$(id -u):$$(id -g)" -v "$$(pwd):/app" -w /app ghcr.io/a-h/templ:$(TEMPL_VERSION) generate; fi

generate: sqlc templ

# go mod tidy necesita el codigo generado para ver todos los imports.
deps: generate
	@echo ">> go mod tidy"
	@if command -v go >/dev/null 2>&1; then go mod tidy; \
	else docker run --rm -v "$$(pwd):/src" -w /src $(GO_IMAGE) go mod tidy; fi

# Ejecucion -----------------------------------------------------------------
clean:
	@echo ">> Borrando contenedores y volumenes"
	@docker compose down -v --remove-orphans

up:
	@echo ">> Compilando (docker) y levantando contenedores"
	@docker compose up -d --build

wait:
	@./scripts/wait-ready.sh

# Genera, levanta Postgres y el servidor, y avisa cuando ya se puede usar.
run: generate up wait
	@echo ">> Listo: $(BASE_URL)"

prepare: generate clean up wait

# Pruebas + limpieza (siempre se ejecuta, pasen o fallen las pruebas) --------
test:
	chmod +x ./scripts/wait-ready.sh
	chmod +x ./scripts/test-web.sh
	@status=0; \
	$(MAKE) --no-print-directory prepare || status=$$?; \
	if [ $$status -eq 0 ]; then ./scripts/test-web.sh || status=$$?; \
	else echo "!! Fallaron las tareas previas"; $(MAKE) --no-print-directory logs; fi; \
	$(MAKE) --no-print-directory clean; \
	exit $$status

logs:
	@docker compose logs --tail=60
