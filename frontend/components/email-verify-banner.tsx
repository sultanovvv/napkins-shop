'use client'

import { useState } from 'react'
import { MailWarning } from 'lucide-react'

import { useAuth } from '@/components/auth-context'
import { requestEmailVerificationRequest } from '@/lib/auth-client'

// Тонкая полоска сверху для залогиненных пользователей с неподтверждённым
// email. Кнопка «Отправить ещё раз» дёргает бэк; SMTP пока нет, поэтому в
// UX-копии это явно «… когда подключим почту». После confirm баннер сам
// исчезает, потому что user.emailVerifiedAt станет не null.
export function EmailVerifyBanner() {
  const { user, status } = useAuth()
  const [state, setState] = useState<'idle' | 'sending' | 'sent'>('idle')

  if (status !== 'authenticated' || !user) return null
  if (user.emailVerifiedAt) return null

  async function handleResend() {
    setState('sending')
    try {
      await requestEmailVerificationRequest()
    } finally {
      setState('sent')
    }
  }

  return (
    <div className="border-b border-amber-300 bg-amber-50 text-amber-900">
      <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-3 px-4 py-2 text-sm">
        <div className="flex items-center gap-2">
          <MailWarning className="size-4" />
          <span>
            Подтвердите email <strong>{user.email}</strong>, чтобы оформлять заказы.
          </span>
        </div>
        {state === 'sent' ? (
          <span className="text-xs text-amber-700">
            Ссылка выпущена. Отправка письма пока не подключена — попросите
            администратора достать токен из БД.
          </span>
        ) : (
          <button
            type="button"
            onClick={handleResend}
            disabled={state === 'sending'}
            className="rounded-md border border-amber-400 px-3 py-1 text-xs text-amber-900 hover:bg-amber-100 disabled:opacity-50"
          >
            {state === 'sending' ? 'Отправляем…' : 'Отправить ещё раз'}
          </button>
        )}
      </div>
    </div>
  )
}
