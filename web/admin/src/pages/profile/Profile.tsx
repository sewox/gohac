import { useState, useEffect } from 'react'
import { Save } from 'lucide-react'
import toast from 'react-hot-toast'
import { useAuth } from '../../context/AuthContext'
import { authAPI } from '../../lib/api'
import '../settings/Settings.css'

export default function Profile() {
  const { user, refreshUser } = useAuth()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (user) {
      setName(user.name || '')
      setEmail(user.email || '')
    }
  }, [user])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (password && password.length < 6) {
      setError('Password must be at least 6 characters long.')
      toast.error('Password must be at least 6 characters long.')
      return
    }

    if (password && password !== confirmPassword) {
      setError('Passwords do not match.')
      toast.error('Passwords do not match.')
      return
    }

    setLoading(true)

    try {
      // Prepare update data
      const updateData: { name?: string; password?: string } = {}
      
      if (name.trim() && name !== user?.name) {
        updateData.name = name.trim()
      }
      
      if (password.trim()) {
        updateData.password = password.trim()
      }

      // Only make API call if there are changes
      if (Object.keys(updateData).length === 0) {
        toast('No changes to save', { icon: 'ℹ️' })
        setLoading(false)
        return
      }

      // Update profile via API
      const updatePromise = authAPI.updateProfile(updateData)

      toast.promise(updatePromise, {
        loading: 'Updating profile...',
        success: 'Profile updated successfully!',
        error: (err: any) => {
          console.error('Profile update error:', err)
          return err.response?.data?.error || 'Failed to update profile'
        },
      })

      await updatePromise

      // Clear password fields
      setPassword('')
      setConfirmPassword('')
      
      // Refresh user data
      if (refreshUser) {
        await refreshUser()
      }
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to update profile'
      setError(errorMsg)
      // Error toast is handled by toast.promise
    } finally {
      setLoading(false)
    }
  }

  if (!user) {
    return (
      <div className="settings-page">
        <div className="loading">Loading profile...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <div className="settings-header">
        <h1>My Profile</h1>
        <p className="settings-description">
          Update your personal information and password.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="settings-form">
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Your Name"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            disabled
            style={{ backgroundColor: '#f7fafc', cursor: 'not-allowed' }}
          />
          <small>Email cannot be changed</small>
        </div>

        <div className="form-group">
          <label htmlFor="role">Role</label>
          <input
            type="text"
            id="role"
            value={user.role || 'editor'}
            disabled
            style={{ backgroundColor: '#f7fafc', cursor: 'not-allowed' }}
          />
          <small>Role is managed by administrators</small>
        </div>

        <div className="form-group">
          <label htmlFor="password">New Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Leave empty to keep current password"
            disabled={loading}
            minLength={6}
          />
          <small>Leave empty to keep your current password</small>
        </div>

        {password && (
          <div className="form-group">
            <label htmlFor="confirmPassword">Confirm New Password</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Confirm new password"
              disabled={loading}
              minLength={6}
            />
          </div>
        )}

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : 'Save Changes'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

