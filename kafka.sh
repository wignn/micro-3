#!/bin/bash


set -e

CONNECT_URL="http://localhost:8083"
KAFKA_CONTAINER=$(docker ps | grep kafka | grep -v zookeeper | grep -v ui | awk '{print $1}')
ACCOUNT_DB_CONTAINER=$(docker ps | grep account_db | awk '{print $1}')
AUTH_DB_CONTAINER=$(docker ps | grep auth_db | awk '{print $1}')

echo "=== KAFKA CONNECT INITIALIZATION ==="
echo "Waiting for services to be ready..."

# Function to wait for service
wait_for_service() {
    local url=$1
    local service_name=$2
    echo "Waiting for $service_name to be ready..."
    until curl -f -s $url > /dev/null 2>&1; do
        echo "  $service_name not ready, waiting 5 seconds..."
        sleep 5
    done
    echo "  $service_name is ready!"
}

# Wait for Kafka Connect
wait_for_service "$CONNECT_URL/connectors" "Kafka Connect"

echo ""
echo "=== STEP 1: CLEANUP EXISTING CONNECTORS ==="

# Get existing connectors
EXISTING_CONNECTORS=$(curl -s $CONNECT_URL/connectors | jq -r '.[]' 2>/dev/null || echo "")

if [ ! -z "$EXISTING_CONNECTORS" ]; then
    echo "Removing existing connectors..."
    for connector in $EXISTING_CONNECTORS; do
        if [[ $connector == *"account"* ]] || [[ $connector == *"auth"* ]]; then
            echo "  Deleting: $connector"
            curl -X DELETE $CONNECT_URL/connectors/$connector
        fi
    done
    sleep 3
else
    echo "No existing connectors found."
fi

echo ""
echo "=== STEP 2: CREATE SOURCE CONNECTOR ==="

# Create PostgreSQL Source Connector dengan snapshot
echo "Creating PostgreSQL Source Connector for accounts table..."
curl -X POST $CONNECT_URL/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "account-postgres-source",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "plugin.name": "pgoutput",
      "database.hostname": "account_db",
      "database.port": "5432",
      "database.user": "wignn",
      "database.password": "123456",
      "database.dbname": "account",
      "database.server.name": "account_server",
      "table.include.list": "public.accounts",
      "topic.prefix": "account",
      "snapshot.mode": "always",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false"
    }
  }'

if [ $? -eq 0 ]; then
    echo "✅ Source connector created successfully"
else
    echo "❌ Failed to create source connector"
    exit 1
fi

echo ""
echo "=== STEP 3: WAIT FOR SOURCE CONNECTOR TO BE READY ==="

# Wait for source connector to be running
sleep 10
for i in {1..12}; do
    STATUS=$(curl -s $CONNECT_URL/connectors/account-postgres-source/status | jq -r '.connector.state' 2>/dev/null || echo "UNKNOWN")
    if [ "$STATUS" = "RUNNING" ]; then
        echo "✅ Source connector is running"
        break
    else
        echo "  Source connector status: $STATUS (attempt $i/12)"
        sleep 5
    fi
done

echo ""
echo "=== STEP 4: CREATE SINK CONNECTOR ==="

# Wait a bit more for topics to be created
sleep 5

echo "Creating PostgreSQL Sink Connector for auth_db..."
curl -f -X POST $CONNECT_URL/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "auth-postgres-sink",
    "config": {
      "connector.class": "io.debezium.connector.jdbc.JdbcSinkConnector",
      "connection.url": "jdbc:postgresql://auth_db:5432/auth",
      "connection.username": "wignn",
      "connection.password": "123456",
      "topics": "account.public.accounts",
      "insert.mode": "upsert",
      "delete.enabled": "true",
      "pk.mode": "record_key",
      "pk.fields": "id",
      "table.name.format": "accounts",
      "auto.create": "true",
      "auto.evolve": "true",
      "transforms": "unwrap",
      "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
      "transforms.unwrap.drop.tombstones": "true",
      "transforms.unwrap.delete.handling.mode": "rewrite",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false"
    }
  }'

if [ $? -eq 0 ]; then
    echo "✅ Sink connector created successfully"
else
    echo "❌ Failed to create sink connector"
    exit 1
fi

echo ""
echo "=== STEP 5: VERIFY SETUP ==="

sleep 5

# Check connector status
echo "Checking connector status..."
echo "Source Connector:"
curl -s $CONNECT_URL/connectors/account-postgres-source/status | jq '.'
echo ""
echo "Sink Connector:"
curl -s $CONNECT_URL/connectors/auth-postgres-sink/status | jq '.'

echo ""
echo "=== STEP 6: CHECK TOPICS ==="

echo "Available topics:"
if [ ! -z "$KAFKA_CONTAINER" ]; then
    docker exec $KAFKA_CONTAINER kafka-topics --bootstrap-server localhost:9092 --list | grep account || echo "No account topics found yet"
else
    echo "Kafka container not found"
fi

echo ""
echo "=== STEP 7: TEST DATA SYNC ==="

# Check current data in source
echo "Current data in account_db:"
if [ ! -z "$ACCOUNT_DB_CONTAINER" ]; then
    docker exec $ACCOUNT_DB_CONTAINER psql -U wignn -d account -c "SELECT COUNT(*) as total_accounts FROM accounts;"
else
    echo "Account DB container not found"
fi

# Wait for sync
sleep 10

# Check data in destination
echo "Data in auth_db after sync:"
if [ ! -z "$AUTH_DB_CONTAINER" ]; then
    docker exec $AUTH_DB_CONTAINER psql -U wignn -d auth -c "SELECT COUNT(*) as total_accounts FROM accounts;" 2>/dev/null || echo "Table not created yet or no data synced"
else
    echo "Auth DB container not found"
fi


# Wait for CDC to process
sleep 5

echo ""
echo "=== STEP 8: FINAL VERIFICATION ==="

echo "Final data count in auth_db:"
if [ ! -z "$AUTH_DB_CONTAINER" ]; then
    docker exec $AUTH_DB_CONTAINER psql -U wignn -d auth -c "SELECT COUNT(*) as total_accounts FROM accounts;" 2>/dev/null || echo "Sync may still be in progress..."
    echo ""
    echo "Latest records in auth_db:"
    docker exec $AUTH_DB_CONTAINER psql -U wignn -d auth -c "SELECT * FROM accounts ORDER BY id DESC LIMIT 3;" 2>/dev/null || echo "No data found yet"
else
    echo "Auth DB container not found"
fi

echo ""
echo "=== INITIALIZATION COMPLETE ==="
echo ""
echo "🎉 Kafka Connect CDC Pipeline Setup Complete!"
echo ""
echo "📊 Monitor your setup:"
echo "   - Kafka UI: http://localhost:4000"
echo "   - Kafka Connect REST API: http://localhost:8083"
echo ""
echo "🔧 Useful commands:"
echo "   - Check connector status: curl http://localhost:8083/connectors/account-postgres-source/status"
echo "   - List topics: docker exec \$KAFKA_CONTAINER kafka-topics --bootstrap-server localhost:9092 --list"
echo "   - View messages: docker exec \$KAFKA_CONTAINER kafka-console-consumer --bootstrap-server localhost:9092 --topic account.accounts --from-beginning"
echo ""
echo "✨ Try inserting more data into account_db to see real-time sync!"