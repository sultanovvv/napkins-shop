'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ChevronRight } from 'lucide-react'
import type { CategoryNode } from '@/lib/shop-data'
import { cn } from '@/lib/utils'

export function CategoryNav({
  tree,
  activeSlug,
}: {
  tree: CategoryNode[]
  activeSlug?: string
}) {
  if (tree.length === 0) {
    return (
      <p className="px-3 py-2 text-xs text-muted-foreground">
        Категории не загружены
      </p>
    )
  }

  return (
    <nav aria-label="Категории товаров" className="text-sm">
      <ul className="flex flex-col gap-1">
        {tree.map((parent) => (
          <CategoryItem key={parent.slug} parent={parent} activeSlug={activeSlug} />
        ))}
      </ul>
    </nav>
  )
}

function CategoryItem({
  parent,
  activeSlug,
}: {
  parent: CategoryNode
  activeSlug?: string
}) {
  const containsActive =
    activeSlug === parent.slug ||
    parent.children.some((c) => c.slug === activeSlug)
  const [open, setOpen] = useState(containsActive)
  const hasChildren = parent.children.length > 0

  return (
    <li className="border-b border-border/60 last:border-0">
      <div
        className={cn(
          'group flex items-center gap-1 transition-colors hover:bg-secondary',
          activeSlug === parent.slug && 'bg-secondary text-primary',
        )}
      >
        {hasChildren ? (
          <button
            type="button"
            onClick={() => setOpen((v) => !v)}
            aria-expanded={open}
            aria-label={open ? 'Свернуть подкатегории' : 'Развернуть подкатегории'}
            className="flex shrink-0 items-center justify-center self-stretch px-2 text-primary"
          >
            <ChevronRight
              className={cn(
                'size-3.5 transition-transform duration-150',
                open && 'rotate-90',
              )}
            />
          </button>
        ) : (
          <span className="w-7 shrink-0" aria-hidden />
        )}

        <Link
          href={`/category/${parent.slug}`}
          className="flex-1 py-2.5 pr-3 font-medium leading-snug"
        >
          {parent.name}
        </Link>
      </div>

      {hasChildren && open && (
        <ul className="ml-3 mb-1 flex flex-col">
          {parent.children.map((child) => (
            <li key={child.slug}>
              <Link
                href={`/category/${child.slug}`}
                className={cn(
                  'block rounded-md px-3 py-1.5 text-xs text-foreground/85 transition-colors hover:bg-secondary hover:text-primary',
                  activeSlug === child.slug && 'bg-secondary font-medium text-primary',
                )}
              >
                {child.name}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </li>
  )
}
