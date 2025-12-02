import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, Save } from 'lucide-react'
import toast from 'react-hot-toast'
import { pagesAPI } from '../../lib/api'
import './PageForm.css'

export default function PageForm() {
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    slug: '',
    title: '',
    status: 'draft' as 'draft' | 'published' | 'archived',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    const createPromise = pagesAPI.create(formData)

    toast.promise(createPromise, {
      loading: 'Creating page...',
      success: 'Page created successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to create page',
    })

    try {
      await createPromise
      navigate('/pages')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to create page'
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page-form">
      <div className="page-form-header">
        <button onClick={() => navigate('/pages')} className="back-button">
          <ArrowLeft size={20} />
          <span>Back to Pages</span>
        </button>
        <h1>Create New Page</h1>
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
            <span>{loading ? 'Creating...' : 'Create Page'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

