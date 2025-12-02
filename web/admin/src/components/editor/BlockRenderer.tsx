import { Block } from '../../types/block'
import HeroBlock from './blocks/HeroBlock'
import TextBlock from './blocks/TextBlock'
import ImageBlock from './blocks/ImageBlock'

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

