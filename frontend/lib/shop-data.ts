import {
  ApiError,
  fetchCategoriesTree,
  fetchProduct,
  fetchProductsList,
  type ApiCategoryNode,
  type ApiProduct,
  type CategoryRef,
  type ProductAttributeValue,
  type ProductImage,
} from './api'

export type CategoryNode = {
  id: number
  slug: string
  name: string
  sortOrder: number
  children: CategoryNode[]
}

export type Product = {
  id: number
  slug: string
  name: string
  attributes: ProductAttributeValue[]
  images: ProductImage[]
  category: CategoryRef | null
}

function toUiProduct(p: ApiProduct): Product {
  return {
    id: p.id,
    slug: p.slug,
    name: p.name,
    attributes: p.attributes,
    images: p.images,
    category: p.category,
  }
}

// primaryImageUrl возвращает URL основной картинки (или первой по сортировке),
// либо локальный placeholder, если у товара нет изображений.
export function primaryImageUrl(p: Product): string {
  if (p.images.length === 0) return '/placeholder.svg'
  const primary = p.images.find((i) => i.isPrimary)
  return (primary ?? p.images[0]).url
}

export async function getAllProducts(categorySlug?: string): Promise<Product[]> {
  const data = await fetchProductsList(categorySlug ? { categorySlug } : {})
  return data.items.map(toUiProduct)
}

// getProduct возвращает null только на честный 404; на любую другую ошибку —
// проброс, чтобы не маскировать сбои бэка.
export async function getProduct(slug: string): Promise<Product | null> {
  try {
    const p = await fetchProduct(slug)
    return toUiProduct(p)
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) return null
    throw e
  }
}

export async function getCategoriesTree(): Promise<CategoryNode[]> {
  const data = await fetchCategoriesTree()
  return data.items.map(toUiNode)
}

function toUiNode(n: ApiCategoryNode): CategoryNode {
  return {
    id: n.id,
    slug: n.slug,
    name: n.name,
    sortOrder: n.sortOrder,
    children: (n.children ?? []).map(toUiNode),
  }
}

export function findCategoryInTree(
  tree: CategoryNode[],
  slug: string,
): CategoryNode | null {
  for (const node of tree) {
    if (node.slug === slug) return node
    const child = findCategoryInTree(node.children, slug)
    if (child) return child
  }
  return null
}
