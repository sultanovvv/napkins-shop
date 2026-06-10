import { ShopShell } from '@/components/shop-shell'

export default function AboutPage() {
  return (
    <ShopShell>
      <h1 className="font-serif text-3xl text-foreground">О магазине</h1>
      <div className="mt-4 flex flex-col gap-4 text-sm leading-relaxed text-muted-foreground">
        <p>
          «Много салфеток» — это интернет-магазин для любителей декупажа и
          коллекционеров. Уже много лет мы собираем большую коллекцию салфеток,
          декупажных карт и рисовой бумаги на самые разные сюжеты.
        </p>
        <p>
          В нашем каталоге вы найдёте цветы, птиц, гномов, города, праздничные и
          новогодние мотивы, продуктовую тему, орнаменты и многое другое. Мы
          постоянно добавляем новинки и бережно относимся к каждому заказу.
        </p>
        <p>
          Спасибо, что выбираете нас для своего творчества!
        </p>
      </div>
    </ShopShell>
  )
}
