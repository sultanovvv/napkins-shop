import Link from 'next/link'
import type { ReactNode } from 'react'

// Узкая колонка-карточка по центру — общий лейаут для всех /auth/* страниц.
// Шапка/футер магазина намеренно не показываются: на форме логина важно
// убрать шум и оставить один CTA.
export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="border-b border-border bg-card">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4">
          <Link href="/" className="font-serif text-2xl italic text-primary">
            Много салфеток
          </Link>
        </div>
      </header>
      <main className="flex flex-1 items-center justify-center px-4 py-10">
        <div className="w-full max-w-sm rounded-lg border border-border bg-card p-6 shadow-sm">
          {children}
        </div>
      </main>
    </div>
  )
}
