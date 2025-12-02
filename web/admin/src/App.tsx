import { useState } from 'react'
import './App.css'

interface LoginFormData {
  email: string
  password: string
}

interface User {
  id: string
  email: string
}

function App() {
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [user, setUser] = useState<User | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Important for cookies
        body: JSON.stringify(formData),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Login failed')
      }

      // Fetch user info after successful login
      const meResponse = await fetch('/api/auth/me', {
        credentials: 'include',
      })
      const meData = await meResponse.json()

      if (meResponse.ok && meData.user) {
        setUser(meData.user)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    // Clear user state
    setUser(null)
    // In a real app, you'd also call a logout endpoint to clear the cookie
  }

  if (user) {
    return (
      <div className="app">
        <div className="container">
          <div className="card">
            <h1>Welcome to Gohac CMS</h1>
            <div className="user-info">
              <p><strong>Email:</strong> {user.email}</p>
              <p><strong>ID:</strong> {user.id}</p>
            </div>
            <button onClick={handleLogout} className="btn btn-secondary">
              Logout
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="app">
      <div className="container">
        <div className="card">
          <h1>Gohac CMS - Admin Login</h1>
          <form onSubmit={handleSubmit}>
            {error && <div className="error">{error}</div>}
            <div className="form-group">
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                value={formData.email}
                onChange={(e) =>
                  setFormData({ ...formData, email: e.target.value })
                }
                required
                placeholder="admin@example.com"
              />
            </div>
            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                value={formData.password}
                onChange={(e) =>
                  setFormData({ ...formData, password: e.target.value })
                }
                required
                placeholder="admin123"
              />
            </div>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </button>
          </form>
          <div className="hint">
            <p>Demo credentials:</p>
            <p>Email: any email</p>
            <p>Password: admin123</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App

