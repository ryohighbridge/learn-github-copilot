import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { CalendarProvider } from '@/contexts/CalendarContext'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: '日本のカレンダー',
  description: '日本の祝日・六曜を表示するカレンダーアプリ',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body className={inter.className}>
        <CalendarProvider>
          {children}
        </CalendarProvider>
      </body>
    </html>
  )
}
