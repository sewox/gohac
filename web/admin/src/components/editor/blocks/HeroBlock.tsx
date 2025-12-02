import { HeroData } from '../../../types/block'
import ImageUpload from '../ImageUpload'
import './Block.css'

interface HeroBlockProps {
  data: HeroData
  onChange: (data: HeroData) => void
}

export default function HeroBlock({ data, onChange }: HeroBlockProps) {
  const handleChange = (field: keyof HeroData, value: string) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Hero Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="hero-title">Title *</label>
          <input
            type="text"
            id="hero-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Hero Title"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="hero-subtitle">Subtitle</label>
          <input
            type="text"
            id="hero-subtitle"
            value={data.subtitle || ''}
            onChange={(e) => handleChange('subtitle', e.target.value)}
            placeholder="Hero Subtitle"
          />
        </div>
        <div className="form-group">
          <ImageUpload
            value={data.image_url}
            onChange={(url) => handleChange('image_url', url)}
            label="Hero Image"
          />
        </div>
      </div>
    </div>
  )
}

