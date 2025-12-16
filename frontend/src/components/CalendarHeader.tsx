'use client'

import { useCalendar } from '@/contexts/CalendarContext'

export default function CalendarHeader() {
  const { currentYear, currentMonth, previousMonth, nextMonth, goToToday } = useCalendar()

  return (
    <div className="flex items-center justify-between mb-6">
      <h2 className="text-2xl font-bold text-gray-800">
        {currentYear}年 {currentMonth}月
      </h2>
      <div className="flex gap-2">
        <button
          onClick={previousMonth}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
        >
          前月
        </button>
        <button
          onClick={goToToday}
          className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition"
        >
          今日
        </button>
        <button
          onClick={nextMonth}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
        >
          次月
        </button>
      </div>
    </div>
  )
}
