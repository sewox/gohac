import { useAuth } from '../context/AuthContext'
import { LayoutDashboard, User, Mail } from 'lucide-react'
import './Dashboard.css'

export default function Dashboard() {
  const { user } = useAuth()

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
            <h3>Pages</h3>
            <p className="stat-value">0</p>
            <p className="stat-label">Total pages</p>
          </div>
          <div className="stat-card">
            <h3>Drafts</h3>
            <p className="stat-value">0</p>
            <p className="stat-label">Draft pages</p>
          </div>
          <div className="stat-card">
            <h3>Published</h3>
            <p className="stat-value">0</p>
            <p className="stat-label">Published pages</p>
          </div>
        </div>
      </div>
    </div>
  )
}

