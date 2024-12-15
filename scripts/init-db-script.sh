#!/bin/bash

db_user=${DB_USER:-}
db_password=${DB_PASSWORD:-}
db_name=${DB_NAME:-}

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER ${db_user} WITH ENCRYPTED PASSWORD '${db_password}';
	CREATE DATABASE ${db_name};
	GRANT ALL PRIVILEGES ON DATABASE ${db_name} TO ${db_user};
  
  \c ${db_name}
  
  CREATE TABLE urls (
      id SERIAL PRIMARY KEY,
      shortened VARCHAR(20) NOT NULL UNIQUE,
      original TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT NOW()
  );
  CREATE INDEX idx_shortened ON urls(shortened);
  GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${db_user};
  GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ${db_user};
EOSQL

if [ $? -eq 0 ]; then
  echo "Database initialized successfully!"
else
  echo "Failed to initialize the database."
fi
