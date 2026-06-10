import { notFound } from 'next/navigation'
import Image from 'next/image'
import Link from 'next/link'
import { Home } from 'lucide-react'
import { ShopShell } from '@/components/shop-shell'
import { ProductCard } from '@/components/product-card'
import { AddToCart } from '@/components/add-to-cart'
import { getAllProducts, getProduct, primaryImageUrl } from '@/lib/shop-data'

export const dynamic = 'force-dynamic'

export default async function ProductPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const product = await getProduct(id)
  if (!product) notFound()

  const all = await getAllProducts()
  const related = all.filter((p) => p.id !== product.id).slice(0, 4)

  return (
    <ShopShell activeCategorySlug={product.category?.slug}>
      <nav className="mb-4 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
        <Link href="/" className="hover:text-primary" aria-label="Главная">
          <Home className="size-4" />
        </Link>
        {product.category && (
          <>
            <span>›</span>
            <Link
              href={`/category/${product.category.slug}`}
              className="hover:text-primary"
            >
              {product.category.name}
            </Link>
          </>
        )}
        <span>›</span>
        <span className="text-foreground">{product.name}</span>
      </nav>

      <div className="grid gap-6 md:grid-cols-2">
        <div className="relative aspect-square overflow-hidden rounded-lg border border-border bg-muted">
          <Image
            src={primaryImageUrl(product)}
            alt={product.name}
            fill
            priority
            sizes="(max-width: 768px) 100vw, 40vw"
            className="object-cover"
          />
        </div>
        <div className="flex flex-col gap-4">
          <h1 className="font-serif text-3xl text-foreground">{product.name}</h1>
          {product.attributes.length > 0 && (
            <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 text-sm">
              {product.attributes.map((v) => (
                <div key={v.attribute.slug} className="contents">
                  <dt className="text-muted-foreground">{v.attribute.name}:</dt>
                  <dd className="text-foreground">
                    {v.valueInt !== null ? v.valueInt : v.valueText}
                  </dd>
                </div>
              ))}
            </dl>
          )}
          <AddToCart product={product} />
        </div>
      </div>

      {related.length > 0 && (
        <section className="mt-10">
          <h2 className="font-serif text-2xl text-primary">Другие товары</h2>
          <div className="mt-4 grid grid-cols-2 gap-4 md:grid-cols-4">
            {related.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        </section>
      )}
    </ShopShell>
  )
}
