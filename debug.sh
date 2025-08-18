#!/bin/bash

echo "=== DEBUGGING KAFKA CONNECT SETUP ==="

# 1. Cek Connectors yang ada
echo -e "\n1. Connectors yang aktif:"
curl -s http://localhost:8083/connectors | jq '.'

# 2. Cek status source connector
echo -e "\n2. Status Source Connector:"
curl -s http://localhost:8083/connectors/account-postgres-source/status | jq '.'

# 3. Cek topics yang dibuat
echo -e "\n3. Topics di Kafka:"
docker exec $(docker ps | grep kafka | grep -v zookeeper | grep -v ui | awk '{print $1}') kafka-topics --bootstrap-server localhost:9092 --list | grep account

# 5. Cek database account_db
echo -e "\n4. Tables di account_db:"
docker exec $(docker ps | grep account_db | awk '{print $1}') psql -U wignn -d account -c '\dt'

echo -e "\n5. Sample data dari account_db:"
docker exec $(docker ps | grep account_db | awk '{print $1}') psql -U wignn -d account -c 'SELECT table_name FROM information_schema.tables WHERE table_schema = '\''public'\'' LIMIT 5;'

# 7. Cek database auth_db
echo -e "\n6. Tables di auth_db:"
docker exec $(docker ps | grep auth_db | awk '{print $1}') psql -U wignn -d auth -c '\dt'

echo -e "\n=== DEBUG SELESAI ==="