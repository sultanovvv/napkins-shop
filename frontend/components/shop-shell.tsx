import type { ReactNode } from 'react'
import { SiteHeader } from '@/components/site-header'
import { SiteFooter } from '@/components/site-footer'
import { CategoryNav } from '@/components/category-nav'
import { Sidebar } from '@/components/sidebar'
import { getCategoriesTree } from '@/lib/shop-data'

export async function ShopShell({
  children,
  activeCategorySlug,
}: {
  children: ReactNode
  activeCategorySlug?: string
}) {
  const tree = await getCategoriesTree()

  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader
        mobileCategoryNav={
          <CategoryNav tree={tree} activeSlug={activeCategorySlug} />
        }
      />
      <div className="mx-auto grid w-full max-w-7xl flex-1 gap-6 px-4 py-6 lg:grid-cols-[260px_1fr_280px]">
        <div className="hidden lg:block">
          <div className="rounded-lg border border-border bg-card">
            <h2 className="border-b border-border px-4 py-3 font-serif text-lg text-primary">
              Категории
            </h2>
            <div className="p-2">
              <CategoryNav tree={tree} activeSlug={activeCategorySlug} />
            </div>
          </div>
        </div>

        <main className="min-w-0">{children}</main>

        <div className="hidden lg:block">
          <Sidebar />
        </div>
      </div>
      <SiteFooter />
    </div>
  )
}
