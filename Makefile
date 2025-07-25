DB_URL=postgres://postgres:password@localhost?sslmode=disable
DB_ENV?=development

.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f .devcontainer/docker-compose.yml down  --remove-orphans || true

.PHONY: docker-compose-up
docker-compose-up:
	NEW_RELIC_API_KEY=$${NEW_RELIC_API_KEY:-dummy-key} docker-compose -f .devcontainer/docker-compose.yml up -d

.PHONY: restart-docker-compose
restart-docker-compose: docker-compose-down docker-compose-up


.PHONY: install-tools-apt-get
install-tools-apt-get:
	sudo apt-get update -y
	sudo apt-get install -y \
		postgresql-client \
		git \
		jq

.PHONY: install-tools-go
install-tools-go:
	go install tool

.PHONY: setup
setup: install-tools-apt-get install-tools-go 

initial-tool-install:
	go get -tool github.com/nats-io/natscli/nats@v0.2.3 
	go get -tool github.com/rubenv/sql-migrate/...@v1.8.0
	go get -tool github.com/santhosh-tekuri/jsonschema/cmd/jv@v0.7.0
	go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
	go get -tool github.com/equinix-labs/otel-cli@v0.4.5

.PHONY: build
build: migrate-db sqlc
	mkdir -p bin
	go build -o bin/beaker $(shell ls *.go | grep -v '_test.go')


.PHONY: terminate-conns
terminate-conns:
	psql "$(DB_URL)" -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'beaker_$(DB_ENV)' AND pid <> pg_backend_pid();"

.PHONY: drop-db
drop-db: terminate-conns
	psql "$(DB_URL)" -c "DROP DATABASE IF EXISTS beaker_$(DB_ENV);"

.PHONY: create-db
create-db:
	 psql "$(DB_URL)" -c "CREATE DATABASE beaker_$(DB_ENV) WITH OWNER postgres ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8' TEMPLATE template0;"

.PHONY: migrate-db
migrate-db:
	sql-migrate up --config dbconfig.yml --env $(DB_ENV)

.PHONY: recreate-db
recreate-db: drop-db create-db migrate-db
	@echo "Database 'beaker_$(DB_ENV)' recreated successfully."

.PHONY: sqlc
sqlc:
	sqlc generate 

.PHONY: run
run: build  
	OTEL_SERVICE_NAME=beaker \
	OTEL_RESOURCE_ATTRIBUTES=service.version=0.1.0,deployment.environment=codespace \
	OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net \
	OTEL_EXPORTER_OTLP_HEADERS=api-key=${NEW_RELIC_API_KEY} \
	OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT=4095 \
	OTEL_EXPORTER_OTLP_COMPRESSION=gzip \
	OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
	OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE=delta \
	bin/beaker --credentials $(HOME)/credentials.txt --postgres "postgres://postgres:password@localhost:5432/beaker_$(DB_ENV)?sslmode=disable"

.PHONY: test-otel
test-otel:
	OTEL_SERVICE_NAME=beaker \
	OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net \
	OTEL_EXPORTER_OTLP_HEADERS=api-key=${NEW_RELIC_API_KEY} \
	OTEL_EXPORTER_OTLP_COMPRESSION=gzip \
	OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
	otel-cli exec --name "curl google" curl https://google.com

# .IGNORE means that the target **will not** be considered failed if it returns a non-zero exit code.
.IGNORE: test-otel-error
.PHONY: test-otel-error
test-otel-error:
	OTEL_SERVICE_NAME=beaker \
	OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net \
	OTEL_EXPORTER_OTLP_HEADERS=api-key=${NEW_RELIC_API_KEY} \
	OTEL_EXPORTER_OTLP_COMPRESSION=gzip \
	OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
	otel-cli exec --name "test error" --attrs "beaker.foo=bar,beaker.baz=qux" false
