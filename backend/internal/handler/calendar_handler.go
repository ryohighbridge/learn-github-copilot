package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/service"
)

type CalendarHandler struct {
	service *service.CalendarService
}

func NewCalendarHandler(service *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

// GetCalendar カレンダー取得
func (h *CalendarHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(vars["month"])
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	calendar, err := h.service.GetCalendar(year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

// GetHolidays 祝日一覧取得
func (h *CalendarHandler) GetHolidays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	holidays := h.service.GetHolidays(year)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holidays)
}
