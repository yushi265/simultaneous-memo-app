import { getAuthHeaders } from './auth-api'
import { Workspace } from './store'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface WorkspaceResponse extends Workspace {
  member_count: number
  user_role: string
}

export interface CreateWorkspaceRequest {
  name: string
  description?: string
}

export interface UpdateWorkspaceRequest {
  name?: string
  description?: string
}

// Helper function to get token from auth store
const getToken = () => {
  if (typeof window !== 'undefined') {
    const authStorage = JSON.parse(localStorage.getItem('auth-storage') || '{}')
    return authStorage.state?.token
  }
  return null
}

export const workspaceApi = {
  async getWorkspaces(): Promise<WorkspaceResponse[]> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces`, {
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to get workspaces')
    }
    
    return response.json()
  },

  async getWorkspace(id: string): Promise<WorkspaceResponse> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces/${id}`, {
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to get workspace')
    }
    
    return response.json()
  },

  async createWorkspace(data: CreateWorkspaceRequest): Promise<WorkspaceResponse> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces`, {
      method: 'POST',
      headers: getAuthHeaders(token),
      body: JSON.stringify(data)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to create workspace')
    }
    
    return response.json()
  },

  async updateWorkspace(id: string, data: UpdateWorkspaceRequest): Promise<WorkspaceResponse> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(token),
      body: JSON.stringify(data)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to update workspace')
    }
    
    return response.json()
  },

  async deleteWorkspace(id: string): Promise<void> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete workspace')
    }
  },

  async switchWorkspace(id: string): Promise<{token: string, workspace: Workspace, message: string}> {
    const token = getToken()
    const response = await fetch(`${API_URL}/api/workspaces/${id}/switch`, {
      method: 'POST',
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to switch workspace')
    }
    
    return response.json()
  }
}