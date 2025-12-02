import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { authAPI } from '../lib/api'

interface User {
  id: string
  email: string
}

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  // Check if user is logged in on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await authAPI.me()
        if (response.data.user) {
          setUser(response.data.user)
        }
      } catch (error) {
        // User is not authenticated
        setUser(null)
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [])

  const login = async (email: string, password: string) => {
    try {
      const response = await authAPI.login(email, password)
      if (response.data.success) {
        // Fetch user info after successful login
        const meResponse = await authAPI.me()
        if (meResponse.data.user) {
          setUser(meResponse.data.user)
        }
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed')
    }
  }

  const logout = () => {
    setUser(null)
    // In a real app, you'd call authAPI.logout() here
    window.location.href = '/login'
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

