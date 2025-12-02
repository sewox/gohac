import { ImageData } from '../../../types/block'
import ImageUpload from '../ImageUpload'
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
          <ImageUpload
            value={data.url}
            onChange={(url) => handleChange('url', url)}
            label="Image"
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
      </div>
    </div>
  )
}

