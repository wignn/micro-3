echo "=== Setup Account Source Connector (CDC) ==="
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "account-cdc-source",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "database.hostname": "account_db",
      "database.port": "5432",
      "database.user": "wignn",
      "database.password": "123456",
      "database.dbname": "account",
      "topic.prefix": "account_events",
      "table.include.list": "public.users,public.profiles",
      "plugin.name": "pgoutput",
      "slot.name": "account_replica_slot",
      "publication.name": "account_replica_pub",
      "database.history.kafka.bootstrap.servers": "kafka:9092",
      "database.history.kafka.topic": "account.history",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false",
      "snapshot.mode": "initial"
    }
  }'

echo -e "\n"

# 3. Setup Sink Connector - Replika ke Auth DB  
echo "=== Setup Auth Replica Sink Connector ==="
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "auth-replica-sink",
    "config": {
      "connector.class": "io.debezium.connector.jdbc.JdbcSinkConnector",
      "connection.url": "jdbc:postgresql://auth_db:5432/auth",
      "connection.user": "wignn",
      "connection.password": "123456",
      "topics": "account_events.public.users,account_events.public.profiles",
      "table.name.format": "replica_${topic.replace(\"account_events.public.\", \"\")}",
      "auto.create": "true",
      "auto.evolve": "true",
      "insert.mode": "upsert",
      "pk.mode": "record_key",
      "pk.fields": "id",
      "delete.enabled": "true",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false",
      "transforms": "unwrap",
      "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
      "transforms.unwrap.drop.tombstones": "false"
    }
  }'

echo -e "\n"

echo "=== Monitoring Commands ==="
echo "# Cek status connectors:"
echo "curl -s http://localhost:8083/connectors/account-cdc-source/status | jq ."
echo "curl -s http://localhost:8083/connectors/auth-replica-sink/status | jq ."
echo ""
echo "# Cek topics:"
echo "curl -s http://localhost:8083/connectors | jq ."
echo ""
echo "# Test data sync:"
echo "# 1. Insert data ke account database"
echo "# 2. Cek apakah muncul di auth database sebagai replica"

echo -e "\n=== ALTERNATIF: Simple Replica Setup ==="
echo "curl -X POST http://localhost:8083/connectors \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{"
echo "    \"name\": \"simple-auth-replica\","
echo "    \"config\": {"
echo "      \"connector.class\": \"io.debezium.connector.jdbc.JdbcSinkConnector\","
echo "      \"connection.url\": \"jdbc:postgresql://auth_db:5432/auth\","
echo "      \"connection.user\": \"wignn\","
echo "      \"connection.password\": \"123456\","
echo "      \"topics\": \"account.public.users\","
echo "      \"table.name.format\": \"users_replica\","
echo "      \"auto.create\": \"true\","
echo "      \"auto.evolve\": \"true\","
echo "      \"insert.mode\": \"upsert\","
echo "      \"pk.mode\": \"record_value\","
echo "      \"pk.fields\": \"id\","
echo "      \"key.converter\": \"org.apache.kafka.connect.json.JsonConverter\","
echo "      \"value.converter\": \"org.apache.kafka.connect.json.JsonConverter\","
echo "      \"key.converter.schemas.enable\": \"false\","
echo "      \"value.converter.schemas.enable\": \"false\","
echo "      \"transforms\": \"unwrap\","
echo "      \"transforms.unwrap.type\": \"io.debezium.transforms.ExtractNewRecordState\""
echo "    }"
echo "  }'"