# Account Service

Account Service is a core component of a larger microservices architecture. It handles all user account-related operations. The service follows a GraphQL → gRPC → Server → Service → Repository → Database flow to ensure separation of concerns, scalability, and maintainability.

---

## 📁 Project Structure

```txt
.
├── client/             # GraphQL layer (queries/mutations)
├── server/             # gRPC server definition
├── service/            # Business logic
├── repository/         # Data access layer
├── proto/              # .proto files for gRPC
├── database/           # DB migrations, seeders, schemas
├── Makefile            # Commands for proto generation and other tools
└── README.md
