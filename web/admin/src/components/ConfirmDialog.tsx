import { ReactNode } from 'react'
import { AlertTriangle, X } from 'lucide-react'
import './ConfirmDialog.css'

interface ConfirmDialogProps {
  isOpen: boolean
  title: string
  message: string
  onConfirm: () => void
  onCancel: () => void
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
}

export default function ConfirmDialog({
  isOpen,
  title,
  message,
  onConfirm,
  onCancel,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'danger',
}: ConfirmDialogProps) {
  if (!isOpen) return null

  return (
    <div className="confirm-dialog-overlay" onClick={onCancel}>
      <div className="confirm-dialog" onClick={(e) => e.stopPropagation()}>
        <div className="confirm-dialog-header">
          <div className={`confirm-icon ${variant}`}>
            <AlertTriangle size={24} />
          </div>
          <h3>{title}</h3>
          <button onClick={onCancel} className="close-button">
            <X size={20} />
          </button>
        </div>
        <div className="confirm-dialog-body">
          <p>{message}</p>
        </div>
        <div className="confirm-dialog-footer">
          <button onClick={onCancel} className="cancel-button">
            {cancelText}
          </button>
          <button
            onClick={onConfirm}
            className={`confirm-button ${variant}`}
          >
            {confirmText}
          </button>
        </div>
      </div>
    </div>
  )
}

