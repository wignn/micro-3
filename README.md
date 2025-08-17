
# Micro-3

<div align="center">

<!-- Technologies Used -->
<img src="https://www.vectorlogo.zone/logos/golang/golang-icon.svg" height="40" />
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/grpc/grpc-original.svg" height="40"/>
<img src="https://www.vectorlogo.zone/logos/graphql/graphql-icon.svg" height="40" />
<img src="https://www.vectorlogo.zone/logos/elastic/elastic-icon.svg" height="40" />
<img src="https://www.vectorlogo.zone/logos/postgresql/postgresql-icon.svg" height="40" />
<img src="https://www.vectorlogo.zone/logos/apache_kafka/apache_kafka-icon.svg" height="40" />

</div>

<div align="center">

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/wignn/micro-3/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](#license)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue)](https://golang.org/)

</div>

---

## Overview

Micro-3 is a polyglot microservices system written in Go. Services communicate internally via gRPC and expose a unified public API through a GraphQL gateway. Data is stored in PostgreSQL and Elasticsearch. Kafka + Debezium enables CDC to stream data changes (for example from Account to Auth).

### Core modules

- Account: user registration and profile
- Auth: login and JWT token issuance
- Catalog: product catalog (Elasticsearch)
- Order: order placement and pricing
- Review: product reviews
- GraphQL Gateway: single entry point for clients
- Kafka + Debezium: CDC and async integration

---

## Services and ports

- GraphQL gateway: <http://localhost:8000> (Playground at /playground, API at /graphql)
- Kafka UI: <http://localhost:4000>
- Kafka Connect REST: <http://localhost:8083>
- Kafka broker: localhost:9092 (inside Docker network: kafka:9092)
- Zookeeper: localhost:2181

Note: gRPC services (account, auth, catalog, order, review) listen on port 8080 inside the Docker network and aren’t exposed directly. Access them through the GraphQL gateway.

---

## Project structure

```text
micro-3/
├── account/    ── gRPC service + Postgres
├── auth/       ── gRPC service + Postgres
├── catalog/    ── gRPC service + Elasticsearch
├── order/      ── gRPC service + Postgres
├── review/     ── gRPC service + Postgres
├── graphql/    ── GraphQL gateway (gqlgen)
├── kafka/      ── Debezium/Kafka configs
├── compose.yml ── Docker Compose stack
└── vendor/     ── vendored dependencies
```

---

## Quick start

Prerequisites: Docker Desktop (Compose v2), ~4GB RAM available for containers.

1. Clone

```powershell
git clone https://github.com/wignn/micro-3.git
cd micro-3
```

1. Start the stack

```powershell
docker compose up --build -d
```

1. Check status

```powershell
docker compose ps
```

1. Open the GraphQL Playground

- <http://localhost:8000/playground> (queries go to /graphql)

---

## GraphQL usage

- Endpoint: <http://localhost:8000/graphql>
- Playground: <http://localhost:8000/playground>

Example queries

Query accounts

```graphql
query {
  accounts(pagination: { skip: 0, take: 10 }) {
    id
    name
    email
  }
}
```

Search products

```graphql
query {
  products(pagination: { skip: 0, take: 12 }, query: "laptop") {
    id
    name
    description
    price
    image
  }
}
```

Mutations

Register account

```graphql
mutation {
  createAccount(account: { name: "Jane", email: "jane@example.com", password: "secret" }) {
    id
    name
    email
  }
}
```

Login and get tokens

```graphql
mutation {
  login(account: { email: "jane@example.com", password: "secret" }) {
    id
    email
    backendToken {
      accessToken
      refreshToken
      expiresIn
    }
  }
}
```

Create an order

```graphql
mutation {
  createOrder(order: { accountId: "<ACCOUNT_ID>", products: [{ id: "<PRODUCT_ID>", quantity: 2 }] }) {
    id
    totalPrice
    status
    products { id name price quantity }
  }
}
```

Notes

- By default, the GraphQL gateway doesn’t require an Authorization header. The Auth service issues tokens that your client app can store and use if you add protected endpoints later.

---

## Kafka + Debezium (CDC)

This repo includes Kafka, Zookeeper, Kafka Connect (Debezium), and a Kafka UI in the Compose stack. After the stack is up, create the source and sink connectors.

1) Ensure Postgres for Account allows logical replication

Add to postgresql.conf (or via a mounted config/ALTER SYSTEM).

```ini
wal_level=logical
max_replication_slots=10
max_wal_senders=10
```

1. Create connectors using the included JSON files

- Source (Account -> Kafka): `kafka/consumer/consumer-account-auth.json` (Debezium Postgres source)
- Sink (Kafka -> Auth DB): `kafka/connectors/sink-account-to-auth.json` (JDBC sink)

Example (using curl). Run after the stack is healthy.

```powershell
# Source connector
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" --data @kafka/consumer/consumer-account-auth.json

# Sink connector
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" --data @kafka/connectors/sink-account-to-auth.json
```

1. Inspect topics and connectors

- Kafka UI: <http://localhost:4000>
- Kafka Connect REST: <http://localhost:8083/connectors>

---

## Development

Code generation

- Protobuf/gRPC (per service), e.g. in `account/`:

```powershell
cd account
make gen
```

- GraphQL schema to Go models/resolvers:

```powershell
cd graphql
make gen
```

Running services locally

- Services default to gRPC port 8080 in containers, and 50051 by default when running locally. Set `PORT` to override. Example for Account:

```powershell
set PORT=8080; set DATABASE_URL=postgres://wignn:123456@localhost:5432/account?sslmode=disable
go run ./account/cmd/account
```

Tip: For local-only runs you’ll need databases and dependencies reachable from your host, or use `docker compose` for the full stack and iterate on code with container restarts.

---

## Configuration

Key environment variables (set by Compose already):

- Account/Auth/Order/Review: `DATABASE_URL`, `PORT`
- Catalog: `DATABASE_URL` (Elasticsearch URL), `PORT`
- Auth: `ACCESS_SECRET_KEY`, `REFRESH_SECRET_KEY`
- GraphQL gateway: `*_SERVICE_URL` for each backend gRPC service

See `compose.yml` for the complete list and defaults.

---

## Troubleshooting

- GraphQL not reachable: ensure `graphql` container is healthy and port 8000 is free.
- Empty results: seed data may be minimal; create resources via mutations.
- Debezium errors: verify Postgres `wal_level=logical` and connector configs in Kafka UI.
- Logs: `docker compose logs -f <service>` (e.g., `graphql`, `account`, `connect`).

---

## License

MIT