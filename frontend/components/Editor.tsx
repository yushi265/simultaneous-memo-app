'use client'

import { useEffect, useRef } from 'react'
import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { ResizableImageExtension } from '@/lib/image-resize-extension'
import { common, createLowlight } from 'lowlight'
import Collaboration from '@tiptap/extension-collaboration'
import CollaborationCursor from '@tiptap/extension-collaboration-cursor'
import * as Y from 'yjs'
import { WebsocketProvider } from 'y-websocket'
import { useStore } from '@/lib/store'
import { api } from '@/lib/api'
import { EditorMenuBar } from './EditorMenuBar'
import { imageUploader } from '@/lib/image-upload'

const lowlight = createLowlight(common)

interface EditorProps {
  pageId: number
}

export function Editor({ pageId }: EditorProps) {
  const { currentPage, updatePage } = useStore()
  const ydocRef = useRef<Y.Doc | null>(null)
  const providerRef = useRef<WebsocketProvider | null>(null)
  const saveTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const editorRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Initialize Yjs
    const ydoc = new Y.Doc()
    ydocRef.current = ydoc

    // Connect to WebSocket
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080'
    const provider = new WebsocketProvider(
      `${wsUrl}/ws/${pageId}`,
      'page-' + pageId,
      ydoc
    )
    providerRef.current = provider

    return () => {
      provider.destroy()
      ydoc.destroy()
    }
  }, [pageId])

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        history: false, // Yjs handles history
        codeBlock: false, // Use CodeBlockLowlight instead
      }),
      Placeholder.configure({
        placeholder: 'ここに入力してください...',
      }),
      CodeBlockLowlight.configure({
        lowlight,
      }),
      ResizableImageExtension.configure({
        HTMLAttributes: {
          class: 'max-w-full h-auto rounded-lg',
        },
        inline: false,
        allowBase64: false,
      }),
      ...(ydocRef.current && providerRef.current ? [
        Collaboration.configure({
          document: ydocRef.current,
        }),
        CollaborationCursor.configure({
          provider: providerRef.current,
          user: {
            name: `User ${Math.floor(Math.random() * 100)}`,
            color: `#${Math.floor(Math.random()*16777215).toString(16)}`,
          },
        }),
      ] : []),
    ],
    content: currentPage?.content || '',
    onUpdate: ({ editor }) => {
      // Debounce save
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current)
      }
      
      saveTimeoutRef.current = setTimeout(() => {
        const content = editor.getJSON()
        saveContent(content)
      }, 1000)
    },
    editorProps: {
      handlePaste: (view, event, slice) => {
        // Handle image paste from clipboard
        const clipboardData = event.clipboardData
        if (clipboardData && editor) {
          imageUploader.uploadFromClipboard(clipboardData, {
            pageId,
            onComplete: (response) => {
              editor.chain().focus().setImage({
                src: response.url,
                alt: response.filename,
                title: response.filename,
                width: response.width,
                height: response.height,
                'data-image-id': response.id.toString(),
                'data-width': response.width.toString(),
                'data-height': response.height.toString(),
              }).run()
            },
            onError: (error) => {
              console.error('Paste image upload failed:', error)
            }
          })
        }
        return false // Let default paste handling continue
      },
    },
  }, [pageId, ydocRef.current, providerRef.current])

  useEffect(() => {
    // Load initial content when page changes
    if (editor && currentPage?.content) {
      editor.commands.setContent(currentPage.content)
    }
  }, [editor, currentPage])

  const saveContent = async (content: any) => {
    try {
      await api.updatePage(pageId, { content })
      updatePage(pageId, { content })
    } catch (error) {
      console.error('Failed to save content:', error)
    }
  }

  const saveTitle = async (title: string) => {
    try {
      await api.updatePage(pageId, { title })
      updatePage(pageId, { title })
    } catch (error) {
      console.error('Failed to save title:', error)
    }
  }

  // Handle drag and drop
  const handleDrop = async (event: React.DragEvent) => {
    event.preventDefault()
    
    const files = imageUploader.uploadFromDrop(event.dataTransfer, {
      pageId,
      onComplete: (response) => {
        if (editor) {
          editor.chain().focus().setImage({
            src: response.url,
            alt: response.filename,
            title: response.filename,
            width: response.width,
            height: response.height,
            'data-image-id': response.id.toString(),
            'data-width': response.width.toString(),
            'data-height': response.height.toString(),
          }).run()
        }
      },
      onError: (error) => {
        console.error('Drop image upload failed:', error)
      }
    })
  }

  const handleDragOver = (event: React.DragEvent) => {
    event.preventDefault()
  }

  if (!currentPage) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-500">
        ページを選択してください
      </div>
    )
  }

  return (
    <div className="flex-1 flex flex-col">
      <div className="p-8 pb-0">
        <input
          type="text"
          value={currentPage.title}
          onChange={(e) => {
            updatePage(pageId, { title: e.target.value })
            saveTitle(e.target.value)
          }}
          className="text-3xl font-bold w-full outline-none border-none"
          placeholder="無題"
        />
      </div>
      
      <EditorMenuBar editor={editor} pageId={pageId} />
      
      <div 
        ref={editorRef}
        className="flex-1 p-8 pt-4"
        onDrop={handleDrop}
        onDragOver={handleDragOver}
      >
        <EditorContent 
          editor={editor} 
          className="prose prose-lg max-w-none focus:outline-none"
        />
      </div>
    </div>
  )
}