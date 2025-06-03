const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const getApiUrl = () => API_URL

export const api = {
  // Pages
  async getPages() {
    const response = await fetch(`${API_URL}/api/pages`)
    if (!response.ok) throw new Error('Failed to fetch pages')
    return response.json()
  },

  async getPage(id: number) {
    const response = await fetch(`${API_URL}/api/pages/${id}`)
    if (!response.ok) throw new Error('Failed to fetch page')
    return response.json()
  },

  async createPage(data: { title: string; content?: any }) {
    const response = await fetch(`${API_URL}/api/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    if (!response.ok) throw new Error('Failed to create page')
    return response.json()
  },

  async updatePage(id: number, data: { title?: string; content?: any }) {
    const response = await fetch(`${API_URL}/api/pages/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    if (!response.ok) throw new Error('Failed to update page')
    return response.json()
  },

  async deletePage(id: number) {
    const response = await fetch(`${API_URL}/api/pages/${id}`, {
      method: 'DELETE'
    })
    if (!response.ok) throw new Error('Failed to delete page')
    return response.json()
  },

  // Image upload
  async uploadFile(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await fetch(`${API_URL}/api/upload`, {
      method: 'POST',
      body: formData
    })
    if (!response.ok) throw new Error('Failed to upload file')
    return response.json()
  },

  // General file upload
  async uploadGeneralFile(file: File, pageId?: number) {
    const formData = new FormData()
    formData.append('file', file)
    if (pageId) {
      formData.append('page_id', pageId.toString())
    }
    
    const response = await fetch(`${API_URL}/api/upload/file`, {
      method: 'POST',
      body: formData
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to upload file')
    }
    return response.json()
  },

  async getFiles(pageId?: number, type?: string, page: number = 1, limit: number = 20) {
    const params = new URLSearchParams()
    if (pageId) params.append('page_id', pageId.toString())
    if (type) params.append('type', type)
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    
    const response = await fetch(`${API_URL}/api/files?${params.toString()}`)
    if (!response.ok) throw new Error('Failed to fetch files')
    return response.json()
  },

  async getFileMetadata(id: number) {
    const response = await fetch(`${API_URL}/api/files/${id}`)
    if (!response.ok) throw new Error('Failed to fetch file metadata')
    return response.json()
  },

  async deleteFile(id: number) {
    const response = await fetch(`${API_URL}/api/files/${id}`, {
      method: 'DELETE'
    })
    if (!response.ok) throw new Error('Failed to delete file')
    return response.json()
  }
}