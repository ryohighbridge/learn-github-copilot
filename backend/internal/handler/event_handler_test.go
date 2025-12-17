package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

// MockEventService はテスト用のモックサービス
type MockEventService struct {
	GetAllEventsFunc func() ([]domain.Event, error)
	GetEventByIDFunc func(id int) (*domain.Event, error)
	CreateEventFunc  func(event *domain.Event) error
	UpdateEventFunc  func(event *domain.Event) error
	DeleteEventFunc  func(id int) error
}

func (m *MockEventService) GetAllEvents() ([]domain.Event, error) {
	if m.GetAllEventsFunc != nil {
		return m.GetAllEventsFunc()
	}
	return []domain.Event{}, nil
}

func (m *MockEventService) GetEventByID(id int) (*domain.Event, error) {
	if m.GetEventByIDFunc != nil {
		return m.GetEventByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventService) CreateEvent(event *domain.Event) error {
	if m.CreateEventFunc != nil {
		return m.CreateEventFunc(event)
	}
	return nil
}

func (m *MockEventService) UpdateEvent(event *domain.Event) error {
	if m.UpdateEventFunc != nil {
		return m.UpdateEventFunc(event)
	}
	return nil
}

func (m *MockEventService) DeleteEvent(id int) error {
	if m.DeleteEventFunc != nil {
		return m.DeleteEventFunc(id)
	}
	return nil
}

func TestNewEventHandler(t *testing.T) {
	service := &MockEventService{}
	handler := NewEventHandler(service)

	if handler == nil {
		t.Error("NewEventHandler should return a non-nil handler")
	}
}

func TestEventHandler_GetAllEvents_Success(t *testing.T) {
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

	service := &MockEventService{
		GetAllEventsFunc: func() ([]domain.Event, error) {
			return mockEvents, nil
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	w := httptest.NewRecorder()

	handler.GetEvents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var events []domain.Event
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestEventHandler_GetAllEvents_Error(t *testing.T) {
	service := &MockEventService{
		GetAllEventsFunc: func() ([]domain.Event, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	w := httptest.NewRecorder()

	handler.GetEvents(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestEventHandler_GetEventByID_Success(t *testing.T) {
	mockEvent := &domain.Event{
		ID:          1,
		Title:       "テストイベント",
		Description: "説明",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		AllDay:      false,
	}

	service := &MockEventService{
		GetEventByIDFunc: func(id int) (*domain.Event, error) {
			if id == 1 {
				return mockEvent, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/events/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetEvent(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var event domain.Event
	if err := json.NewDecoder(w.Body).Decode(&event); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if event.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", event.ID)
	}
}

func TestEventHandler_GetEventByID_NotFound(t *testing.T) {
	service := &MockEventService{
		GetEventByIDFunc: func(id int) (*domain.Event, error) {
			return nil, domain.ErrNotFound
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/events/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetEvent(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEventHandler_GetEventByID_InvalidID(t *testing.T) {
	service := &MockEventService{}
	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/events/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	w := httptest.NewRecorder()

	handler.GetEvent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestEventHandler_CreateEvent_Success(t *testing.T) {
	service := &MockEventService{
		CreateEventFunc: func(event *domain.Event) error {
			event.ID = 1
			event.CreatedAt = time.Now()
			event.UpdatedAt = time.Now()
			return nil
		},
	}

	handler := NewEventHandler(service)

	eventData := map[string]interface{}{
		"title":       "新しいイベント",
		"description": "説明",
		"start_date":  time.Now().Format(time.RFC3339),
		"end_date":    time.Now().Add(time.Hour).Format(time.RFC3339),
		"all_day":     false,
	}

	body, _ := json.Marshal(eventData)
	req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var event domain.Event
	if err := json.NewDecoder(w.Body).Decode(&event); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if event.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", event.ID)
	}
}

func TestEventHandler_CreateEvent_InvalidInput(t *testing.T) {
	service := &MockEventService{
		CreateEventFunc: func(event *domain.Event) error {
			return domain.ErrInvalidInput
		},
	}

	handler := NewEventHandler(service)

	eventData := map[string]interface{}{
		"title":       "",
		"description": "説明",
		"start_date":  time.Now().Format(time.RFC3339),
		"end_date":    time.Now().Add(time.Hour).Format(time.RFC3339),
		"all_day":     false,
	}

	body, _ := json.Marshal(eventData)
	req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestEventHandler_UpdateEvent_Success(t *testing.T) {
	service := &MockEventService{
		UpdateEventFunc: func(event *domain.Event) error {
			event.UpdatedAt = time.Now()
			return nil
		},
	}

	handler := NewEventHandler(service)

	eventData := map[string]interface{}{
		"title":       "更新されたイベント",
		"description": "更新された説明",
		"start_date":  time.Now().Format(time.RFC3339),
		"end_date":    time.Now().Add(time.Hour).Format(time.RFC3339),
		"all_day":     false,
	}

	body, _ := json.Marshal(eventData)
	req := httptest.NewRequest(http.MethodPut, "/api/events/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateEvent(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestEventHandler_UpdateEvent_NotFound(t *testing.T) {
	service := &MockEventService{
		UpdateEventFunc: func(event *domain.Event) error {
			return domain.ErrNotFound
		},
	}

	handler := NewEventHandler(service)

	eventData := map[string]interface{}{
		"title":       "更新されたイベント",
		"description": "更新された説明",
		"start_date":  time.Now().Format(time.RFC3339),
		"end_date":    time.Now().Add(time.Hour).Format(time.RFC3339),
		"all_day":     false,
	}

	body, _ := json.Marshal(eventData)
	req := httptest.NewRequest(http.MethodPut, "/api/events/999", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateEvent(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEventHandler_DeleteEvent_Success(t *testing.T) {
	service := &MockEventService{
		DeleteEventFunc: func(id int) error {
			return nil
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/api/events/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteEvent(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestEventHandler_DeleteEvent_NotFound(t *testing.T) {
	service := &MockEventService{
		DeleteEventFunc: func(id int) error {
			return domain.ErrNotFound
		},
	}

	handler := NewEventHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/api/events/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.DeleteEvent(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
