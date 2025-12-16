'use client'

import { useEffect } from 'react'
import { useCalendar } from '@/contexts/CalendarContext'
import CalendarGrid from './CalendarGrid'
import CalendarHeader from './CalendarHeader'

export default function Calendar() {
  const { currentYear, currentMonth, fetchCalendar, fetchEvents } = useCalendar()

  useEffect(() => {
    fetchCalendar(currentYear, currentMonth)
    fetchEvents()
  }, [currentYear, currentMonth, fetchCalendar, fetchEvents])

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <CalendarHeader />
      <CalendarGrid />
    </div>
  )
}
