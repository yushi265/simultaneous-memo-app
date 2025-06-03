'use client'

import React, { useState, useRef, useCallback } from 'react'
import { api } from '@/lib/api'
import { 
  FileIcon, 
  UploadIcon, 
  FileTextIcon, 
  FileArchiveIcon,
  FileCodeIcon,
  TrashIcon,
  DownloadIcon,
  ExternalLinkIcon
} from '@radix-ui/react-icons'

interface FileMetadata {
  id: number
  filename: string
  original_name: string
  content_type: string
  size: number
  url: string
  page_id?: number
  created_at: string
}

interface FileUploadProps {
  pageId?: number
  onFileUploaded?: (file: FileMetadata) => void
  onFileDeleted?: (fileId: number) => void
  showUploadArea?: boolean
}

export default function FileUpload({ pageId, onFileUploaded, onFileDeleted, showUploadArea = true }: FileUploadProps) {
  const [isUploading, setIsUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [files, setFiles] = useState<FileMetadata[]>([])
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Load files on component mount
  React.useEffect(() => {
    loadFiles()
  }, [pageId])

  const loadFiles = async () => {
    try {
      const data = await api.getFiles(pageId)
      // Handle both old format (array) and new format (paginated response)
      if (Array.isArray(data)) {
        setFiles(data)
      } else if (data.files) {
        setFiles(data.files)
      }
    } catch (err) {
      console.error('Failed to load files:', err)
    }
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = e.target.files
    if (selectedFiles && selectedFiles.length > 0) {
      uploadFiles(Array.from(selectedFiles))
    }
  }

  const uploadFiles = async (filesToUpload: File[]) => {
    setIsUploading(true)
    setError(null)

    for (const file of filesToUpload) {
      try {
        const result = await api.uploadGeneralFile(file, pageId)
        setFiles(prev => [...prev, result])
        if (onFileUploaded) {
          onFileUploaded(result)
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Upload failed')
      }
    }

    setIsUploading(false)
    setUploadProgress(0)
  }

  const handleDelete = async (fileId: number) => {
    if (!confirm('このファイルを削除してもよろしいですか？')) {
      return
    }

    try {
      await api.deleteFile(fileId)
      setFiles(prev => prev.filter(f => f.id !== fileId))
      if (onFileDeleted) {
        onFileDeleted(fileId)
      }
    } catch (err) {
      setError('Failed to delete file')
    }
  }

  const handleDragEnter = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }, [])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)

    const droppedFiles = Array.from(e.dataTransfer.files)
    if (droppedFiles.length > 0) {
      uploadFiles(droppedFiles)
    }
  }, [pageId])

  const getFileIcon = (contentType: string) => {
    if (contentType.includes('text') || contentType.includes('document')) {
      return <FileTextIcon className="w-4 h-4" />
    } else if (contentType.includes('zip') || contentType.includes('compressed')) {
      return <FileArchiveIcon className="w-4 h-4" />
    } else if (contentType.includes('javascript') || contentType.includes('json') || contentType.includes('xml')) {
      return <FileCodeIcon className="w-4 h-4" />
    }
    return <FileIcon className="w-4 h-4" />
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    else if (bytes < 1048576) return Math.round(bytes / 1024) + ' KB'
    else return (bytes / 1048576).toFixed(2) + ' MB'
  }

  // Show nothing if no upload area and no files
  if (!showUploadArea && files.length === 0) {
    return null
  }

  return (
    <div className="space-y-4">
      {/* Upload Area */}
      {showUploadArea && (
        <div
          onDragEnter={handleDragEnter}
          onDragLeave={handleDragLeave}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
          className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
            isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400'
          }`}
        >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={handleFileSelect}
          className="hidden"
          accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.csv,.rtf,.zip,.rar,.7z,.tar,.gz,.js,.ts,.json,.xml,.html,.css,.py,.go,.java,.cpp,.c,.sh,.md"
        />
        
        <UploadIcon className="mx-auto h-12 w-12 text-gray-400" />
        <p className="mt-2 text-sm text-gray-600">
          ドラッグ&ドロップまたは
          <button
            onClick={() => fileInputRef.current?.click()}
            className="mx-1 text-blue-600 hover:text-blue-800 underline"
            disabled={isUploading}
          >
            クリックしてファイルを選択
          </button>
        </p>
        <p className="text-xs text-gray-500 mt-1">
          最大50MB (PDF, ドキュメント, アーカイブ, コードファイル)
        </p>
        </div>
      )}

      {/* Upload Progress */}
      {isUploading && (
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-blue-700">アップロード中...</span>
            <span className="text-sm text-blue-700">{uploadProgress}%</span>
          </div>
          <div className="w-full bg-blue-200 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${uploadProgress}%` }}
            />
          </div>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-3">
          <p className="text-sm text-red-700">{error}</p>
        </div>
      )}

      {/* File List */}
      {files.length > 0 && (
        <div className="border rounded-lg overflow-hidden">
          <div className="bg-gray-50 px-4 py-2 border-b">
            <h3 className="text-sm font-medium text-gray-700">
              アップロードされたファイル ({files.length})
            </h3>
          </div>
          <div className="divide-y">
            {files.map((file) => (
              <div key={file.id} className="px-4 py-3 hover:bg-gray-50 transition-colors">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3 flex-1 min-w-0">
                    <div className="flex-shrink-0">
                      {getFileIcon(file.content_type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {file.original_name}
                      </p>
                      <p className="text-xs text-gray-500">
                        {formatFileSize(file.size)} • {new Date(file.created_at).toLocaleString('ja-JP')}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2 ml-4">
                    <button
                      onClick={() => window.open(file.url, '_blank')}
                      className="p-1 text-gray-500 hover:text-blue-600 transition-colors"
                      title="開く"
                    >
                      <ExternalLinkIcon className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => {
                        const link = document.createElement('a')
                        link.href = `${file.url}?download=true`
                        link.download = file.original_name
                        document.body.appendChild(link)
                        link.click()
                        document.body.removeChild(link)
                      }}
                      className="p-1 text-gray-500 hover:text-green-600 transition-colors"
                      title="ダウンロード"
                    >
                      <DownloadIcon className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(file.id)}
                      className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                      title="削除"
                    >
                      <TrashIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}