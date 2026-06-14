import { apiFetch } from './auth-client'

export type OrderItemResponse = {
  productId: number
  name: string
  quantity: number
}

export type OrderResponse = {
  id: number
  items: OrderItemResponse[]
  status: string
  createdAt: string
}

export type CreateOrderItemInput = {
  productId: number
  quantity: number
}

export function createOrder(items: CreateOrderItemInput[]): Promise<OrderResponse> {
  return apiFetch<OrderResponse>('/api/v1/public/createOrder', {
    method: 'POST',
    body: JSON.stringify({ items }),
  })
}
