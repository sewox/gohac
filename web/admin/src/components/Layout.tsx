import { ReactNode } from 'react'
import { useAuth } from '../context/AuthContext'
import { LogOut, LayoutDashboard, FileText, Settings } from 'lucide-react'
import './Layout.css'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { user, logout } = useAuth()

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2>Gohac CMS</h2>
        </div>

        <nav className="sidebar-nav">
          <a href="/" className="nav-item active">
            <LayoutDashboard size={20} />
            <span>Dashboard</span>
          </a>
          <a href="/pages" className="nav-item">
            <FileText size={20} />
            <span>Pages</span>
          </a>
          <a href="/settings" className="nav-item">
            <Settings size={20} />
            <span>Settings</span>
          </a>
        </nav>

        <div className="sidebar-footer">
          <div className="user-profile">
            <div className="user-avatar">
              {user?.email.charAt(0).toUpperCase()}
            </div>
            <div className="user-details">
              <div className="user-name">{user?.email}</div>
              <div className="user-role">Admin</div>
            </div>
          </div>
          <button onClick={logout} className="logout-button">
            <LogOut size={18} />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      <main className="main-content">
        {children}
      </main>
    </div>
  )
}

