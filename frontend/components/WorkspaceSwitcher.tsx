'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { ChevronDownIcon, PlusIcon } from '@radix-ui/react-icons'
import { useAuthStore, useStore } from '@/lib/store'
import { workspaceApi, WorkspaceResponse } from '@/lib/workspace-api'
import { api } from '@/lib/api'

interface WorkspaceSwitcherProps {
  onCreateWorkspace?: () => void
}

export default function WorkspaceSwitcher({ onCreateWorkspace }: WorkspaceSwitcherProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [workspaces, setWorkspaces] = useState<WorkspaceResponse[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const { currentWorkspace, login, user, token } = useAuthStore()
  const { setPages, setCurrentPage } = useStore()
  const router = useRouter()

  useEffect(() => {
    loadWorkspaces()
  }, [])

  const loadWorkspaces = async () => {
    if (!token) return
    
    try {
      setIsLoading(true)
      const fetchedWorkspaces = await workspaceApi.getWorkspaces()
      setWorkspaces(fetchedWorkspaces)
    } catch (error) {
      console.error('Failed to load workspaces:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleSwitchWorkspace = async (workspaceId: string) => {
    if (!user || workspaceId === currentWorkspace?.id) return

    try {
      setIsLoading(true)
      console.log('Switching to workspace:', workspaceId)
      
      const response = await workspaceApi.switchWorkspace(workspaceId)
      console.log('Switch response:', response)
      
      // Update auth store with new token and workspace
      login(response.token, user, response.workspace)
      console.log('Updated store with new token and workspace')
      
      // Clear current page and load pages for new workspace
      setCurrentPage(null)
      setPages([])
      
      setIsOpen(false)
      
      // Reload workspace list to reflect changes
      loadWorkspaces()
      
      // Load pages for the new workspace in the background
      setTimeout(async () => {
        try {
          const newPages = await api.getPages()
          setPages(newPages)
          console.log('Loaded pages for new workspace:', newPages.length)
        } catch (pageError) {
          console.error('Failed to load pages for new workspace:', pageError)
        }
      }, 500)
      
    } catch (error) {
      console.error('Failed to switch workspace:', error)
      alert('ワークスペースの切り替えに失敗しました: ' + (error instanceof Error ? error.message : String(error)))
    } finally {
      setIsLoading(false)
    }
  }

  if (!currentWorkspace) return null

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors min-w-[200px] justify-between"
        disabled={isLoading}
      >
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${currentWorkspace.is_personal ? 'bg-blue-500' : 'bg-green-500'}`}></div>
          <span className="truncate">{currentWorkspace.name}</span>
        </div>
        <ChevronDownIcon className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>

      {isOpen && (
        <>
          {/* Backdrop */}
          <div 
            className="fixed inset-0 z-40" 
            onClick={() => setIsOpen(false)}
          />
          
          {/* Dropdown */}
          <div className="absolute left-0 mt-2 w-72 bg-white border border-gray-200 rounded-md shadow-lg z-50">
            <div className="p-3 border-b border-gray-100">
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                ワークスペース
              </div>
            </div>
            
            <div className="py-1 max-h-64 overflow-y-auto">
              {workspaces.map((workspace) => (
                <button
                  key={workspace.id}
                  onClick={() => handleSwitchWorkspace(workspace.id)}
                  className={`w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center justify-between ${
                    workspace.id === currentWorkspace?.id ? 'bg-blue-50 text-blue-700' : 'text-gray-700'
                  }`}
                  disabled={isLoading}
                >
                  <div className="flex items-center gap-3">
                    <div className={`w-2 h-2 rounded-full ${workspace.is_personal ? 'bg-blue-500' : 'bg-green-500'}`}></div>
                    <div>
                      <div className="font-medium truncate">{workspace.name}</div>
                      <div className="text-xs text-gray-500">
                        {workspace.member_count}名 • {workspace.user_role}
                        {workspace.is_personal && ' • 個人'}
                      </div>
                    </div>
                  </div>
                  {workspace.id === currentWorkspace?.id && (
                    <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                  )}
                </button>
              ))}
            </div>
            
            <div className="border-t border-gray-100 p-2">
              <button
                onClick={() => {
                  setIsOpen(false)
                  onCreateWorkspace?.()
                }}
                className="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md"
              >
                <PlusIcon className="w-4 h-4" />
                新しいワークスペース
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  )
}