import { FAQData, FAQItem } from '../../../types/block'
import Repeater from '../Repeater'
import './Block.css'

interface FAQBlockProps {
  data: FAQData
  onChange: (data: FAQData) => void
}

export default function FAQBlock({ data, onChange }: FAQBlockProps) {
  const handleChange = (field: keyof FAQData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  const handleItemChange = (index: number, field: keyof FAQItem, value: string) => {
    const updatedItems = [...(data.items || [])]
    updatedItems[index] = {
      ...updatedItems[index],
      [field]: value,
    }
    handleChange('items', updatedItems)
  }

  const handleAddItem = () => {
    const newItem: FAQItem = {
      question: '',
      answer: '',
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
        <h3>FAQ Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="faq-title">Title</label>
          <input
            type="text"
            id="faq-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Frequently Asked Questions"
          />
        </div>
        <div className="form-group">
          <label>FAQ Items</label>
          <Repeater
            items={data.items || []}
            renderItem={(item, index) => (
              <div className="faq-item-form">
                <div className="form-group">
                  <label>Question *</label>
                  <input
                    type="text"
                    value={item.question || ''}
                    onChange={(e) => handleItemChange(index, 'question', e.target.value)}
                    placeholder="What is your question?"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Answer *</label>
                  <textarea
                    value={item.answer || ''}
                    onChange={(e) => handleItemChange(index, 'answer', e.target.value)}
                    placeholder="Answer to the question"
                    rows={4}
                    required
                  />
                </div>
              </div>
            )}
            onAdd={handleAddItem}
            onRemove={handleRemoveItem}
            onChange={(index, newItem) => handleChange('items', data.items.map((item, i) => i === index ? newItem : item))}
            addButtonText="Add FAQ Item"
            emptyMessage="No FAQ items yet. Add your first question below."
          />
        </div>
      </div>
    </div>
  )
}

