package service

import (
	"testing"
	"time"
)

func TestNewCalendarService(t *testing.T) {
	service := NewCalendarService()

	if service == nil {
		t.Error("NewCalendarService should return a non-nil service")
	}
}

func TestCalendarService_GetCalendar(t *testing.T) {
	service := NewCalendarService()
	calendar, err := service.GetCalendar(2025, 12)

	if err != nil {
		t.Errorf("GetCalendar should not return error: %v", err)
	}

	if calendar == nil {
		t.Error("GetCalendar should return a non-nil calendar")
		return
	}

	if calendar.Year != 2025 {
		t.Errorf("Expected year 2025, got %d", calendar.Year)
	}

	if calendar.Month != 12 {
		t.Errorf("Expected month 12, got %d", calendar.Month)
	}

	// 12月は31日まである
	if len(calendar.Days) != 31 {
		t.Errorf("Expected 31 days in December, got %d", len(calendar.Days))
	}

	// 最初の日付が12月1日であることを確認
	firstDay := calendar.Days[0]
	if firstDay.Day != 1 {
		t.Errorf("Expected first day to be 1, got %d", firstDay.Day)
	}

	// 各日が適切な曜日を持っていることを確認
	for _, day := range calendar.Days {
		if day.Weekday == "" {
			t.Error("Day should have a weekday")
		}
		if day.Rokuyo == "" {
			t.Error("Day should have a rokuyo")
		}
	}
}

func TestCalendarService_GetHolidays(t *testing.T) {
	service := NewCalendarService()
	holidays := service.GetHolidays(2025)

	if len(holidays) == 0 {
		t.Error("GetHolidays should return at least one holiday")
	}

	// 主要な祝日が含まれているか確認
	foundNewYear := false
	foundCultureDay := false

	for _, holiday := range holidays {
		if holiday.Name == "元日" {
			foundNewYear = true
			if holiday.Date.Month() != time.January || holiday.Date.Day() != 1 {
				t.Errorf("元日 should be January 1st, got %v", holiday.Date)
			}
		}
		if holiday.Name == "文化の日" {
			foundCultureDay = true
			if holiday.Date.Month() != time.November || holiday.Date.Day() != 3 {
				t.Errorf("文化の日 should be November 3rd, got %v", holiday.Date)
			}
		}
	}

	if !foundNewYear {
		t.Error("元日 should be included in holidays")
	}
	if !foundCultureDay {
		t.Error("文化の日 should be included in holidays")
	}
}

func TestCalendarService_GetNthWeekday(t *testing.T) {
	service := NewCalendarService()

	// 2025年1月の第2月曜日（成人の日）
	seijinNoHi := service.getNthWeekday(2025, 1, time.Monday, 2)

	if seijinNoHi.Month() != time.January {
		t.Errorf("Expected January, got %v", seijinNoHi.Month())
	}

	if seijinNoHi.Weekday() != time.Monday {
		t.Errorf("Expected Monday, got %v", seijinNoHi.Weekday())
	}

	// 2025年1月の第2月曜日は13日
	if seijinNoHi.Day() != 13 {
		t.Errorf("Expected 13th, got %d", seijinNoHi.Day())
	}
}

func TestCalendarService_CalculateShunbun(t *testing.T) {
	service := NewCalendarService()

	// 2025年の春分の日
	shunbun := service.calculateShunbun(2025)

	if shunbun.Month() != time.March {
		t.Errorf("春分の日 should be in March, got %v", shunbun.Month())
	}

	// 春分の日は3月20日前後
	if shunbun.Day() < 19 || shunbun.Day() > 21 {
		t.Errorf("春分の日 should be around March 20th, got March %d", shunbun.Day())
	}
}

func TestCalendarService_CalculateShubun(t *testing.T) {
	service := NewCalendarService()

	// 2025年の秋分の日
	shubun := service.calculateShubun(2025)

	if shubun.Month() != time.September {
		t.Errorf("秋分の日 should be in September, got %v", shubun.Month())
	}

	// 秋分の日は9月22日前後
	if shubun.Day() < 21 || shubun.Day() > 24 {
		t.Errorf("秋分の日 should be around September 23rd, got September %d", shubun.Day())
	}
}

func TestCalendarService_CalculateRokuyo(t *testing.T) {
	service := NewCalendarService()

	date := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	rokuyo := service.calculateRokuyo(date)

	validRokuyo := map[string]bool{
		"大安": true,
		"赤口": true,
		"先勝": true,
		"友引": true,
		"先負": true,
		"仏滅": true,
	}

	if !validRokuyo[rokuyo] {
		t.Errorf("Invalid rokuyo: %s", rokuyo)
	}
}

func TestCalendarService_GetWeekdayJapanese(t *testing.T) {
	service := NewCalendarService()

	tests := []struct {
		weekday  time.Weekday
		expected string
	}{
		{time.Sunday, "日"},
		{time.Monday, "月"},
		{time.Tuesday, "火"},
		{time.Wednesday, "水"},
		{time.Thursday, "木"},
		{time.Friday, "金"},
		{time.Saturday, "土"},
	}

	for _, test := range tests {
		result := service.getWeekdayJapanese(test.weekday)
		if result != test.expected {
			t.Errorf("For %v, expected %s, got %s", test.weekday, test.expected, result)
		}
	}
}

func TestCalendarService_HolidayIntegrity(t *testing.T) {
	service := NewCalendarService()

	// 2025年のカレンダーと祝日を取得
	calendar, err := service.GetCalendar(2025, 1)
	if err != nil {
		t.Fatalf("Failed to get calendar: %v", err)
	}

	holidays := service.GetHolidays(2025)

	// 1月1日が祝日としてマークされているか確認
	janFirst := calendar.Days[0]
	if !janFirst.IsHoliday {
		t.Error("January 1st should be marked as a holiday")
	}
	if janFirst.Holiday != "元日" {
		t.Errorf("Expected holiday name '元日', got '%s'", janFirst.Holiday)
	}

	// 祝日の総数が妥当か確認（日本の祝日は年間15-17日程度）
	if len(holidays) < 10 || len(holidays) > 20 {
		t.Errorf("Unexpected number of holidays: %d", len(holidays))
	}
}
