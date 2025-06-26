'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { workspaceApi } from '@/lib/workspace-api'
import { Header } from '@/components/Header'
import { CheckCircledIcon, CrossCircledIcon, ExclamationTriangleIcon } from '@radix-ui/react-icons'

interface InvitePageProps {
  params: {
    token: string
  }
}

export default function InvitePage({ params }: InvitePageProps) {
  const { isAuthenticated, user } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [workspaceName, setWorkspaceName] = useState('')
  const router = useRouter()

  useEffect(() => {
    if (!isAuthenticated) {
      // 未認証の場合は招待トークンを保持してログインページへ
      router.push(`/login?invite=${params.token}`)
      return
    }
  }, [isAuthenticated, params.token, router])

  const handleAcceptInvitation = async () => {
    if (!isAuthenticated) return

    try {
      setIsLoading(true)
      setError('')
      
      const response = await workspaceApi.acceptInvitation(params.token)
      
      setSuccess(`「${response.workspace.name}」への参加が完了しました！`)
      setWorkspaceName(response.workspace.name)
      
      // 3秒後にワークスペースに移動
      setTimeout(() => {
        router.push('/')
      }, 3000)
      
    } catch (err) {
      setError(err instanceof Error ? err.message : '招待の受諾に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const handleDecline = () => {
    router.push('/')
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">認証情報を確認しています...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      <div className="max-w-md mx-auto py-20 px-4">
        <div className="bg-white rounded-lg shadow-xl p-8 text-center">
          {success ? (
            <div>
              <CheckCircledIcon className="w-16 h-16 text-green-500 mx-auto mb-4" />
              <h1 className="text-2xl font-bold text-gray-900 mb-4">
                参加完了！
              </h1>
              <p className="text-gray-600 mb-6">
                {success}
              </p>
              <p className="text-sm text-gray-500">
                まもなくワークスペースにリダイレクトします...
              </p>
            </div>
          ) : error ? (
            <div>
              <CrossCircledIcon className="w-16 h-16 text-red-500 mx-auto mb-4" />
              <h1 className="text-2xl font-bold text-gray-900 mb-4">
                招待の受諾に失敗
              </h1>
              <p className="text-red-600 mb-6">
                {error}
              </p>
              <div className="space-y-3">
                <button
                  onClick={handleAcceptInvitation}
                  className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                  disabled={isLoading}
                >
                  再試行
                </button>
                <button
                  onClick={handleDecline}
                  className="w-full px-4 py-2 text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                >
                  ホームに戻る
                </button>
              </div>
            </div>
          ) : (
            <div>
              <ExclamationTriangleIcon className="w-16 h-16 text-blue-500 mx-auto mb-4" />
              <h1 className="text-2xl font-bold text-gray-900 mb-4">
                ワークスペースへの招待
              </h1>
              <p className="text-gray-600 mb-2">
                {user?.name} さん、
              </p>
              <p className="text-gray-600 mb-6">
                ワークスペースに参加しますか？
              </p>
              
              <div className="space-y-3">
                <button
                  onClick={handleAcceptInvitation}
                  className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
                  disabled={isLoading}
                >
                  {isLoading ? '処理中...' : '参加する'}
                </button>
                <button
                  onClick={handleDecline}
                  className="w-full px-4 py-2 text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                  disabled={isLoading}
                >
                  辞退する
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}