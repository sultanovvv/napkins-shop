'use client'

import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import { useState, type FormEvent } from 'react'

import { ApiError, confirmPasswordResetRequest } from '@/lib/auth-client'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export default function ResetPasswordPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  // secret попадает в URL ссылкой из письма: /auth/reset-password?token=<…>
  const tokenFromQuery = searchParams.get('token') ?? ''

  const [secret, setSecret] = useState(tokenFromQuery)
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [done, setDone] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setIsSubmitting(true)
    try {
      await confirmPasswordResetRequest(secret.trim(), password)
      setDone(true)
      // Все активные refresh-токены пользователя бэк отзывает; перенаправляем
      // на логин с явным flash-параметром, чтобы фронт мог показать «вход
      // с новым паролем», если когда-нибудь захочет.
      setTimeout(() => router.push('/auth/login?reset=ok'), 1200)
    } catch (err) {
      if (err instanceof ApiError && err.status === 400) {
        setError('Ссылка недействительна или устарела. Запросите новую.')
      } else {
        setError('Не удалось сменить пароль. Попробуйте ещё раз.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  if (done) {
    return (
      <div className="flex flex-col gap-3">
        <h1 className="font-serif text-2xl text-primary">Пароль изменён</h1>
        <p className="text-sm text-muted-foreground">
          Сейчас перебросим на страницу входа…
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-serif text-2xl text-primary">Новый пароль</h1>

      <form className="flex flex-col gap-3" onSubmit={handleSubmit}>
        {!tokenFromQuery && (
          <label className="flex flex-col gap-1 text-sm">
            <span className="text-foreground">Код из письма</span>
            <Input
              required
              value={secret}
              onChange={(e) => setSecret(e.target.value)}
            />
          </label>
        )}

        <label className="flex flex-col gap-1 text-sm">
          <span className="text-foreground">Новый пароль</span>
          <Input
            type="password"
            autoComplete="new-password"
            required
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <span className="text-xs text-muted-foreground">Не менее 8 символов</span>
        </label>

        {error && <p className="text-sm text-destructive">{error}</p>}

        <Button type="submit" disabled={isSubmitting || !secret}>
          {isSubmitting ? 'Сохраняем…' : 'Сменить пароль'}
        </Button>
      </form>

      <Link href="/auth/login" className="text-sm text-muted-foreground hover:text-primary">
        ← Вернуться ко входу
      </Link>
    </div>
  )
}
