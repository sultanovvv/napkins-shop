'use client'

import Image from 'next/image'
import Link from 'next/link'
import { useState } from 'react'
import { Minus, Plus, Trash2, ShoppingBasket, CheckCircle2, AlertCircle } from 'lucide-react'

import { useCart } from '@/components/cart-context'
import { primaryImageUrl } from '@/lib/shop-data'
import { ApiError } from '@/lib/auth-client'
import { createOrder, type OrderResponse } from '@/lib/orders-client'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from '@/components/ui/sheet'

// Тонкая state-машина для футера корзины:
// idle — обычный вид с кнопкой «Оформить заказ».
// submitting — спиннер на кнопке, ввод заблокирован.
// success — показываем номер заказа и предлагаем закрыть.
// error — показываем сообщение, кнопка снова активна.
type Submission =
  | { kind: 'idle' }
  | { kind: 'submitting' }
  | { kind: 'success'; order: OrderResponse }
  | { kind: 'error'; message: string; emailNotVerified: boolean }

export function CartDrawer({
  open,
  onOpenChange,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const { items, totalItems, removeItem, setQuantity, clear } = useCart()
  const [submission, setSubmission] = useState<Submission>({ kind: 'idle' })

  // При закрытии drawer'а сбрасываем результат, чтобы при следующем открытии
  // не висело старое «спасибо за заказ».
  function handleOpenChange(next: boolean) {
    if (!next) setSubmission({ kind: 'idle' })
    onOpenChange(next)
  }

  async function handleCheckout() {
    if (items.length === 0) return
    setSubmission({ kind: 'submitting' })
    try {
      const order = await createOrder(
        items.map((i) => ({ productId: i.product.id, quantity: i.quantity })),
      )
      // Сервер заказ создал — локально чистим корзину. clear() сходит на
      // /cart/clear и обновит счётчик в шапке.
      await clear()
      setSubmission({ kind: 'success', order })
    } catch (err) {
      if (err instanceof ApiError && err.code === 'EMAIL_NOT_VERIFIED') {
        setSubmission({
          kind: 'error',
          message: 'Подтвердите email, прежде чем оформить заказ.',
          emailNotVerified: true,
        })
        return
      }
      if (err instanceof ApiError) {
        setSubmission({
          kind: 'error',
          message: err.message || 'Не удалось оформить заказ',
          emailNotVerified: false,
        })
        return
      }
      setSubmission({
        kind: 'error',
        message: 'Не удалось оформить заказ. Попробуйте ещё раз.',
        emailNotVerified: false,
      })
    }
  }

  const isSubmitting = submission.kind === 'submitting'

  return (
    <Sheet open={open} onOpenChange={handleOpenChange}>
      <SheetContent className="flex w-full flex-col gap-0 p-0 sm:max-w-md">
        <SheetHeader className="border-b border-border px-5 py-4">
          <SheetTitle className="font-serif text-xl text-primary">
            Корзина
          </SheetTitle>
        </SheetHeader>

        {submission.kind === 'success' ? (
          <SuccessState order={submission.order} onClose={() => handleOpenChange(false)} />
        ) : items.length === 0 ? (
          <EmptyState />
        ) : (
          <>
            <div className="flex-1 overflow-y-auto px-5 py-4">
              <ul className="flex flex-col gap-4">
                {items.map((item) => (
                  <li key={item.product.id} className="flex gap-3">
                    <div className="relative size-16 shrink-0 overflow-hidden rounded-md bg-muted">
                      <Image
                        src={primaryImageUrl(item.product)}
                        alt={item.product.name}
                        fill
                        sizes="64px"
                        className="object-cover"
                      />
                    </div>
                    <div className="flex flex-1 flex-col">
                      <span className="font-serif text-sm text-foreground">
                        {item.product.name}
                      </span>
                      <div className="mt-1 flex items-center gap-2">
                        <div className="flex items-center rounded-md border border-input">
                          <button
                            type="button"
                            aria-label="Уменьшить"
                            disabled={isSubmitting}
                            onClick={() =>
                              setQuantity(item.product.id, item.quantity - 1)
                            }
                            className="p-1.5 text-foreground hover:text-primary disabled:opacity-50"
                          >
                            <Minus className="size-3.5" />
                          </button>
                          <span className="min-w-6 text-center text-sm">
                            {item.quantity}
                          </span>
                          <button
                            type="button"
                            aria-label="Увеличить"
                            disabled={isSubmitting}
                            onClick={() =>
                              setQuantity(item.product.id, item.quantity + 1)
                            }
                            className="p-1.5 text-foreground hover:text-primary disabled:opacity-50"
                          >
                            <Plus className="size-3.5" />
                          </button>
                        </div>
                        <button
                          type="button"
                          aria-label="Удалить"
                          disabled={isSubmitting}
                          onClick={() => removeItem(item.product.id)}
                          className="ml-auto p-1.5 text-muted-foreground hover:text-destructive disabled:opacity-50"
                        >
                          <Trash2 className="size-4" />
                        </button>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>

            <SheetFooter className="border-t border-border px-5 py-4">
              <div className="flex w-full flex-col gap-3">
                <div className="flex items-center justify-between font-serif text-lg">
                  <span>Товаров:</span>
                  <span className="font-semibold text-primary">{totalItems}</span>
                </div>

                {submission.kind === 'error' && (
                  <ErrorState
                    message={submission.message}
                    emailNotVerified={submission.emailNotVerified}
                    onCloseDrawer={() => handleOpenChange(false)}
                  />
                )}

                <button
                  type="button"
                  onClick={handleCheckout}
                  disabled={isSubmitting}
                  className="w-full rounded-md bg-primary px-4 py-3 text-center text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {isSubmitting ? 'Оформляем…' : 'Оформить заказ'}
                </button>
                <button
                  type="button"
                  disabled={isSubmitting}
                  onClick={() => {
                    void clear()
                  }}
                  className="text-xs text-muted-foreground hover:text-destructive disabled:opacity-50"
                >
                  Очистить корзину
                </button>
              </div>
            </SheetFooter>
          </>
        )}
      </SheetContent>
    </Sheet>
  )
}

function EmptyState() {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center">
      <ShoppingBasket className="size-12 text-muted-foreground/50" />
      <p className="text-muted-foreground">Ваша корзина пуста</p>
    </div>
  )
}

function SuccessState({
  order,
  onClose,
}: {
  order: OrderResponse
  onClose: () => void
}) {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-4 p-8 text-center">
      <CheckCircle2 className="size-12 text-emerald-500" />
      <div>
        <p className="font-serif text-lg text-foreground">Заказ оформлен</p>
        <p className="mt-1 text-sm text-muted-foreground">
          Номер заказа: <span className="font-medium text-primary">#{order.id}</span>
        </p>
        <p className="mt-1 text-xs text-muted-foreground">
          Статус: {order.status}
        </p>
      </div>
      <button
        type="button"
        onClick={onClose}
        className="rounded-md bg-primary px-5 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        Закрыть
      </button>
    </div>
  )
}

function ErrorState({
  message,
  emailNotVerified,
  onCloseDrawer,
}: {
  message: string
  emailNotVerified: boolean
  onCloseDrawer: () => void
}) {
  return (
    <div className="flex items-start gap-2 rounded-md border border-destructive/40 bg-destructive/5 px-3 py-2 text-sm text-destructive">
      <AlertCircle className="mt-0.5 size-4 shrink-0" />
      <div className="flex flex-col gap-1">
        <span>{message}</span>
        {emailNotVerified && (
          <Link
            href="/auth/verify-email"
            onClick={onCloseDrawer}
            className="text-xs underline hover:no-underline"
          >
            Перейти к подтверждению email →
          </Link>
        )}
      </div>
    </div>
  )
}
