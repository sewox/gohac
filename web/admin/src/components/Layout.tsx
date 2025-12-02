import { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { LogOut, LayoutDashboard, FileText, Settings } from 'lucide-react'
import './Layout.css'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { user, logout } = useAuth()
  const location = useLocation()

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2>Gohac CMS</h2>
        </div>

        <nav className="sidebar-nav">
          <Link
            to="/admin"
            className={`nav-item ${location.pathname === '/admin' ? 'active' : ''}`}
          >
            <LayoutDashboard size={20} />
            <span>Dashboard</span>
          </Link>
          <Link
            to="/admin/pages"
            className={`nav-item ${location.pathname.startsWith('/admin/pages') ? 'active' : ''}`}
          >
            <FileText size={20} />
            <span>Pages</span>
          </Link>
          <Link
            to="/admin/settings"
            className={`nav-item ${location.pathname === '/admin/settings' ? 'active' : ''}`}
          >
            <Settings size={20} />
            <span>Settings</span>
          </Link>
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

