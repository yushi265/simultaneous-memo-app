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
  const [hasHydrated, setHasHydrated] = useState(false)
  const { isAuthenticated, token, login, logout } = useAuthStore()
  const router = useRouter()
  const pathname = usePathname()

  // Wait for Zustand to hydrate from localStorage
  useEffect(() => {
    const unsubscribe = useAuthStore.persist.onFinishHydration(() => {
      setHasHydrated(true)
    })
    
    // If already hydrated, set immediately
    if (useAuthStore.persist.hasHydrated()) {
      setHasHydrated(true)
    }
    
    return unsubscribe
  }, [])

  useEffect(() => {
    if (!hasHydrated) {
      return // Wait for hydration
    }
    
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
        console.log('Checking token validity...')
        const userInfo = await authApi.me(token)
        console.log('Token is valid, user info:', userInfo)
        
        // ユーザー情報を更新（既存のワークスペース情報を保持）
        login(token, userInfo.user, userInfo.currentWorkspace)
        setIsLoading(false)
      } catch (error) {
        console.error('Token validation failed:', error)
        
        // 429エラーの場合は一時的な問題として扱い、ログアウトしない
        const errorMessage = error instanceof Error ? error.message : String(error)
        if (errorMessage.includes('Rate limited') || errorMessage.includes('429')) {
          console.warn('Authentication check temporarily failed due to rate limiting, keeping user logged in')
          setIsLoading(false)
          return // ログアウトしない
        }
        
        // その他のエラーの場合はログアウト
        logout()
        router.push('/login')
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [pathname, isAuthenticated, token, login, logout, router, hasHydrated])

  // ハイドレーション完了まで待機
  if (!hasHydrated || isLoading) {
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