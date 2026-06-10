import Link from 'next/link'

export function SiteFooter() {
  return (
    <footer className="mt-12 border-t border-border bg-secondary">
      <div className="mx-auto grid max-w-7xl gap-8 px-4 py-10 md:grid-cols-3">
        <div>
          <p className="font-serif text-2xl italic text-primary">Много салфеток</p>
          <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
            Интернет-магазин салфеток, карт и рисовой бумаги для декупажа.
            Большой выбор сюжетов для творчества и коллекционирования.
          </p>
        </div>
        <div>
          <h3 className="mb-3 font-serif text-lg text-foreground">Покупателям</h3>
          <ul className="flex flex-col gap-2 text-sm text-muted-foreground">
            <li>
              <Link href="/how-to-buy" className="hover:text-primary">
                Как покупать
              </Link>
            </li>
            <li>
              <Link href="/about" className="hover:text-primary">
                О магазине
              </Link>
            </li>
            <li>
              <Link href="/contacts" className="hover:text-primary">
                Контакты
              </Link>
            </li>
            <li>
              <Link href="/sitemap" className="hover:text-primary">
                Карта сайта
              </Link>
            </li>
          </ul>
        </div>
        <div>
          <h3 className="mb-3 font-serif text-lg text-foreground">Контакты</h3>
          <ul className="flex flex-col gap-2 text-sm text-muted-foreground">
            <li>Телефон: +7 (000) 000-00-00</li>
            <li>E-mail: info@mnogosalfetok.ru</li>
            <li>Работаем ежедневно с 10:00 до 20:00</li>
          </ul>
        </div>
      </div>
      <div className="border-t border-border py-4 text-center text-xs text-muted-foreground">
        © {new Date().getFullYear()} Много салфеток. Все права защищены.
      </div>
    </footer>
  )
}
