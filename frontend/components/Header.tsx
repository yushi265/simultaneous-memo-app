'use client'

import { PlusIcon, ExitIcon, PersonIcon, GearIcon } from '@radix-ui/react-icons'
import { useStore, useAuthStore } from '@/lib/store'
import { api } from '@/lib/api'
import { authApi } from '@/lib/auth-api'
import { Logo } from './Logo'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import WorkspaceSwitcher from './WorkspaceSwitcher'
import CreateWorkspaceModal from './CreateWorkspaceModal'

export function Header() {
  const { setCurrentPage, addPage, currentPage } = useStore()
  const { user, currentWorkspace, logout, token } = useAuthStore()
  const [showUserMenu, setShowUserMenu] = useState(false)
  const [showCreateWorkspace, setShowCreateWorkspace] = useState(false)
  const router = useRouter()
  
  const handleNewPage = async () => {
    try {
      const newPage = await api.createPage({
        title: '無題のページ',
        content: { blocks: [] }
      })
      addPage(newPage)
      setCurrentPage(newPage)
    } catch (error) {
      console.error('Failed to create page:', error)
    }
  }

  return (
    <header className="h-14 border-b border-gray-200 bg-white px-4 flex items-center justify-between">
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Logo className="w-8 h-8 text-gray-700" />
            <h1 className="text-xl font-semibold">リアルタイムメモ</h1>
          </div>
          
          <WorkspaceSwitcher onCreateWorkspace={() => setShowCreateWorkspace(true)} />
        </div>
        
        <button
          onClick={handleNewPage}
          className="flex items-center gap-1 px-3 py-1.5 text-sm bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
        >
          <PlusIcon className="w-4 h-4" />
          新規ページ
        </button>
      </div>
      
      <div className="flex items-center gap-4">
        {currentPage && (
          <span className="text-sm text-gray-500">
            編集中: {currentPage.title}
          </span>
        )}
        
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <span className="text-sm text-gray-600">接続中</span>
        </div>
        
        {/* User Menu */}
        <div className="relative">
          <button
            onClick={() => setShowUserMenu(!showUserMenu)}
            className="flex items-center gap-2 px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
          >
            <PersonIcon className="w-4 h-4" />
            <span>{user?.name}</span>
          </button>
          
          {showUserMenu && (
            <div className="absolute right-0 mt-2 w-64 bg-white border border-gray-200 rounded-md shadow-lg z-50">
              <div className="p-3 border-b border-gray-100">
                <div className="text-sm font-medium text-gray-900">{user?.name}</div>
                <div className="text-xs text-gray-500">{user?.email}</div>
                <div className="text-xs text-gray-500 mt-1">{currentWorkspace?.name}</div>
              </div>
              
              <div className="py-1">
                <button
                  onClick={() => {
                    setShowUserMenu(false)
                    router.push('/workspace/settings')
                  }}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  <GearIcon className="w-4 h-4" />
                  ワークスペース設定
                </button>
                
                <button
                  onClick={async () => {
                    try {
                      if (token) {
                        await authApi.logout(token)
                      }
                    } catch (error) {
                      console.error('Logout error:', error)
                    } finally {
                      logout()
                      router.push('/login')
                    }
                  }}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  <ExitIcon className="w-4 h-4" />
                  ログアウト
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
      
      <CreateWorkspaceModal
        isOpen={showCreateWorkspace}
        onClose={() => setShowCreateWorkspace(false)}
        onSuccess={() => {
          // Close modal and let WorkspaceSwitcher reload its list
          setShowCreateWorkspace(false)
        }}
      />
    </header>
  )
}