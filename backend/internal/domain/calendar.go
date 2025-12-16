package domain

import "time"

// Event イベントドメインモデル
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	AllDay      bool      `json:"all_day"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CalendarDay カレンダーの1日分のデータ
type CalendarDay struct {
	Date      time.Time `json:"date"`
	Day       int       `json:"day"`
	Weekday   string    `json:"weekday"`
	IsHoliday bool      `json:"is_holiday"`
	Holiday   string    `json:"holiday,omitempty"`
	Rokuyo    string    `json:"rokuyo"`
	Events    []Event   `json:"events"`
}

// Calendar カレンダー情報
type Calendar struct {
	Year  int           `json:"year"`
	Month int           `json:"month"`
	Days  []CalendarDay `json:"days"`
}

// Holiday 祝日情報
type Holiday struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
}
