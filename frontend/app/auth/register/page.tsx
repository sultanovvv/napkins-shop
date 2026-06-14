'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState, type FormEvent } from 'react'

import { useAuth } from '@/components/auth-context'
import { ApiError } from '@/lib/auth-client'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export default function RegisterPage() {
  const router = useRouter()
  const { register } = useAuth()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setIsSubmitting(true)
    try {
      await register(email.trim(), password)
      router.push('/')
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setError('Пользователь с таким email уже зарегистрирован')
      } else if (err instanceof ApiError && err.status === 400) {
        setError(err.message)
      } else {
        setError('Не удалось зарегистрироваться. Попробуйте ещё раз.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-serif text-2xl text-primary">Регистрация</h1>

      <form className="flex flex-col gap-3" onSubmit={handleSubmit}>
        <label className="flex flex-col gap-1 text-sm">
          <span className="text-foreground">Email</span>
          <Input
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>

        <label className="flex flex-col gap-1 text-sm">
          <span className="text-foreground">Пароль</span>
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

        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Создаём аккаунт…' : 'Создать аккаунт'}
        </Button>
      </form>

      <div className="text-center text-sm">
        Уже есть аккаунт?{' '}
        <Link href="/auth/login" className="text-primary hover:underline">
          Войти
        </Link>
      </div>
    </div>
  )
}
