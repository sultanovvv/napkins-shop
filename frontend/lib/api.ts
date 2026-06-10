const API_BASE = process.env.API_BASE_URL ?? 'http://127.0.0.1:9888'

export type AttributeValueType = 'string' | 'int'

export type Attribute = {
  id: number
  slug: string
  name: string
  valueType: AttributeValueType
  sortOrder: number
}

export type ProductAttributeValue = {
  attribute: Attribute
  valueText: string | null
  valueInt: number | null
}

export type CategoryRef = {
  id: number
  slug: string
  name: string
}

export type ProductImage = {
  url: string
  isPrimary: boolean
  sortOrder: number
}

export type ApiProduct = {
  id: number
  slug: string
  name: string
  attributes: ProductAttributeValue[]
  images: ProductImage[]
  category: CategoryRef | null
}

export type ApiProductListResponse = {
  items: ApiProduct[]
  totalCount: number
}

export type ApiCategoryNode = {
  id: number
  slug: string
  name: string
  sortOrder: number
  children: ApiCategoryNode[]
}

export type ApiCategoriesTreeResponse = {
  items: ApiCategoryNode[]
}

async function get<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { cache: 'no-store', ...init })
  if (!res.ok) {
    throw new ApiError(res.status, `${init?.method ?? 'GET'} ${path} → HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}

export class ApiError extends Error {
  constructor(public readonly status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function fetchProductsList(
  options: { categorySlug?: string } = {},
): Promise<ApiProductListResponse> {
  const qs = options.categorySlug
    ? `?category=${encodeURIComponent(options.categorySlug)}`
    : ''
  return get<ApiProductListResponse>(`/getProductsList${qs}`)
}

export async function fetchProduct(slug: string): Promise<ApiProduct> {
  return get<ApiProduct>(`/getProduct/${encodeURIComponent(slug)}`)
}

export async function fetchCategoriesTree(): Promise<ApiCategoriesTreeResponse> {
  return get<ApiCategoriesTreeResponse>('/getCategoriesTree')
}
