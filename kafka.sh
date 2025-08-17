
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
  "name": "account-source-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "plugin.name": "pgoutput",
    "database.hostname": "account_db",
    "database.port": "5432",
    "database.user": "wignn",
    "database.password": "123456",
    "database.dbname": "account",
    "database.server.name": "account_db",
    "slot.name": "account_slot",
    "publication.name": "account_pub",
    "table.include.list": "public.accounts",
    "key.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "key.converter.schemas.enable": false,
    "value.converter.schemas.enable": false,
    "transforms": "unwrap",
    "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
    "transforms.unwrap.drop.tombstones": "false",
    "snapshot.mode": "initial"
  }
}'

echo -e "\n"

echo "=== Setup Auth Replica Sink Connector ==="
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
  "name": "account-to-auth-sink",
  "config": {
    "connector.class": "io.debezium.connector.jdbc.JdbcSinkConnector",
    "tasks.max": "1",
    "topics": "account_db.public.accounts",
    "connection.url": "jdbc:postgresql://auth_db:5432/auth?user=wignn&password=123456",
    "auto.create": "false",
    "insert.mode": "upsert",
    "pk.mode": "record_key",
    "pk.fields": "id",
    "delete.enabled": "true"
  }
}'