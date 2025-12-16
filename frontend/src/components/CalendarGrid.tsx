'use client'

import { useCalendar } from '@/contexts/CalendarContext'
import { format, parseISO } from 'date-fns'
import { ja } from 'date-fns/locale'

export default function CalendarGrid() {
  const { calendarData, loading, error } = useCalendar()

  if (loading) {
    return (
      <div className="flex justify-center items-center h-96">
        <div className="text-gray-500">読み込み中...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex justify-center items-center h-96">
        <div className="text-red-500">{error}</div>
      </div>
    )
  }

  if (!calendarData) {
    return null
  }

  const weekdays = ['日', '月', '火', '水', '木', '金', '土']
  const firstDayWeekday = parseISO(calendarData.days[0].date).getDay()
  const emptyCells = Array(firstDayWeekday).fill(null)

  const getCellColor = (weekday: string, isHoliday: boolean) => {
    if (isHoliday) return 'text-holiday'
    if (weekday === '日') return 'text-sunday'
    if (weekday === '土') return 'text-saturday'
    return 'text-gray-700'
  }

  const isToday = (dateStr: string) => {
    const today = new Date()
    const date = parseISO(dateStr)
    return (
      date.getDate() === today.getDate() &&
      date.getMonth() === today.getMonth() &&
      date.getFullYear() === today.getFullYear()
    )
  }

  return (
    <div className="border border-gray-200 rounded-lg overflow-hidden">
      {/* 曜日ヘッダー */}
      <div className="grid grid-cols-7 bg-gray-100">
        {weekdays.map((day, index) => (
          <div
            key={day}
            className={`py-3 text-center font-semibold ${
              index === 0 ? 'text-sunday' : index === 6 ? 'text-saturday' : 'text-gray-700'
            }`}
          >
            {day}
          </div>
        ))}
      </div>

      {/* カレンダーグリッド */}
      <div className="grid grid-cols-7">
        {/* 空白セル */}
        {emptyCells.map((_, index) => (
          <div key={`empty-${index}`} className="border-t border-r border-gray-200 min-h-24 bg-gray-50" />
        ))}

        {/* 日付セル */}
        {calendarData.days.map((day) => (
          <div
            key={day.date}
            className={`border-t border-r border-gray-200 min-h-24 p-2 ${
              isToday(day.date) ? 'bg-blue-50' : 'bg-white'
            } hover:bg-gray-50 transition`}
          >
            <div className="flex flex-col h-full">
              <div className={`text-right font-bold ${getCellColor(day.weekday, day.is_holiday)}`}>
                {day.day}
              </div>
              {day.is_holiday && (
                <div className="text-xs text-holiday mt-1 truncate">{day.holiday}</div>
              )}
              <div className="text-xs text-gray-500 mt-1">{day.rokuyo}</div>
              {day.events && day.events.length > 0 && (
                <div className="mt-auto">
                  {day.events.map((event) => (
                    <div
                      key={event.id}
                      className="text-xs bg-blue-100 text-blue-800 rounded px-1 py-0.5 mb-1 truncate"
                    >
                      {event.title}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
