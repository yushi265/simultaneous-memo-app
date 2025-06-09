'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { workspaceApi, WorkspaceResponse } from '@/lib/workspace-api'
import { Header } from '@/components/Header'

export default function WorkspaceSettingsPage() {
  const { currentWorkspace, user } = useAuthStore()
  const [workspace, setWorkspace] = useState<WorkspaceResponse | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const router = useRouter()

  useEffect(() => {
    if (!currentWorkspace) {
      router.push('/')
      return
    }
    
    loadWorkspace()
  }, [currentWorkspace, router])

  const loadWorkspace = async () => {
    if (!currentWorkspace) return

    try {
      setIsLoading(true)
      const workspaceData = await workspaceApi.getWorkspace(currentWorkspace.id)
      setWorkspace(workspaceData)
      setName(workspaceData.name)
      setDescription(workspaceData.description || '')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ワークスペースの取得に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!workspace || !name.trim()) return

    try {
      setIsSaving(true)
      setError('')
      setSuccess('')

      await workspaceApi.updateWorkspace(workspace.id, {
        name: name.trim(),
        description: description.trim() || undefined
      })

      setSuccess('ワークスペース設定を更新しました')
      
      // Reload workspace data
      setTimeout(() => {
        loadWorkspace()
      }, 1000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ワークスペースの更新に失敗しました')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!workspace || workspace.is_personal) return

    const confirmed = window.confirm(
      `「${workspace.name}」を削除してもよろしいですか？この操作は取り消せません。`
    )

    if (!confirmed) return

    try {
      setIsLoading(true)
      await workspaceApi.deleteWorkspace(workspace.id)
      
      // Redirect to home after deletion
      router.push('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ワークスペースの削除に失敗しました')
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <div className="flex items-center justify-center py-20">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  if (!workspace) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <div className="max-w-2xl mx-auto py-20 px-4">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            ワークスペースが見つかりません
          </div>
        </div>
      </div>
    )
  }

  const canEdit = workspace.user_role === 'owner' || workspace.user_role === 'admin'
  const canDelete = workspace.user_role === 'owner' && !workspace.is_personal

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      <div className="max-w-2xl mx-auto py-8 px-4">
        <div className="bg-white rounded-lg shadow">
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-xl font-semibold text-gray-900">
              ワークスペース設定
            </h1>
            <p className="mt-1 text-sm text-gray-600">
              ワークスペースの基本情報を管理します
            </p>
          </div>

          {/* Content */}
          <div className="p-6">
            {error && (
              <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                {error}
              </div>
            )}

            {success && (
              <div className="mb-6 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
                {success}
              </div>
            )}

            <form onSubmit={handleSave} className="space-y-6">
              {/* Basic Info */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">基本情報</h3>
                
                <div className="space-y-4">
                  <div>
                    <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                      ワークスペース名
                    </label>
                    <input
                      id="name"
                      type="text"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      required
                      disabled={!canEdit || isSaving}
                    />
                  </div>

                  <div>
                    <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                      説明
                    </label>
                    <textarea
                      id="description"
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      rows={3}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      disabled={!canEdit || isSaving}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-gray-500">タイプ:</span>
                      <span className="ml-2">
                        {workspace.is_personal ? '個人ワークスペース' : 'チームワークスペース'}
                      </span>
                    </div>
                    <div>
                      <span className="text-gray-500">メンバー数:</span>
                      <span className="ml-2">{workspace.member_count}名</span>
                    </div>
                    <div>
                      <span className="text-gray-500">あなたの権限:</span>
                      <span className="ml-2">{workspace.user_role}</span>
                    </div>
                    <div>
                      <span className="text-gray-500">作成日:</span>
                      <span className="ml-2">
                        {new Date(workspace.created_at).toLocaleDateString('ja-JP')}
                      </span>
                    </div>
                  </div>
                </div>

                {canEdit && (
                  <div className="mt-6 flex justify-end">
                    <button
                      type="submit"
                      className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
                      disabled={isSaving || !name.trim()}
                    >
                      {isSaving ? '保存中...' : '設定を保存'}
                    </button>
                  </div>
                )}
              </div>

              {/* Danger Zone */}
              {canDelete && (
                <div className="border-t border-gray-200 pt-6">
                  <h3 className="text-lg font-medium text-red-900 mb-4">危険な操作</h3>
                  
                  <div className="bg-red-50 border border-red-200 rounded-md p-4">
                    <h4 className="text-sm font-medium text-red-800 mb-2">
                      ワークスペースを削除
                    </h4>
                    <p className="text-sm text-red-700 mb-4">
                      ワークスペースを削除すると、すべてのページ、ファイル、メンバーシップが永久に削除されます。この操作は取り消せません。
                    </p>
                    <button
                      type="button"
                      onClick={handleDelete}
                      className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
                      disabled={isLoading}
                    >
                      ワークスペースを削除
                    </button>
                  </div>
                </div>
              )}
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}