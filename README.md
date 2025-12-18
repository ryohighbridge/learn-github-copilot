# 日本のカレンダーWebアプリ

日本人向けのカレンダーWebアプリケーションです。

## 技術スタック

### バックエンド
- **言語**: Go 1.21
- **フレームワーク**: Gorilla Mux (REST API)
- **データベース**: PostgreSQL 16
- **アーキテクチャ**: クリーンアーキテクチャ（Handler → Service → Repository）

### フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS
- **状態管理**: React Context API
- **日付処理**: date-fns

### インフラ
- **コンテナ**: Docker & Docker Compose
- **開発環境**: ホットリロード対応
- **マイグレーション**: golang-migrate

## 機能

### 実装済み機能
- ✅ カレンダー表示（月次ビュー）
- ✅ 日本の祝日表示（国民の祝日法対応）
- ✅ 六曜表示（大安、赤口、先勝、友引、先負、仏滅）
- ✅ イベントCRUD機能
- ✅ 前月・次月ナビゲーション
- ✅ 今日へ移動機能
- ✅ レスポンシブデザイン

### API エンドポイント

**カレンダーAPI**
- `GET /api/calendar/{year}/{month}` - カレンダーデータ取得
- `GET /api/holidays/{year}` - 祝日一覧取得

**イベントAPI**
- `GET /api/events` - イベント一覧取得
- `POST /api/events` - イベント作成
- `GET /api/events/{id}` - イベント詳細取得
- `PUT /api/events/{id}` - イベント更新
- `DELETE /api/events/{id}` - イベント削除

## セットアップ

### 前提条件
- Docker
- Docker Compose v2以降

### インストール手順

1. **リポジトリのクローン**
```bash
git clone <repository-url>
cd learn-github-copilot
```

2. **環境変数ファイルの作成**
```bash
cp .env.sample .env
```

必要に応じて`.env`ファイルを編集してください：
```bash
# Backend API URL (for frontend)
NEXT_PUBLIC_API_URL=http://localhost:8080

# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=calendar_user
DB_PASSWORD=calendar_pass
DB_NAME=calendar_db
```

3. **Docker Composeで起動**
```bash
docker compose up --build
```

初回起動時は、イメージのビルドとデータベースの初期化が実行されます。
compose.yaml                # Docker Compose設定
├── .env.sample                 # 環境変数サンプル
├── .gitignore                  # Git除外設定
├── architecture.md             # アーキテクチャ設計書
├── README.md                   # このファイル
4. **アプリケーションにアクセス**
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080
- PostgreSQL: localhost:5432

### 使用方法

#### 基本操作

1. **カレンダーの表示**
   - ブラウザで http://localhost:3000 にアクセス
   - 現在の月のカレンダーが表示されます

2. **月の移動**
   - 「前月」ボタン：前月のカレンダーを表示
   - 「次月」ボタン：次月のカレンダーを表示
   - 「今日」ボタン：現在の月に戻る

3. **祝日と六曜の確認**
   - 祝日は赤色で表示され、祝日名が表示されます
   - 各日付に六曜（大安、友引など）が表示されます

4. **イベントの管理**（将来実装予定）
   - カレンダー上の日付をクリックしてイベントを作成
   - イベント詳細の編集・削除

#### APIの使用

バックエンドAPIを直接呼び出すこともできます：

```bash
# カレンダーデータの取得
curl http://localhost:8080/api/calendar/2025/12

# 祝日一覧の取得
curl http://localhost:8080/api/holidays/2025

# イベント一覧の取得
curl http://localhost:8080/api/events

# イベントの作成
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "会議",
    "description": "プロジェクト定例会議",
    "start_date": "2025-12-20T10:00:00Z",
    "end_date": "2025-12-20T11:00:00Z",
    "all_day": false
  }'
```

### 開発モード

各サービスは自動的にホットリロードされます。

**フロントエンド**
- `frontend/src`配下のファイル変更時に自動でリロード
- Next.js開発サーバーが起動

**バックエンド**
- 開発環境では`Dockerfile.dev`を使用してGoツールチェーンを含むイメージを使用
- `backend`配下のファイル変更時にコンテナを再起動
- ソースコードは`/app`にマウントされます
- テストは`make test-backend`または`cd backend && go test ./...`で実行可能

**データベース**
- データは`db-data`ボリュームに永続化されます
- マイグレーションは`db/migrations`ディレクトリで管理
- コンテナ起動時に自動でマイグレーションが実行されます

> **注意**: 本番環境では`backend/Dockerfile`（マルチステージビルド）を使用してください。開発環境では`backend/Dockerfile.dev`を使用してGoのツールチェーンを含めています。

## データベースマイグレーション

このプロジェクトでは[golang-migrate](https://github.com/golang-migrate/migrate)を使用してデータベーススキーマを管理しています。

### マイグレーションファイルの構造

マイグレーションファイルは`db/migrations`ディレクトリに配置されます：

```
db/migrations/
├── 000001_create_events_table.up.sql    # マイグレーション適用
└── 000001_create_events_table.down.sql  # マイグレーションロールバック
```

### Makefileを使用したマイグレーション管理

便利なMakeコマンドを用意しています：

```bash
# マイグレーションを最新バージョンまで適用
make migrate-up

# マイグレーションを1つロールバック
make migrate-down

# 現在のマイグレーションバージョンを確認
make migrate-version

# 新しいマイグレーションファイルを作成
make migrate-create
# 例: "add_users_table" という名前を入力すると
#   000002_add_users_table.up.sql
#   000002_add_users_table.down.sql
# が作成されます

# マイグレーションバージョンを強制設定（エラー時の回復用）
make migrate-force
```

### 手動でのマイグレーション実行

Makefileを使わずに直接実行することもできます：

```bash
# コンテナ内でマイグレーションを実行
docker compose exec db migrate \
  -path /migrations \
  -database "postgres://calendar_user:calendar_pass@localhost:5432/calendar_db?sslmode=disable" \
  up

# マイグレーションを1つ戻す
docker compose exec db migrate \
  -path /migrations \
  -database "postgres://calendar_user:calendar_pass@localhost:5432/calendar_db?sslmode=disable" \
  down 1

# 特定のバージョンまでマイグレーション
docker compose exec db migrate \
  -path /migrations \
  -database "postgres://calendar_user:calendar_pass@localhost:5432/calendar_db?sslmode=disable" \
  goto 1
```

### 新しいマイグレーションの作成手順

1. **Makefileを使用する場合（推奨）**

```bash
make migrate-create
# マイグレーション名を入力: add_categories_table
```

2. **手動で作成する場合**

```bash
# upファイル（適用用）
touch db/migrations/000002_add_categories_table.up.sql

# downファイル（ロールバック用）
touch db/migrations/000002_add_categories_table.down.sql
```

3. **マイグレーションファイルの記述**

`000002_add_categories_table.up.sql`:
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

`000002_add_categories_table.down.sql`:
```sql
DROP TABLE IF EXISTS categories;
```

4. **マイグレーションの適用**

```bash
make migrate-up
```

### マイグレーショントラブルシューティング

#### マイグレーションが失敗した場合

```bash
# 現在のバージョンを確認
make migrate-version

# エラーを確認
docker compose logs db

# 必要に応じて強制的にバージョンを設定
make migrate-force
# 例: バージョン1に強制設定
```

#### マイグレーションをやり直す場合

```bash
# 全てのマイグレーションをロールバック
docker compose exec db migrate \
  -path /migrations \
  -database "postgres://calendar_user:calendar_pass@localhost:5432/calendar_db?sslmode=disable" \
  down -all

# 再度適用
make migrate-up
```

#### データベースを完全にリセット

```bash
# コンテナとボリュームを削除
docker compose down -v

# 再起動（マイグレーションが自動実行される）
docker compose up -d
```

### コンテナの管理

```bash
# バックグラウンドで起動
docker compose up -d

# ログの確認
docker compose logs -f

# 特定のサービスのログ確認
docker compose logs -f backend

# コンテナの停止
docker compose down

# コンテナとボリュームの削除（データベースもリセット）
docker compose down -v

# 特定のサービスの再起動
docker compose restart backend
```

### トラブルシューティング

#### ポートが既に使用されている場合

```bash
# 使用中のポートを確認
sudo lsof -i :3000
sudo lsof -i :8080
sudo lsof -i :5432

# compose.yamlのポート設定を変更
# 例: "3001:3000" に変更
```

#### データベース接続エラー

```bash
# データベースコンテナのログを確認
docker compose logs db

# データベースのヘルスチェック
docker compose ps
```

#### フロントエンドのビルドエラー

```bash
# node_modulesを再インストール
docker compose down
docker compose up --build
```

## プロジェクト構造

```
.
├── compose.yaml                # Docker Compose設定
├── .env.sample                 # 環境変数サンプル
├── .gitignore                  # Git除外設定
├── architecture.md             # アーキテクチャ設計書
├── Makefile                    # マイグレーション管理用
├── README.md                   # このファイル
├── frontend/                   # Next.jsフロントエンド
│   ├── Dockerfile
│   ├── src/
│   │   ├── app/               # App Router
│   │   ├── components/        # Reactコンポーネント
│   │   ├── contexts/          # React Context
│   │   ├── lib/              # API クライアント
│   │   └── types/            # TypeScript型定義
│   └── package.json
├── backend/                    # Goバックエンド
│   ├── Dockerfile
│   ├── cmd/
│   │   └── api/              # エントリーポイント
│   ├── internal/
│   │   ├── domain/           # ドメインモデル
│   │   ├── handler/          # HTTPハンドラー
│   │   ├── service/          # ビジネスロジック
│   │   └── repository/       # データアクセス層
│   └── go.mod
└── db/
    ├── Dockerfile              # カスタムPostgreSQLイメージ
    ├── docker-entrypoint.sh   # マイグレーション自動実行スクリプト
    └── migrations/             # マイグレーションファイル
        ├── 000001_create_events_table.up.sql
        └── 000001_create_events_table.down.sql
```

## データベーススキーマ

### events テーブル
| カラム | 型 | 説明 |
|--------|-----|------|
| id | SERIAL | 主キー |
| title | VARCHAR(255) | イベントタイトル |
| description | TEXT | イベント詳細 |
| start_date | TIMESTAMP | 開始日時 |
| end_date | TIMESTAMP | 終了日時 |
| all_day | BOOLEAN | 終日フラグ |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

環境変数は`.env`ファイルで管理します。`.env.sample`をコピーして使用してください。

**バックエンド関連**
- `DB_HOST`: データベースホスト（デフォルト: `db`）
- `DB_PORT`: データベースポート（デフォルト: `5432`）
- `DB_USER`: データベースユーザー（デフォルト: `calendar_user`）
- `DB_PASSWORD`: データベースパスワード（デフォルト: `calendar_pass`）
- `DB_NAME`: データベース名（デフォルト: `calendar_db`）

**フロントエンド関連**
- `NEXT_PUBLIC_API_URL`: バックエンドAPIのURL（デフォルト: `http://backend:8080`）

### データベースのバックアップとリストア

```bash
# データベースのバックアップ
docker compose exec db pg_dump -U calendar_user calendar_db > backup.sql

# データベースのリストア
docker compose exec -T db psql -U calendar_user calendar_db < backup.sql
```
- `DB_NAME`: データベース名（デフォルト: calendar_db）

**フロントエンド**
- `NEXT_PUBLIC_API_URL`: バックエンドAPIのURL（デフォルト: http://backend:8080）

## テスト

このプロジェクトでは、バックエンドとフロントエンドの両方で包括的なユニットテストを実装しています。

### Makefileを使用したテスト実行

便利なMakeコマンドを用意しています：

```bash
# 全てのテストを実行（バックエンド + フロントエンド）
make test

# バックエンドのテストのみ実行
make test-backend

# フロントエンドのテストのみ実行
make test-frontend

# フロントエンドのカバレッジレポートを生成
make test-coverage
```

### バックエンドのテスト

Goの標準テストフレームワークを使用しています。

**テストの実行**

```bash
# Makefileを使用（推奨）
make test-backend

# または直接実行
docker compose exec backend go test -v ./...

# カバレッジ付きでテスト
docker compose exec backend go test -cover ./...

# 詳細なカバレッジレポート
docker compose exec backend go test -coverprofile=coverage.out ./...
docker compose exec backend go tool cover -html=coverage.out

# 特定のパッケージのみテスト
docker compose exec backend go test -v ./internal/service/...
```

**テスト構成**

- **Service層**: ビジネスロジックのテスト
  - `backend/internal/service/event_service_test.go`
  - `backend/internal/service/calendar_service_test.go`
  
- **Handler層**: HTTPハンドラーのテスト
  - `backend/internal/handler/event_handler_test.go`
  - `backend/internal/handler/calendar_handler_test.go`
  
- **Repository層**: データアクセス層の統合テスト
  - `backend/internal/repository/event_repository_test.go`
  - 注意: Repository層のテストは実際のPostgreSQLデータベースが必要なため、`-short`フラグでスキップされます

**短いテストの実行（統合テストをスキップ）**

```bash
docker compose exec backend go test -short ./...
```

### フロントエンドのテスト

Jest と React Testing Library を使用しています。

**テストの実行**

```bash
# Makefileを使用（推奨）
make test-frontend

# または直接実行
docker compose exec frontend npm test

# ウォッチモード
docker compose exec frontend npm run test:watch

# カバレッジレポート生成
make test-coverage
# または
docker compose exec frontend npm run test:coverage
```

**テスト構成**

- **コンポーネントテスト**:
  - `frontend/src/components/__tests__/CalendarHeader.test.tsx`
  - `frontend/src/components/__tests__/CalendarGrid.test.tsx`
  
- **Contextテスト**:
  - `frontend/src/contexts/__tests__/CalendarContext.test.tsx`
  
- **APIクライアントテスト**:
  - `frontend/src/lib/__tests__/api.test.ts`

**テスト設定ファイル**

- `frontend/jest.config.ts`: Jest設定
- `frontend/jest.setup.ts`: テストセットアップ（Testing Libraryの設定）

### テストのベストプラクティス

**バックエンド**

1. モックを使用してテストを独立させる
2. テーブル駆動テストを活用する
3. エラーケースも必ずテストする
4. 統合テストは`testing.Short()`でスキップできるようにする

**フロントエンド**

1. ユーザーの操作をシミュレートする
2. 実装の詳細ではなく、動作をテストする
3. APIはモックする
4. アクセシビリティを考慮したクエリを使用する

### 継続的インテグレーション

将来的にはGitHub Actionsなどを使用して、プルリクエスト時に自動的にテストを実行する予定です。

## 今後の拡張可能性

- 📅 週次・日次ビュー
- 🔔 イベント通知機能
- 👥 ユーザー認証・マルチユーザー対応
- 📱 PWA対応
- 🌙 旧暦表示
- 🎌 二十四節気表示
- 📊 イベント統計・分析
- 📤 カレンダーエクスポート（iCal形式）

## コントリビューション

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。

1. フォークする
2. フィーチャーブランチを作成 (`git switch -c feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

このプロジェクトは Apache License 2.0 のもとで公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 関連ドキュメント

- [アーキテクチャ設計書](architecture.md) - システムアーキテクチャの詳細
- [API仕様](architecture.md#api設計) - REST APIエンドポイントの詳細
