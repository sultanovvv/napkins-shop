import Link from 'next/link'
import Image from 'next/image'
import { getAllProducts, primaryImageUrl } from '@/lib/shop-data'

export async function Sidebar() {
  const all = await getAllProducts()
  const popular = all.slice(0, 4)

  if (popular.length === 0) return null

  return (
    <aside className="flex flex-col gap-6">
      <section className="rounded-lg border border-border bg-card p-4">
        <h2 className="mb-3 font-serif text-lg text-primary">Популярные товары</h2>
        <ul className="flex flex-col gap-3">
          {popular.map((p, i) => (
            <li key={p.id}>
              <Link
                href={`/product/${p.slug}`}
                className="flex items-center gap-3 rounded-md p-1.5 hover:bg-secondary"
              >
                <span className="flex size-6 shrink-0 items-center justify-center rounded bg-primary text-xs font-semibold text-primary-foreground">
                  {i + 1}
                </span>
                <div className="relative size-12 shrink-0 overflow-hidden rounded bg-muted">
                  <Image
                    src={primaryImageUrl(p)}
                    alt={p.name}
                    fill
                    sizes="48px"
                    className="object-cover"
                  />
                </div>
                <span className="font-serif text-sm leading-snug text-foreground">
                  {p.name}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </section>
    </aside>
  )
}
