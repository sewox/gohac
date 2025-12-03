import { useState, useEffect } from 'react'
import { Copy, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import { mediaAPI } from '../../lib/api'
import './MediaLibrary.css'

interface MediaItem {
  name: string
  url: string
  size: number
  type: string
}

export default function MediaLibrary() {
  const [mediaItems, setMediaItems] = useState<MediaItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [copiedUrl, setCopiedUrl] = useState<string | null>(null)

  useEffect(() => {
    fetchMedia()
  }, [])

  const fetchMedia = async () => {
    try {
      setLoading(true)
      const response = await mediaAPI.list()
      setMediaItems(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load media'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = async (url: string) => {
    try {
      const fullUrl = url.startsWith('http') ? url : `${window.location.origin}${url}`
      await navigator.clipboard.writeText(fullUrl)
      setCopiedUrl(url)
      toast.success('URL copied to clipboard!')
      setTimeout(() => setCopiedUrl(null), 2000)
    } catch (err) {
      toast.error('Failed to copy URL')
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  const isImage = (mimeType: string): boolean => {
    return mimeType.startsWith('image/')
  }

  if (loading) {
    return (
      <div className="media-library">
        <div className="loading">Loading media...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="media-library">
        <div className="error-message">{error}</div>
      </div>
    )
  }

  return (
    <div className="media-library">
      <div className="media-library-header">
        <h1>Media Library</h1>
        <p className="media-library-description">
          Browse and manage your uploaded media files. Click on an image to copy its URL.
        </p>
      </div>

      {mediaItems.length === 0 ? (
        <div className="empty-state">
          <p>No media files found. Upload files to see them here.</p>
        </div>
      ) : (
        <div className="media-grid">
          {mediaItems.map((item) => (
            <div key={item.name} className="media-item">
              {isImage(item.type) ? (
                <div className="media-item-image-container">
                  <img
                    src={item.url}
                    alt={item.name}
                    className="media-item-image"
                    onClick={() => copyToClipboard(item.url)}
                  />
                  <div className="media-item-overlay">
                    <button
                      className="copy-button"
                      onClick={() => copyToClipboard(item.url)}
                      title="Copy URL"
                    >
                      {copiedUrl === item.url ? (
                        <Check size={20} />
                      ) : (
                        <Copy size={20} />
                      )}
                    </button>
                  </div>
                </div>
              ) : (
                <div className="media-item-file">
                  <div className="media-item-file-icon">📄</div>
                  <button
                    className="copy-button"
                    onClick={() => copyToClipboard(item.url)}
                    title="Copy URL"
                  >
                    {copiedUrl === item.url ? (
                      <Check size={20} />
                    ) : (
                      <Copy size={20} />
                    )}
                  </button>
                </div>
              )}
              <div className="media-item-info">
                <div className="media-item-name" title={item.name}>
                  {item.name}
                </div>
                <div className="media-item-meta">
                  <span>{formatFileSize(item.size)}</span>
                  <span className="media-item-type">{item.type}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

