package service

import (
	"testing"
	"time"

	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

// MockEventRepository はテスト用のモックリポジトリ
type MockEventRepository struct {
	GetAllFunc         func() ([]domain.Event, error)
	GetByIDFunc        func(id int) (*domain.Event, error)
	GetByDateRangeFunc func(start, end time.Time) ([]domain.Event, error)
	CreateFunc         func(event *domain.Event) error
	UpdateFunc         func(event *domain.Event) error
	DeleteFunc         func(id int) error
}

func (m *MockEventRepository) GetAll() ([]domain.Event, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return []domain.Event{}, nil
}

func (m *MockEventRepository) GetByID(id int) (*domain.Event, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventRepository) GetByDateRange(start, end time.Time) ([]domain.Event, error) {
	if m.GetByDateRangeFunc != nil {
		return m.GetByDateRangeFunc(start, end)
	}
	return []domain.Event{}, nil
}

func (m *MockEventRepository) Create(event *domain.Event) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(event)
	}
	return nil
}

func (m *MockEventRepository) Update(event *domain.Event) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(event)
	}
	return nil
}

func (m *MockEventRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestNewEventService(t *testing.T) {
	repo := &MockEventRepository{}
	service := NewEventService(repo)

	if service == nil {
		t.Error("NewEventService should return a non-nil service")
	}
}

func TestEventService_GetAllEvents(t *testing.T) {
	mockEvents := []domain.Event{
		{
			ID:          1,
			Title:       "テストイベント1",
			Description: "説明1",
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(time.Hour),
			AllDay:      false,
		},
		{
			ID:          2,
			Title:       "テストイベント2",
			Description: "説明2",
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(2 * time.Hour),
			AllDay:      true,
		},
	}

	repo := &MockEventRepository{
		GetAllFunc: func() ([]domain.Event, error) {
			return mockEvents, nil
		},
	}

	service := NewEventService(repo)
	events, err := service.GetAllEvents()

	if err != nil {
		t.Errorf("GetAllEvents should not return error: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}

	if events[0].Title != "テストイベント1" {
		t.Errorf("Expected first event title 'テストイベント1', got '%s'", events[0].Title)
	}
}

func TestEventService_CreateEvent_Success(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(time.Hour)

	event := &domain.Event{
		Title:       "新しいイベント",
		Description: "説明",
		StartDate:   startDate,
		EndDate:     endDate,
		AllDay:      false,
	}

	repo := &MockEventRepository{
		CreateFunc: func(e *domain.Event) error {
			e.ID = 1
			e.CreatedAt = time.Now()
			e.UpdatedAt = time.Now()
			return nil
		},
	}

	service := NewEventService(repo)
	err := service.CreateEvent(event)

	if err != nil {
		t.Errorf("CreateEvent should not return error: %v", err)
	}

	if event.ID != 1 {
		t.Errorf("Expected event ID to be set to 1, got %d", event.ID)
	}
}

func TestEventService_CreateEvent_EmptyTitle(t *testing.T) {
	event := &domain.Event{
		Title:       "",
		Description: "説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	repo := &MockEventRepository{}
	service := NewEventService(repo)
	err := service.CreateEvent(event)

	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

func TestEventService_CreateEvent_InvalidDateRange(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(-time.Hour) // 終了日が開始日より前

	event := &domain.Event{
		Title:       "イベント",
		Description: "説明",
		StartDate:   startDate,
		EndDate:     endDate,
		AllDay:      false,
	}

	repo := &MockEventRepository{}
	service := NewEventService(repo)
	err := service.CreateEvent(event)

	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput for invalid date range, got %v", err)
	}
}

func TestEventService_UpdateEvent_Success(t *testing.T) {
	existingEvent := &domain.Event{
		ID:          1,
		Title:       "既存イベント",
		Description: "既存の説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	updatedEvent := &domain.Event{
		ID:          1,
		Title:       "更新されたイベント",
		Description: "更新された説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(2 * time.Hour),
		AllDay:      true,
	}

	repo := &MockEventRepository{
		GetByIDFunc: func(id int) (*domain.Event, error) {
			if id == 1 {
				return existingEvent, nil
			}
			return nil, nil
		},
		UpdateFunc: func(e *domain.Event) error {
			e.UpdatedAt = time.Now()
			return nil
		},
	}

	service := NewEventService(repo)
	err := service.UpdateEvent(updatedEvent)

	if err != nil {
		t.Errorf("UpdateEvent should not return error: %v", err)
	}
}

func TestEventService_UpdateEvent_NotFound(t *testing.T) {
	event := &domain.Event{
		ID:          999,
		Title:       "存在しないイベント",
		Description: "説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	repo := &MockEventRepository{
		GetByIDFunc: func(id int) (*domain.Event, error) {
			return nil, nil
		},
	}

	service := NewEventService(repo)
	err := service.UpdateEvent(event)

	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestEventService_DeleteEvent_Success(t *testing.T) {
	repo := &MockEventRepository{
		GetByIDFunc: func(id int) (*domain.Event, error) {
			if id == 1 {
				return &domain.Event{ID: 1, Title: "削除するイベント"}, nil
			}
			return nil, nil
		},
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	service := NewEventService(repo)
	err := service.DeleteEvent(1)

	if err != nil {
		t.Errorf("DeleteEvent should not return error: %v", err)
	}
}

func TestEventService_DeleteEvent_NotFound(t *testing.T) {
	repo := &MockEventRepository{
		GetByIDFunc: func(id int) (*domain.Event, error) {
			return nil, nil
		},
	}

	service := NewEventService(repo)
	err := service.DeleteEvent(999)

	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
