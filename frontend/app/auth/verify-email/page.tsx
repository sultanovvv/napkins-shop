'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'

import { useAuth } from '@/components/auth-context'
import { ApiError, confirmEmailVerificationRequest } from '@/lib/auth-client'

type State = 'idle' | 'pending' | 'ok' | 'error'

export default function VerifyEmailPage() {
  const searchParams = useSearchParams()
  const token = searchParams.get('token')?.trim() ?? ''
  const { setUser, user } = useAuth()

  const [state, setState] = useState<State>('idle')
  const [error, setError] = useState<string | null>(null)

  // Подтверждаем токен один раз на mount; повторно не пытаемся, потому что
  // токен одноразовый — повторный confirm всё равно вернёт 400.
  useEffect(() => {
    if (!token || state !== 'idle') return
    setState('pending')
    confirmEmailVerificationRequest(token)
      .then((u) => {
        setUser(u)
        setState('ok')
      })
      .catch((err: unknown) => {
        if (err instanceof ApiError && err.status === 400) {
          setError('Ссылка недействительна или устарела.')
        } else {
          setError('Не удалось подтвердить email. Попробуйте ещё раз позже.')
        }
        setState('error')
      })
  }, [token, state, setUser])

  if (!token) {
    return (
      <div className="flex flex-col gap-3">
        <h1 className="font-serif text-2xl text-primary">Подтверждение email</h1>
        <p className="text-sm text-muted-foreground">
          В адресе должен быть параметр <code>?token=…</code> из письма.
        </p>
        <Link href="/" className="text-sm text-primary hover:underline">
          ← На главную
        </Link>
      </div>
    )
  }

  if (state === 'pending' || state === 'idle') {
    return (
      <div className="flex flex-col gap-3">
        <h1 className="font-serif text-2xl text-primary">Подтверждаем email…</h1>
        <p className="text-sm text-muted-foreground">Это займёт секунду.</p>
      </div>
    )
  }

  if (state === 'ok') {
    return (
      <div className="flex flex-col gap-3">
        <h1 className="font-serif text-2xl text-primary">Email подтверждён</h1>
        <p className="text-sm text-muted-foreground">
          {user?.email ? `${user.email} — теперь активный аккаунт.` : 'Готово.'}
        </p>
        <Link href="/" className="text-sm text-primary hover:underline">
          ← На главную
        </Link>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      <h1 className="font-serif text-2xl text-primary">Не получилось</h1>
      <p className="text-sm text-destructive">{error}</p>
      <Link href="/" className="text-sm text-muted-foreground hover:text-primary">
        ← На главную
      </Link>
    </div>
  )
}
