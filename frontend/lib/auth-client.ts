// auth-client: тонкий слой между React-кодом и Go-бэком.
//
// Access JWT хранится В ПАМЯТИ модуля. Это значит: hard-reload страницы
// = токен теряется, но на старте AuthProvider дёргает /auth/refresh
// (HttpOnly cookie уйдёт сама), и сессия восстанавливается. Никакого
// localStorage/sessionStorage — таким образом XSS не может украсть токен.
//
// `apiFetch` ниже умеет одну вещь, которую сырой fetch не умеет: при 401
// один раз пробует /auth/refresh и повторяет исходный запрос. Параллельные
// запросы делят одну refresh-Promise через `pendingRefresh`, чтобы не
// устроить шторм из refresh'ей.

// Пустая строка по умолчанию = ходим относительно текущего origin. Под
// капотом Next.js rewrites проксируют /api/* на Go-бэкенд (см. next.config.mjs).
// Это делает запросы same-origin и решает все cookie/CORS проблемы.
//
// Переопределить можно через NEXT_PUBLIC_API_BASE_URL (например, при деплое
// фронта и бэка на разные origin'ы).
export const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? ''

let accessToken: string | null = null
let pendingRefresh: Promise<boolean> | null = null

const onAuthChangeListeners = new Set<() => void>()

export function getAccessToken(): string | null {
  return accessToken
}

export function setAccessToken(token: string | null): void {
  accessToken = token
  for (const fn of onAuthChangeListeners) fn()
}

// onAuthChange — для AuthContext: «токен сбросили» = надо переснять user
// state и убрать пользователя из UI. Дёргается также из самого contex'та
// при login/logout.
export function onAuthChange(fn: () => void): () => void {
  onAuthChangeListeners.add(fn)
  return () => {
    onAuthChangeListeners.delete(fn)
  }
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
    public readonly code?: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

type ApiFetchOptions = RequestInit & {
  // Если true — на 401 НЕ пытаемся /refresh (используется самим /refresh
  // и /logout, иначе будет бесконечная рекурсия).
  skipAuthRefresh?: boolean
}

export async function apiFetch<T>(
  path: string,
  options: ApiFetchOptions = {},
): Promise<T> {
  const res = await rawFetch(path, options)
  if (res.status === 401 && !options.skipAuthRefresh) {
    const refreshed = await tryRefresh()
    if (refreshed) {
      const retry = await rawFetch(path, options)
      return handleResponse<T>(retry, path, options.method)
    }
  }
  return handleResponse<T>(res, path, options.method)
}

async function rawFetch(path: string, options: ApiFetchOptions): Promise<Response> {
  const headers = new Headers(options.headers)
  if (!headers.has('Content-Type') && options.body && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }
  if (accessToken && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${accessToken}`)
  }
  return fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  })
}

async function handleResponse<T>(res: Response, path: string, method?: string): Promise<T> {
  if (res.status === 204) {
    return undefined as T
  }
  if (!res.ok) {
    let message = `${method ?? 'GET'} ${path} → HTTP ${res.status}`
    let code: string | undefined
    try {
      const body = await res.json()
      if (body?.error?.message) message = body.error.message
      if (body?.error?.code) code = body.error.code
    } catch {
      // не JSON — оставляем сухое сообщение
    }
    throw new ApiError(res.status, message, code)
  }
  return (await res.json()) as T
}

// tryRefresh — единая точка обращения к /auth/refresh. Если несколько
// запросов получили 401 одновременно, refresh выполняется один раз; все
// ждут одну и ту же Promise.
async function tryRefresh(): Promise<boolean> {
  if (!pendingRefresh) {
    pendingRefresh = refreshNow().finally(() => {
      pendingRefresh = null
    })
  }
  return pendingRefresh
}

async function refreshNow(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/public/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!res.ok) {
      setAccessToken(null)
      return false
    }
    const data = (await res.json()) as { accessToken: string }
    setAccessToken(data.accessToken)
    return true
  } catch {
    setAccessToken(null)
    return false
  }
}

// bootstrapAuth — вызывается AuthProvider'ом на mount. Возвращает user или
// null. Под капотом — тот же refresh: cookie уйдёт автоматически, если есть.
export type AuthUser = {
  id: string
  email: string
  emailVerifiedAt: string | null
}

type LoginResponse = {
  user: AuthUser
  accessToken: string
  accessExpiresAt: string
  refreshExpiresAt: string
}

export async function bootstrapAuth(): Promise<AuthUser | null> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/public/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!res.ok) {
      setAccessToken(null)
      return null
    }
    const data = (await res.json()) as LoginResponse
    setAccessToken(data.accessToken)
    return data.user
  } catch {
    setAccessToken(null)
    return null
  }
}

export async function loginRequest(email: string, password: string): Promise<AuthUser> {
  const data = await apiFetch<LoginResponse>('/api/v1/public/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
    skipAuthRefresh: true,
  })
  setAccessToken(data.accessToken)
  return data.user
}

export async function registerRequest(email: string, password: string): Promise<AuthUser> {
  const data = await apiFetch<LoginResponse>('/api/v1/public/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
    skipAuthRefresh: true,
  })
  setAccessToken(data.accessToken)
  return data.user
}

export async function logoutRequest(): Promise<void> {
  await apiFetch<void>('/api/v1/public/auth/logout', {
    method: 'POST',
    skipAuthRefresh: true,
  })
  setAccessToken(null)
}

export async function logoutAllRequest(): Promise<void> {
  await apiFetch<void>('/api/v1/public/auth/logoutAll', {
    method: 'POST',
  })
  setAccessToken(null)
}

export async function requestPasswordResetRequest(email: string): Promise<void> {
  await apiFetch<void>('/api/v1/public/auth/requestPasswordReset', {
    method: 'POST',
    body: JSON.stringify({ email }),
    skipAuthRefresh: true,
  })
}

export async function confirmPasswordResetRequest(
  secret: string,
  newPassword: string,
): Promise<void> {
  await apiFetch<void>('/api/v1/public/auth/confirmPasswordReset', {
    method: 'POST',
    body: JSON.stringify({ secret, newPassword }),
    skipAuthRefresh: true,
  })
}

export function meRequest(): Promise<AuthUser> {
  return apiFetch<AuthUser>('/api/v1/public/auth/me')
}

export async function requestEmailVerificationRequest(): Promise<void> {
  await apiFetch<void>('/api/v1/public/auth/requestEmailVerification', {
    method: 'POST',
  })
}

export function confirmEmailVerificationRequest(secret: string): Promise<AuthUser> {
  return apiFetch<AuthUser>('/api/v1/public/auth/confirmEmailVerification', {
    method: 'POST',
    body: JSON.stringify({ secret }),
    skipAuthRefresh: true,
  })
}
