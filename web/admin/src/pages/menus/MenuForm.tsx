import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save, Plus, Trash2, ChevronDown, ChevronUp } from 'lucide-react'
import toast from 'react-hot-toast'
import { menusAPI } from '../../lib/api'
import '../settings/Settings.css'

interface MenuItem {
  label: string
  url: string
  target?: string
  children?: MenuItem[]
}

export default function MenuForm() {
  const navigate = useNavigate()
  const { id } = useParams<{ id: string }>()
  const isEdit = !!id

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [items, setItems] = useState<MenuItem[]>([])
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(isEdit)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isEdit) {
      fetchMenu()
    }
  }, [id])

  const fetchMenu = async () => {
    if (!id) return

    try {
      setFetching(true)
      const response = await menusAPI.get(id)
      setName(response.data.name)
      setDescription(response.data.description || '')
      setItems(response.data.items || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load menu'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const addMenuItem = () => {
    const newItem: MenuItem = {
      label: '',
      url: '',
      target: '_self',
    }
    setItems([...items, newItem])
  }

  const updateMenuItem = (index: number, field: keyof MenuItem, value: string) => {
    const newItems = [...items]
    newItems[index] = { ...newItems[index], [field]: value }
    setItems(newItems)
  }

  const removeMenuItem = (index: number) => {
    const newItems = items.filter((_, i) => i !== index)
    setItems(newItems)
  }

  const moveMenuItem = (index: number, direction: 'up' | 'down') => {
    const newItems = [...items]
    const targetIndex = direction === 'up' ? index - 1 : index + 1

    if (targetIndex < 0 || targetIndex >= newItems.length) {
      return
    }

    ;[newItems[index], newItems[targetIndex]] = [newItems[targetIndex], newItems[index]]
    setItems(newItems)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!name.trim()) {
      toast.error('Menu name is required')
      return
    }

    // Validate menu items
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      if (!item.label || !item.url) {
        toast.error(`Menu item ${i + 1} is missing label or URL`)
        return
      }
    }

    setLoading(true)

    try {
      if (isEdit && id) {
        const updatePromise = menusAPI.update(id, {
          name,
          description,
          items,
        })

        toast.promise(updatePromise, {
          loading: 'Updating menu...',
          success: 'Menu updated successfully!',
          error: (err: any) => err.response?.data?.error || 'Failed to update menu',
        })

        await updatePromise
      } else {
        const createPromise = menusAPI.create({
          name,
          description,
          items,
        })

        toast.promise(createPromise, {
          loading: 'Creating menu...',
          success: 'Menu created successfully!',
          error: (err: any) => err.response?.data?.error || 'Failed to create menu',
        })

        await createPromise
      }

      navigate('/admin/menus')
    } catch (err: any) {
      console.error('Menu creation error:', err)
      const errorMsg = err.response?.data?.error || err.message || `Failed to ${isEdit ? 'update' : 'create'} menu`
      setError(errorMsg)
      // Don't show toast here as toast.promise already handles it
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="settings-page">
        <div className="loading">Loading menu...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <div className="settings-header">
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <button
            onClick={() => navigate('/admin/menus')}
            className="back-button"
            type="button"
          >
            <ArrowLeft size={20} />
          </button>
          <h1>{isEdit ? 'Edit Menu' : 'Create New Menu'}</h1>
        </div>
        <p className="settings-description">
          {isEdit ? 'Update menu details and items.' : 'Create a reusable menu that can be used in headers, footers, or as a block.'}
        </p>
      </div>

      <form onSubmit={handleSubmit} className="settings-form">
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="name">Menu Name *</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            placeholder="Main Navigation"
            disabled={loading}
          />
          <small>A descriptive name for this menu (e.g., "Main Navigation", "Footer Links")</small>
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Optional description for this menu"
            rows={3}
            disabled={loading}
            style={{ width: '100%', padding: '12px 16px', border: '2px solid #e2e8f0', borderRadius: '8px', fontFamily: 'inherit', fontSize: '16px' }}
          />
        </div>

        <div className="form-group">
          <label>Menu Items</label>
          <div className="menu-items-list">
            {items.length === 0 ? (
              <div className="empty-state">
                <p>No menu items yet. Click "Add Item" to get started.</p>
              </div>
            ) : (
              items.map((item, index) => (
                <div key={index} className="menu-item-card">
                  <div className="menu-item-header">
                    <span className="menu-item-number">#{index + 1}</span>
                    <div className="menu-item-actions">
                      <button
                        type="button"
                        onClick={() => moveMenuItem(index, 'up')}
                        disabled={index === 0}
                        className="menu-action-button"
                        title="Move up"
                      >
                        <ChevronUp size={16} />
                      </button>
                      <button
                        type="button"
                        onClick={() => moveMenuItem(index, 'down')}
                        disabled={index === items.length - 1}
                        className="menu-action-button"
                        title="Move down"
                      >
                        <ChevronDown size={16} />
                      </button>
                      <button
                        type="button"
                        onClick={() => removeMenuItem(index)}
                        className="menu-action-button delete"
                        title="Remove"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </div>

                  <div className="menu-item-fields">
                    <div className="form-group">
                      <label>Label *</label>
                      <input
                        type="text"
                        value={item.label}
                        onChange={(e) => updateMenuItem(index, 'label', e.target.value)}
                        placeholder="Home"
                        required
                        disabled={loading}
                      />
                    </div>

                    <div className="form-group">
                      <label>URL *</label>
                      <input
                        type="text"
                        value={item.url}
                        onChange={(e) => updateMenuItem(index, 'url', e.target.value)}
                        placeholder="/home"
                        required
                        disabled={loading}
                      />
                    </div>

                    <div className="form-group">
                      <label>Target</label>
                      <select
                        value={item.target || '_self'}
                        onChange={(e) => updateMenuItem(index, 'target', e.target.value)}
                        disabled={loading}
                      >
                        <option value="_self">Same Window</option>
                        <option value="_blank">New Window</option>
                      </select>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="menu-actions">
            <button
              type="button"
              onClick={addMenuItem}
              className="add-menu-item-button"
              disabled={loading}
            >
              <Plus size={18} />
              <span>Add Item</span>
            </button>
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : isEdit ? 'Update Menu' : 'Create Menu'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

