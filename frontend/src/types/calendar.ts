export interface Event {
  id: number
  title: string
  description: string
  start_date: string
  end_date: string
  all_day: boolean
  created_at: string
  updated_at: string
}

export interface CalendarDay {
  date: string
  day: number
  weekday: string
  is_holiday: boolean
  holiday?: string
  rokuyo: string
  events: Event[]
}

export interface CalendarData {
  year: number
  month: number
  days: CalendarDay[]
}

export interface Holiday {
  date: string
  name: string
}
