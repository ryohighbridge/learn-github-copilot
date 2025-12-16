# Webアプリケーションアーキテクチャ設計書

## 概要

本ドキュメントは、日本人向けカレンダーWebアプリケーションのアーキテクチャ設計を定義します。

## システムアーキテクチャ

### 全体構成図

```
┌─────────────────────────────────────────┐
│      ユーザー（ブラウザ）                │
└──────────────┬──────────────────────────┘
               │ HTTP/HTTPS
               │
┌──────────────▼──────────────────────────┐
│      Next.js Frontend (Port 3000)       │
│  ┌────────────────────────────────────┐ │
│  │  App Router (React 18)             │ │
│  │  - Server Components               │ │
│  │  - Client Components               │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  State Management                  │ │
│  │  - React Context API               │ │
│  │  - CalendarContext                 │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  UI Layer                          │ │
│  │  - Tailwind CSS                    │ │
│  │  - Responsive Design               │ │
│  └────────────────────────────────────┘ │
└──────────────┬──────────────────────────┘
               │ REST API (JSON)
               │
┌──────────────▼──────────────────────────┐
│      Go Backend API (Port 8080)         │
│  ┌────────────────────────────────────┐ │
│  │  Handler Layer                     │ │
│  │  - HTTP Request/Response           │ │
│  │  - Routing (Gorilla Mux)           │ │
│  │  - CORS Middleware                 │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  Service Layer                     │ │
│  │  - Business Logic                  │ │
│  │  - Validation                      │ │
│  │  - Holiday/Rokuyo Calculation      │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  Repository Layer                  │ │
│  │  - Data Access                     │ │
│  │  - DB Connection Pool              │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  Domain Layer                      │ │
│  │  - Domain Models                   │ │
│  │  - Business Rules                  │ │
│  └────────────────────────────────────┘ │
└──────────────┬──────────────────────────┘
               │ SQL (PostgreSQL Protocol)
               │
┌──────────────▼──────────────────────────┐
│      PostgreSQL Database (Port 5432)    │
│  - Events Table                         │
│  - Indexes                              │
│  - Triggers                             │
└─────────────────────────────────────────┘
```

## 技術スタック

### フロントエンド

| 技術 | バージョン | 用途 |
|------|-----------|------|
| Next.js | 14.x | Reactフレームワーク（App Router） |
| React | 18.x | UIライブラリ |
| TypeScript | 5.x | 型安全な開発 |
| Tailwind CSS | 3.x | スタイリング |
| date-fns | 3.x | 日付操作 |

### バックエンド

| 技術 | バージョン | 用途 |
|------|-----------|------|
| Go | 1.21 | プログラミング言語 |
| Gorilla Mux | 1.8 | HTTPルーティング |
| lib/pq | 1.10 | PostgreSQLドライバ |
| rs/cors | 1.10 | CORS対応 |

### インフラ

| 技術 | バージョン | 用途 |
|------|-----------|------|
| Docker | Latest | コンテナ化 |
| Docker Compose | Latest | マルチコンテナ管理 |
| PostgreSQL | 16 | リレーショナルデータベース |

## アーキテクチャパターン

### バックエンド: クリーンアーキテクチャ

```
┌─────────────────────────────────────────────────────┐
│                   Handler Layer                     │
│  - HTTPリクエスト受信                                │
│  - レスポンス生成                                    │
│  - ルーティング                                      │
│  依存: Service Interface                            │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                  Service Layer                      │
│  - ビジネスロジック                                  │
│  - バリデーション                                    │
│  - トランザクション制御                              │
│  依存: Repository Interface                         │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                 Repository Layer                    │
│  - データアクセス                                    │
│  - SQL実行                                          │
│  - データマッピング                                  │
│  依存: Database, Domain Models                      │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                   Domain Layer                      │
│  - エンティティ                                      │
│  - ビジネスルール                                    │
│  - カスタムエラー                                    │
│  依存: なし（最も内側の層）                          │
└─────────────────────────────────────────────────────┘
```

#### 依存関係のルール

1. **依存は内側に向かう**: 外側の層は内側の層に依存できるが、逆は不可
2. **インターフェースによる抽象化**: 層間はインターフェースを介して通信
3. **ドメイン層の独立性**: ドメイン層は他の層に依存しない

### フロントエンド: レイヤードアーキテクチャ

```
┌─────────────────────────────────────────────────────┐
│                  Presentation Layer                 │
│  - Pages (App Router)                               │
│  - Components                                       │
│  - UI State Management                              │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                Application Layer                    │
│  - Context Providers                                │
│  - Custom Hooks                                     │
│  - Business Logic                                   │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                  Data Access Layer                  │
│  - API Client                                       │
│  - HTTP Communication                               │
│  - Data Transformation                              │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                    Types Layer                      │
│  - TypeScript Interfaces                            │
│  - Type Definitions                                 │
│  - Constants                                        │
└─────────────────────────────────────────────────────┘
```

## ディレクトリ構造

### バックエンド

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/                     # ドメイン層
│   │   ├── calendar.go            # カレンダードメインモデル
│   │   └── errors.go              # カスタムエラー定義
│   ├── repository/                 # リポジトリ層
│   │   ├── db.go                  # DB接続管理
│   │   └── event_repository.go    # イベントデータアクセス
│   ├── service/                    # サービス層
│   │   ├── calendar_service.go    # カレンダービジネスロジック
│   │   └── event_service.go       # イベントビジネスロジック
│   └── handler/                    # ハンドラー層
│       ├── calendar_handler.go    # カレンダーHTTPハンドラー
│       └── event_handler.go       # イベントHTTPハンドラー
├── pkg/                            # 公開パッケージ（将来の拡張用）
├── test/                           # テストコード
│   ├── integration/               # 統合テスト
│   └── testutil/                  # テストヘルパー
├── go.mod
└── go.sum
```

### フロントエンド

```
frontend/
├── src/
│   ├── app/                        # App Router
│   │   ├── layout.tsx             # ルートレイアウト
│   │   ├── page.tsx               # ホームページ
│   │   └── globals.css            # グローバルスタイル
│   ├── components/                 # UIコンポーネント
│   │   ├── Calendar.tsx           # カレンダーメインコンポーネント
│   │   ├── CalendarGrid.tsx       # カレンダーグリッド
│   │   └── CalendarHeader.tsx     # カレンダーヘッダー
│   ├── contexts/                   # React Context
│   │   └── CalendarContext.tsx    # カレンダー状態管理
│   ├── lib/                        # ユーティリティ
│   │   └── api.ts                 # APIクライアント
│   └── types/                      # 型定義
│       └── calendar.ts            # カレンダー型定義
├── public/                         # 静的ファイル
├── package.json
└── tsconfig.json
```

## API設計

### REST APIエンドポイント

#### カレンダーAPI

| メソッド | パス | 説明 | レスポンス |
|---------|------|------|-----------|
| GET | `/api/calendar/{year}/{month}` | 月次カレンダー取得 | Calendar |
| GET | `/api/holidays/{year}` | 年次祝日一覧取得 | []Holiday |

#### イベントAPI

| メソッド | パス | 説明 | レスポンス |
|---------|------|------|-----------|
| GET | `/api/events` | イベント一覧取得 | []Event |
| POST | `/api/events` | イベント作成 | Event |
| GET | `/api/events/{id}` | イベント詳細取得 | Event |
| PUT | `/api/events/{id}` | イベント更新 | Event |
| DELETE | `/api/events/{id}` | イベント削除 | 204 No Content |

#### ヘルスチェック

| メソッド | パス | 説明 | レスポンス |
|---------|------|------|-----------|
| GET | `/health` | サービス稼働確認 | "OK" |

### データモデル

#### Calendar

```typescript
{
  year: number,
  month: number,
  days: CalendarDay[]
}
```

#### CalendarDay

```typescript
{
  date: string,          // ISO 8601形式
  day: number,           // 日（1-31）
  weekday: string,       // 曜日（日本語）
  is_holiday: boolean,   // 祝日フラグ
  holiday?: string,      // 祝日名
  rokuyo: string,        // 六曜
  events: Event[]        // その日のイベント
}
```

#### Event

```typescript
{
  id: number,
  title: string,
  description: string,
  start_date: string,    // ISO 8601形式
  end_date: string,      // ISO 8601形式
  all_day: boolean,
  created_at: string,
  updated_at: string
}
```

#### Holiday

```typescript
{
  date: string,          // ISO 8601形式
  name: string           // 祝日名
}
```

## 状態管理

### React Context API

```typescript
interface CalendarContextType {
  // State
  currentYear: number
  currentMonth: number
  calendarData: CalendarData | null
  events: Event[]
  loading: boolean
  error: string | null
  
  // Actions
  setCurrentDate: (year: number, month: number) => void
  fetchCalendar: (year: number, month: number) => Promise<void>
  fetchEvents: () => Promise<void>
  createEvent: (event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => Promise<void>
  updateEvent: (id: number, event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => Promise<void>
  deleteEvent: (id: number) => Promise<void>
  nextMonth: () => void
  previousMonth: () => void
  goToToday: () => void
}
```

### データフロー

```
User Action
    ↓
Component Event Handler
    ↓
Context Action (useCalendar hook)
    ↓
API Call (lib/api.ts)
    ↓
Backend REST API
    ↓
Database
    ↓
Response
    ↓
Context State Update
    ↓
Component Re-render
```

## セキュリティ

### CORS設定

- 許可オリジン: `http://localhost:3000` (開発環境)
- 許可メソッド: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- 許可ヘッダー: `Content-Type`, `Authorization`

### 入力バリデーション

#### バックエンド

- タイトル: 必須、最大255文字
- 日付: RFC3339形式、終了日 >= 開始日
- パスパラメータ: 正規表現による型チェック

#### フロントエンド

- クライアントサイドバリデーション
- TypeScriptによる型安全性

### データベース

- プリペアドステートメント使用（SQLインジェクション対策）
- 接続プール管理
- トランザクション制御

## パフォーマンス最適化

### フロントエンド

1. **Server Side Rendering (SSR)**
   - Next.js App Routerによる初期表示高速化

2. **コンポーネント最適化**
   - React.memo による不要な再レンダリング防止
   - useCallback/useMemo によるメモ化

3. **バンドルサイズ最適化**
   - Tree shaking
   - Code splitting

### バックエンド

1. **データベース最適化**
   - インデックス設定（start_date, end_date）
   - コネクションプール

2. **キャッシング**
   - 祝日データのメモリキャッシュ（年単位）
   - HTTPキャッシュヘッダー設定

3. **非同期処理**
   - Goのgoroutineによる並行処理

## テスト戦略

### ユニットテストの容易性を考慮した設計

#### 1. インターフェース駆動開発

```go
// サービス層インターフェース
type CalendarServiceInterface interface {
    GetCalendar(year, month int) (*domain.Calendar, error)
    GetHolidays(year int) []domain.Holiday
}

type EventServiceInterface interface {
    GetAllEvents() ([]domain.Event, error)
    GetEventByID(id int) (*domain.Event, error)
    CreateEvent(event *domain.Event) error
    UpdateEvent(event *domain.Event) error
    DeleteEvent(id int) error
}

// リポジトリ層インターフェース
type EventRepositoryInterface interface {
    GetAll() ([]domain.Event, error)
    GetByID(id int) (*domain.Event, error)
    GetByDateRange(start, end time.Time) ([]domain.Event, error)
    Create(event *domain.Event) error
    Update(event *domain.Event) error
    Delete(id int) error
}
```

#### 2. 依存性注入

- コンストラクタインジェクション
- インターフェース型の依存
- モック・スタブの容易な差し替え

#### 3. 時刻依存コードの抽象化

```go
type TimeProvider interface {
    Now() time.Time
}

// 本番環境
type RealTimeProvider struct{}

// テスト環境
type MockTimeProvider struct {
    FixedTime time.Time
}
```

#### 4. 計算ロジックの分離

```go
type HolidayCalculator interface {
    Calculate(year int) []domain.Holiday
}

type RokuyoCalculator interface {
    Calculate(date time.Time) string
}
```

### テストレベル

#### バックエンド

1. **ユニットテスト**
   - 各関数・メソッドの単体テスト
   - モックを使用した依存関係の分離
   - カバレッジ目標: 80%以上

2. **統合テスト**
   - API エンドポイントのテスト
   - データベーステスト（テストコンテナ使用）
   - ビルドタグで分離: `// +build integration`

3. **E2Eテスト**
   - Docker Compose環境での全体テスト

#### フロントエンド

1. **コンポーネントテスト**
   - React Testing Library
   - スナップショットテスト

2. **統合テスト**
   - ユーザーインタラクションのテスト

3. **E2Eテスト**
   - Playwright/Cypress

### テストツール

| 層 | ツール | 用途 |
|----|--------|------|
| Backend Unit | Go標準testing | ユニットテスト |
| Backend Mock | gomock/testify | モック生成 |
| Backend DB | testcontainers-go | DBテスト |
| Frontend Unit | Jest | ユニットテスト |
| Frontend Component | React Testing Library | コンポーネントテスト |
| E2E | Playwright | E2Eテスト |

## 拡張性

### 将来の機能拡張

1. **認証・認可**
   - JWT認証の実装
   - ユーザーごとのイベント管理

2. **通知機能**
   - WebSocket/Server-Sent Events
   - プッシュ通知

3. **カレンダービュー拡張**
   - 週次ビュー
   - 日次ビュー
   - 年次ビュー

4. **データエクスポート**
   - iCal形式
   - CSV形式

5. **国際化 (i18n)**
   - 多言語対応
   - タイムゾーン対応

6. **モバイルアプリ**
   - React Native
   - 同一バックエンドAPI使用

### スケーラビリティ

1. **水平スケーリング**
   - ステートレスなAPI設計
   - ロードバランサー導入

2. **キャッシュ層**
   - Redis導入
   - カレンダーデータのキャッシュ

3. **CDN**
   - 静的アセットの配信最適化

## デプロイメント

### 環境

- **開発環境**: Docker Compose
- **ステージング**: Kubernetes (GKE/EKS/AKS)
- **本番環境**: Kubernetes (GKE/EKS/AKS)

### CI/CD

```
Git Push
    ↓
GitHub Actions
    ↓
Lint & Test
    ↓
Build Docker Images
    ↓
Push to Container Registry
    ↓
Deploy to Kubernetes
    ↓
Health Check
```

### モニタリング

- **ログ**: 構造化ログ（JSON形式）
- **メトリクス**: Prometheus + Grafana
- **トレーシング**: OpenTelemetry
- **アラート**: PagerDuty/Slack

## まとめ

本アーキテクチャは以下の原則に基づいて設計されています：

1. **関心の分離**: 各層が明確な責務を持つ
2. **テスタビリティ**: インターフェース駆動で容易にテスト可能
3. **保守性**: クリーンアーキテクチャによる変更容易性
4. **拡張性**: 新機能追加が容易な設計
5. **パフォーマンス**: 適切なキャッシング戦略
6. **型安全性**: TypeScript/Goによる静的型チェック

これらの設計原則により、長期的に保守・拡張可能なシステムを実現します。
