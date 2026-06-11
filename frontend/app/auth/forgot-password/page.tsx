'use client'

import Link from 'next/link'
import { useState, type FormEvent } from 'react'

import { requestPasswordResetRequest } from '@/lib/auth-client'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitted, setSubmitted] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      // Бэк ВСЕГДА отдаёт 204, даже если email не зарегистрирован — это
      // намеренная защита от энумерации. Поэтому показываем одинаковое
      // сообщение в любом исходе и не пытаемся различать "нашли/не нашли".
      await requestPasswordResetRequest(email.trim())
    } finally {
      setIsSubmitting(false)
      setSubmitted(true)
    }
  }

  if (submitted) {
    return (
      <div className="flex flex-col gap-4">
        <h1 className="font-serif text-2xl text-primary">Письмо отправлено</h1>
        <p className="text-sm text-muted-foreground">
          Если такой email зарегистрирован, в течение нескольких минут на него
          придёт письмо с ссылкой для сброса пароля. Проверьте папку «Спам».
        </p>
        <p className="text-xs text-muted-foreground">
          Email-рассылка ещё не подключена в этом окружении. Токен сейчас
          сохраняется только в БД — используйте админ-интерфейс или CLI, чтобы
          получить его.
        </p>
        <Link href="/auth/login" className="text-sm text-primary hover:underline">
          ← Вернуться ко входу
        </Link>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-serif text-2xl text-primary">Восстановление пароля</h1>
      <p className="text-sm text-muted-foreground">
        Укажите email, на который зарегистрирован аккаунт. Мы отправим ссылку
        для сброса пароля.
      </p>

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

        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Отправляем…' : 'Отправить'}
        </Button>
      </form>

      <Link href="/auth/login" className="text-sm text-muted-foreground hover:text-primary">
        ← Вернуться ко входу
      </Link>
    </div>
  )
}
