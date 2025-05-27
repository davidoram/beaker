
.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f .devcontainer/docker-compose.yml down || true

.PHONY: docker-compose-up
docker-compose-up:
	docker-compose -f .devcontainer/docker-compose.yml up -d

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
	go get -tool github.com/nats-io/natscli/nats@latest 
	go get -tool github.com/rubenv/sql-migrate/...@latest
	go get -tool github.com/santhosh-tekuri/jsonschema/cmd/jv@latest
	go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: run
run: 
	OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
	OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
	OTEL_SERVICE_NAME=beaker \
	OTEL_RESOURCE_ATTRIBUTES=service.version=0.1.0,deployment.environment=development \
	OTEL_LOGS_EXPORTER=otlp \
	OTEL_METRICS_EXPORTER=otlp \
	OTEL_TRACES_EXPORTER=otlp \
	go run $(shell ls *.go | grep -v '_test.go') --credentials credentials.txt --postgres postgres://postgres:password@localhost:5432/beaker_dev?sslmode=disable

.PHONY: run-debug
run-debug:
	OTEL_EXPORTER_TYPE=stdout \
	OTEL_SERVICE_NAME=beaker \
	OTEL_RESOURCE_ATTRIBUTES=service.version=0.1.0,deployment.environment=development \
	OTEL_LOGS_EXPORTER=stdout \
	OTEL_METRICS_EXPORTER=stdout \
	OTEL_TRACES_EXPORTER=stdout \
	go run $(shell ls *.go | grep -v '_test.go') --credentials credentials.txt --postgres postgres://postgres:password@localhost:5432/beaker_dev?sslmode=disable 