/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  // Прокидываем /api/* на Go-бэкенд, чтобы для браузера запросы выглядели
  // same-origin. Иначе SameSite=Lax cookie (gct/rft) с 127.0.0.1:9888 не
  // отправляются на fetch с localhost:3000 — корзина «забывается» между
  // запросами и каждый раз создаётся новая гостевая.
  async rewrites() {
    const backend = process.env.BACKEND_URL ?? 'http://127.0.0.1:9888'
    return [{ source: '/api/:path*', destination: `${backend}/api/:path*` }]
  },
}

export default nextConfig
