'use client'

import Image from 'next/image'
import { Minus, Plus, Trash2, ShoppingBasket } from 'lucide-react'
import { useCart } from '@/components/cart-context'
import { primaryImageUrl } from '@/lib/shop-data'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from '@/components/ui/sheet'

export function CartDrawer({
  open,
  onOpenChange,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const { items, totalItems, removeItem, setQuantity, clear } = useCart()

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="flex w-full flex-col gap-0 p-0 sm:max-w-md">
        <SheetHeader className="border-b border-border px-5 py-4">
          <SheetTitle className="font-serif text-xl text-primary">
            Корзина
          </SheetTitle>
        </SheetHeader>

        {items.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center">
            <ShoppingBasket className="size-12 text-muted-foreground/50" />
            <p className="text-muted-foreground">Ваша корзина пуста</p>
          </div>
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
                            onClick={() =>
                              setQuantity(item.product.id, item.quantity - 1)
                            }
                            className="p-1.5 text-foreground hover:text-primary"
                          >
                            <Minus className="size-3.5" />
                          </button>
                          <span className="min-w-6 text-center text-sm">
                            {item.quantity}
                          </span>
                          <button
                            type="button"
                            aria-label="Увеличить"
                            onClick={() =>
                              setQuantity(item.product.id, item.quantity + 1)
                            }
                            className="p-1.5 text-foreground hover:text-primary"
                          >
                            <Plus className="size-3.5" />
                          </button>
                        </div>
                        <button
                          type="button"
                          aria-label="Удалить"
                          onClick={() => removeItem(item.product.id)}
                          className="ml-auto p-1.5 text-muted-foreground hover:text-destructive"
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
                <button
                  type="button"
                  className="w-full rounded-md bg-primary px-4 py-3 text-center text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
                  onClick={() => alert('Оформление заказа — демо-версия')}
                >
                  Оформить заказ
                </button>
                <button
                  type="button"
                  onClick={clear}
                  className="text-xs text-muted-foreground hover:text-destructive"
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
