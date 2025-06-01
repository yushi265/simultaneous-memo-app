'use client'

import { useEffect } from 'react'
import { TrashIcon } from '@radix-ui/react-icons'
import { useStore, Page } from '@/lib/store'
import { api } from '@/lib/api'

export function Sidebar() {
  const { pages, currentPage, setPages, setCurrentPage, deletePage, setLoading } = useStore()

  useEffect(() => {
    loadPages()
  }, [])

  const loadPages = async () => {
    try {
      setLoading(true)
      const data = await api.getPages()
      setPages(data)
    } catch (error) {
      console.error('Failed to load pages:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSelectPage = async (page: Page) => {
    try {
      const fullPage = await api.getPage(page.id)
      setCurrentPage(fullPage)
    } catch (error) {
      console.error('Failed to load page:', error)
    }
  }

  const handleDeletePage = async (e: React.MouseEvent, pageId: number) => {
    e.stopPropagation()
    if (!confirm('このページを削除しますか？')) return
    
    try {
      await api.deletePage(pageId)
      deletePage(pageId)
    } catch (error) {
      console.error('Failed to delete page:', error)
    }
  }

  return (
    <aside className="w-64 bg-gray-50 border-r border-gray-200 p-4 overflow-y-auto">
      <h2 className="text-sm font-semibold text-gray-600 mb-4">ページ一覧</h2>
      
      <div className="space-y-1">
        {pages.map((page) => (
          <div
            key={page.id}
            onClick={() => handleSelectPage(page)}
            className={`
              group flex items-center justify-between px-3 py-2 rounded-md cursor-pointer transition-colors
              ${currentPage?.id === page.id 
                ? 'bg-blue-100 text-blue-700' 
                : 'hover:bg-gray-100 text-gray-700'
              }
            `}
          >
            <span className="text-sm truncate flex-1">{page.title}</span>
            <button
              onClick={(e) => handleDeletePage(e, page.id)}
              className="opacity-0 group-hover:opacity-100 p-1 hover:bg-gray-200 rounded transition-opacity"
            >
              <TrashIcon className="w-4 h-4 text-gray-500" />
            </button>
          </div>
        ))}
        
        {pages.length === 0 && (
          <p className="text-sm text-gray-500 text-center py-4">
            ページがありません
          </p>
        )}
      </div>
    </aside>
  )
}