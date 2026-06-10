import Image from 'next/image'
import Link from 'next/link'
import { ShopShell } from '@/components/shop-shell'
import { ProductCard } from '@/components/product-card'
import { getAllProducts, getCategoriesTree } from '@/lib/shop-data'

export const dynamic = 'force-dynamic'

export default async function HomePage() {
  const [all, tree] = await Promise.all([getAllProducts(), getCategoriesTree()])

  const featured = all.slice(0, 8)
  // 6 листьев из первых корней — для блока "Популярные категории".
  const topCategories = tree.flatMap((parent) => parent.children).slice(0, 6)
  const heroCategory = topCategories[0]?.slug ?? tree[0]?.slug ?? 'cvety'

  return (
    <ShopShell>
      <section className="overflow-hidden rounded-lg border border-border bg-card">
        <div className="relative aspect-[16/7] w-full">
          <Image
            src="/hero-decoupage.png"
            alt="Салфетки и материалы для декупажа"
            fill
            priority
            sizes="(max-width: 1024px) 100vw, 60vw"
            className="object-cover"
          />
          <div className="absolute inset-0 flex flex-col justify-center bg-gradient-to-r from-card/85 via-card/40 to-transparent p-6 md:p-10">
            <h1 className="max-w-md text-balance font-serif text-3xl italic text-primary md:text-4xl">
              Салфетки для декупажа на любой сюжет
            </h1>
            <p className="mt-3 max-w-sm text-pretty text-sm text-foreground/80 md:text-base">
              Цветы, птицы, гномы, города и праздники — тысячи мотивов для вашего
              творчества и коллекции.
            </p>
            <Link
              href={`/category/${heroCategory}`}
              className="mt-5 w-fit rounded-md bg-primary px-5 py-2.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
            >
              Перейти в каталог
            </Link>
          </div>
        </div>
      </section>

      {topCategories.length > 0 && (
        <section className="mt-8">
          <h2 className="font-serif text-2xl text-primary">Популярные категории</h2>
          <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3">
            {topCategories.map((cat) => (
              <Link
                key={cat.slug}
                href={`/category/${cat.slug}`}
                className="flex items-center justify-between rounded-lg border border-border bg-card px-4 py-3 transition-colors hover:border-primary hover:bg-secondary"
              >
                <span className="font-serif text-sm leading-snug text-foreground">
                  {cat.name}
                </span>
              </Link>
            ))}
          </div>
        </section>
      )}

      <section className="mt-8">
        <div className="flex items-center justify-between">
          <h2 className="font-serif text-2xl text-primary">Новинки и хиты</h2>
        </div>
        <div className="mt-4 grid grid-cols-2 gap-4 md:grid-cols-3 xl:grid-cols-4">
          {featured.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>
    </ShopShell>
  )
}
