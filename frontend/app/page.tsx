'use client'

import { Header } from '@/components/Header'
import { Sidebar } from '@/components/Sidebar'
import { Editor } from '@/components/Editor'
import { useStore } from '@/lib/store'

export default function Home() {
  const { currentPage } = useStore()

  return (
    <div className="h-screen flex flex-col">
      <Header />
      
      <div className="flex-1 flex overflow-hidden">
        <Sidebar />
        
        <main className="flex-1 bg-white overflow-y-auto">
          {currentPage ? (
            <Editor key={currentPage.id} pageId={currentPage.id} />
          ) : (
            <div className="h-full flex items-center justify-center text-gray-500">
              <div className="text-center">
                <p className="text-lg">ページを選択または作成してください</p>
                <p className="text-sm mt-2">左のサイドバーから既存のページを選択するか、新規ページボタンをクリックしてください</p>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}