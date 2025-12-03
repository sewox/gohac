import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Plus, Edit, Trash2 } from 'lucide-react'
import toast from 'react-hot-toast'
import { usersAPI } from '../../lib/api'
import ConfirmDialog from '../../components/ConfirmDialog'
import '../../pages/pages/PageList.css'

interface User {
  id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export default function UserList() {
  const navigate = useNavigate()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; user: User | null }>({
    open: false,
    user: null,
  })

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const response = await usersAPI.list()
      setUsers(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load users'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (user: User) => {
    if (!user.id) return

    const deletePromise = usersAPI.delete(user.id)

    toast.promise(deletePromise, {
      loading: 'Deleting user...',
      success: 'User deleted successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to delete user',
    })

    try {
      await deletePromise
      setDeleteDialog({ open: false, user: null })
      fetchUsers() // Refresh the list
    } catch (err: any) {
      console.error('Failed to delete user:', err)
    }
  }

  if (loading) {
    return (
      <div className="page-list">
        <div className="loading">Loading users...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="page-list">
        <div className="error-message">{error}</div>
      </div>
    )
  }

  return (
    <div className="page-list">
      <div className="page-list-header">
        <div className="header-title">
          <h1>Users</h1>
        </div>
        <Link to="/admin/users/new" className="create-button">
          <Plus size={20} />
          <span>Add User</span>
        </Link>
      </div>

      {users.length === 0 ? (
        <div className="empty-state">
          <p>No users yet. Create your first user to get started.</p>
          <Link to="/admin/users/new" className="create-button">
            <Plus size={20} />
            <span>Add User</span>
          </Link>
        </div>
      ) : (
        <div className="page-list-table-container">
          <table className="page-list-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.name}</td>
                  <td>{user.email}</td>
                  <td>
                    <span style={{
                      display: 'inline-block',
                      padding: '0.25rem 0.75rem',
                      borderRadius: '12px',
                      fontSize: '0.75rem',
                      fontWeight: 600,
                      textTransform: 'uppercase',
                      backgroundColor: user.role === 'admin' ? '#fed7d7' : '#bee3f8',
                      color: user.role === 'admin' ? '#c53030' : '#2c5282'
                    }}>
                      {user.role}
                    </span>
                  </td>
                  <td>{new Date(user.created_at).toLocaleDateString()}</td>
                  <td className="actions">
                    <button
                      onClick={() => navigate(`/admin/users/${user.id}/edit`)}
                      className="action-button edit"
                      title="Edit User"
                    >
                      <Edit size={18} />
                    </button>
                    <button
                      onClick={() => setDeleteDialog({ open: true, user })}
                      className="action-button delete"
                      title="Delete User"
                    >
                      <Trash2 size={18} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <ConfirmDialog
        isOpen={deleteDialog.open}
        title="Delete User"
        message={`Are you sure you want to delete "${deleteDialog.user?.name}"? This action cannot be undone.`}
        variant="danger"
        onConfirm={() => deleteDialog.user && handleDelete(deleteDialog.user)}
        onCancel={() => setDeleteDialog({ open: false, user: null })}
      />
    </div>
  )
}

