import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

export interface User {
  id: string
  email: string
  name: string
  avatar_url: string
  created_at: string
  updated_at: string
}

export interface Workspace {
  id: string
  name: string
  slug: string
  description?: string
  is_personal: boolean
  owner_id: string
  created_at: string
  updated_at: string
}

export interface Page {
  id: string
  workspace_id: string
  title: string
  content: any
  created_by: string
  last_edited_by: string
  is_public: boolean
  created_at: string
  updated_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  currentWorkspace: Workspace | null
  workspaces: Array<{id: string, name: string, role: string}>
  isAuthenticated: boolean
  isLoading: boolean
  
  // Actions
  login: (token: string, user: User, workspace: Workspace) => void
  logout: () => void
  setCurrentWorkspace: (workspace: Workspace) => void
  setUser: (user: User) => void
  setWorkspaces: (workspaces: Array<{id: string, name: string, role: string}>) => void
}

interface AppState {
  pages: Page[]
  currentPage: Page | null
  isLoading: boolean
  error: string | null
  
  // Actions
  setPages: (pages: Page[]) => void
  setCurrentPage: (page: Page | null) => void
  addPage: (page: Page) => void
  updatePage: (id: string, updates: Partial<Page>) => void
  deletePage: (id: string) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
}

export const useStore = create<AppState>((set) => ({
  pages: [],
  currentPage: null,
  isLoading: false,
  error: null,
  
  setPages: (pages) => set({ pages }),
  setCurrentPage: (page) => set({ currentPage: page }),
  addPage: (page) => set((state) => ({ pages: [page, ...state.pages] })),
  updatePage: (id, updates) => set((state) => ({
    pages: state.pages.map(p => p.id === id ? { ...p, ...updates } : p),
    currentPage: state.currentPage?.id === id ? { ...state.currentPage, ...updates } : state.currentPage
  })),
  deletePage: (id) => set((state) => ({
    pages: state.pages.filter(p => p.id !== id),
    currentPage: state.currentPage?.id === id ? null : state.currentPage
  })),
  setLoading: (loading) => set({ isLoading: loading }),
  setError: (error) => set({ error })
}))

// Auth store with persistence
export const useAuthStore = create<AuthState>()(persist(
  (set) => ({
    user: null,
    token: null,
    currentWorkspace: null,
    workspaces: [],
    isAuthenticated: false,
    isLoading: false,
    
    login: (token, user, workspace) => set({
      token,
      user,
      currentWorkspace: workspace,
      isAuthenticated: true,
      workspaces: [{ id: workspace.id, name: workspace.name, role: 'owner' }]
    }),
    
    logout: () => set({
      user: null,
      token: null,
      currentWorkspace: null,
      workspaces: [],
      isAuthenticated: false
    }),
    
    setCurrentWorkspace: (workspace) => set({ currentWorkspace: workspace }),
    
    setUser: (user) => set({ user }),
    
    setWorkspaces: (workspaces) => set({ workspaces })
  }),
  {
    name: 'auth-storage',
    storage: createJSONStorage(() => {
      // Only use localStorage on client side
      if (typeof window !== 'undefined') {
        return localStorage
      }
      // Return a no-op storage for SSR
      return {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {}
      }
    }),
    partialize: (state) => ({
      user: state.user,
      token: state.token,
      currentWorkspace: state.currentWorkspace,
      workspaces: state.workspaces,
      isAuthenticated: state.isAuthenticated
    })
  }
))