import { ShopShell } from '@/components/shop-shell'
import { Mail, Phone, Clock } from 'lucide-react'

export default function ContactsPage() {
  return (
    <ShopShell>
      <h1 className="font-serif text-3xl text-foreground">Контакты</h1>
      <div className="mt-6 flex flex-col gap-4">
        <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
          <Phone className="size-5 text-primary" />
          <span className="text-sm text-foreground">+7 (000) 000-00-00</span>
        </div>
        <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
          <Mail className="size-5 text-primary" />
          <span className="text-sm text-foreground">info@mnogosalfetok.ru</span>
        </div>
        <div className="flex items-center gap-3 rounded-lg border border-border bg-card p-4">
          <Clock className="size-5 text-primary" />
          <span className="text-sm text-foreground">
            Ежедневно с 10:00 до 20:00
          </span>
        </div>
      </div>
    </ShopShell>
  )
}
