.PHONY: start-docker-compose
start-docker-compose:
	docker-compose -f .devcontainer/docker-compose.yml up -d


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
	go run $(shell ls *.go | grep -v '_test.go') --credentials credentials.txt --postgres postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable 