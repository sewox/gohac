import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Plus, Edit, Trash2 } from 'lucide-react'
import toast from 'react-hot-toast'
import { menusAPI } from '../../lib/api'
import ConfirmDialog from '../../components/ConfirmDialog'
import '../pages/PageList.css'

interface Menu {
  id: string
  name: string
  description?: string
  items: Array<{
    label: string
    url: string
  }>
  created_at: string
  updated_at: string
}

export default function MenuList() {
  const navigate = useNavigate()
  const [menus, setMenus] = useState<Menu[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; menu: Menu | null }>({
    open: false,
    menu: null,
  })

  useEffect(() => {
    fetchMenus()
  }, [])

  const fetchMenus = async () => {
    try {
      setLoading(true)
      const response = await menusAPI.list()
      setMenus(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load menus'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (menu: Menu) => {
    const deletePromise = menusAPI.delete(menu.id)

    toast.promise(deletePromise, {
      loading: 'Deleting menu...',
      success: 'Menu deleted successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to delete menu',
    })

    try {
      await deletePromise
      setDeleteDialog({ open: false, menu: null })
      fetchMenus()
    } catch (err: any) {
      // Error handled by toast
    }
  }

  if (loading) {
    return (
      <div className="page-list">
        <div className="loading">Loading menus...</div>
      </div>
    )
  }

  return (
    <div className="page-list">
      <div className="page-list-header">
        <h1>Menus</h1>
        <Link to="/admin/menus/new" className="create-button">
          <Plus size={18} />
          <span>Create New Menu</span>
        </Link>
      </div>

      {error && <div className="error-message">{error}</div>}

      {menus.length === 0 ? (
        <div className="empty-state">
          <p>No menus yet. Create your first menu to get started.</p>
          <Link to="/admin/menus/new" className="create-button">
            <Plus size={18} />
            <span>Create New Menu</span>
          </Link>
        </div>
      ) : (
        <div className="table-container">
          <table className="data-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Items</th>
                <th>Updated</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {menus.map((menu) => (
                <tr key={menu.id}>
                  <td>
                    <strong>{menu.name}</strong>
                  </td>
                  <td>{menu.description || '-'}</td>
                  <td>{menu.items?.length || 0} items</td>
                  <td>{new Date(menu.updated_at).toLocaleDateString()}</td>
                  <td>
                    <div className="action-buttons">
                      <button
                        onClick={() => navigate(`/admin/menus/${menu.id}/edit`)}
                        className="action-button edit"
                        title="Edit"
                      >
                        <Edit size={16} />
                      </button>
                      <button
                        onClick={() => setDeleteDialog({ open: true, menu })}
                        className="action-button delete"
                        title="Delete"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <ConfirmDialog
        isOpen={deleteDialog.open}
        title="Delete Menu"
        message={`Are you sure you want to delete "${deleteDialog.menu?.name}"? This action cannot be undone.`}
        variant="danger"
        onConfirm={() => deleteDialog.menu && handleDelete(deleteDialog.menu)}
        onCancel={() => setDeleteDialog({ open: false, menu: null })}
      />
    </div>
  )
}

