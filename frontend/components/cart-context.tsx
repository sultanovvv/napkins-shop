'use client'

import {
  createContext,
  useContext,
  useMemo,
  useReducer,
  type ReactNode,
} from 'react'
import type { Product } from '@/lib/shop-data'

export type CartItem = {
  product: Product
  quantity: number
}

type CartState = {
  items: CartItem[]
}

type Action =
  | { type: 'add'; product: Product; quantity?: number }
  | { type: 'remove'; id: number }
  | { type: 'setQty'; id: number; quantity: number }
  | { type: 'clear' }

function reducer(state: CartState, action: Action): CartState {
  switch (action.type) {
    case 'add': {
      const qty = action.quantity ?? 1
      const existing = state.items.find((i) => i.product.id === action.product.id)
      if (existing) {
        return {
          items: state.items.map((i) =>
            i.product.id === action.product.id
              ? { ...i, quantity: i.quantity + qty }
              : i,
          ),
        }
      }
      return { items: [...state.items, { product: action.product, quantity: qty }] }
    }
    case 'remove':
      return { items: state.items.filter((i) => i.product.id !== action.id) }
    case 'setQty':
      return {
        items: state.items
          .map((i) =>
            i.product.id === action.id
              ? { ...i, quantity: Math.max(1, action.quantity) }
              : i,
          )
          .filter((i) => i.quantity > 0),
      }
    case 'clear':
      return { items: [] }
    default:
      return state
  }
}

type CartContextValue = {
  items: CartItem[]
  addItem: (product: Product, quantity?: number) => void
  removeItem: (id: number) => void
  setQuantity: (id: number, quantity: number) => void
  clear: () => void
  totalItems: number
}

const CartContext = createContext<CartContextValue | null>(null)

export function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(reducer, { items: [] })

  const value = useMemo<CartContextValue>(() => {
    const totalItems = state.items.reduce((sum, i) => sum + i.quantity, 0)
    return {
      items: state.items,
      addItem: (product, quantity) => dispatch({ type: 'add', product, quantity }),
      removeItem: (id) => dispatch({ type: 'remove', id }),
      setQuantity: (id, quantity) => dispatch({ type: 'setQty', id, quantity }),
      clear: () => dispatch({ type: 'clear' }),
      totalItems,
    }
  }, [state.items])

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>
}

export function useCart() {
  const ctx = useContext(CartContext)
  if (!ctx) throw new Error('useCart must be used within CartProvider')
  return ctx
}
