package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

// MockCalendarService はテスト用のモックサービス
type MockCalendarService struct {
	GetCalendarFunc func(year, month int) (*domain.Calendar, error)
	GetHolidaysFunc func(year int) []domain.Holiday
}

func (m *MockCalendarService) GetCalendar(year, month int) (*domain.Calendar, error) {
	if m.GetCalendarFunc != nil {
		return m.GetCalendarFunc(year, month)
	}
	return nil, nil
}

func (m *MockCalendarService) GetHolidays(year int) []domain.Holiday {
	if m.GetHolidaysFunc != nil {
		return m.GetHolidaysFunc(year)
	}
	return []domain.Holiday{}
}

func TestNewCalendarHandler(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	if handler == nil {
		t.Error("NewCalendarHandler should return a non-nil handler")
	}
}

func TestCalendarHandler_GetCalendar_Success(t *testing.T) {
	mockCalendar := &domain.Calendar{
		Year:  2025,
		Month: 12,
		Days: []domain.CalendarDay{
			{
				Date:      time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				Day:       1,
				Weekday:   "月",
				Rokuyo:    "大安",
				IsHoliday: false,
				Holiday:   "",
			},
		},
	}

	service := &MockCalendarService{
		GetCalendarFunc: func(year, month int) (*domain.Calendar, error) {
			if year == 2025 && month == 12 {
				return mockCalendar, nil
			}
			return nil, errors.New("invalid parameters")
		},
	}

	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar/2025/12", nil)
	req = mux.SetURLVars(req, map[string]string{"year": "2025", "month": "12"})
	w := httptest.NewRecorder()

	handler.GetCalendar(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var calendar domain.Calendar
	if err := json.NewDecoder(w.Body).Decode(&calendar); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calendar.Year != 2025 {
		t.Errorf("Expected year 2025, got %d", calendar.Year)
	}

	if calendar.Month != 12 {
		t.Errorf("Expected month 12, got %d", calendar.Month)
	}

	if len(calendar.Days) != 1 {
		t.Errorf("Expected 1 day, got %d", len(calendar.Days))
	}
}

func TestCalendarHandler_GetCalendar_InvalidYear(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar/invalid/12", nil)
	req = mux.SetURLVars(req, map[string]string{"year": "invalid", "month": "12"})
	w := httptest.NewRecorder()

	handler.GetCalendar(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCalendarHandler_GetCalendar_InvalidMonth(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar/2025/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"year": "2025", "month": "invalid"})
	w := httptest.NewRecorder()

	handler.GetCalendar(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCalendarHandler_GetCalendar_MonthOutOfRange(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	tests := []struct {
		month      string
		shouldFail bool
	}{
		{"0", true},
		{"1", false},
		{"12", false},
		{"13", true},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/calendar/2025/"+test.month, nil)
		req = mux.SetURLVars(req, map[string]string{"year": "2025", "month": test.month})
		w := httptest.NewRecorder()

		handler.GetCalendar(w, req)

		if test.shouldFail {
			if w.Code != http.StatusBadRequest {
				t.Errorf("Month %s: Expected status code %d, got %d", test.month, http.StatusBadRequest, w.Code)
			}
		}
	}
}

func TestCalendarHandler_GetHolidays_Success(t *testing.T) {
	mockHolidays := []domain.Holiday{
		{
			Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Name: "元日",
		},
		{
			Date: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
			Name: "文化の日",
		},
	}

	service := &MockCalendarService{
		GetHolidaysFunc: func(year int) []domain.Holiday {
			if year == 2025 {
				return mockHolidays
			}
			return []domain.Holiday{}
		},
	}

	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/holidays/2025", nil)
	req = mux.SetURLVars(req, map[string]string{"year": "2025"})
	w := httptest.NewRecorder()

	handler.GetHolidays(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var holidays []domain.Holiday
	if err := json.NewDecoder(w.Body).Decode(&holidays); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(holidays) != 2 {
		t.Errorf("Expected 2 holidays, got %d", len(holidays))
	}

	if holidays[0].Name != "元日" {
		t.Errorf("Expected first holiday '元日', got '%s'", holidays[0].Name)
	}
}

func TestCalendarHandler_GetHolidays_InvalidYear(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/holidays/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"year": "invalid"})
	w := httptest.NewRecorder()

	handler.GetHolidays(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCalendarHandler_GetHolidays_EmptyYear(t *testing.T) {
	service := &MockCalendarService{}
	handler := NewCalendarHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/holidays/", nil)
	req = mux.SetURLVars(req, map[string]string{"year": ""})
	w := httptest.NewRecorder()

	handler.GetHolidays(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
