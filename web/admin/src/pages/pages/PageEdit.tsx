import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save } from 'lucide-react'
import toast from 'react-hot-toast'
import { pagesAPI } from '../../lib/api'
import './PageForm.css'

interface Page {
  id: string
  slug: string
  title: string
  status: 'draft' | 'published' | 'archived'
}

export default function PageEdit() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    slug: '',
    title: '',
    status: 'draft' as 'draft' | 'published' | 'archived',
  })
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      fetchPage()
    }
  }, [id])

  const fetchPage = async () => {
    try {
      setFetching(true)
      const response = await pagesAPI.getById(id!)
      const page = response.data
      setFormData({
        slug: page.slug,
        title: page.title,
        status: page.status,
      })
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load page'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    const updatePromise = pagesAPI.update(id!, formData)

    toast.promise(updatePromise, {
      loading: 'Updating page...',
      success: 'Page updated successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to update page',
    })

    try {
      await updatePromise
      navigate('/pages')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to update page'
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="page-form">
        <div className="loading">Loading page...</div>
      </div>
    )
  }

  return (
    <div className="page-form">
      <div className="page-form-header">
        <button onClick={() => navigate('/pages')} className="back-button">
          <ArrowLeft size={20} />
          <span>Back to Pages</span>
        </button>
        <h1>Edit Page</h1>
      </div>

      <form onSubmit={handleSubmit} className="form">
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="title">Title *</label>
          <input
            type="text"
            id="title"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            required
            placeholder="My Page Title"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="slug">Slug *</label>
          <input
            type="text"
            id="slug"
            value={formData.slug}
            onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
            required
            placeholder="my-page-slug"
            disabled={loading}
          />
          <small>URL-friendly identifier (e.g., "about-us")</small>
        </div>

        <div className="form-group">
          <label htmlFor="status">Status</label>
          <select
            id="status"
            value={formData.status}
            onChange={(e) =>
              setFormData({
                ...formData,
                status: e.target.value as 'draft' | 'published' | 'archived',
              })
            }
            disabled={loading}
          >
            <option value="draft">Draft</option>
            <option value="published">Published</option>
            <option value="archived">Archived</option>
          </select>
        </div>

        <div className="form-actions">
          <button
            type="button"
            onClick={() => navigate('/pages')}
            className="cancel-button"
            disabled={loading}
          >
            Cancel
          </button>
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Updating...' : 'Update Page'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

