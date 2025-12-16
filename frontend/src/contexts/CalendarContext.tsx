'use client'

import React, { createContext, useContext, useState, useCallback } from 'react'
import { CalendarData, Event } from '@/types/calendar'
import { api } from '@/lib/api'

interface CalendarContextType {
  currentYear: number
  currentMonth: number
  calendarData: CalendarData | null
  events: Event[]
  loading: boolean
  error: string | null
  setCurrentDate: (year: number, month: number) => void
  fetchCalendar: (year: number, month: number) => Promise<void>
  fetchEvents: () => Promise<void>
  createEvent: (event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => Promise<void>
  updateEvent: (id: number, event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => Promise<void>
  deleteEvent: (id: number) => Promise<void>
  nextMonth: () => void
  previousMonth: () => void
  goToToday: () => void
}

const CalendarContext = createContext<CalendarContextType | undefined>(undefined)

export function CalendarProvider({ children }: { children: React.ReactNode }) {
  const now = new Date()
  const [currentYear, setCurrentYear] = useState(now.getFullYear())
  const [currentMonth, setCurrentMonth] = useState(now.getMonth() + 1)
  const [calendarData, setCalendarData] = useState<CalendarData | null>(null)
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const setCurrentDate = useCallback((year: number, month: number) => {
    setCurrentYear(year)
    setCurrentMonth(month)
  }, [])

  const fetchCalendar = useCallback(async (year: number, month: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await api.getCalendar(year, month)
      setCalendarData(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'カレンダーの取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchEvents = useCallback(async () => {
    try {
      const data = await api.getEvents()
      setEvents(data)
    } catch (err) {
      console.error('Failed to fetch events:', err)
    }
  }, [])

  const createEvent = useCallback(async (event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      await api.createEvent(event)
      await fetchEvents()
      await fetchCalendar(currentYear, currentMonth)
    } catch (err) {
      throw new Error('イベントの作成に失敗しました')
    }
  }, [currentYear, currentMonth, fetchEvents, fetchCalendar])

  const updateEvent = useCallback(async (id: number, event: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      await api.updateEvent(id, event)
      await fetchEvents()
      await fetchCalendar(currentYear, currentMonth)
    } catch (err) {
      throw new Error('イベントの更新に失敗しました')
    }
  }, [currentYear, currentMonth, fetchEvents, fetchCalendar])

  const deleteEvent = useCallback(async (id: number) => {
    try {
      await api.deleteEvent(id)
      await fetchEvents()
      await fetchCalendar(currentYear, currentMonth)
    } catch (err) {
      throw new Error('イベントの削除に失敗しました')
    }
  }, [currentYear, currentMonth, fetchEvents, fetchCalendar])

  const nextMonth = useCallback(() => {
    if (currentMonth === 12) {
      setCurrentYear(currentYear + 1)
      setCurrentMonth(1)
    } else {
      setCurrentMonth(currentMonth + 1)
    }
  }, [currentYear, currentMonth])

  const previousMonth = useCallback(() => {
    if (currentMonth === 1) {
      setCurrentYear(currentYear - 1)
      setCurrentMonth(12)
    } else {
      setCurrentMonth(currentMonth - 1)
    }
  }, [currentYear, currentMonth])

  const goToToday = useCallback(() => {
    const now = new Date()
    setCurrentYear(now.getFullYear())
    setCurrentMonth(now.getMonth() + 1)
  }, [])

  const value: CalendarContextType = {
    currentYear,
    currentMonth,
    calendarData,
    events,
    loading,
    error,
    setCurrentDate,
    fetchCalendar,
    fetchEvents,
    createEvent,
    updateEvent,
    deleteEvent,
    nextMonth,
    previousMonth,
    goToToday,
  }

  return (
    <CalendarContext.Provider value={value}>
      {children}
    </CalendarContext.Provider>
  )
}

export function useCalendar() {
  const context = useContext(CalendarContext)
  if (context === undefined) {
    throw new Error('useCalendar must be used within a CalendarProvider')
  }
  return context
}
