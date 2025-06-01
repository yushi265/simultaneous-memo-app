const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

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

  // File upload
  async uploadFile(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await fetch(`${API_URL}/api/upload`, {
      method: 'POST',
      body: formData
    })
    if (!response.ok) throw new Error('Failed to upload file')
    return response.json()
  }
}