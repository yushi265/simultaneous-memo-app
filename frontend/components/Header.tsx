'use client'

import { PlusIcon } from '@radix-ui/react-icons'
import { useStore } from '@/lib/store'
import { api } from '@/lib/api'

export function Header() {
  const { setCurrentPage, addPage, currentPage } = useStore()
  
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
        <h1 className="text-xl font-semibold">Notion-like Memo</h1>
        <button
          onClick={handleNewPage}
          className="flex items-center gap-1 px-3 py-1.5 text-sm bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
        >
          <PlusIcon className="w-4 h-4" />
          新規ページ
        </button>
      </div>
      
      <div className="flex items-center gap-2">
        {currentPage && (
          <span className="text-sm text-gray-500">
            編集中: {currentPage.title}
          </span>
        )}
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <span className="text-sm text-gray-600">接続中</span>
        </div>
      </div>
    </header>
  )
}