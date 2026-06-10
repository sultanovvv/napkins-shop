import { notFound } from 'next/navigation'
import Link from 'next/link'
import { Home } from 'lucide-react'
import { ShopShell } from '@/components/shop-shell'
import { ProductCard } from '@/components/product-card'
import {
  findCategoryInTree,
  getAllProducts,
  getCategoriesTree,
} from '@/lib/shop-data'

export const dynamic = 'force-dynamic'

export default async function CategoryPage({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params

  const [tree, items] = await Promise.all([
    getCategoriesTree(),
    getAllProducts(slug),
  ])

  const category = findCategoryInTree(tree, slug)
  if (!category) notFound()

  // Хлебные крошки: попробуем найти родителя в дереве.
  const parent = tree.find((p) =>
    p.children.some((c) => c.slug === category.slug),
  )

  return (
    <ShopShell activeCategorySlug={slug}>
      <nav className="mb-4 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
        <Link href="/" className="hover:text-primary" aria-label="Главная">
          <Home className="size-4" />
        </Link>
        {parent && (
          <>
            <span>›</span>
            <Link
              href={`/category/${parent.slug}`}
              className="hover:text-primary"
            >
              {parent.name}
            </Link>
          </>
        )}
        <span>›</span>
        <span className="text-foreground">{category.name}</span>
      </nav>

      <h1 className="font-serif text-2xl text-foreground md:text-3xl">
        {category.name}
      </h1>
      <div className="mt-3 rounded-md bg-secondary px-4 py-2 text-sm text-secondary-foreground">
        Найдено товаров: {items.length}
      </div>

      {category.children.length > 0 && (
        <div className="mt-4 flex flex-wrap gap-2">
          {category.children.map((child) => (
            <Link
              key={child.slug}
              href={`/category/${child.slug}`}
              className="rounded-full border border-border bg-card px-3 py-1 text-xs text-foreground transition-colors hover:border-primary hover:text-primary"
            >
              {child.name}
            </Link>
          ))}
        </div>
      )}

      {items.length === 0 ? (
        <p className="mt-8 rounded-md border border-dashed border-border px-4 py-8 text-center text-sm text-muted-foreground">
          В этой категории пока нет товаров.
        </p>
      ) : (
        <div className="mt-6 grid grid-cols-2 gap-4 md:grid-cols-3 xl:grid-cols-4">
          {items.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </ShopShell>
  )
}
