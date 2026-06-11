'use client'

import Link from 'next/link'
import { useState, type ReactNode } from 'react'
import { LogOut, Search, User, ShoppingBasket, Menu } from 'lucide-react'
import { useCart } from '@/components/cart-context'
import { useAuth } from '@/components/auth-context'
import { CartDrawer } from '@/components/cart-drawer'
import { EmailVerifyBanner } from '@/components/email-verify-banner'
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
} from '@/components/ui/sheet'

export function SiteHeader({ mobileCategoryNav }: { mobileCategoryNav?: ReactNode }) {
  const { totalItems } = useCart()
  const { user, status, logout } = useAuth()
  const [cartOpen, setCartOpen] = useState(false)

  return (
    <header className="border-b border-border bg-card">
      <EmailVerifyBanner />
      {/* Top utility bar */}
      <div className="bg-secondary text-secondary-foreground">
        <div className="mx-auto flex max-w-7xl items-center justify-end gap-4 px-4 py-2 text-sm">
          <Link href="/contacts" className="hover:text-primary">
            Контакты
          </Link>
          <span className="text-border">|</span>
          <Link href="/sitemap" className="hover:text-primary">
            Карта сайта
          </Link>
        </div>
      </div>

      {/* Logo + actions */}
      <div className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-5">
        <div className="flex items-center gap-3">
          {/* Mobile nav trigger */}
          <Sheet>
            <SheetTrigger
              className="rounded-md p-2 text-primary hover:bg-secondary lg:hidden"
              aria-label="Открыть меню категорий"
            >
              <Menu className="size-6" />
            </SheetTrigger>
            <SheetContent side="left" className="w-80 overflow-y-auto p-0">
              <SheetTitle className="px-5 pt-5 font-serif text-xl text-primary">
                Категории
              </SheetTitle>
              <div className="p-3">{mobileCategoryNav}</div>
            </SheetContent>
          </Sheet>

          <Link href="/" className="flex flex-col leading-none">
            <span className="font-serif text-3xl italic text-primary md:text-4xl">
              Много салфеток
            </span>
            <span className="mt-1 text-xs tracking-wide text-muted-foreground">
              для любителей декупажа и коллекционеров
            </span>
          </Link>
        </div>

        <div className="flex items-center gap-2">
          {status === 'authenticated' && user ? (
            <div className="hidden items-center gap-1 sm:flex">
              <span
                className="flex items-center gap-2 rounded-md px-3 py-2 text-sm text-foreground"
                title={user.email}
              >
                <User className="size-5 text-primary" />
                <span className="max-w-[160px] truncate">{user.email}</span>
              </span>
              <button
                type="button"
                onClick={() => {
                  void logout()
                }}
                className="rounded-md p-2 text-muted-foreground hover:bg-secondary hover:text-destructive"
                aria-label="Выйти"
                title="Выйти"
              >
                <LogOut className="size-4" />
              </button>
            </div>
          ) : (
            <Link
              href="/auth/login"
              className="hidden items-center gap-2 rounded-md px-3 py-2 text-sm text-foreground hover:bg-secondary sm:flex"
            >
              <User className="size-5 text-primary" />
              {status === 'loading' ? '…' : 'Вход'}
            </Link>
          )}
          <button
            type="button"
            onClick={() => setCartOpen(true)}
            className="relative flex items-center gap-2 rounded-md bg-primary px-3 py-2 text-sm text-primary-foreground transition-colors hover:bg-primary/90"
          >
            <ShoppingBasket className="size-5" />
            <span className="hidden sm:inline">Корзина</span>
            {totalItems > 0 && (
              <span className="absolute -right-1.5 -top-1.5 flex size-5 items-center justify-center rounded-full bg-accent-foreground text-xs font-semibold text-card">
                {totalItems}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Main nav + search */}
      <div className="border-t border-border bg-card">
        <div className="mx-auto flex max-w-7xl flex-col items-stretch gap-3 px-4 py-3 md:flex-row md:items-center md:justify-between">
          <nav className="flex items-center gap-1 font-serif text-lg">
            <Link
              href="/"
              className="rounded-md px-3 py-1.5 italic text-foreground hover:text-primary"
            >
              Главная
            </Link>
            <Link
              href="/about"
              className="rounded-md px-3 py-1.5 italic text-foreground hover:text-primary"
            >
              О магазине
            </Link>
            <Link
              href="/how-to-buy"
              className="rounded-md px-3 py-1.5 italic text-foreground hover:text-primary"
            >
              Как покупать
            </Link>
          </nav>

          <form
            action="/search"
            className="flex w-full items-center overflow-hidden rounded-md border border-input bg-background md:w-80"
          >
            <input
              type="search"
              name="q"
              placeholder="Поиск по сюжетам…"
              className="w-full bg-transparent px-3 py-2 text-sm outline-none placeholder:text-muted-foreground"
            />
            <button
              type="submit"
              className="flex items-center gap-1 bg-primary px-3 py-2 text-sm text-primary-foreground"
              aria-label="Искать"
            >
              <Search className="size-4" />
            </button>
          </form>
        </div>
      </div>

      <CartDrawer open={cartOpen} onOpenChange={setCartOpen} />
    </header>
  )
}
