DB_CONTAINER ?= gestion_futbol_db
BASE_URL     ?= http://localhost:8080
TIMEOUT      ?= 90

export DB_CONTAINER BASE_URL TIMEOUT

.PHONY: help test sqlc up wait prepare clean logs

help:
	@echo "make test     Ejecuta todo: sqlc, build, contenedores, pruebas y limpieza"
	@echo "make prepare  Solo tareas previas (sqlc, build, levantar y esperar)"
	@echo "make sqlc     Regenera el codigo con sqlc"
	@echo "make clean    Borra contenedores y volumenes"
	@echo "make logs     Muestra los logs de los contenedores"
	@echo "Variables: BASE_URL=$(BASE_URL) DB_CONTAINER=$(DB_CONTAINER) TIMEOUT=$(TIMEOUT)"

sqlc:
	@echo ">> sqlc generate"
	@if command -v sqlc >/dev/null 2>&1; then sqlc generate; \
	else docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate; fi

clean:
	@echo ">> Borrando contenedores y volumenes"
	@docker compose down -v --remove-orphans

up:
	@echo ">> Compilando (docker) y levantando contenedores"
	@docker compose up -d --build

wait:
	@./scripts/wait-ready.sh

prepare: sqlc clean up wait

# 2) Pruebas + 3) Limpieza (siempre se ejecuta, pasen o fallen las pruebas) --
test:
	chmod +x ./scripts/wait-ready.sh
	chmod +x ./scripts/test-api.sh
	@status=0; \
	$(MAKE) --no-print-directory prepare || status=$$?; \
	if [ $$status -eq 0 ]; then ./scripts/test-api.sh || status=$$?; \
	else echo "!! Fallaron las tareas previas"; $(MAKE) --no-print-directory logs; fi; \
	$(MAKE) --no-print-directory clean; \
	exit $$status

logs:
	@docker compose logs --tail=60