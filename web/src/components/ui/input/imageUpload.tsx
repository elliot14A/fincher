import { Loader2, UploadCloud, X } from 'lucide-preact'
import { useRef, useState } from 'preact/hooks'
import { deleteUploadsById, postUploads } from '#/lib/api'
import {
  dropzone,
  dropzoneActive,
  errorText,
  hiddenFileInput,
  previewContainer,
  previewDetails,
  previewFilename,
  previewSize,
  removeButton,
  thumbnail,
  uploadHint,
  uploadIcon,
  uploadPrompt,
} from './imageUpload.css'

const MAX_SIZE_BYTES = 1 * 1024 * 1024 // Strict 1MB

function extractUploadId(url: string | undefined): string | null {
  if (!url) return null
  const prefix = '/api/uploads/'
  if (url.startsWith(prefix)) {
    return url.slice(prefix.length)
  }
  return null
}

export type ImageUploadProps = {
  value?: string
  onChange: (url: string | undefined) => void
  disabled?: boolean
}

export function ImageUpload({ value, onChange, disabled }: ImageUploadProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [isDragging, setIsDragging] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [fileName, setFileName] = useState<string | null>(null)
  const [fileSize, setFileSize] = useState<number | null>(null)

  const handleFile = async (file: File) => {
    setErrorMessage(null)

    if (file.size > MAX_SIZE_BYTES) {
      setErrorMessage('File size exceeds the 1MB limit. Please choose a smaller image.')
      return
    }

    const allowedTypes = ['image/png', 'image/jpeg', 'image/webp', 'image/gif']
    if (!allowedTypes.includes(file.type)) {
      setErrorMessage('Unsupported file format. Please upload a PNG, JPEG, WebP, or GIF.')
      return
    }

    const previousUploadId = extractUploadId(value)

    try {
      setIsUploading(true)
      setFileName(file.name)
      setFileSize(file.size)

      const { data, error } = await postUploads({
        body: {
          file,
        },
      })

      if (error) {
        throw new Error(error.message || 'Upload failed')
      }

      if (data?.url) {
        // Clean up previously replaced upload if present
        if (previousUploadId) {
          deleteUploadsById({ path: { id: previousUploadId } }).catch(() => {})
        }
        onChange(data.url)
      }
    } catch (err) {
      setErrorMessage(err instanceof Error ? err.message : 'Failed to upload image.')
      onChange(undefined)
      setFileName(null)
      setFileSize(null)
    } finally {
      setIsUploading(false)
    }
  }

  const handleRemove = () => {
    const uploadId = extractUploadId(value)
    if (uploadId) {
      deleteUploadsById({ path: { id: uploadId } }).catch(() => {})
    }
    onChange(undefined)
    setFileName(null)
    setFileSize(null)
    setErrorMessage(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
  }

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/png,image/jpeg,image/webp,image/gif"
        class={hiddenFileInput}
        disabled={disabled || isUploading}
        onChange={(e) => {
          const files = (e.target as HTMLInputElement).files
          if (files?.[0]) {
            handleFile(files[0])
          }
        }}
      />

      {value ? (
        <div class={previewContainer}>
          <img src={value} alt="Preview" class={thumbnail} />
          <div class={previewDetails}>
            <span class={previewFilename}>{fileName ?? 'Uploaded image'}</span>
            <span class={previewSize}>
              {fileSize ? formatSize(fileSize) : 'Stored in database'}
            </span>
          </div>
          <button
            type="button"
            class={removeButton}
            onClick={handleRemove}
            aria-label="Remove image"
            disabled={disabled || isUploading}
          >
            <X size={14} />
          </button>
        </div>
      ) : (
        <button
          type="button"
          class={isDragging ? `${dropzone} ${dropzoneActive}` : dropzone}
          onClick={() => fileInputRef.current?.click()}
          onDragOver={(e) => {
            e.preventDefault()
            setIsDragging(true)
          }}
          onDragLeave={() => setIsDragging(false)}
          onDrop={(e) => {
            e.preventDefault()
            setIsDragging(false)
            if (e.dataTransfer?.files?.[0]) {
              handleFile(e.dataTransfer.files[0])
            }
          }}
        >
          {isUploading ? (
            <>
              <Loader2 size={22} class={uploadIcon} />
              <span class={uploadPrompt}>Uploading image to database...</span>
            </>
          ) : (
            <>
              <UploadCloud size={22} class={uploadIcon} />
              <span class={uploadPrompt}>Click to upload or drag and drop</span>
              <span class={uploadHint}>PNG, JPEG, WebP, or GIF (Strict max 1MB)</span>
            </>
          )}
        </button>
      )}

      {errorMessage ? <div class={errorText}>{errorMessage}</div> : null}
    </div>
  )
}
