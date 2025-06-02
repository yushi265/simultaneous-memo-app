export interface ImageUploadResponse {
  id: number
  filename: string
  size: number
  originalSize: number
  url: string
  thumbnailUrl: string
  contentType: string
  width: number
  height: number
  pageId?: number
  uploadedAt: string
}

export interface ImageUploadOptions {
  pageId?: number
  onProgress?: (progress: number) => void
  onStart?: () => void
  onComplete?: (response: ImageUploadResponse) => void
  onError?: (error: string) => void
}

export class ImageUploader {
  private baseUrl: string

  constructor(baseUrl = 'http://localhost:8080') {
    this.baseUrl = baseUrl
  }

  async uploadImage(file: File, options: ImageUploadOptions = {}): Promise<ImageUploadResponse> {
    const { pageId, onProgress, onStart, onComplete, onError } = options

    try {
      onStart?.()

      // Validate file type
      if (!file.type.startsWith('image/')) {
        throw new Error('画像ファイルのみアップロード可能です')
      }

      // Validate file size (10MB)
      const maxSize = 10 * 1024 * 1024
      if (file.size > maxSize) {
        throw new Error(`ファイルサイズが大きすぎます。最大サイズは${maxSize / 1024 / 1024}MBです`)
      }

      const formData = new FormData()
      formData.append('file', file)
      
      if (pageId) {
        formData.append('page_id', pageId.toString())
      }

      const xhr = new XMLHttpRequest()

      return new Promise((resolve, reject) => {
        xhr.upload.addEventListener('progress', (event) => {
          if (event.lengthComputable) {
            const progress = Math.round((event.loaded / event.total) * 100)
            onProgress?.(progress)
          }
        })

        xhr.addEventListener('load', () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            try {
              const response: ImageUploadResponse = JSON.parse(xhr.responseText)
              onComplete?.(response)
              resolve(response)
            } catch (error) {
              const errorMsg = 'レスポンスの解析に失敗しました'
              onError?.(errorMsg)
              reject(new Error(errorMsg))
            }
          } else {
            let errorMsg = `アップロードに失敗しました (${xhr.status})`
            try {
              const errorResponse = JSON.parse(xhr.responseText)
              errorMsg = errorResponse.error || errorMsg
            } catch (e) {
              // Use default error message
            }
            onError?.(errorMsg)
            reject(new Error(errorMsg))
          }
        })

        xhr.addEventListener('error', () => {
          const errorMsg = 'ネットワークエラーが発生しました'
          onError?.(errorMsg)
          reject(new Error(errorMsg))
        })

        xhr.addEventListener('abort', () => {
          const errorMsg = 'アップロードがキャンセルされました'
          onError?.(errorMsg)
          reject(new Error(errorMsg))
        })

        xhr.open('POST', `${this.baseUrl}/api/upload`)
        xhr.send(formData)
      })
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : 'アップロードに失敗しました'
      onError?.(errorMsg)
      throw error
    }
  }

  /**
   * Get optimized image URL with optional size parameter
   */
  getImageUrl(path: string, size?: 'thumbnail' | 'original'): string {
    const url = `${this.baseUrl}/api/img${path}`
    if (size === 'thumbnail') {
      return `${url}?size=thumbnail`
    }
    return url
  }

  /**
   * Upload image from clipboard
   */
  async uploadFromClipboard(clipboardData: DataTransfer, options: ImageUploadOptions = {}): Promise<ImageUploadResponse | null> {
    const items = clipboardData.items
    
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      
      if (item.type.startsWith('image/')) {
        const file = item.getAsFile()
        if (file) {
          return this.uploadImage(file, options)
        }
      }
    }
    
    return null
  }

  /**
   * Upload image from drag and drop
   */
  async uploadFromDrop(dataTransfer: DataTransfer, options: ImageUploadOptions = {}): Promise<ImageUploadResponse[]> {
    const files = Array.from(dataTransfer.files).filter(file => file.type.startsWith('image/'))
    const results: ImageUploadResponse[] = []
    
    for (const file of files) {
      try {
        const result = await this.uploadImage(file, options)
        results.push(result)
      } catch (error) {
        console.error('Failed to upload file:', file.name, error)
        // Continue with other files even if one fails
      }
    }
    
    return results
  }
}

// Global instance
export const imageUploader = new ImageUploader()