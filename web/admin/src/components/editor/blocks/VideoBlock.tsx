import { VideoData } from '../../../types/block'
import './Block.css'

interface VideoBlockProps {
  data: VideoData
  onChange: (data: VideoData) => void
}

export default function VideoBlock({ data, onChange }: VideoBlockProps) {
  const handleChange = (field: keyof VideoData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Video Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="video-url">Video URL *</label>
          <input
            type="url"
            id="video-url"
            value={data.url || ''}
            onChange={(e) => handleChange('url', e.target.value)}
            placeholder="https://www.youtube.com/watch?v=..."
            required
          />
          <small>Supports YouTube, Vimeo, or direct video URLs</small>
        </div>
        <div className="form-group">
          <label htmlFor="video-title">Title</label>
          <input
            type="text"
            id="video-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Video Title"
          />
        </div>
        <div className="form-group">
          <label htmlFor="video-description">Description</label>
          <textarea
            id="video-description"
            value={data.description || ''}
            onChange={(e) => handleChange('description', e.target.value)}
            placeholder="Video description"
            rows={3}
          />
        </div>
        <div className="form-group">
          <label>
            <input
              type="checkbox"
              checked={data.autoplay || false}
              onChange={(e) => handleChange('autoplay', e.target.checked)}
            />
            {' '}Autoplay
          </label>
        </div>
        <div className="form-group">
          <label>
            <input
              type="checkbox"
              checked={data.loop || false}
              onChange={(e) => handleChange('loop', e.target.checked)}
            />
            {' '}Loop
          </label>
        </div>
        {data.url && (
          <div className="video-preview">
            <p>Preview:</p>
            <div className="video-preview-placeholder">
              Video will be embedded here
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

