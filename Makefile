# Makefileでマイグレーション管理とテストを簡単にする

.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version test test-backend test-frontend test-coverage

# マイグレーションを実行（最新まで）
migrate-up:
	docker compose exec db migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" up

# マイグレーションを1つロールバック
migrate-down:
	docker compose exec db migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" down 1

# 新しいマイグレーションファイルを作成
migrate-create:
	@read -p "Enter migration name: " name; \
	docker run --rm -v $(PWD)/db/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $$name

# マイグレーションバージョンを強制設定
migrate-force:
	@read -p "Enter version to force: " version; \
	docker compose exec db migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" force $$version

# 現在のマイグレーションバージョンを確認
migrate-version:
	docker compose exec db migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" version

# 全てのテストを実行
test: test-backend test-frontend

# バックエンドのテストを実行
test-backend:
	@echo "Running backend tests..."
	cd backend && go test -v ./...

# フロントエンドのテストを実行
test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm test

# フロントエンドのカバレッジレポートを生成
test-coverage:
	@echo "Running frontend tests with coverage..."
	cd frontend && npm run test:coverage

# 環境変数を.envから読み込み
include .env
export

