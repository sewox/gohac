import { useState, useRef } from 'react'
import { Upload, X, Download, Loader } from 'lucide-react'
import api from '../../lib/api'
import toast from 'react-hot-toast'
import './ImageUpload.css'

interface ImageUploadProps {
  value?: string
  onChange: (url: string) => void
  label?: string
  required?: boolean
}

export default function ImageUpload({
  value,
  onChange,
  label = 'Image',
  required = false,
}: ImageUploadProps) {
  const [uploading, setUploading] = useState(false)
  const [dragActive, setDragActive] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      toast.error('Please select an image file')
      return
    }

    setUploading(true)
    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await api.post('/v1/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      onChange(response.data.url)
      toast.success('Image uploaded successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to upload image')
    } finally {
      setUploading(false)
    }
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      handleFileSelect(file)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setDragActive(false)
    const file = e.dataTransfer.files?.[0]
    if (file) {
      handleFileSelect(file)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setDragActive(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setDragActive(false)
  }

  const handleDownloadFromURL = async () => {
    const url = prompt('Enter image URL:')
    if (!url) return

    setUploading(true)
    try {
      const response = await api.post('/v1/upload/from-url', { url })
      onChange(response.data.url)
      toast.success('Image downloaded and saved successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to download image')
    } finally {
      setUploading(false)
    }
  }

  const handleRemove = () => {
    onChange('')
  }

  return (
    <div className="image-upload">
      {label && (
        <label className="image-upload-label">
          {label}
          {required && <span className="required">*</span>}
        </label>
      )}

      {value ? (
        <div className="image-preview-container">
          <div className="image-preview-wrapper">
            <img src={value} alt="Preview" className="image-preview-img" />
            <button
              type="button"
              onClick={handleRemove}
              className="image-remove-button"
              title="Remove image"
            >
              <X size={18} />
            </button>
          </div>
          <div className="image-preview-actions">
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="image-action-button"
              disabled={uploading}
            >
              <Upload size={16} />
              Replace Image
            </button>
            <button
              type="button"
              onClick={handleDownloadFromURL}
              className="image-action-button"
              disabled={uploading}
            >
              <Download size={16} />
              Download from URL
            </button>
          </div>
        </div>
      ) : (
        <div
          className={`image-upload-dropzone ${dragActive ? 'drag-active' : ''} ${uploading ? 'uploading' : ''}`}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onClick={() => !uploading && fileInputRef.current?.click()}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileInputChange}
            className="image-upload-input"
            disabled={uploading}
          />
          {uploading ? (
            <div className="image-upload-loading">
              <Loader size={24} className="spinner" />
              <span>Uploading...</span>
            </div>
          ) : (
            <>
              <Upload size={32} />
              <p className="image-upload-text">
                <strong>Click to upload</strong> or drag and drop
              </p>
              <p className="image-upload-hint">PNG, JPG, GIF up to 10MB</p>
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation()
                  handleDownloadFromURL()
                }}
                className="image-upload-url-button"
              >
                <Download size={16} />
                Or download from URL
              </button>
            </>
          )}
        </div>
      )}
    </div>
  )
}

