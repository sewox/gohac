import { useState, useEffect } from 'react'
import { Save } from 'lucide-react'
import toast from 'react-hot-toast'
import { settingsAPI } from '../../lib/api'
import ImageUpload from '../../components/editor/ImageUpload'
import './Settings.css'

interface GlobalSettings {
  site_name: string
  logo: string
  favicon: string
  contact_email: string
}

export default function GeneralSettings() {
  const [settings, setSettings] = useState<GlobalSettings>({
    site_name: '',
    logo: '',
    favicon: '',
    contact_email: '',
  })
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchSettings()
  }, [])

  const fetchSettings = async () => {
    try {
      setFetching(true)
      const response = await settingsAPI.get()
      setSettings(response.data)
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load settings'
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

    const updatePromise = settingsAPI.update(settings)

    toast.promise(updatePromise, {
      loading: 'Saving settings...',
      success: 'Settings saved successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to save settings',
    })

    try {
      await updatePromise
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to save settings'
      setError(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="settings-page">
        <div className="loading">Loading settings...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <div className="settings-header">
        <h1>General Settings</h1>
        <p className="settings-description">
          Configure your site's global settings, logo, and branding.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="settings-form">
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="site_name">Site Name *</label>
          <input
            type="text"
            id="site_name"
            value={settings.site_name}
            onChange={(e) => setSettings({ ...settings, site_name: e.target.value })}
            required
            placeholder="My Website"
            disabled={loading}
          />
          <small>The name of your website displayed in the header</small>
        </div>

        <div className="form-group">
          <label htmlFor="contact_email">Contact Email</label>
          <input
            type="email"
            id="contact_email"
            value={settings.contact_email}
            onChange={(e) => setSettings({ ...settings, contact_email: e.target.value })}
            placeholder="contact@example.com"
            disabled={loading}
          />
          <small>Email address for contact forms and inquiries</small>
        </div>

        <div className="form-group">
          <label htmlFor="logo">Logo</label>
          <ImageUpload
            value={settings.logo}
            onChange={(url) => setSettings({ ...settings, logo: url })}
            label="Site Logo"
          />
          <small>Logo displayed in the header (recommended: 200x50px)</small>
        </div>

        <div className="form-group">
          <label htmlFor="favicon">Favicon</label>
          <ImageUpload
            value={settings.favicon}
            onChange={(url) => setSettings({ ...settings, favicon: url })}
            label="Favicon"
          />
          <small>Favicon displayed in browser tabs (recommended: 32x32px or 16x16px)</small>
        </div>

        <div className="form-actions">
          <button type="submit" className="save-button" disabled={loading}>
            <Save size={18} />
            <span>{loading ? 'Saving...' : 'Save Settings'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}

