'use client'

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react'

import type { Product } from '@/lib/shop-data'
import * as cartApi from '@/lib/cart-client'
import { onAuthChanged, useAuth } from '@/components/auth-context'

// Внешний API остаётся тем же, что был у локального reducer'а: items с
// product + quantity, и три метода мутации. Внутри теперь — серверная
// корзина, мутации идут через apiFetch и заменяют локальный snapshot
// тем, что вернул бэк. Никакого optimistic-state: один источник правды.

export type CartItem = {
  product: Product
  quantity: number
}

type CartContextValue = {
  items: CartItem[]
  totalItems: number
  isLoading: boolean
  addItem: (product: Product, quantity?: number) => Promise<void>
  removeItem: (productId: number) => Promise<void>
  setQuantity: (productId: number, quantity: number) => Promise<void>
  clear: () => Promise<void>
}

const CartContext = createContext<CartContextValue | null>(null)

export function CartProvider({ children }: { children: ReactNode }) {
  const { status } = useAuth()
  const [items, setItems] = useState<CartItem[]>([])
  const [totalItems, setTotalItems] = useState(0)
  const [isLoading, setIsLoading] = useState(true)

  // Чтобы быстрые повторные клики на + / – не устроили гонку: мутации
  // сериализуем — каждая ждёт предыдущую. Так UI не «дёргается» из-за
  // ответа, пришедшего после следующего клика.
  const pending = useRef<Promise<unknown>>(Promise.resolve())

  const applySnapshot = useCallback((r: cartApi.CartResponse) => {
    setItems(toCartItems(r))
    setTotalItems(r.totalQuantity)
  }, [])

  const fetchSnapshot = useCallback(async () => {
    try {
      const r = await cartApi.getCart()
      applySnapshot(r)
    } catch {
      // Сетевые сбои не должны валить страницу — оставим как пустую
      // корзину, фронт сам ретрайнет на следующем действии.
      setItems([])
      setTotalItems(0)
    } finally {
      setIsLoading(false)
    }
  }, [applySnapshot])

  // Стартуем getCart только когда auth-bootstrap завершился: тогда
  // первый же запрос пойдёт с правильным Bearer (или без — для гостя),
  // вместо «сначала как гость, потом как юзер».
  //
  // onAuthChanged продолжаем слушать — это нужно для дополнительной
  // перезагрузки корзины при login/register/logout уже в рантайме.
  useEffect(() => {
    if (status === 'loading') return
    void fetchSnapshot()
    const off = onAuthChanged(() => {
      void fetchSnapshot()
    })
    return off
  }, [status, fetchSnapshot])

  const runMutation = useCallback(
    async (fn: () => Promise<cartApi.CartResponse>) => {
      const next = pending.current.then(fn, fn)
      pending.current = next.catch(() => undefined)
      const r = await next
      applySnapshot(r)
    },
    [applySnapshot],
  )

  const addItem = useCallback(
    (product: Product, quantity = 1) =>
      runMutation(() => cartApi.addItem(product.id, quantity)),
    [runMutation],
  )

  const removeItem = useCallback(
    (productId: number) => runMutation(() => cartApi.removeItem(productId)),
    [runMutation],
  )

  const setQuantity = useCallback(
    (productId: number, quantity: number) =>
      runMutation(() => cartApi.setQuantity(productId, quantity)),
    [runMutation],
  )

  const clear = useCallback(
    () => runMutation(() => cartApi.clearCart()),
    [runMutation],
  )

  const value = useMemo<CartContextValue>(
    () => ({ items, totalItems, isLoading, addItem, removeItem, setQuantity, clear }),
    [items, totalItems, isLoading, addItem, removeItem, setQuantity, clear],
  )

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>
}

export function useCart() {
  const ctx = useContext(CartContext)
  if (!ctx) throw new Error('useCart must be used within CartProvider')
  return ctx
}

// toCartItems отбрасывает позиции, у которых бэк не прислал product
// (битая запись каталога). Альтернатива — показать «товар недоступен», но
// корзина без названия товара бесполезна для пользователя; для MVP скрываем.
function toCartItems(r: cartApi.CartResponse): CartItem[] {
  const out: CartItem[] = []
  for (const it of r.items) {
    if (!it.product) continue
    out.push({
      product: {
        id: it.product.id,
        slug: it.product.slug,
        name: it.product.name,
        attributes: it.product.attributes,
        images: it.product.images,
        category: it.product.category,
      },
      quantity: it.quantity,
    })
  }
  return out
}
