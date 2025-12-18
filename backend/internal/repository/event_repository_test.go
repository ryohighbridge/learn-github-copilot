package repository

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

// setupTestDB はテスト用のインメモリDBをセットアップ
// 注意: 実際の統合テストではPostgreSQLのテストコンテナを使用することを推奨
func setupTestDB(t *testing.T) *sql.DB {
	// このテストでは実際のDBが必要なため、
	// 実際にはDockerコンテナを使用したintegration testとして実装するのが望ましい
	// ここではrepositoryのインターフェースの構造が正しいことを確認
	t.Skip("Integration test - requires PostgreSQL database")
	return nil
}

func TestNewEventRepository(t *testing.T) {
	// モックDBを使用しない場合はスキップ
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)
	if repo == nil {
		t.Error("NewEventRepository should return a non-nil repository")
	}
}

// 以下のテストは実際のPostgreSQLデータベースが必要なため、
// 統合テストとして実装されています。
// 実際のテスト実行時にはテストコンテナを使用することを推奨します。

func TestEventRepository_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	event := &domain.Event{
		Title:       "統合テストイベント",
		Description: "テスト説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	err := repo.Create(event)
	if err != nil {
		t.Errorf("Create should not return error: %v", err)
	}

	if event.ID == 0 {
		t.Error("Event ID should be set after creation")
	}
}

func TestEventRepository_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	// まずイベントを作成
	event := &domain.Event{
		Title:       "取得テストイベント",
		Description: "テスト説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	err := repo.Create(event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// IDで取得
	retrieved, err := repo.GetByID(event.ID)
	if err != nil {
		t.Errorf("GetByID should not return error: %v", err)
	}

	if retrieved == nil {
		t.Error("GetByID should return an event")
		return
	}

	if retrieved.Title != event.Title {
		t.Errorf("Expected title '%s', got '%s'", event.Title, retrieved.Title)
	}
}

func TestEventRepository_GetAll_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	events, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error: %v", err)
	}

	if events == nil {
		t.Error("GetAll should return a non-nil slice")
	}
}

func TestEventRepository_Update_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	// まずイベントを作成
	event := &domain.Event{
		Title:       "更新前のイベント",
		Description: "更新前の説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	err := repo.Create(event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// イベントを更新
	event.Title = "更新後のイベント"
	event.Description = "更新後の説明"

	err = repo.Update(event)
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}

	// 更新されたイベントを取得
	updated, err := repo.GetByID(event.ID)
	if err != nil {
		t.Fatalf("Failed to get updated event: %v", err)
	}

	if updated.Title != "更新後のイベント" {
		t.Errorf("Expected title '更新後のイベント', got '%s'", updated.Title)
	}
}

func TestEventRepository_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	// まずイベントを作成
	event := &domain.Event{
		Title:       "削除するイベント",
		Description: "テスト説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	err := repo.Create(event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// イベントを削除
	err = repo.Delete(event.ID)
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	// 削除されたことを確認
	deleted, err := repo.GetByID(event.ID)
	if err != nil {
		t.Errorf("GetByID after delete should not return error: %v", err)
	}

	if deleted != nil {
		t.Error("Event should be deleted")
	}
}

func TestEventRepository_GetByDateRange_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewEventRepository(db)

	// テストイベントを作成
	now := time.Now()
	event1 := &domain.Event{
		Title:       "期間内イベント1",
		Description: "テスト説明",
		StartDate:   now,
		EndDate:     now.Add(time.Hour),
		AllDay:      false,
	}

	event2 := &domain.Event{
		Title:       "期間内イベント2",
		Description: "テスト説明",
		StartDate:   now.Add(2 * time.Hour),
		EndDate:     now.Add(3 * time.Hour),
		AllDay:      false,
	}

	repo.Create(event1)
	repo.Create(event2)

	// 期間内のイベントを取得
	start := now.Add(-1 * time.Hour)
	end := now.Add(4 * time.Hour)
	events, err := repo.GetByDateRange(start, end)

	if err != nil {
		t.Errorf("GetByDateRange should not return error: %v", err)
	}

	if len(events) < 2 {
		t.Errorf("Expected at least 2 events in range, got %d", len(events))
	}
}
