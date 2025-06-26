import { User, Workspace } from './store'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// Retry configuration for authentication calls
const RETRY_ATTEMPTS = 3
const RETRY_DELAY = 2000 // 2 seconds for auth calls

// Exponential backoff retry logic for auth calls
async function retryAuthFetch(url: string, options: RequestInit, attempts: number = RETRY_ATTEMPTS): Promise<Response> {
  try {
    const response = await fetch(url, options)
    
    // If 429 error, retry with exponential backoff
    if (response.status === 429 && attempts > 0) {
      const delay = RETRY_DELAY * (RETRY_ATTEMPTS - attempts + 1)
      console.log(`Auth request rate limited, retrying in ${delay}ms... (${attempts} attempts left)`)
      await new Promise(resolve => setTimeout(resolve, delay))
      return retryAuthFetch(url, options, attempts - 1)
    }
    
    return response
  } catch (error) {
    if (attempts > 0) {
      const delay = RETRY_DELAY * (RETRY_ATTEMPTS - attempts + 1)
      console.log(`Auth request failed, retrying in ${delay}ms... (${attempts} attempts left)`)
      await new Promise(resolve => setTimeout(resolve, delay))
      return retryAuthFetch(url, options, attempts - 1)
    }
    throw error
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  token: string
  user: User
  workspace: Workspace
}

export interface UserWithWorkspaces {
  user: User
  currentWorkspace: Workspace
  workspaces: Array<{id: string, name: string, role: string}>
}

// Helper function to get auth headers
export const getAuthHeaders = (token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  return headers
}

export const authApi = {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_URL}/api/auth/register`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Registration failed')
    }
    
    return response.json()
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Login failed')
    }
    
    return response.json()
  },

  async me(token: string): Promise<UserWithWorkspaces> {
    const response = await retryAuthFetch(`${API_URL}/api/auth/me`, {
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      // Don't treat 429 as an authentication failure if we've exhausted retries
      if (response.status === 429) {
        console.warn('Authentication check failed due to rate limiting')
        throw new Error('Rate limited - please try again in a moment')
      }
      
      const error = await response.json()
      throw new Error(error.error || 'Failed to get user info')
    }
    
    return response.json()
  },

  async logout(token: string): Promise<void> {
    const response = await fetch(`${API_URL}/api/auth/logout`, {
      method: 'POST',
      headers: getAuthHeaders(token)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Logout failed')
    }
  }
}