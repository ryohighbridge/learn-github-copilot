package service

import (
	"strconv"
	"time"

	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

type CalendarService struct{}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}

// GetCalendar 指定月のカレンダー情報を取得
func (s *CalendarService) GetCalendar(year, month int) (*domain.Calendar, error) {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	calendar := &domain.Calendar{
		Year:  year,
		Month: month,
		Days:  []domain.CalendarDay{},
	}

	holidays := s.GetHolidays(year)
	holidayMap := make(map[string]string)
	for _, h := range holidays {
		key := h.Date.Format("2006-01-02")
		holidayMap[key] = h.Name
	}

	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		holidayName, isHoliday := holidayMap[dateKey]

		day := domain.CalendarDay{
			Date:      d,
			Day:       d.Day(),
			Weekday:   s.getWeekdayJapanese(d.Weekday()),
			IsHoliday: isHoliday,
			Holiday:   holidayName,
			Rokuyo:    s.calculateRokuyo(d),
			Events:    []domain.Event{},
		}

		calendar.Days = append(calendar.Days, day)
	}

	return calendar, nil
}

// GetHolidays 指定年の祝日一覧を取得
func (s *CalendarService) GetHolidays(year int) []domain.Holiday {
	holidays := []domain.Holiday{}

	// 固定祝日
	fixedHolidays := map[string]string{
		"01-01": "元日",
		"02-11": "建国記念の日",
		"02-23": "天皇誕生日",
		"04-29": "昭和の日",
		"05-03": "憲法記念日",
		"05-04": "みどりの日",
		"05-05": "こどもの日",
		"08-11": "山の日",
		"11-03": "文化の日",
		"11-23": "勤労感謝の日",
	}

	for dateStr, name := range fixedHolidays {
		date, _ := time.Parse("2006-01-02", strconv.Itoa(year)+"-"+dateStr)
		holidays = append(holidays, domain.Holiday{Date: date, Name: name})
	}

	// 成人の日（1月第2月曜日）
	holidays = append(holidays, domain.Holiday{
		Date: s.getNthWeekday(year, 1, time.Monday, 2),
		Name: "成人の日",
	})

	// 海の日（7月第3月曜日）
	holidays = append(holidays, domain.Holiday{
		Date: s.getNthWeekday(year, 7, time.Monday, 3),
		Name: "海の日",
	})

	// 敬老の日（9月第3月曜日）
	holidays = append(holidays, domain.Holiday{
		Date: s.getNthWeekday(year, 9, time.Monday, 3),
		Name: "敬老の日",
	})

	// スポーツの日（10月第2月曜日）
	holidays = append(holidays, domain.Holiday{
		Date: s.getNthWeekday(year, 10, time.Monday, 2),
		Name: "スポーツの日",
	})

	// 春分の日・秋分の日
	holidays = append(holidays, domain.Holiday{
		Date: s.calculateShunbun(year),
		Name: "春分の日",
	})
	holidays = append(holidays, domain.Holiday{
		Date: s.calculateShubun(year),
		Name: "秋分の日",
	})

	return holidays
}

// getNthWeekday 指定月のN番目の曜日を取得
func (s *CalendarService) getNthWeekday(year, month int, weekday time.Weekday, n int) time.Time {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// 最初の指定曜日を見つける
	daysUntilWeekday := int(weekday - firstDay.Weekday())
	if daysUntilWeekday < 0 {
		daysUntilWeekday += 7
	}

	// N番目の曜日を計算
	return firstDay.AddDate(0, 0, daysUntilWeekday+(n-1)*7)
}

// calculateShunbun 春分の日を計算（簡易版）
func (s *CalendarService) calculateShunbun(year int) time.Time {
	// 簡易計算式（2000-2099年まで有効）
	day := 20
	if year >= 2000 && year <= 2099 {
		day = int(20.8431+0.242194*(float64(year)-1980)) - (year-1980)/4
	}
	return time.Date(year, 3, day, 0, 0, 0, 0, time.UTC)
}

// calculateShubun 秋分の日を計算（簡易版）
func (s *CalendarService) calculateShubun(year int) time.Time {
	// 簡易計算式（2000-2099年まで有効）
	day := 23
	if year >= 2000 && year <= 2099 {
		day = int(23.2488+0.242194*(float64(year)-1980)) - (year-1980)/4
	}
	return time.Date(year, 9, day, 0, 0, 0, 0, time.UTC)
}

// calculateRokuyo 六曜を計算
func (s *CalendarService) calculateRokuyo(date time.Time) string {
	rokuyo := []string{"大安", "赤口", "先勝", "友引", "先負", "仏滅"}

	// 旧暦計算の簡易版（正確な旧暦変換ではなく近似値）
	month := int(date.Month())
	day := date.Day()

	index := (month + day) % 6
	return rokuyo[index]
}

// getWeekdayJapanese 曜日の日本語名を取得
func (s *CalendarService) getWeekdayJapanese(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "日",
		time.Monday:    "月",
		time.Tuesday:   "火",
		time.Wednesday: "水",
		time.Thursday:  "木",
		time.Friday:    "金",
		time.Saturday:  "土",
	}
	return weekdays[weekday]
}
