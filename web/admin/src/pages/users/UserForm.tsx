import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Save, ArrowLeft } from 'lucide-react'
import toast from 'react-hot-toast'
import { usersAPI } from '../../lib/api'
import '../settings/Settings.css'

export default function UserForm() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEdit = !!id

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<'admin' | 'editor'>('editor')
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isEdit && id) {
      fetchUser()
    } else {
      setFetching(false)
    }
  }, [isEdit, id])

  const fetchUser = async () => {
    try {
      setFetching(true)
      const response = await usersAPI.get(id!)
      const user = response.data
      setName(user.name)
      setEmail(user.email)
      setRole(user.role)
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load user'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!name.trim() || !email.trim()) {
      setError('Name and Email are required.')
      toast.error('Name and Email are required.')
      return
    }

    if (!isEdit && !password.trim()) {
      setError('Password is required for new users.')
      toast.error('Password is required for new users.')
      return
    }

    if (password.trim() && password.length < 6) {
      setError('Password must be at least 6 characters.')
      toast.error('Password must be at least 6 characters.')
      return
    }

    setLoading(true)

    try {
      const data: any = {
        name,
        email,
        role,
      }

      if (password.trim()) {
        data.password = password
      }

      if (isEdit && id) {
        const updatePromise = usersAPI.update(id, data)

        toast.promise(updatePromise, {
          loading: 'Updating user...',
          success: 'User updated successfully!',
          error: (err: any) => err.response?.data?.error || 'Failed to update user',
        })

        await updatePromise
      } else {
        const createPromise = usersAPI.create(data)

        toast.promise(createPromise, {
          loading: 'Creating user...',
          success: 'User created successfully!',
          error: (err: any) => err.response?.data?.error || 'Failed to create user',
        })

        await createPromise
      }

      navigate('/admin/users')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || `Failed to ${isEdit ? 'update' : 'create'} user`
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="settings-page">
        <div className="loading">Loading user...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <div className="settings-header">
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <button
            onClick={() => navigate('/admin/users')}
            className="back-button"
            type="button"
          >
            <ArrowLeft size={20} />
          </button>
          <h1>{isEdit ? 'Edit User' : 'Add User'}</h1>
        </div>
        <p className="settings-description">
          {isEdit ? 'Update user information and permissions.' : 'Create a new user account.'}
        </p>
      </div>

      <form onSubmit={handleSubmit} className="settings-form">
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="name">Name *</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            placeholder="John Doe"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="email">Email *</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            placeholder="user@example.com"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">
            Password {isEdit ? '(leave empty to keep current)' : '*'}
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required={!isEdit}
            placeholder={isEdit ? 'Leave empty to keep current password' : 'Minimum 6 characters'}
            disabled={loading}
            minLength={isEdit ? 0 : 6}
          />
          <small>Password must be at least 6 characters long</small>
        </div>

        <div className="form-group">
          <label htmlFor="role">Role *</label>
          <select
            id="role"
            value={role}
            onChange={(e) => setRole(e.target.value as 'admin' | 'editor')}
            disabled={loading}
          >
            <option value="editor">Editor</option>
            <option value="admin">Admin</option>
          </select>
          <small>Admin users can manage other users and system settings</small>
        </div>

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : (isEdit ? 'Save User' : 'Create User')}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

