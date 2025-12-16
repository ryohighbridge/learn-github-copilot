#!/bin/bash
set -e

# このスクリプトはPostgreSQLの初期化時に実行される
# /docker-entrypoint-initdb.d/ に配置されたスクリプトは
# データベースが初めて作成されたときにのみ実行される

echo "Running database migrations..."

# マイグレーションを実行
migrate -path /migrations \
  -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
  up || {
    echo "Migration failed. Checking version..."
    # エラーが発生した場合、dirtyバージョンをクリーンアップ
    CURRENT_VERSION=$(migrate -path /migrations \
      -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
      version 2>&1 | grep -oP '\d+' | head -1 || echo "0")
    
    if [ "$CURRENT_VERSION" != "0" ]; then
      echo "Forcing version to $CURRENT_VERSION..."
      migrate -path /migrations \
        -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
        force $CURRENT_VERSION
      
      echo "Retrying migration..."
      migrate -path /migrations \
        -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" \
        up
    fi
  }

echo "Migration completed successfully!"
