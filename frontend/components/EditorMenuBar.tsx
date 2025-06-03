'use client'

import { Editor } from '@tiptap/react'
import {
  FontBoldIcon,
  FontItalicIcon,
  CodeIcon,
  ListBulletIcon,
  TextIcon,
  ImageIcon,
  FileIcon,
} from '@radix-ui/react-icons'
import { useRef, useState } from 'react'
import { imageUploader } from '@/lib/image-upload'
import { getFullImageUrl } from '@/lib/image-utils'

interface EditorMenuBarProps {
  editor: Editor | null
  pageId?: number
  onFileUploadClick?: () => void
}

export function EditorMenuBar({ editor, pageId, onFileUploadClick }: EditorMenuBarProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [isUploading, setIsUploading] = useState(false)

  if (!editor) {
    return null
  }

  const handleImageUpload = async (file: File) => {
    try {
      setIsUploading(true)
      setUploadProgress(0)

      const response = await imageUploader.uploadImage(file, {
        pageId,
        onProgress: (progress) => setUploadProgress(progress),
        onStart: () => setUploadProgress(0),
      })

      // Insert image into editor using custom command
      const imageUrl = getFullImageUrl(response.url)
        
      editor.chain().focus().setImage({
        src: imageUrl,
        alt: response.filename,
        title: response.filename,
        width: response.width,
        height: response.height,
        // Store image metadata for reference tracking
        'data-image-id': response.id.toString(),
        'data-width': response.width.toString(),
        'data-height': response.height.toString(),
      }).run()

    } catch (error) {
      console.error('Image upload failed:', error)
      alert(error instanceof Error ? error.message : '画像のアップロードに失敗しました')
    } finally {
      setIsUploading(false)
      setUploadProgress(0)
    }
  }

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      handleImageUpload(file)
      // Clear the input so the same file can be selected again
      event.target.value = ''
    }
  }

  return (
    <div className="border-b border-gray-200 p-2 flex items-center gap-1">
      <button
        onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('heading', { level: 1 }) ? 'bg-gray-200' : ''
        }`}
        title="見出し1"
      >
        <span className="font-bold text-lg">H1</span>
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('heading', { level: 2 }) ? 'bg-gray-200' : ''
        }`}
        title="見出し2"
      >
        <span className="font-bold">H2</span>
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('heading', { level: 3 }) ? 'bg-gray-200' : ''
        }`}
        title="見出し3"
      >
        <span className="font-bold text-sm">H3</span>
      </button>

      <div className="w-px h-6 bg-gray-300 mx-1" />

      <button
        onClick={() => editor.chain().focus().toggleBold().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('bold') ? 'bg-gray-200' : ''
        }`}
        title="太字"
      >
        <FontBoldIcon className="w-4 h-4" />
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleItalic().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('italic') ? 'bg-gray-200' : ''
        }`}
        title="斜体"
      >
        <FontItalicIcon className="w-4 h-4" />
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleCode().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('code') ? 'bg-gray-200' : ''
        }`}
        title="インラインコード"
      >
        <CodeIcon className="w-4 h-4" />
      </button>

      <div className="w-px h-6 bg-gray-300 mx-1" />

      <button
        onClick={() => editor.chain().focus().toggleBulletList().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('bulletList') ? 'bg-gray-200' : ''
        }`}
        title="箇条書き"
      >
        <ListBulletIcon className="w-4 h-4" />
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('orderedList') ? 'bg-gray-200' : ''
        }`}
        title="番号付きリスト"
      >
        <span className="text-sm font-medium">1.</span>
      </button>
      
      <button
        onClick={() => editor.chain().focus().toggleCodeBlock().run()}
        className={`p-2 rounded hover:bg-gray-100 ${
          editor.isActive('codeBlock') ? 'bg-gray-200' : ''
        }`}
        title="コードブロック"
      >
        <span className="text-sm font-mono">&lt;/&gt;</span>
      </button>

      <div className="w-px h-6 bg-gray-300 mx-1" />

      <button
        onClick={() => fileInputRef.current?.click()}
        className="p-2 rounded hover:bg-gray-100 relative"
        title="画像をアップロード"
        disabled={isUploading}
      >
        <ImageIcon className="w-4 h-4" />
        {isUploading && (
          <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 rounded">
            <div className="text-xs text-blue-600">{uploadProgress}%</div>
          </div>
        )}
      </button>
      
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        className="hidden"
      />

      <button
        onClick={onFileUploadClick}
        className="p-2 rounded hover:bg-gray-100"
        title="ファイルをアップロード"
      >
        <FileIcon className="w-4 h-4" />
      </button>

      <div className="w-px h-6 bg-gray-300 mx-1" />

      <button
        onClick={() => editor.chain().focus().setHorizontalRule().run()}
        className="p-2 rounded hover:bg-gray-100"
        title="区切り線"
      >
        <span className="text-gray-500">—</span>
      </button>
    </div>
  )
}