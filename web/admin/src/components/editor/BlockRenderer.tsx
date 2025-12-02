import { Block } from '../../types/block'
import HeroBlock from './blocks/HeroBlock'
import TextBlock from './blocks/TextBlock'
import ImageBlock from './blocks/ImageBlock'
import FeaturesBlock from './blocks/FeaturesBlock'
import PricingBlock from './blocks/PricingBlock'
import FAQBlock from './blocks/FAQBlock'
import TestimonialBlock from './blocks/TestimonialBlock'
import VideoBlock from './blocks/VideoBlock'
import CTABlock from './blocks/CTABlock'

interface BlockRendererProps {
  block: Block
  onChange: (data: Block['data']) => void
}

export default function BlockRenderer({ block, onChange }: BlockRendererProps) {
  switch (block.type) {
    case 'hero':
      return (
        <HeroBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'text':
      return (
        <TextBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'image':
      return (
        <ImageBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'features':
      return (
        <FeaturesBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'pricing':
      return (
        <PricingBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'faq':
      return (
        <FAQBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'testimonial':
      return (
        <TestimonialBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'video':
      return (
        <VideoBlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    case 'cta':
      return (
        <CTABlock
          data={block.data as any}
          onChange={(data) => onChange(data)}
        />
      )
    default:
      return (
        <div className="block-editor">
          <div className="block-header">
            <h3>Unknown Block Type: {block.type}</h3>
          </div>
        </div>
      )
  }
}

