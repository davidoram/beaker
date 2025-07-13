
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
	go get -tool github.com/nats-io/natscli/nats@latest 
	go get -tool github.com/rubenv/sql-migrate/...@latest
	go get -tool github.com/santhosh-tekuri/jsonschema/cmd/jv@latest
	go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go get -tool github.com/equinix-labs/otel-cli@latest

.PHONY: build
build: 
	mkdir -p bin
	go build -o bin/beaker $(shell ls *.go | grep -v '_test.go')

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
	bin/beaker --credentials credentials.txt --postgres "postgres://postgres:password@db:5432/beaker_dev?sslmode=disable"

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
