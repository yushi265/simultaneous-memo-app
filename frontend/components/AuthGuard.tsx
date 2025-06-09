'use client'

import { useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { authApi } from '@/lib/auth-api'

interface AuthGuardProps {
  children: React.ReactNode
}

const publicRoutes = ['/login', '/register']

export default function AuthGuard({ children }: AuthGuardProps) {
  const [isLoading, setIsLoading] = useState(true)
  const { isAuthenticated, token, login, logout } = useAuthStore()
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    const checkAuth = async () => {
      // Public routesなら認証チェックをスキップ
      if (publicRoutes.includes(pathname)) {
        // 既にログイン済みならホームにリダイレクト
        if (isAuthenticated && token) {
          router.push('/')
          return
        }
        setIsLoading(false)
        return
      }

      // トークンがない場合はログインページへ
      if (!token) {
        router.push('/login')
        setIsLoading(false)
        return
      }

      // トークンがある場合は有効性をチェック
      try {
        const userInfo = await authApi.me(token)
        // ユーザー情報を更新
        login(token, userInfo.user, userInfo.currentWorkspace)
        setIsLoading(false)
      } catch (error) {
        // トークンが無効な場合はログアウト
        logout()
        router.push('/login')
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [pathname, isAuthenticated, token, login, logout, router])

  // ローディング中は何も表示しない
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  // Public routesまたは認証済みの場合のみ子コンポーネントを表示
  if (publicRoutes.includes(pathname) || (isAuthenticated && token)) {
    return <>{children}</>
  }

  // その他の場合は何も表示しない（リダイレクト処理中）
  return null
}