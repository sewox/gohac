import { useState, useEffect } from 'react'
import { Plus, Trash2, ChevronUp, ChevronDown, GripVertical } from 'lucide-react'
import { Block } from '../../types/block'
import BlockRenderer from './BlockRenderer'
import './BlockEditor.css'

interface BlockEditorProps {
  initialBlocks?: Block[]
  onChange: (blocks: Block[]) => void
}

export default function BlockEditor({
  initialBlocks = [],
  onChange,
}: BlockEditorProps) {
  const [blocks, setBlocks] = useState<Block[]>(initialBlocks)
  const [showAddMenu, setShowAddMenu] = useState(false)

  useEffect(() => {
    setBlocks(initialBlocks)
  }, [initialBlocks])

  useEffect(() => {
    onChange(blocks)
  }, [blocks, onChange])

  const addBlock = (type: 'hero' | 'text' | 'image') => {
    const newBlock: Block = {
      id: `block-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      type,
      data: getDefaultDataForType(type),
    }
    setBlocks([...blocks, newBlock])
  }

  const getDefaultDataForType = (
    type: 'hero' | 'text' | 'image'
  ): Block['data'] => {
    switch (type) {
      case 'hero':
        return { title: '', subtitle: '' }
      case 'text':
        return { content: '' }
      case 'image':
        return { url: '', alt: '' }
    }
  }

  const updateBlock = (index: number, data: Block['data']) => {
    const updatedBlocks = [...blocks]
    updatedBlocks[index] = {
      ...updatedBlocks[index],
      data,
    }
    setBlocks(updatedBlocks)
  }

  const deleteBlock = (index: number) => {
    setBlocks(blocks.filter((_, i) => i !== index))
  }

  const moveBlock = (index: number, direction: 'up' | 'down') => {
    const newIndex = direction === 'up' ? index - 1 : index + 1
    if (newIndex < 0 || newIndex >= blocks.length) return

    const updatedBlocks = [...blocks]
    ;[updatedBlocks[index], updatedBlocks[newIndex]] = [
      updatedBlocks[newIndex],
      updatedBlocks[index],
    ]
    setBlocks(updatedBlocks)
  }

  return (
    <div className="block-editor-container">
      <div className="blocks-list">
        {blocks.length === 0 ? (
          <div className="empty-blocks">
            <p>No blocks yet. Add your first block below.</p>
          </div>
        ) : (
          blocks.map((block, index) => (
            <div key={block.id} className="block-wrapper">
              <div className="block-controls">
                <div className="block-handle">
                  <GripVertical size={16} />
                  <span className="block-number">{index + 1}</span>
                </div>
                <div className="block-actions">
                  <button
                    onClick={() => moveBlock(index, 'up')}
                    disabled={index === 0}
                    className="block-action-button"
                    title="Move up"
                  >
                    <ChevronUp size={16} />
                  </button>
                  <button
                    onClick={() => moveBlock(index, 'down')}
                    disabled={index === blocks.length - 1}
                    className="block-action-button"
                    title="Move down"
                  >
                    <ChevronDown size={16} />
                  </button>
                  <button
                    onClick={() => deleteBlock(index)}
                    className="block-action-button delete"
                    title="Delete block"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
              <BlockRenderer
                block={block}
                onChange={(data) => updateBlock(index, data)}
              />
            </div>
          ))
        )}
      </div>

      <div className="add-block-section">
        <div className="add-block-dropdown">
          <button
            type="button"
            className="add-block-button"
            onClick={() => setShowAddMenu(!showAddMenu)}
          >
            <Plus size={18} />
            <span>Add Block</span>
            <ChevronDown size={16} className={showAddMenu ? 'rotate' : ''} />
          </button>
          {showAddMenu && (
            <>
              <div
                className="add-block-overlay"
                onClick={() => setShowAddMenu(false)}
              />
              <div className="add-block-menu">
                <button
                  type="button"
                  onClick={() => {
                    addBlock('hero')
                    setShowAddMenu(false)
                  }}
                  className="add-block-option"
                >
                  <span className="block-type-icon">🎯</span>
                  <div>
                    <div className="block-type-name">Hero</div>
                    <div className="block-type-desc">Title and subtitle</div>
                  </div>
                </button>
                <button
                  type="button"
                  onClick={() => {
                    addBlock('text')
                    setShowAddMenu(false)
                  }}
                  className="add-block-option"
                >
                  <span className="block-type-icon">📝</span>
                  <div>
                    <div className="block-type-name">Text</div>
                    <div className="block-type-desc">Rich text content</div>
                  </div>
                </button>
                <button
                  type="button"
                  onClick={() => {
                    addBlock('image')
                    setShowAddMenu(false)
                  }}
                  className="add-block-option"
                >
                  <span className="block-type-icon">🖼️</span>
                  <div>
                    <div className="block-type-name">Image</div>
                    <div className="block-type-desc">Image with caption</div>
                  </div>
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

