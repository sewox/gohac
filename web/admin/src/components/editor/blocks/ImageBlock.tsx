import { ImageData } from '../../../types/block'
import './Block.css'

interface ImageBlockProps {
  data: ImageData
  onChange: (data: ImageData) => void
}

export default function ImageBlock({ data, onChange }: ImageBlockProps) {
  const handleChange = (field: keyof ImageData, value: string) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Image Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="image-url">Image URL *</label>
          <input
            type="url"
            id="image-url"
            value={data.url || ''}
            onChange={(e) => handleChange('url', e.target.value)}
            placeholder="https://example.com/image.jpg"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="image-alt">Alt Text</label>
          <input
            type="text"
            id="image-alt"
            value={data.alt || ''}
            onChange={(e) => handleChange('alt', e.target.value)}
            placeholder="Descriptive alt text"
          />
        </div>
        <div className="form-group">
          <label htmlFor="image-caption">Caption</label>
          <input
            type="text"
            id="image-caption"
            value={data.caption || ''}
            onChange={(e) => handleChange('caption', e.target.value)}
            placeholder="Image caption"
          />
        </div>
        {data.url && (
          <div className="image-preview">
            <img src={data.url} alt={data.alt || 'Preview'} />
          </div>
        )}
      </div>
    </div>
  )
}

