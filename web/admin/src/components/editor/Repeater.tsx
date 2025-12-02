import { Plus, Trash2 } from 'lucide-react'
import './Repeater.css'

interface RepeaterProps<T> {
  items: T[]
  renderItem: (item: T, index: number) => React.ReactNode
  onAdd: () => void
  onRemove: (index: number) => void
  onChange: (index: number, newItem: T) => void
  addButtonText?: string
  emptyMessage?: string
}

export default function Repeater<T>({
  items,
  renderItem,
  onAdd,
  onRemove,
  onChange,
  addButtonText = 'Add Item',
  emptyMessage = 'No items yet. Add your first item below.',
}: RepeaterProps<T>) {
  return (
    <div className="repeater">
      {items.length === 0 ? (
        <div className="repeater-empty">
          <p>{emptyMessage}</p>
        </div>
      ) : (
        <div className="repeater-items">
          {items.map((item, index) => (
            <div key={index} className="repeater-item">
              <div className="repeater-item-header">
                <span className="repeater-item-number">#{index + 1}</span>
                <button
                  type="button"
                  onClick={() => onRemove(index)}
                  className="repeater-remove-button"
                  title="Remove item"
                >
                  <Trash2 size={16} />
                </button>
              </div>
              <div className="repeater-item-content">
                {renderItem(item, index)}
              </div>
            </div>
          ))}
        </div>
      )}
      <button
        type="button"
        onClick={onAdd}
        className="repeater-add-button"
      >
        <Plus size={18} />
        <span>{addButtonText}</span>
      </button>
    </div>
  )
}

