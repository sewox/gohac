import { FeaturesData, FeatureItem } from '../../../types/block'
import Repeater from '../Repeater'
import './Block.css'

interface FeaturesBlockProps {
  data: FeaturesData
  onChange: (data: FeaturesData) => void
}

export default function FeaturesBlock({ data, onChange }: FeaturesBlockProps) {
  const handleChange = (field: keyof FeaturesData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  const handleItemChange = (index: number, field: keyof FeatureItem, value: string) => {
    const updatedItems = [...(data.items || [])]
    updatedItems[index] = {
      ...updatedItems[index],
      [field]: value,
    }
    handleChange('items', updatedItems)
  }

  const handleAddItem = () => {
    const newItem: FeatureItem = {
      title: '',
      description: '',
      icon: '',
    }
    handleChange('items', [...(data.items || []), newItem])
  }

  const handleRemoveItem = (index: number) => {
    const updatedItems = [...(data.items || [])]
    updatedItems.splice(index, 1)
    handleChange('items', updatedItems)
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Features Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="features-title">Title</label>
          <input
            type="text"
            id="features-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Features Section Title"
          />
        </div>
        <div className="form-group">
          <label htmlFor="features-subtitle">Subtitle</label>
          <input
            type="text"
            id="features-subtitle"
            value={data.subtitle || ''}
            onChange={(e) => handleChange('subtitle', e.target.value)}
            placeholder="Features Section Subtitle"
          />
        </div>
        <div className="form-group">
          <label htmlFor="features-columns">Columns</label>
          <select
            id="features-columns"
            value={data.columns || 3}
            onChange={(e) => handleChange('columns', parseInt(e.target.value) as 2 | 3 | 4)}
          >
            <option value={2}>2 Columns</option>
            <option value={3}>3 Columns</option>
            <option value={4}>4 Columns</option>
          </select>
        </div>
        <div className="form-group">
          <label>Feature Items</label>
          <Repeater
            items={data.items || []}
            renderItem={(item, index) => (
              <div className="feature-item-form">
                <div className="form-group">
                  <label>Title *</label>
                  <input
                    type="text"
                    value={item.title || ''}
                    onChange={(e) => handleItemChange(index, 'title', e.target.value)}
                    placeholder="Feature Title"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Description</label>
                  <textarea
                    value={item.description || ''}
                    onChange={(e) => handleItemChange(index, 'description', e.target.value)}
                    placeholder="Feature Description"
                    rows={3}
                  />
                </div>
                <div className="form-group">
                  <label>Icon</label>
                  <input
                    type="text"
                    value={item.icon || ''}
                    onChange={(e) => handleItemChange(index, 'icon', e.target.value)}
                    placeholder="Icon name or URL"
                  />
                  <small>e.g., "star", "check", or icon URL</small>
                </div>
              </div>
            )}
            onAdd={handleAddItem}
            onRemove={handleRemoveItem}
            onChange={(index, newItem) => handleChange('items', data.items.map((item, i) => i === index ? newItem : item))}
            addButtonText="Add Feature"
            emptyMessage="No features yet. Add your first feature below."
          />
        </div>
      </div>
    </div>
  )
}

