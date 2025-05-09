.PHONY: setup
setup:
	sudo apt-get update -y
	sudo apt-get install -y \
		postgresql-client \
		git \
		jq
	go install tool

initial-tool-install:
	go get -tool github.com/nats-io/natscli/nats@latest 
	go get -tool github.com/rubenv/sql-migrate/...@latest
	go get -tool github.com/santhosh-tekuri/jsonschema/cmd/jv@latest
	go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest