import { getAuthHeaders } from './auth-api'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const getApiUrl = () => API_URL

// Retry configuration
const RETRY_ATTEMPTS = 3
const RETRY_DELAY = 1000 // 1 second

// Exponential backoff retry logic
async function retryFetch(url: string, options: RequestInit, attempts: number = RETRY_ATTEMPTS): Promise<Response> {
  try {
    const response = await fetch(url, options)
    
    // If 429 error, retry with exponential backoff
    if (response.status === 429 && attempts > 0) {
      const delay = RETRY_DELAY * (RETRY_ATTEMPTS - attempts + 1)
      console.log(`Rate limited, retrying in ${delay}ms... (${attempts} attempts left)`)
      await new Promise(resolve => setTimeout(resolve, delay))
      return retryFetch(url, options, attempts - 1)
    }
    
    return response
  } catch (error) {
    if (attempts > 0) {
      const delay = RETRY_DELAY * (RETRY_ATTEMPTS - attempts + 1)
      console.log(`Request failed, retrying in ${delay}ms... (${attempts} attempts left)`)
      await new Promise(resolve => setTimeout(resolve, delay))
      return retryFetch(url, options, attempts - 1)
    }
    throw error
  }
}

// Helper function to get token from auth store
const getToken = () => {
  if (typeof window !== 'undefined') {
    const authStorage = JSON.parse(localStorage.getItem('auth-storage') || '{}')
    return authStorage.state?.token
  }
  return null
}

export const api = {
  // Pages
  async getPages() {
    const token = getToken()
    const response = await retryFetch(`${API_URL}/api/pages`, {
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to fetch pages')
    return response.json()
  },

  async getPage(id: string) {
    const token = getToken()
    const response = await retryFetch(`${API_URL}/api/pages/${id}`, {
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to fetch page')
    return response.json()
  },

  async createPage(data: { title: string; content?: any }) {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/pages`, {
      method: 'POST',
      headers: getAuthHeaders(token),
      body: JSON.stringify(data)
    })
    if (!response.ok) throw new Error('Failed to create page')
    return response.json()
  },

  async updatePage(id: string, data: { title?: string; content?: any }) {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/pages/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(token),
      body: JSON.stringify(data)
    })
    if (!response.ok) throw new Error('Failed to update page')
    return response.json()
  },

  async deletePage(id: string) {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/pages/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to delete page')
    return response.json()
  },

  // Image upload
  async uploadFile(file: File) {
    const token = getToken()
    const formData = new FormData()
    formData.append('file', file)
    
    const headers: Record<string, string> = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const response = await fetch(`${API_URL}/api/upload`, {
      method: 'POST',
      headers,
      body: formData
    })
    if (!response.ok) throw new Error('Failed to upload file')
    return response.json()
  },

  // General file upload
  async uploadGeneralFile(file: File, pageId?: string) {
    const token = getToken()
    const formData = new FormData()
    formData.append('file', file)
    if (pageId) {
      formData.append('page_id', pageId)
    }
    
    const headers: Record<string, string> = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const response = await fetch(`${API_URL}/api/upload/file`, {
      method: 'POST',
      headers,
      body: formData
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to upload file')
    }
    return response.json()
  },

  async getFiles(pageId?: string, type?: string, page: number = 1, limit: number = 20) {
    const token = getToken()
    const params = new URLSearchParams()
    if (pageId) params.append('page_id', pageId)
    if (type) params.append('type', type)
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    
    const response = await retryFetch(`${API_URL}/api/files?${params.toString()}`, {
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to fetch files')
    return response.json()
  },

  async getFileMetadata(id: number) {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/files/${id}`, {
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to fetch file metadata')
    return response.json()
  },

  async deleteFile(id: number) {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/files/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(token)
    })
    if (!response.ok) throw new Error('Failed to delete file')
    return response.json()
  }
}