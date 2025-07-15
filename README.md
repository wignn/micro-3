
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

**Micro-3** is a microservices architecture built using **Go**, leveraging **gRPC** for internal service communication and **GraphQL** as an external API gateway. The system is backed by **PostgreSQL** and **Elasticsearch** for data storage, and integrates **Kafka + Debezium** for real-time *Change Data Capture (CDC)* across services such as `account` and `auth`.

### Core Modules

- **Account Service**: Handles user registration and profile management.
- **Authentication Service**: Handles login, authentication, and JWT-based authorization.
- **Catalog Service**: Provides available product or item data.
- **Order Service**: Manages order processing and delivery status.
- **Review Service**: Collects and manages user reviews.
- **GraphQL Gateway**: Public API interface to access all services.
- **Kafka Integration**: Facilitates service-to-service communication and data syncing.
- **Real-time CDC**: Enables database synchronization using Kafka + Debezium.

---

## Technology Stack

| Component         | Technology              |
|------------------|--------------------------|
| Language          | Go                      |
| Communication     | gRPC, GraphQL           |
| Database          | PostgreSQL, Elasticsearch |
| CDC               | Kafka + Debezium        |
| Containerization  | Docker + Docker Compose |
| Infrastructure    | Makefile, Dockerfile    |

---

## Project Structure

```
micro-3/
├── account/        # User account management
├── auth/           # Authentication and authorization
├── catalog/        # Product catalog service
├── order/          # Order processing
├── review/         # Review system
├── graphql/        # GraphQL API gateway
├── kafka/          # Kafka & Debezium configuration
├── compose.yml     # Docker Compose definition
├── go.mod          # Go module definition
├── go.sum          # Go dependencies checksum
└── README.md       # Project documentation
```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/wignn/micro-3.git
cd micro-3
```

### 2. Build and Start All Services

```bash
docker compose up --build -d
```

### 3. Verify Running Containers

```bash
docker compose ps
```
