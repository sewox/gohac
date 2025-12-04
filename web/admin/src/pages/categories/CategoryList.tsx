import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Plus, Edit, Trash2, Tag } from 'lucide-react'
import toast from 'react-hot-toast'
import { categoriesAPI } from '../../lib/api'
import ConfirmDialog from '../../components/ConfirmDialog'
import '../pages/PageList.css'

interface Category {
  id: string
  name: string
  slug: string
  description: string
  created_at: string
  updated_at: string
}

export default function CategoryList() {
  const navigate = useNavigate()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; category: Category | null }>({
    open: false,
    category: null,
  })

  useEffect(() => {
    fetchCategories()
  }, [])

  const fetchCategories = async () => {
    try {
      setLoading(true)
      const response = await categoriesAPI.list()
      setCategories(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load categories'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (category: Category) => {
    if (!category.id) return

    const deletePromise = categoriesAPI.delete(category.id)

    toast.promise(deletePromise, {
      loading: 'Deleting category...',
      success: 'Category deleted successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to delete category',
    })

    try {
      await deletePromise
      setDeleteDialog({ open: false, category: null })
      fetchCategories() // Refresh the list
    } catch (err: any) {
      console.error('Failed to delete category:', err)
    }
  }

  if (loading) {
    return (
      <div className="page-list">
        <div className="loading">Loading categories...</div>
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
          <h1>Categories</h1>
        </div>
        <Link to="/admin/categories/new" className="create-button">
          <Plus size={20} />
          <span>New Category</span>
        </Link>
      </div>

      {categories.length === 0 ? (
        <div className="empty-state">
          <p>No categories yet. Create your first category to get started.</p>
          <Link to="/admin/categories/new" className="create-button">
            <Plus size={20} />
            <span>New Category</span>
          </Link>
        </div>
      ) : (
        <div className="page-list-table-container">
          <table className="page-list-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Slug</th>
                <th>Description</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {categories.map((category) => (
                <tr key={category.id}>
                  <td>
                    <strong>{category.name}</strong>
                  </td>
                  <td>
                    <code style={{ fontSize: '12px' }}>{category.slug}</code>
                  </td>
                  <td>{category.description || '-'}</td>
                  <td>{new Date(category.created_at).toLocaleDateString()}</td>
                  <td className="actions">
                    <button
                      onClick={() => navigate(`/admin/categories/${category.id}/edit`)}
                      className="action-button edit"
                      title="Edit Category"
                    >
                      <Edit size={18} />
                    </button>
                    <button
                      onClick={() => setDeleteDialog({ open: true, category })}
                      className="action-button delete"
                      title="Delete Category"
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
        title="Delete Category"
        message={`Are you sure you want to delete "${deleteDialog.category?.name}"? This action cannot be undone.`}
        variant="danger"
        onConfirm={() => deleteDialog.category && handleDelete(deleteDialog.category)}
        onCancel={() => setDeleteDialog({ open: false, category: null })}
      />
    </div>
  )
}

