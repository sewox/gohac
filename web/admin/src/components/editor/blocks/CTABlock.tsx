import { CTAData } from '../../../types/block'
import './Block.css'

interface CTABlockProps {
  data: CTAData
  onChange: (data: CTAData) => void
}

export default function CTABlock({ data, onChange }: CTABlockProps) {
  const handleChange = (field: keyof CTAData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Call to Action Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="cta-title">Title *</label>
          <input
            type="text"
            id="cta-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Ready to get started?"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="cta-subtitle">Subtitle</label>
          <input
            type="text"
            id="cta-subtitle"
            value={data.subtitle || ''}
            onChange={(e) => handleChange('subtitle', e.target.value)}
            placeholder="Join thousands of happy customers"
          />
        </div>
        <div className="form-group">
          <label htmlFor="cta-button-text">Button Text *</label>
          <input
            type="text"
            id="cta-button-text"
            value={data.button_text || ''}
            onChange={(e) => handleChange('button_text', e.target.value)}
            placeholder="Get Started"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="cta-button-url">Button URL *</label>
          <input
            type="url"
            id="cta-button-url"
            value={data.button_url || ''}
            onChange={(e) => handleChange('button_url', e.target.value)}
            placeholder="https://example.com/signup"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="cta-button-style">Button Style</label>
          <select
            id="cta-button-style"
            value={data.button_style || 'primary'}
            onChange={(e) => handleChange('button_style', e.target.value as 'primary' | 'secondary' | 'outline')}
          >
            <option value="primary">Primary</option>
            <option value="secondary">Secondary</option>
            <option value="outline">Outline</option>
          </select>
        </div>
        <div className="form-group">
          <label htmlFor="cta-background">Background</label>
          <input
            type="text"
            id="cta-background"
            value={data.background || ''}
            onChange={(e) => handleChange('background', e.target.value)}
            placeholder="Color or gradient (e.g., #667eea or linear-gradient(...))"
          />
        </div>
      </div>
    </div>
  )
}

