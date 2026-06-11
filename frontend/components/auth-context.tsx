'use client'

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'

import {
  bootstrapAuth,
  loginRequest,
  logoutAllRequest,
  logoutRequest,
  registerRequest,
  type AuthUser,
} from '@/lib/auth-client'

type Status = 'loading' | 'authenticated' | 'anonymous'

type AuthContextValue = {
  user: AuthUser | null
  status: Status
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  logoutAll: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

// Сигнализирует подписчикам (CartProvider), что состояние пользователя
// изменилось — пора переподтянуть корзину. Использую CustomEvent на window,
// чтобы не плодить взаимных импортов между двумя контекстами.
const AUTH_CHANGED_EVENT = 'napkins:auth-changed'

export function emitAuthChanged() {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent(AUTH_CHANGED_EVENT))
  }
}

export function onAuthChanged(handler: () => void): () => void {
  if (typeof window === 'undefined') return () => {}
  window.addEventListener(AUTH_CHANGED_EVENT, handler)
  return () => window.removeEventListener(AUTH_CHANGED_EVENT, handler)
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [status, setStatus] = useState<Status>('loading')

  // На mount — пробуем восстановить сессию через /auth/refresh. HttpOnly
  // cookie уйдёт автоматически: если она валидна, придёт новый access; если
  // нет — пользователь остаётся гостем, никаких ошибок не показываем.
  //
  // emitAuthChanged() ЗДЕСЬ НЕ ВЫЗЫВАЕМ — подписчики (CartProvider) ждут
  // изменения `status` через useAuth() и сами решают, когда фетчить. Иначе
  // на каждом mount получаем два запроса cart-а.
  useEffect(() => {
    let cancelled = false
    bootstrapAuth().then((u) => {
      if (cancelled) return
      setUser(u)
      setStatus(u ? 'authenticated' : 'anonymous')
    })
    return () => {
      cancelled = true
    }
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const u = await loginRequest(email, password)
    setUser(u)
    setStatus('authenticated')
    emitAuthChanged()
  }, [])

  const register = useCallback(async (email: string, password: string) => {
    const u = await registerRequest(email, password)
    setUser(u)
    setStatus('authenticated')
    emitAuthChanged()
  }, [])

  const logout = useCallback(async () => {
    try {
      await logoutRequest()
    } finally {
      setUser(null)
      setStatus('anonymous')
      emitAuthChanged()
    }
  }, [])

  const logoutAll = useCallback(async () => {
    try {
      await logoutAllRequest()
    } finally {
      setUser(null)
      setStatus('anonymous')
      emitAuthChanged()
    }
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({ user, status, login, register, logout, logoutAll }),
    [user, status, login, register, logout, logoutAll],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
