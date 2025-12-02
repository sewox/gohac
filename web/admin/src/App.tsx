import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Layout from './components/Layout'
import RequireAuth from './components/RequireAuth'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import PageList from './pages/pages/PageList'
import PageForm from './pages/pages/PageForm'
import PageEdit from './pages/pages/PageEdit'
import './App.css'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public route */}
          <Route path="/admin/login" element={<Login />} />

          {/* Protected routes */}
          <Route
            path="/admin"
            element={
              <RequireAuth>
                <Layout>
                  <Dashboard />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/pages"
            element={
              <RequireAuth>
                <Layout>
                  <PageList />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/pages/new"
            element={
              <RequireAuth>
                <Layout>
                  <PageForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/pages/:id/edit"
            element={
              <RequireAuth>
                <Layout>
                  <PageEdit />
                </Layout>
              </RequireAuth>
            }
          />

          {/* Redirect root to admin */}
          <Route path="/" element={<Navigate to="/admin" replace />} />

          {/* Redirect unknown routes to admin */}
          <Route path="*" element={<Navigate to="/admin" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
