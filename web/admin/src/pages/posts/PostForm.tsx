import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save, FileText } from 'lucide-react'
import toast from 'react-hot-toast'
import { postsAPI, categoriesAPI } from '../../lib/api'
import BlockEditor from '../../components/editor/BlockEditor'
import ImageUpload from '../../components/editor/ImageUpload'
import { Block } from '../../types/block'
import '../pages/PageForm.css'

interface Post {
  id: string
  slug: string
  title: string
  excerpt: string
  content: string
  featured_image: string
  status: 'draft' | 'published' | 'archived'
  category_ids?: string[]
  categories?: Array<{ id: string; name: string }>
}

interface Category {
  id: string
  name: string
  slug: string
}

export default function PostForm() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEdit = !!id

  const [formData, setFormData] = useState({
    slug: '',
    title: '',
    excerpt: '',
    featured_image: '',
    status: 'draft' as 'draft' | 'published' | 'archived',
  })
  const [blocks, setBlocks] = useState<Block[]>([])
  const [selectedCategories, setSelectedCategories] = useState<string[]>([])
  const [availableCategories, setAvailableCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(isEdit)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCategories()
    if (isEdit && id) {
      fetchPost()
    }
  }, [isEdit, id])

  const fetchCategories = async () => {
    try {
      const response = await categoriesAPI.list()
      setAvailableCategories(response.data.data || [])
    } catch (err: any) {
      console.error('Failed to load categories:', err)
      // Don't show error toast for empty categories, it's normal
      if (err.response?.status !== 401) {
        toast.error('Failed to load categories')
      }
    }
  }

  const fetchPost = async () => {
    try {
      setFetching(true)
      const response = await postsAPI.getById(id!)
      const post = response.data

      setFormData({
        slug: post.slug,
        title: post.title,
        excerpt: post.excerpt || '',
        featured_image: post.featured_image || '',
        status: post.status,
      })

      // Parse content blocks from JSON
      if (post.content) {
        try {
          let parsedBlocks: Block[] = []
          if (typeof post.content === 'string') {
            parsedBlocks = JSON.parse(post.content)
          } else if (Array.isArray(post.content)) {
            parsedBlocks = post.content
          }
          parsedBlocks = parsedBlocks.map((block) => ({
            id: block.id || `block-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            type: block.type,
            data: block.data || {},
          }))
          setBlocks(parsedBlocks)
        } catch (e) {
          console.error('Failed to parse content blocks:', e)
          setBlocks([])
        }
      } else {
        setBlocks([])
      }

      // Set selected categories
      if (post.categories && Array.isArray(post.categories)) {
        setSelectedCategories(post.categories.map((cat: any) => cat.id))
      } else if (post.category_ids && Array.isArray(post.category_ids)) {
        setSelectedCategories(post.category_ids)
      }

      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load post'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!formData.title.trim() || !formData.slug.trim()) {
      setError('Title and slug are required.')
      toast.error('Title and slug are required.')
      return
    }

    setLoading(true)

    const postData: any = {
      ...formData,
      content: blocks.length > 0 ? JSON.stringify(blocks) : '', // Convert blocks to JSON string, empty string if no blocks
      category_ids: selectedCategories.length > 0 ? selectedCategories : [], // Ensure it's always an array
    }

    try {
      if (isEdit && id) {
        const updatePromise = postsAPI.update(id, postData)

        toast.promise(updatePromise, {
          loading: 'Updating post...',
          success: 'Post updated successfully!',
          error: (err: any) => {
            console.error('Post update error:', err)
            return err.response?.data?.error || 'Failed to update post'
          },
        })

        await updatePromise
      } else {
        const createPromise = postsAPI.create(postData)

        toast.promise(createPromise, {
          loading: 'Creating post...',
          success: 'Post created successfully!',
          error: (err: any) => {
            console.error('Post creation error:', err)
            return err.response?.data?.error || 'Failed to create post'
          },
        })

        await createPromise
      }

      navigate('/admin/posts')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || `Failed to ${isEdit ? 'update' : 'create'} post`
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="page-form">
        <div className="loading">Loading post...</div>
      </div>
    )
  }

  return (
    <div className="page-form">
      <div className="page-form-header">
        <button onClick={() => navigate('/admin/posts')} className="back-button">
          <ArrowLeft size={20} />
          <span>Back to Posts</span>
        </button>
        <h1>{isEdit ? 'Edit Post' : 'Create New Post'}</h1>
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
            placeholder="My Blog Post Title"
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
            placeholder="my-blog-post-slug"
            disabled={loading}
          />
          <small>URL-friendly identifier (e.g., "hello-world")</small>
        </div>

        <div className="form-group">
          <label htmlFor="excerpt">Excerpt</label>
          <textarea
            id="excerpt"
            value={formData.excerpt}
            onChange={(e) => setFormData({ ...formData, excerpt: e.target.value })}
            placeholder="A brief summary of this post"
            rows={3}
            disabled={loading}
          />
          <small>Short description shown in post listings</small>
        </div>

        <div className="form-group">
          <label htmlFor="featured_image">Featured Image</label>
          <ImageUpload
            value={formData.featured_image}
            onChange={(url) => setFormData({ ...formData, featured_image: url })}
            label="Featured Image"
          />
          <small>Main image displayed with this post</small>
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

        <div className="form-group">
          <label htmlFor="categories">Categories</label>
          {availableCategories.length > 0 ? (
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginTop: '8px' }}>
              {availableCategories.map((category) => (
                <label
                  key={category.id}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    padding: '6px 12px',
                    border: selectedCategories.includes(category.id)
                      ? '2px solid #4299e1'
                      : '2px solid #e2e8f0',
                    borderRadius: '6px',
                    cursor: 'pointer',
                    backgroundColor: selectedCategories.includes(category.id) ? '#ebf8ff' : '#fff',
                  }}
                >
                  <input
                    type="checkbox"
                    checked={selectedCategories.includes(category.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedCategories([...selectedCategories, category.id])
                      } else {
                        setSelectedCategories(selectedCategories.filter((id) => id !== category.id))
                      }
                    }}
                    disabled={loading}
                    style={{ marginRight: '6px' }}
                  />
                  <span>{category.name}</span>
                </label>
              ))}
            </div>
          ) : (
            <div style={{ marginTop: '8px' }}>
              <small style={{ color: '#999', display: 'block', marginBottom: '8px' }}>
                No categories available. Create categories first.
              </small>
              <button
                type="button"
                onClick={() => navigate('/admin/categories/new')}
                style={{
                  padding: '6px 12px',
                  backgroundColor: '#4299e1',
                  color: '#fff',
                  border: 'none',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontSize: '14px',
                }}
              >
                Create Category
              </button>
            </div>
          )}
        </div>

        <div className="form-group">
          <label>Content Blocks</label>
          <BlockEditor initialBlocks={blocks} onChange={setBlocks} />
        </div>

        <div className="form-actions">
          <button
            type="button"
            onClick={() => navigate('/admin/posts')}
            className="cancel-button"
            disabled={loading}
          >
            Cancel
          </button>
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : isEdit ? 'Update Post' : 'Create Post'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

