import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Save, Plus, Trash2, ChevronDown, ChevronUp, Settings } from 'lucide-react'
import toast from 'react-hot-toast'
import { settingsAPI } from '../../lib/api'
import './Settings.css'

type MenuPosition = 'header' | 'footer'

interface MenuItem {
  label: string
  url: string
  target?: string
  children?: MenuItem[]
}

interface Menu {
  position: MenuPosition
  items: MenuItem[]
}

export default function MenuManager() {
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<MenuPosition>('header')
  const [headerMenu, setHeaderMenu] = useState<Menu>({ position: 'header', items: [] })
  const [footerMenu, setFooterMenu] = useState<Menu>({ position: 'footer', items: [] })
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchMenus()
  }, [])

  const fetchMenus = async () => {
    try {
      setFetching(true)
      const [headerRes, footerRes] = await Promise.all([
        settingsAPI.getMenu('header'),
        settingsAPI.getMenu('footer'),
      ])
      setHeaderMenu({ position: 'header', items: headerRes.data.items || [] })
      setFooterMenu({ position: 'footer', items: footerRes.data.items || [] })
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load menus'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setFetching(false)
    }
  }

  const getCurrentMenu = () => {
    return activeTab === 'header' ? headerMenu : footerMenu
  }

  const setCurrentMenu = (menu: Menu) => {
    if (activeTab === 'header') {
      setHeaderMenu(menu)
    } else {
      setFooterMenu(menu)
    }
  }

  const addMenuItem = () => {
    const menu = getCurrentMenu()
    const newItem: MenuItem = {
      label: '',
      url: '',
      target: '_self',
    }
    setCurrentMenu({
      ...menu,
      items: [...menu.items, newItem],
    })
  }

  const updateMenuItem = (index: number, field: keyof MenuItem, value: string) => {
    const menu = getCurrentMenu()
    const newItems = [...menu.items]
    newItems[index] = { ...newItems[index], [field]: value }
    setCurrentMenu({ ...menu, items: newItems })
  }

  const removeMenuItem = (index: number) => {
    const menu = getCurrentMenu()
    const newItems = menu.items.filter((_, i) => i !== index)
    setCurrentMenu({ ...menu, items: newItems })
  }

  const moveMenuItem = (index: number, direction: 'up' | 'down') => {
    const menu = getCurrentMenu()
    const newItems = [...menu.items]
    const targetIndex = direction === 'up' ? index - 1 : index + 1

    if (targetIndex < 0 || targetIndex >= newItems.length) {
      return
    }

    ;[newItems[index], newItems[targetIndex]] = [newItems[targetIndex], newItems[index]]
    setCurrentMenu({ ...menu, items: newItems })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    const menu = getCurrentMenu()

    // Validate menu items
    for (let i = 0; i < menu.items.length; i++) {
      const item = menu.items[i]
      if (!item.label || !item.url) {
        toast.error(`Menu item ${i + 1} is missing label or URL`)
        setLoading(false)
        return
      }
    }

    const updatePromise = settingsAPI.updateMenu(menu.position, menu.items)

    toast.promise(updatePromise, {
      loading: `Saving ${menu.position} menu...`,
      success: `${menu.position.charAt(0).toUpperCase() + menu.position.slice(1)} menu saved successfully!`,
      error: (err: any) => err.response?.data?.error || 'Failed to save menu',
    })

    try {
      await updatePromise
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to save menu'
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="settings-page">
        <div className="loading">Loading menus...</div>
      </div>
    )
  }

  const currentMenu = getCurrentMenu()

  return (
    <div className="settings-page">
      <div className="settings-header">
        <h1>Menu Manager</h1>
        <p className="settings-description">
          Manage navigation menus for your site's header and footer.
        </p>
      </div>

      {/* Tab Navigation - Settings Pages */}
      <div className="tab-navigation">
        <button
          type="button"
          className="tab-button"
          onClick={() => navigate('/admin/settings')}
        >
          <Settings size={16} />
          <span>General</span>
        </button>
        <button
          type="button"
          className="tab-button active"
          disabled
        >
          <span>Menus</span>
        </button>
      </div>

      {/* Tab Navigation - Menu Positions */}
      <div className="tab-navigation" style={{ marginTop: '16px' }}>
        <button
          type="button"
          className={`tab-button ${activeTab === 'header' ? 'active' : ''}`}
          onClick={() => setActiveTab('header')}
        >
          <span>Header Menu</span>
        </button>
        <button
          type="button"
          className={`tab-button ${activeTab === 'footer' ? 'active' : ''}`}
          onClick={() => setActiveTab('footer')}
        >
          <span>Footer Menu</span>
        </button>
      </div>

      <form onSubmit={handleSubmit} className="settings-form">
        {error && <div className="error-message">{error}</div>}

        <div className="menu-items-list">
          {currentMenu.items.length === 0 ? (
            <div className="empty-state">
              <p>No menu items yet. Click "Add Item" to get started.</p>
            </div>
          ) : (
            currentMenu.items.map((item, index) => (
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
                      disabled={index === currentMenu.items.length - 1}
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

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : `Save ${activeTab.charAt(0).toUpperCase() + activeTab.slice(1)} Menu`}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

