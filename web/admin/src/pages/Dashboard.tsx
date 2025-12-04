import { useEffect, useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { LayoutDashboard, User, Mail, FileText, BookOpen, Tag, Image as ImageIcon, Users } from 'lucide-react'
import { dashboardAPI } from '../lib/api'
import './Dashboard.css'

interface DashboardStats {
  pages: number
  users: number
  media: number
  posts: number
  categories: number
}

export default function Dashboard() {
  const { user } = useAuth()
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const response = await dashboardAPI.getStats()
      setStats(response.data.stats)
    } catch (err: any) {
      console.error('Failed to load dashboard stats:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <div className="dashboard-title">
          <LayoutDashboard size={24} />
          <h1>Dashboard</h1>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="welcome-card">
          <div className="welcome-icon">
            <User size={48} />
          </div>
          <h2>Welcome back!</h2>
          <p className="welcome-message">
            You are logged in as <strong>{user?.email}</strong>
          </p>
          <div className="user-info-card">
            <div className="user-info-item">
              <Mail size={20} />
              <div>
                <span className="label">Email</span>
                <span className="value">{user?.email}</span>
              </div>
            </div>
            <div className="user-info-item">
              <User size={20} />
              <div>
                <span className="label">User ID</span>
                <span className="value">{user?.id}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="dashboard-stats">
          <div className="stat-card">
            <FileText size={24} />
            <h3>Pages</h3>
            <p className="stat-value">{loading ? '...' : stats?.pages || 0}</p>
            <p className="stat-label">Total pages</p>
          </div>
          <div className="stat-card">
            <BookOpen size={24} />
            <h3>Posts</h3>
            <p className="stat-value">{loading ? '...' : stats?.posts || 0}</p>
            <p className="stat-label">Blog posts</p>
          </div>
          <div className="stat-card">
            <Tag size={24} />
            <h3>Categories</h3>
            <p className="stat-value">{loading ? '...' : stats?.categories || 0}</p>
            <p className="stat-label">Post categories</p>
          </div>
          <div className="stat-card">
            <Users size={24} />
            <h3>Users</h3>
            <p className="stat-value">{loading ? '...' : stats?.users || 0}</p>
            <p className="stat-label">Total users</p>
          </div>
          <div className="stat-card">
            <ImageIcon size={24} />
            <h3>Media</h3>
            <p className="stat-value">{loading ? '...' : stats?.media || 0}</p>
            <p className="stat-label">Media files</p>
          </div>
        </div>
      </div>
    </div>
  )
}

