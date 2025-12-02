import { TextData } from '../../../types/block'
import './Block.css'

interface TextBlockProps {
  data: TextData
  onChange: (data: TextData) => void
}

export default function TextBlock({ data, onChange }: TextBlockProps) {
  const handleChange = (field: keyof TextData, value: string) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Text Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="text-content">Content *</label>
          <textarea
            id="text-content"
            value={data.content || ''}
            onChange={(e) => handleChange('content', e.target.value)}
            placeholder="Enter your text content here..."
            rows={6}
            required
          />
          <small>Supports HTML and Markdown</small>
        </div>
        <div className="form-group">
          <label htmlFor="text-align">Alignment</label>
          <select
            id="text-align"
            value={data.align || 'left'}
            onChange={(e) =>
              handleChange('align', e.target.value as 'left' | 'center' | 'right')
            }
          >
            <option value="left">Left</option>
            <option value="center">Center</option>
            <option value="right">Right</option>
          </select>
        </div>
      </div>
    </div>
  )
}

