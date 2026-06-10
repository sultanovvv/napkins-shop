import { ShopShell } from '@/components/shop-shell'

const steps = [
  {
    title: '1. Выберите товары',
    text: 'Найдите нужные салфетки через каталог категорий или поиск и добавьте их в корзину.',
  },
  {
    title: '2. Оформите заказ',
    text: 'Перейдите в корзину, проверьте состав заказа и укажите данные для доставки.',
  },
  {
    title: '3. Оплата',
    text: 'Оплатите заказ удобным способом: банковской картой или при получении.',
  },
  {
    title: '4. Доставка',
    text: 'Мы аккуратно упакуем салфетки и отправим их Почтой России или курьерской службой.',
  },
]

export default function HowToBuyPage() {
  return (
    <ShopShell>
      <h1 className="font-serif text-3xl text-foreground">Как покупать</h1>
      <div className="mt-6 grid gap-4 sm:grid-cols-2">
        {steps.map((s) => (
          <div
            key={s.title}
            className="rounded-lg border border-border bg-card p-5"
          >
            <h2 className="font-serif text-lg text-primary">{s.title}</h2>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              {s.text}
            </p>
          </div>
        ))}
      </div>
    </ShopShell>
  )
}
