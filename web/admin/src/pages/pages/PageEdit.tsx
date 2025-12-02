import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save, FileText, Search } from 'lucide-react'
import toast from 'react-hot-toast'
import { pagesAPI } from '../../lib/api'
import BlockEditor from '../../components/editor/BlockEditor'
import ImageUpload from '../../components/editor/ImageUpload'
import { Block } from '../../types/block'
import './PageForm.css'
import './PageEdit.css'

type TabType = 'content' | 'seo'

interface PageMeta {
  meta_title?: string
  meta_description?: string
  og_image?: string
  no_index?: boolean
}

interface Page {
  id: string
  slug: string
  title: string
  status: 'draft' | 'published' | 'archived'
  blocks?: Block[]
  meta?: PageMeta | string
}

export default function PageEdit() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<TabType>('content')
  const [formData, setFormData] = useState({
    slug: '',
    title: '',
    status: 'draft' as 'draft' | 'published' | 'archived',
  })
  const [blocks, setBlocks] = useState<Block[]>([])
  const [meta, setMeta] = useState<PageMeta>({
    meta_title: '',
    meta_description: '',
    og_image: '',
    no_index: false,
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

      // Parse blocks from JSON
      if (page.blocks) {
        try {
          let parsedBlocks: Block[] = []
          if (Array.isArray(page.blocks)) {
            parsedBlocks = page.blocks
          } else if (typeof page.blocks === 'string') {
            parsedBlocks = JSON.parse(page.blocks)
          } else {
            parsedBlocks = page.blocks as any
          }
          parsedBlocks = parsedBlocks.map((block) => ({
            id: block.id || `block-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            type: block.type,
            data: block.data || {},
          }))
          setBlocks(parsedBlocks)
        } catch (e) {
          console.error('Failed to parse blocks:', e)
          setBlocks([])
        }
      } else {
        setBlocks([])
      }

      // Parse meta from JSON
      if (page.meta) {
        try {
          let parsedMeta: PageMeta = {}
          if (typeof page.meta === 'string') {
            parsedMeta = JSON.parse(page.meta)
          } else if (typeof page.meta === 'object') {
            parsedMeta = page.meta as PageMeta
          }
          setMeta({
            meta_title: parsedMeta.meta_title || '',
            meta_description: parsedMeta.meta_description || '',
            og_image: parsedMeta.og_image || '',
            no_index: parsedMeta.no_index || false,
          })
        } catch (e) {
          console.error('Failed to parse meta:', e)
          setMeta({
            meta_title: '',
            meta_description: '',
            og_image: '',
            no_index: false,
          })
        }
      }

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

    // Construct meta JSON object
    const metaData: PageMeta = {}
    if (meta.meta_title) metaData.meta_title = meta.meta_title
    if (meta.meta_description) metaData.meta_description = meta.meta_description
    if (meta.og_image) metaData.og_image = meta.og_image
    if (meta.no_index) metaData.no_index = meta.no_index

    const updateData = {
      ...formData,
      blocks,
      meta: Object.keys(metaData).length > 0 ? metaData : null,
    }

    const updatePromise = pagesAPI.update(id!, updateData)

    toast.promise(updatePromise, {
      loading: 'Updating page...',
      success: 'Page updated successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to update page',
    })

    try {
      await updatePromise
      navigate('/admin/pages')
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
        <button onClick={() => navigate('/admin/pages')} className="back-button">
          <ArrowLeft size={20} />
          <span>Back to Pages</span>
        </button>
        <h1>Edit Page</h1>
      </div>

      {/* Tab Navigation */}
      <div className="tab-navigation">
        <button
          type="button"
          className={`tab-button ${activeTab === 'content' ? 'active' : ''}`}
          onClick={() => setActiveTab('content')}
        >
          <FileText size={18} />
          <span>Content</span>
        </button>
        <button
          type="button"
          className={`tab-button ${activeTab === 'seo' ? 'active' : ''}`}
          onClick={() => setActiveTab('seo')}
        >
          <Search size={18} />
          <span>SEO & Settings</span>
        </button>
      </div>

      <form onSubmit={handleSubmit} className="form">
        {error && <div className="error-message">{error}</div>}

        {/* Basic Fields - Always Visible */}
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

        {/* Tab Content */}
        {activeTab === 'content' && (
          <div className="form-group">
            <label>Content Blocks</label>
            <BlockEditor initialBlocks={blocks} onChange={setBlocks} />
          </div>
        )}

        {activeTab === 'seo' && (
          <div className="seo-tab-content">
            <div className="form-group">
              <label htmlFor="meta_title">Meta Title</label>
              <input
                type="text"
                id="meta_title"
                value={meta.meta_title || ''}
                onChange={(e) => setMeta({ ...meta, meta_title: e.target.value })}
                placeholder={formData.title || 'Page Title (fallback)'}
                disabled={loading}
              />
              <small>
                SEO title for search engines. If empty, page title will be used.
              </small>
            </div>

            <div className="form-group">
              <label htmlFor="meta_description">Meta Description</label>
              <textarea
                id="meta_description"
                value={meta.meta_description || ''}
                onChange={(e) => setMeta({ ...meta, meta_description: e.target.value })}
                placeholder="A brief description of this page for search engines"
                rows={4}
                disabled={loading}
              />
              <small>Recommended: 150-160 characters</small>
            </div>

            <div className="form-group">
              <label htmlFor="og_image">Open Graph Image</label>
              <ImageUpload
                value={meta.og_image}
                onChange={(url) => setMeta({ ...meta, og_image: url })}
                label="Social Media Preview Image"
              />
              <small>
                Image shown when sharing on social media (Facebook, Twitter, etc.)
              </small>
            </div>

            <div className="form-group">
              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={meta.no_index || false}
                  onChange={(e) => setMeta({ ...meta, no_index: e.target.checked })}
                  disabled={loading}
                />
                <span>Hide from search engines (noindex)</span>
              </label>
              <small>
                When enabled, search engines will not index this page
              </small>
            </div>
          </div>
        )}

        <div className="form-actions">
          <button
            type="button"
            onClick={() => navigate('/admin/pages')}
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
