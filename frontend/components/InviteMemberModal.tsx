'use client'

import { useState } from 'react'
import { Cross2Icon, PersonIcon, EnvelopeClosedIcon } from '@radix-ui/react-icons'
import { workspaceApi } from '@/lib/workspace-api'

interface InviteMemberModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
  workspaceId: string
  workspaceName: string
}

export default function InviteMemberModal({ 
  isOpen, 
  onClose, 
  onSuccess, 
  workspaceId, 
  workspaceName 
}: InviteMemberModalProps) {
  const [email, setEmail] = useState('')
  const [role, setRole] = useState('member')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [inviteUrl, setInviteUrl] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim()) return

    try {
      setIsLoading(true)
      setError('')
      setSuccess('')
      setInviteUrl('')

      const response = await workspaceApi.inviteMember(workspaceId, email.trim(), role)
      
      setSuccess(response.message)
      setInviteUrl(response.invite_url)
      setEmail('')
      
      // Call success callback to refresh member/invitation lists
      onSuccess()
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'メンバーの招待に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const copyInviteUrl = () => {
    const fullUrl = `${window.location.origin}${inviteUrl}`
    navigator.clipboard.writeText(fullUrl)
    alert('招待URLをクリップボードにコピーしました')
  }

  const resetForm = () => {
    setEmail('')
    setRole('member')
    setError('')
    setSuccess('')
    setInviteUrl('')
  }

  const handleClose = () => {
    resetForm()
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div className="flex items-center gap-3">
            <PersonIcon className="w-5 h-5 text-blue-600" />
            <h2 className="text-lg font-semibold text-gray-900">
              メンバーを招待
            </h2>
          </div>
          <button
            onClick={handleClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <Cross2Icon className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          <p className="text-sm text-gray-600 mb-4">
            「{workspaceName}」に新しいメンバーを招待します
          </p>

          {error && (
            <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
              {error}
            </div>
          )}

          {success && (
            <div className="mb-4 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
              <p className="font-medium">{success}</p>
              {inviteUrl && (
                <div className="mt-2">
                  <p className="text-sm">招待URL:</p>
                  <div className="flex items-center gap-2 mt-1">
                    <input
                      type="text"
                      value={`${window.location.origin}${inviteUrl}`}
                      readOnly
                      className="flex-1 text-xs bg-gray-50 border border-gray-200 rounded px-2 py-1"
                    />
                    <button
                      onClick={copyInviteUrl}
                      className="text-xs bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
                    >
                      コピー
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                メールアドレス
              </label>
              <div className="relative">
                <EnvelopeClosedIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="user@example.com"
                  required
                  disabled={isLoading}
                />
              </div>
            </div>

            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-700 mb-1">
                権限
              </label>
              <select
                id="role"
                value={role}
                onChange={(e) => setRole(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                disabled={isLoading}
              >
                <option value="member">メンバー</option>
                <option value="admin">管理者</option>
                <option value="viewer">閲覧者</option>
              </select>
              <p className="mt-1 text-xs text-gray-500">
                {role === 'admin' && '管理者: ワークスペースの設定とメンバー管理が可能'}
                {role === 'member' && 'メンバー: ページの作成・編集が可能'}
                {role === 'viewer' && '閲覧者: ページの閲覧のみ可能'}
              </p>
            </div>

            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={handleClose}
                className="flex-1 px-4 py-2 text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
                disabled={isLoading}
              >
                キャンセル
              </button>
              <button
                type="submit"
                className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
                disabled={isLoading || !email.trim()}
              >
                {isLoading ? '送信中...' : '招待を送信'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}