import { create } from 'zustand'

export interface Page {
  id: number
  title: string
  content: any
  created_at: string
  updated_at: string
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
  updatePage: (id: number, updates: Partial<Page>) => void
  deletePage: (id: number) => void
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