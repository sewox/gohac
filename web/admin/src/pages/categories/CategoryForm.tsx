import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Save, ArrowLeft } from 'lucide-react'
import toast from 'react-hot-toast'
import { categoriesAPI } from '../../lib/api'
import '../settings/Settings.css'

interface CategoryFormFields {
  name: string
  slug: string
  description: string
}

export default function CategoryForm() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEdit = !!id

  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(isEdit)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isEdit) {
      fetchCategory()
    } else {
      setFetching(false)
    }
  }, [isEdit, id])

  // Auto-generate slug from name
  useEffect(() => {
    if (!isEdit && name && !slug) {
      const generatedSlug = name
        .toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '')
      setSlug(generatedSlug)
    }
  }, [name, isEdit, slug])

  const fetchCategory = async () => {
    try {
      setFetching(true)
      const response = await categoriesAPI.getById(id!)
      const category = response.data
      setName(category.name)
      setSlug(category.slug)
      setDescription(category.description || '')
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load category'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!name.trim() || !slug.trim()) {
      setError('Name and slug are required.')
      toast.error('Name and slug are required.')
      return
    }

    setLoading(true)

    const categoryData: CategoryFormFields = {
      name: name.trim(),
      slug: slug.trim(),
      description: description.trim(),
    }

    try {
      if (isEdit && id) {
        const updatePromise = categoriesAPI.update(id, categoryData)

        toast.promise(updatePromise, {
          loading: 'Updating category...',
          success: 'Category updated successfully!',
          error: (err: any) => {
            console.error('Category update error:', err)
            return err.response?.data?.error || 'Failed to update category'
          },
        })

        await updatePromise
      } else {
        const createPromise = categoriesAPI.create(categoryData)

        toast.promise(createPromise, {
          loading: 'Creating category...',
          success: 'Category created successfully!',
          error: (err: any) => {
            console.error('Category creation error:', err)
            return err.response?.data?.error || 'Failed to create category'
          },
        })

        await createPromise
      }

      navigate('/admin/categories')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || `Failed to ${isEdit ? 'update' : 'create'} category`
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="settings-page">
        <div className="loading">Loading category...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <div className="settings-header">
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <button
            onClick={() => navigate('/admin/categories')}
            className="back-button"
            type="button"
          >
            <ArrowLeft size={20} />
          </button>
          <h1>{isEdit ? 'Edit Category' : 'Create New Category'}</h1>
        </div>
        <p className="settings-description">
          {isEdit
            ? 'Update category details.'
            : 'Create a new category to organize your blog posts.'}
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
            placeholder="Tech"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="slug">Slug *</label>
          <input
            type="text"
            id="slug"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            required
            placeholder="tech"
            disabled={loading}
          />
          <small>URL-friendly identifier (auto-generated from name)</small>
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="A brief description of this category"
            rows={4}
            disabled={loading}
          />
        </div>

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : (isEdit ? 'Save Category' : 'Create Category')}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

