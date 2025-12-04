import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Layout from './components/Layout'
import RequireAuth from './components/RequireAuth'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import PageList from './pages/pages/PageList'
import PageForm from './pages/pages/PageForm'
import PageEdit from './pages/pages/PageEdit'
import GeneralSettings from './pages/settings/GeneralSettings'
import MenuList from './pages/menus/MenuList'
import MenuForm from './pages/menus/MenuForm'
import UserList from './pages/users/UserList'
import UserForm from './pages/users/UserForm'
import MediaLibrary from './pages/media/MediaLibrary'
import Profile from './pages/profile/Profile'
import PostList from './pages/posts/PostList'
import PostForm from './pages/posts/PostForm'
import CategoryList from './pages/categories/CategoryList'
import CategoryForm from './pages/categories/CategoryForm'
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

          <Route
            path="/admin/menus"
            element={
              <RequireAuth>
                <Layout>
                  <MenuList />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/menus/new"
            element={
              <RequireAuth>
                <Layout>
                  <MenuForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/menus/:id/edit"
            element={
              <RequireAuth>
                <Layout>
                  <MenuForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/users"
            element={
              <RequireAuth>
                <Layout>
                  <UserList />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/users/new"
            element={
              <RequireAuth>
                <Layout>
                  <UserForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/users/:id/edit"
            element={
              <RequireAuth>
                <Layout>
                  <UserForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/media"
            element={
              <RequireAuth>
                <Layout>
                  <MediaLibrary />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/settings"
            element={
              <RequireAuth>
                <Layout>
                  <GeneralSettings />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/profile"
            element={
              <RequireAuth>
                <Layout>
                  <Profile />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/posts"
            element={
              <RequireAuth>
                <Layout>
                  <PostList />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/posts/new"
            element={
              <RequireAuth>
                <Layout>
                  <PostForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/posts/:id/edit"
            element={
              <RequireAuth>
                <Layout>
                  <PostForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/categories"
            element={
              <RequireAuth>
                <Layout>
                  <CategoryList />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/categories/new"
            element={
              <RequireAuth>
                <Layout>
                  <CategoryForm />
                </Layout>
              </RequireAuth>
            }
          />

          <Route
            path="/admin/categories/:id/edit"
            element={
              <RequireAuth>
                <Layout>
                  <CategoryForm />
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
