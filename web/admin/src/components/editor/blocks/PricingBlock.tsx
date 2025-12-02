import { PricingData, PricingPlan } from '../../../types/block'
import Repeater from '../Repeater'
import './Block.css'

interface PricingBlockProps {
  data: PricingData
  onChange: (data: PricingData) => void
}

export default function PricingBlock({ data, onChange }: PricingBlockProps) {
  const handleChange = (field: keyof PricingData, value: any) => {
    onChange({
      ...data,
      [field]: value,
    })
  }

  const handlePlanChange = (index: number, field: keyof PricingPlan, value: any) => {
    const updatedPlans = [...(data.plans || [])]
    updatedPlans[index] = {
      ...updatedPlans[index],
      [field]: value,
    }
    handleChange('plans', updatedPlans)
  }

  const handleFeatureChange = (planIndex: number, featureIndex: number, value: string) => {
    const updatedPlans = [...(data.plans || [])]
    const updatedFeatures = [...(updatedPlans[planIndex].features || [])]
    updatedFeatures[featureIndex] = value
    updatedPlans[planIndex] = {
      ...updatedPlans[planIndex],
      features: updatedFeatures,
    }
    handleChange('plans', updatedPlans)
  }

  const handleAddFeature = (planIndex: number) => {
    const updatedPlans = [...(data.plans || [])]
    updatedPlans[planIndex] = {
      ...updatedPlans[planIndex],
      features: [...(updatedPlans[planIndex].features || []), ''],
    }
    handleChange('plans', updatedPlans)
  }

  const handleRemoveFeature = (planIndex: number, featureIndex: number) => {
    const updatedPlans = [...(data.plans || [])]
    const updatedFeatures = [...(updatedPlans[planIndex].features || [])]
    updatedFeatures.splice(featureIndex, 1)
    updatedPlans[planIndex] = {
      ...updatedPlans[planIndex],
      features: updatedFeatures,
    }
    handleChange('plans', updatedPlans)
  }

  const handleAddPlan = () => {
    const newPlan: PricingPlan = {
      name: '',
      price: '',
      description: '',
      features: [],
      button_text: 'Get Started',
      button_url: '',
      highlighted: false,
    }
    handleChange('plans', [...(data.plans || []), newPlan])
  }

  const handleRemovePlan = (index: number) => {
    const updatedPlans = [...(data.plans || [])]
    updatedPlans.splice(index, 1)
    handleChange('plans', updatedPlans)
  }

  return (
    <div className="block-editor">
      <div className="block-header">
        <h3>Pricing Block</h3>
      </div>
      <div className="block-content">
        <div className="form-group">
          <label htmlFor="pricing-title">Title</label>
          <input
            type="text"
            id="pricing-title"
            value={data.title || ''}
            onChange={(e) => handleChange('title', e.target.value)}
            placeholder="Pricing Section Title"
          />
        </div>
        <div className="form-group">
          <label htmlFor="pricing-subtitle">Subtitle</label>
          <input
            type="text"
            id="pricing-subtitle"
            value={data.subtitle || ''}
            onChange={(e) => handleChange('subtitle', e.target.value)}
            placeholder="Pricing Section Subtitle"
          />
        </div>
        <div className="form-group">
          <label>Pricing Plans</label>
          <Repeater
            items={data.plans || []}
            renderItem={(plan, planIndex) => (
              <div className="pricing-plan-form">
                <div className="form-group">
                  <label>Plan Name *</label>
                  <input
                    type="text"
                    value={plan.name || ''}
                    onChange={(e) => handlePlanChange(planIndex, 'name', e.target.value)}
                    placeholder="e.g., Starter, Pro, Enterprise"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Price *</label>
                  <input
                    type="text"
                    value={plan.price || ''}
                    onChange={(e) => handlePlanChange(planIndex, 'price', e.target.value)}
                    placeholder="e.g., $99/month or $999/year"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Description</label>
                  <textarea
                    value={plan.description || ''}
                    onChange={(e) => handlePlanChange(planIndex, 'description', e.target.value)}
                    placeholder="Plan description"
                    rows={2}
                  />
                </div>
                <div className="form-group">
                  <label>
                    <input
                      type="checkbox"
                      checked={plan.highlighted || false}
                      onChange={(e) => handlePlanChange(planIndex, 'highlighted', e.target.checked)}
                    />
                    {' '}Highlighted Plan
                  </label>
                </div>
                <div className="form-group">
                  <label>Features</label>
                  <div className="features-list">
                    {plan.features?.map((feature, featureIndex) => (
                      <div key={featureIndex} className="feature-input-row">
                        <input
                          type="text"
                          value={feature}
                          onChange={(e) => handleFeatureChange(planIndex, featureIndex, e.target.value)}
                          placeholder="Feature name"
                        />
                        <button
                          type="button"
                          onClick={() => handleRemoveFeature(planIndex, featureIndex)}
                          className="remove-feature-button"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    <button
                      type="button"
                      onClick={() => handleAddFeature(planIndex)}
                      className="add-feature-button"
                    >
                      + Add Feature
                    </button>
                  </div>
                </div>
                <div className="form-group">
                  <label>Button Text</label>
                  <input
                    type="text"
                    value={plan.button_text || ''}
                    onChange={(e) => handlePlanChange(planIndex, 'button_text', e.target.value)}
                    placeholder="Get Started"
                  />
                </div>
                <div className="form-group">
                  <label>Button URL</label>
                  <input
                    type="url"
                    value={plan.button_url || ''}
                    onChange={(e) => handlePlanChange(planIndex, 'button_url', e.target.value)}
                    placeholder="https://example.com/signup"
                  />
                </div>
              </div>
            )}
            onAdd={handleAddPlan}
            onRemove={handleRemovePlan}
            onChange={(index, newPlan) => handleChange('plans', data.plans.map((plan, i) => i === index ? newPlan : plan))}
            addButtonText="Add Pricing Plan"
            emptyMessage="No pricing plans yet. Add your first plan below."
          />
        </div>
      </div>
    </div>
  )
}

