'use client'

import { useState } from 'react'
import { Minus, Plus, ShoppingBasket, Check } from 'lucide-react'
import type { Product } from '@/lib/shop-data'
import { useCart } from '@/components/cart-context'

export function AddToCart({ product }: { product: Product }) {
  const { addItem } = useCart()
  const [qty, setQty] = useState(1)
  const [added, setAdded] = useState(false)

  function handleAdd() {
    addItem(product, qty)
    setAdded(true)
    setTimeout(() => setAdded(false), 1800)
  }

  return (
    <div className="flex flex-wrap items-center gap-3">
      <div className="flex items-center rounded-md border border-input">
        <button
          type="button"
          aria-label="Уменьшить"
          onClick={() => setQty((q) => Math.max(1, q - 1))}
          className="p-2.5 text-foreground hover:text-primary"
        >
          <Minus className="size-4" />
        </button>
        <span className="min-w-8 text-center font-medium">{qty}</span>
        <button
          type="button"
          aria-label="Увеличить"
          onClick={() => setQty((q) => q + 1)}
          className="p-2.5 text-foreground hover:text-primary"
        >
          <Plus className="size-4" />
        </button>
      </div>
      <button
        type="button"
        onClick={handleAdd}
        className="flex items-center gap-2 rounded-md bg-primary px-6 py-2.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
      >
        {added ? (
          <>
            <Check className="size-4" /> Добавлено
          </>
        ) : (
          <>
            <ShoppingBasket className="size-4" /> В корзину
          </>
        )}
      </button>
    </div>
  )
}
