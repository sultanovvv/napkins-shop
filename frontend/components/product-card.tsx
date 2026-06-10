'use client'

import Image from 'next/image'
import Link from 'next/link'
import { ShoppingBasket } from 'lucide-react'
import { primaryImageUrl, type Product } from '@/lib/shop-data'
import { useCart } from '@/components/cart-context'

export function ProductCard({ product }: { product: Product }) {
  const { addItem } = useCart()

  return (
    <article className="group flex flex-col overflow-hidden rounded-lg border border-border bg-card transition-shadow hover:shadow-md">
      <Link
        href={`/product/${product.slug}`}
        className="relative block aspect-square overflow-hidden bg-muted"
      >
        <Image
          src={primaryImageUrl(product)}
          alt={product.name}
          fill
          sizes="(max-width: 768px) 50vw, 25vw"
          className="object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </Link>
      <div className="flex flex-1 flex-col gap-2 p-3">
        <Link
          href={`/product/${product.slug}`}
          className="font-serif text-base leading-snug text-foreground hover:text-primary"
        >
          {product.name}
        </Link>
        {product.attributes.length > 0 && (
          <p className="line-clamp-2 text-xs text-muted-foreground">
            {product.attributes
              .map((a) => (a.valueInt !== null ? a.valueInt : a.valueText))
              .filter(Boolean)
              .join(' · ')}
          </p>
        )}
        <div className="mt-auto flex items-center justify-end pt-2">
          <button
            type="button"
            onClick={() => addItem(product)}
            className="flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            <ShoppingBasket className="size-4" />В корзину
          </button>
        </div>
      </div>
    </article>
  )
}
