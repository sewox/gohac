import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Edit, Trash2, FileText } from 'lucide-react'
import toast from 'react-hot-toast'
import { pagesAPI } from '../../lib/api'
import ConfirmDialog from '../../components/ConfirmDialog'
import './PageList.css'

interface Page {
  id: string
  slug: string
  title: string
  status: 'draft' | 'published' | 'archived'
  created_at: string
  updated_at: string
}

export default function PageList() {
  const [pages, setPages] = useState<Page[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteConfirm, setDeleteConfirm] = useState<{
    isOpen: boolean
    pageId: string | null
    pageTitle: string
  }>({
    isOpen: false,
    pageId: null,
    pageTitle: '',
  })

  useEffect(() => {
    fetchPages()
  }, [])

  const fetchPages = async () => {
    try {
      setLoading(true)
      const response = await pagesAPI.list()
      setPages(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load pages'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteClick = (id: string, title: string) => {
    setDeleteConfirm({
      isOpen: true,
      pageId: id,
      pageTitle: title,
    })
  }

  const handleDeleteConfirm = async () => {
    if (!deleteConfirm.pageId) return

    const deletePromise = pagesAPI.delete(deleteConfirm.pageId)

    toast.promise(deletePromise, {
      loading: 'Deleting page...',
      success: 'Page deleted successfully',
      error: (err: any) => err.response?.data?.error || 'Failed to delete page',
    })

    try {
      await deletePromise
      setDeleteConfirm({ isOpen: false, pageId: null, pageTitle: '' })
      fetchPages() // Refresh list
    } catch (err) {
      // Error already handled by toast.promise
    }
  }

  const handleDeleteCancel = () => {
    setDeleteConfirm({ isOpen: false, pageId: null, pageTitle: '' })
  }

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'published':
        return 'status-badge published'
      case 'draft':
        return 'status-badge draft'
      case 'archived':
        return 'status-badge archived'
      default:
        return 'status-badge'
    }
  }

  if (loading) {
    return (
      <div className="page-list">
        <div className="loading">Loading pages...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="page-list">
        <div className="error-message">{error}</div>
        <button onClick={fetchPages} className="retry-button">
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="page-list">
      <div className="page-list-header">
        <div className="header-title">
          <FileText size={24} />
          <h1>Pages</h1>
        </div>
        <Link to="/admin/pages/new" className="create-button">
          <Plus size={20} />
          <span>Create New Page</span>
        </Link>
      </div>

      {pages.length === 0 ? (
        <div className="empty-state">
          <FileText size={48} />
          <h2>No pages yet</h2>
          <p>Create your first page to get started</p>
          <Link to="/admin/pages/new" className="create-button">
            <Plus size={20} />
            <span>Create Page</span>
          </Link>
        </div>
      ) : (
        <div className="pages-table-container">
          <table className="pages-table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Slug</th>
                <th>Status</th>
                <th>Updated</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {pages.map((page) => (
                <tr key={page.id}>
                  <td className="title-cell">
                    <strong>{page.title}</strong>
                  </td>
                  <td className="slug-cell">
                    <code>{page.slug}</code>
                  </td>
                  <td>
                    <span className={getStatusBadgeClass(page.status)}>
                      {page.status}
                    </span>
                  </td>
                  <td className="date-cell">
                    {new Date(page.updated_at).toLocaleDateString()}
                  </td>
                  <td className="actions-cell">
                    <Link
                      to={`/admin/pages/${page.id}/edit`}
                      className="action-button edit"
                      title="Edit"
                    >
                      <Edit size={16} />
                    </Link>
                    <button
                      onClick={() => handleDeleteClick(page.id, page.title)}
                      className="action-button delete"
                      title="Delete"
                    >
                      <Trash2 size={16} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <ConfirmDialog
        isOpen={deleteConfirm.isOpen}
        title="Delete Page"
        message={`Are you sure you want to delete "${deleteConfirm.pageTitle}"? This action cannot be undone.`}
        onConfirm={handleDeleteConfirm}
        onCancel={handleDeleteCancel}
        confirmText="Delete"
        cancelText="Cancel"
        variant="danger"
      />
    </div>
  )
}

