import { TestimonialData, TestimonialItem } from '../../../types/block'
import Repeater from '../Repeater'
import ImageUpload from '../ImageUpload'
import './Block.css'

interface TestimonialBlockProps {
  data: TestimonialData
  onChange: (data: TestimonialData) => void
}

export default function TestimonialBlock({ data, onChange }: TestimonialBlockProps) {
  const handleChange = (field: keyof TestimonialData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  const handleItemChange = (index: number, field: keyof TestimonialItem, value: string) => {
    const updatedTestimonials = [...(data.testimonials || [])]
    updatedTestimonials[index] = {
      ...updatedTestimonials[index],
      [field]: value,
    }
    handleChange('testimonials', updatedTestimonials)
  }

  const handleAddItem = () => {
    const newItem: TestimonialItem = {
      quote: '',
      author: '',
      avatar_url: '',
      role: '',
    }
    handleChange('testimonials', [...(data.testimonials || []), newItem])
  }

  const handleRemoveItem = (index: number) => {
    const updatedTestimonials = [...(data.testimonials || [])]
    updatedTestimonials.splice(index, 1)
    handleChange('testimonials', updatedTestimonials)
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Testimonial Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="testimonial-title">Title</label>
          <input
            type="text"
            id="testimonial-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="What Our Customers Say"
          />
        </div>
        <div className="form-group">
          <label htmlFor="testimonial-subtitle">Subtitle</label>
          <input
            type="text"
            id="testimonial-subtitle"
            value={data.subtitle || ''}
            onChange={(e) => handleChange('subtitle', e.target.value)}
            placeholder="Testimonials Section Subtitle"
          />
        </div>
        <div className="form-group">
          <label>Testimonials</label>
          <Repeater
            items={data.testimonials || []}
            renderItem={(item, index) => (
              <div className="testimonial-item-form">
                <div className="form-group">
                  <label>Quote *</label>
                  <textarea
                    value={item.quote || ''}
                    onChange={(e) => handleItemChange(index, 'quote', e.target.value)}
                    placeholder="Customer testimonial quote"
                    rows={4}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Author Name *</label>
                  <input
                    type="text"
                    value={item.author || ''}
                    onChange={(e) => handleItemChange(index, 'author', e.target.value)}
                    placeholder="John Doe"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Role/Company</label>
                  <input
                    type="text"
                    value={item.role || ''}
                    onChange={(e) => handleItemChange(index, 'role', e.target.value)}
                    placeholder="CEO, Company Name"
                  />
                </div>
                <div className="form-group">
                  <ImageUpload
                    value={item.avatar_url}
                    onChange={(url) => handleItemChange(index, 'avatar_url', url)}
                    label="Avatar Image"
                  />
                </div>
              </div>
            )}
            onAdd={handleAddItem}
            onRemove={handleRemoveItem}
            onChange={(index, newItem) => handleChange('testimonials', data.testimonials.map((item, i) => i === index ? newItem : item))}
            addButtonText="Add Testimonial"
            emptyMessage="No testimonials yet. Add your first testimonial below."
          />
        </div>
      </div>
    </div>
  )
}

