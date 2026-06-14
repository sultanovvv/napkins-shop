'use client'

import Image from 'next/image'
import { useCallback, useEffect, useState } from 'react'
import { ChevronLeft, ChevronRight, X, ZoomIn } from 'lucide-react'

import type { ProductImage } from '@/lib/api'

// ProductGallery — главная картинка + ряд миниатюр + полноэкранный лайтбокс.
//
// Поведение:
// — Клик по главной картинке открывает лайтбокс на текущем индексе.
// — Клик по миниатюре переключает главную картинку (не открывает лайтбокс).
// — В лайтбоксе: «‹»/«›», стрелки клавиатуры, Esc и клик по фону — закрытие.
// — При открытом лайтбоксе скролл body заблокирован, чтобы фон не уезжал.
//
// Если изображений нет — fallback на /placeholder.svg, лайтбокс не открывается.
// Если изображение одно — лайтбокс работает, кнопки prev/next спрятаны.
export function ProductGallery({
  images,
  alt,
}: {
  images: ProductImage[]
  alt: string
}) {
  const [index, setIndex] = useState(0)
  const [lightboxOpen, setLightboxOpen] = useState(false)

  const hasImages = images.length > 0
  const safeIndex = Math.min(index, Math.max(0, images.length - 1))
  const current = hasImages ? images[safeIndex] : null

  const next = useCallback(() => {
    if (images.length <= 1) return
    setIndex((i) => (i + 1) % images.length)
  }, [images.length])

  const prev = useCallback(() => {
    if (images.length <= 1) return
    setIndex((i) => (i - 1 + images.length) % images.length)
  }, [images.length])

  const close = useCallback(() => setLightboxOpen(false), [])

  // Клавиатура + блокировка скролла активны только при открытом лайтбоксе.
  useEffect(() => {
    if (!lightboxOpen) return
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') close()
      else if (e.key === 'ArrowLeft') prev()
      else if (e.key === 'ArrowRight') next()
    }
    window.addEventListener('keydown', onKey)
    const prevOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      window.removeEventListener('keydown', onKey)
      document.body.style.overflow = prevOverflow
    }
  }, [lightboxOpen, next, prev, close])

  return (
    <div className="flex flex-col gap-3">
      {/* Главная картинка */}
      <button
        type="button"
        onClick={() => hasImages && setLightboxOpen(true)}
        disabled={!hasImages}
        aria-label={hasImages ? 'Открыть увеличенное изображение' : 'Нет изображений'}
        className="group relative aspect-square overflow-hidden rounded-lg border border-border bg-muted enabled:cursor-zoom-in"
      >
        <Image
          src={current?.url ?? '/placeholder.svg'}
          alt={alt}
          fill
          priority
          sizes="(max-width: 768px) 100vw, 40vw"
          className="object-cover"
        />
        {hasImages && (
          <span
            aria-hidden
            className="absolute right-3 top-3 flex items-center gap-1 rounded-md bg-black/55 px-2 py-1 text-xs text-white opacity-0 transition-opacity group-hover:opacity-100"
          >
            <ZoomIn className="size-3.5" />
            Увеличить
          </span>
        )}
      </button>

      {/* Миниатюры — показываем только при 2+ */}
      {images.length > 1 && (
        <ul className="flex flex-wrap gap-2">
          {images.map((img, i) => {
            const active = i === safeIndex
            return (
              <li key={img.url}>
                <button
                  type="button"
                  onClick={() => setIndex(i)}
                  aria-label={`Изображение ${i + 1}`}
                  aria-current={active}
                  className={
                    'relative size-16 overflow-hidden rounded-md border bg-muted transition-colors ' +
                    (active
                      ? 'border-primary ring-2 ring-primary/30'
                      : 'border-border hover:border-primary/60')
                  }
                >
                  <Image
                    src={img.url}
                    alt=""
                    fill
                    sizes="64px"
                    className="object-cover"
                  />
                </button>
              </li>
            )
          })}
        </ul>
      )}

      {lightboxOpen && current && (
        <Lightbox
          src={current.url}
          alt={alt}
          index={safeIndex}
          total={images.length}
          onClose={close}
          onPrev={prev}
          onNext={next}
        />
      )}
    </div>
  )
}

function Lightbox({
  src,
  alt,
  index,
  total,
  onClose,
  onPrev,
  onNext,
}: {
  src: string
  alt: string
  index: number
  total: number
  onClose: () => void
  onPrev: () => void
  onNext: () => void
}) {
  const hasNav = total > 1
  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="Просмотр изображения"
      onClick={onClose}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/85"
    >
      {/* Сама картинка — отдельный «остров», клики внутри не закрывают модалку. */}
      <div
        className="relative h-[85vh] w-[90vw] max-w-5xl"
        onClick={(e) => e.stopPropagation()}
      >
        <Image
          src={src}
          alt={alt}
          fill
          priority
          sizes="90vw"
          className="object-contain"
        />
      </div>

      {/* Кнопка закрытия. */}
      <button
        type="button"
        onClick={onClose}
        aria-label="Закрыть"
        className="absolute right-4 top-4 rounded-full bg-white/10 p-2 text-white hover:bg-white/20"
      >
        <X className="size-5" />
      </button>

      {hasNav && (
        <>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation()
              onPrev()
            }}
            aria-label="Предыдущее"
            className="absolute left-4 top-1/2 -translate-y-1/2 rounded-full bg-white/10 p-2 text-white hover:bg-white/20"
          >
            <ChevronLeft className="size-6" />
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation()
              onNext()
            }}
            aria-label="Следующее"
            className="absolute right-4 top-1/2 -translate-y-1/2 rounded-full bg-white/10 p-2 text-white hover:bg-white/20"
          >
            <ChevronRight className="size-6" />
          </button>
          <span className="absolute bottom-4 left-1/2 -translate-x-1/2 rounded-md bg-black/50 px-3 py-1 text-sm text-white">
            {index + 1} / {total}
          </span>
        </>
      )}
    </div>
  )
}
