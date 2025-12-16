const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const api = {
  async getCalendar(year: number, month: number) {
    const response = await fetch(`${API_BASE_URL}/api/calendar/${year}/${month}`)
    if (!response.ok) {
      throw new Error('Failed to fetch calendar')
    }
    return response.json()
  },

  async getHolidays(year: number) {
    const response = await fetch(`${API_BASE_URL}/api/holidays/${year}`)
    if (!response.ok) {
      throw new Error('Failed to fetch holidays')
    }
    return response.json()
  },

  async getEvents() {
    const response = await fetch(`${API_BASE_URL}/api/events`)
    if (!response.ok) {
      throw new Error('Failed to fetch events')
    }
    return response.json()
  },

  async createEvent(event: {
    title: string
    description: string
    start_date: string
    end_date: string
    all_day: boolean
  }) {
    const response = await fetch(`${API_BASE_URL}/api/events`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(event),
    })
    if (!response.ok) {
      throw new Error('Failed to create event')
    }
    return response.json()
  },

  async updateEvent(id: number, event: {
    title: string
    description: string
    start_date: string
    end_date: string
    all_day: boolean
  }) {
    const response = await fetch(`${API_BASE_URL}/api/events/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(event),
    })
    if (!response.ok) {
      throw new Error('Failed to update event')
    }
    return response.json()
  },

  async deleteEvent(id: number) {
    const response = await fetch(`${API_BASE_URL}/api/events/${id}`, {
      method: 'DELETE',
    })
    if (!response.ok) {
      throw new Error('Failed to delete event')
    }
  },
}
