include .env
export
POSTGRES_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bot-postgres:5432/${POSTGRES_DB}?sslmode=disable
CLICKHOUSE_URL=clickhouse://${CLICKHOUSE_HOST}?username=${CLICKHOUSE_USER}&password=${CLICKHOUSE_PASSWORD}&database=${CLICKHOUSE_DB}&x-multi-statement=true&x-migrations-table-engine=MergeTree&x-migrations-table=clickhouse_migrations

env-up:
	@$(MAKE) redis-up
	@$(MAKE) clickhouse-up
	@$(MAKE) postgres-up

env-down:
	@$(MAKE) redis-down
	@$(MAKE) clickhouse-down
	@$(MAKE) postgres-down

postgres-up:
	@docker compose up -d postgres

postgres-down:
	@docker compose down ostgres

redis-up:
	@docker compose up -d redis

redis-down:
	@docker compose down redis

clickhouse-up:
	mkdir -p out/clickhousedata
	sudo chown -R 101:101 out/clickhousedata
	@docker compose up -d clickhouse

clickhouse-down:
	@docker compose down clickhouse

migrate-postgres-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-postgres-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root postgres-migrate \
		create \
		-ext sql \
		-dir /migrations/postgres \
		-seq "$(seq)"; \
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations/postgres

migrate-postgres-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-postgres-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root postgres-migrate \
	-path /migrations/postgres \
	-database "$(POSTGRES_URL)" \
	"$(action)"

migrate-postgres-up:
	@$(MAKE) migrate-postgres-action action=up 

migrate-postgres-down:
	@$(MAKE) migrate-postgres-action action=down

migrate-clickhouse-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-clickhouse-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations/clickhouse \
		-seq "$(seq)"; \
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations/clickhouse

migrate-clickhouse-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-clickhouse-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm --user root bot-postgres-migrate \
	-path /migrations/clickhouse \
	-database "$(CLICKHOUSE_URL)" \
	"$(action)"

migrate-clickhouse-up:
	@$(MAKE) migrate-clickhouse-action action=up 

migrate-clickhouse-down:
	@$(MAKE) migrate-clickhouse-action action=down