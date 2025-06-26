'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { workspaceApi, WorkspaceResponse } from '@/lib/workspace-api'
import { Header } from '@/components/Header'
import InviteMemberModal from '@/components/InviteMemberModal'
import { PersonIcon, EnvelopeClosedIcon, TrashIcon, Pencil1Icon } from '@radix-ui/react-icons'

export default function WorkspaceSettingsPage() {
  const { currentWorkspace, user } = useAuthStore()
  const [workspace, setWorkspace] = useState<WorkspaceResponse | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [members, setMembers] = useState<any[]>([])
  const [invitations, setInvitations] = useState<any[]>([])
  const [showInviteModal, setShowInviteModal] = useState(false)
  const [activeTab, setActiveTab] = useState<'settings' | 'members'>('settings')
  const router = useRouter()

  useEffect(() => {
    if (!currentWorkspace) {
      router.push('/')
      return
    }
    
    loadWorkspace()
    if (activeTab === 'members') {
      loadMembers()
      loadInvitations()
    }
  }, [currentWorkspace, router, activeTab])

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

  const loadMembers = async () => {
    if (!currentWorkspace) return

    try {
      const response = await workspaceApi.getMembers(currentWorkspace.id)
      setMembers(response.members)
    } catch (err) {
      console.error('Failed to load members:', err)
    }
  }

  const loadInvitations = async () => {
    if (!currentWorkspace) return

    try {
      const response = await workspaceApi.getInvitations(currentWorkspace.id)
      setInvitations(response.invitations)
    } catch (err) {
      console.error('Failed to load invitations:', err)
    }
  }

  const handleInviteSuccess = () => {
    loadInvitations()
    setSuccess('メンバーの招待を送信しました')
    setTimeout(() => setSuccess(''), 3000)
  }

  const handleCancelInvitation = async (invitationId: string) => {
    if (!currentWorkspace) return
    
    try {
      await workspaceApi.cancelInvitation(currentWorkspace.id, invitationId)
      loadInvitations()
      setSuccess('招待をキャンセルしました')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err instanceof Error ? err.message : '招待のキャンセルに失敗しました')
    }
  }

  const handleUpdateMemberRole = async (memberId: string, newRole: string) => {
    if (!currentWorkspace) return
    
    try {
      await workspaceApi.updateMemberRole(currentWorkspace.id, memberId, newRole)
      loadMembers()
      setSuccess('メンバーの権限を更新しました')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'メンバー権限の更新に失敗しました')
    }
  }

  const handleRemoveMember = async (memberId: string, memberName: string) => {
    if (!currentWorkspace) return
    
    if (!confirm(`${memberName}をワークスペースから削除してもよろしいですか？`)) {
      return
    }
    
    try {
      await workspaceApi.removeMember(currentWorkspace.id, memberId)
      loadMembers()
      setSuccess('メンバーを削除しました')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'メンバーの削除に失敗しました')
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
      
      <div className="max-w-4xl mx-auto py-8 px-4">
        <div className="bg-white rounded-lg shadow">
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-xl font-semibold text-gray-900">
              ワークスペース設定
            </h1>
            <p className="mt-1 text-sm text-gray-600">
              ワークスペースの基本情報とメンバーを管理します
            </p>
          </div>

          {/* Tabs */}
          <div className="px-6 py-0 border-b border-gray-200">
            <nav className="flex space-x-8">
              <button
                onClick={() => setActiveTab('settings')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'settings'
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                基本設定
              </button>
              {!workspace?.is_personal && (
                <button
                  onClick={() => setActiveTab('members')}
                  className={`py-4 px-1 border-b-2 font-medium text-sm ${
                    activeTab === 'members'
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  メンバー管理
                </button>
              )}
            </nav>
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

            {activeTab === 'settings' && (
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
            )}

            {activeTab === 'members' && (
              <div className="space-y-6">
                {/* Member Management Header */}
                <div className="flex items-center justify-between">
                  <h3 className="text-lg font-medium text-gray-900">メンバー管理</h3>
                  {canEdit && (
                    <button
                      onClick={() => setShowInviteModal(true)}
                      className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                    >
                      <PersonIcon className="w-4 h-4" />
                      メンバーを招待
                    </button>
                  )}
                </div>

                {/* Members List */}
                <div>
                  <h4 className="text-md font-medium text-gray-900 mb-3">
                    メンバー ({members.length}名)
                  </h4>
                  <div className="bg-gray-50 rounded-lg divide-y divide-gray-200">
                    {members.map((member) => (
                      <div key={member.id} className="p-4 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                            <PersonIcon className="w-4 h-4 text-blue-600" />
                          </div>
                          <div>
                            <p className="font-medium text-gray-900">{member.user.name}</p>
                            <p className="text-sm text-gray-500">{member.user.email}</p>
                          </div>
                        </div>
                        <div className="flex items-center gap-3">
                          <select
                            value={member.role}
                            onChange={(e) => handleUpdateMemberRole(member.user.id, e.target.value)}
                            className="text-sm border border-gray-300 rounded px-2 py-1"
                            disabled={member.role === 'owner' || !canEdit}
                          >
                            <option value="viewer">閲覧者</option>
                            <option value="member">メンバー</option>
                            <option value="admin">管理者</option>
                            <option value="owner">オーナー</option>
                          </select>
                          {member.role !== 'owner' && canEdit && (
                            <button
                              onClick={() => handleRemoveMember(member.user.id, member.user.name)}
                              className="text-red-600 hover:text-red-800"
                            >
                              <TrashIcon className="w-4 h-4" />
                            </button>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Pending Invitations */}
                {invitations.length > 0 && (
                  <div>
                    <h4 className="text-md font-medium text-gray-900 mb-3">
                      保留中の招待 ({invitations.length}件)
                    </h4>
                    <div className="bg-yellow-50 rounded-lg divide-y divide-yellow-200">
                      {invitations.map((invitation) => (
                        <div key={invitation.id} className="p-4 flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
                              <EnvelopeClosedIcon className="w-4 h-4 text-yellow-600" />
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">{invitation.email}</p>
                              <p className="text-sm text-gray-500">
                                {invitation.role} • {invitation.inviter.name}が招待
                              </p>
                            </div>
                          </div>
                          {canEdit && (
                            <button
                              onClick={() => handleCancelInvitation(invitation.id)}
                              className="text-red-600 hover:text-red-800 text-sm"
                            >
                              キャンセル
                            </button>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Invite Member Modal */}
      {workspace && (
        <InviteMemberModal
          isOpen={showInviteModal}
          onClose={() => setShowInviteModal(false)}
          onSuccess={handleInviteSuccess}
          workspaceId={workspace.id}
          workspaceName={workspace.name}
        />
      )}
    </div>
  )
}