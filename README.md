# Overview

Beaker is a demonstration project that shows how to build a production-ready microservice architecture using [NATS](https://nats.io/). It includes structured request handling, streaming support, service discovery, health checks, and observability features, giving teams a practical foundation for building scalable systems.

# Project Goals

- Showcase the components needed to build and run a modern production quality microservice.
- Build our system in [`go`](https://go.dev/) using the [NATS service](https://docs.nats.io/using-nats/developer/services) protocol, and a [Postgres](https://www.postgresql.org/) database.
- Discuss use of standards, and how they help guide our development.
- Demonstrate tooling used to help accelerate developer productivity when building with these tools
- Provide observability hooks for health checks and metrics.
- Show best practice developer-friendly local development via Github Codespaces and devcontainers.

Each episode covers one of the following topics

- [High Level architecture](./docs/architecture.md)
- [Devlopment environment](./docs/dev-environment.md)
- [Open Telemetry Setup](./docs/otel.md)
- [Database Layer](./docs/db.md)
