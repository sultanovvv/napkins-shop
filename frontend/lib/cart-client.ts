import { apiFetch } from './auth-client'
import type { ApiProduct } from './api'

export type CartItemApi = {
  productId: number
  quantity: number
  product: ApiProduct | null
}

export type CartResponse = {
  items: CartItemApi[]
  totalQuantity: number
}

const base = '/api/v1/public/cart'

export function getCart(): Promise<CartResponse> {
  return apiFetch<CartResponse>(`${base}/getCart`)
}

export function addItem(productId: number, quantity: number): Promise<CartResponse> {
  return apiFetch<CartResponse>(`${base}/addItem`, {
    method: 'POST',
    body: JSON.stringify({ productId, quantity }),
  })
}

export function setQuantity(productId: number, quantity: number): Promise<CartResponse> {
  return apiFetch<CartResponse>(`${base}/setQuantity`, {
    method: 'POST',
    body: JSON.stringify({ productId, quantity }),
  })
}

export function removeItem(productId: number): Promise<CartResponse> {
  return apiFetch<CartResponse>(`${base}/removeItem`, {
    method: 'POST',
    body: JSON.stringify({ productId }),
  })
}

export function clearCart(): Promise<CartResponse> {
  return apiFetch<CartResponse>(`${base}/clear`, {
    method: 'POST',
  })
}
