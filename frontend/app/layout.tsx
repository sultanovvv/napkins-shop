import { Analytics } from '@vercel/analytics/next'
import type { Metadata } from 'next'
import { Geist } from 'next/font/google'
import { Playfair_Display } from 'next/font/google'
import './globals.css'
import { CartProvider } from '@/components/cart-context'

const geistSans = Geist({ variable: '--font-geist-sans', subsets: ['latin'] })
const playfair = Playfair_Display({
  variable: '--font-playfair',
  subsets: ['latin', 'cyrillic'],
  weight: ['400', '500', '600', '700'],
  style: ['normal', 'italic'],
})

export const metadata: Metadata = {
  title: 'Много салфеток — интернет-магазин для декупажа',
  description:
    'Декупажные салфетки, карты и рисовая бумага для творчества и коллекционеров. Большой выбор сюжетов: цветы, птицы, гномы, города и многое другое.',
  generator: 'v0.app',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html
      lang="ru"
      className={`${geistSans.variable} ${playfair.variable} bg-background`}
    >
      <body className="font-sans antialiased">
        <CartProvider>{children}</CartProvider>
        {process.env.NODE_ENV === 'production' && <Analytics />}
      </body>
    </html>
  )
}
