import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Button } from '#/components/ui/button'
import { FormField, ImageUpload, SelectInput, TextInput } from '#/components/ui/input'
import { Modal } from '#/components/ui/modal'
import { vendorsKeys } from '#/features/vendors/queryKeys'
import { postVendors } from '#/lib/api'
import { slugify } from '#/lib/utils/slugify'
import { bareInput, form, inputWithPrefix, prefix } from './modals.css'

export type CreateVendorModalProps = {
  isOpen: boolean
  onClose: () => void
}

export function CreateVendorModal({ isOpen, onClose }: CreateVendorModalProps) {
  const queryClient = useQueryClient()

  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [isSlugManual, setIsSlugManual] = useState(false)
  const [specialty, setSpecialty] = useState('AUDIO_DUBBING')
  const [hourlyRate, setHourlyRate] = useState('120')
  const [turnaround, setTurnaround] = useState('24')
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)

  const handleNameChange = (val: string) => {
    setName(val)
    if (!isSlugManual) {
      setSlug(slugify(val))
    }
  }

  const handleSlugChange = (val: string) => {
    setIsSlugManual(true)
    setSlug(val.toLowerCase().replace(/[^a-z0-9-]/g, ''))
  }

  const resetForm = () => {
    setName('')
    setSlug('')
    setIsSlugManual(false)
    setSpecialty('AUDIO_DUBBING')
    setHourlyRate('120')
    setTurnaround('24')
    setAvatarUrl(undefined)
  }

  const mutation = useMutation({
    mutationFn: async () => {
      const finalId = slug ? `vendor-${slug}` : `vendor-${slugify(name)}`
      if (!name.trim()) throw new Error('Vendor name is required')
      if (!finalId.trim()) throw new Error('Vendor identifier is required')

      const { data, error } = await postVendors({
        body: {
          id: finalId,
          name: name.trim(),
          specialty,
          hourly_rate_usd: hourlyRate ? parseFloat(hourlyRate) : 0,
          turnaround_hours: turnaround ? parseInt(turnaround, 10) : 24,
          metadata: avatarUrl ? { avatar_url: avatarUrl } : {},
        },
      })

      if (error) {
        throw new Error(error.message || 'Failed to register vendor')
      }

      return data
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: vendorsKeys.all })
      toast.success(`Vendor "${data?.name ?? name}" registered successfully`)
      resetForm()
      onClose()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to register vendor')
    },
  })

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Register Localization Vendor"
      description="Onboard an external dubbing studio, subtitling facility, or QC inspection lab."
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose} disabled={mutation.isPending}>
            Cancel
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => mutation.mutate()}
            disabled={mutation.isPending || !name.trim()}
          >
            <Plus size={14} />
            <span>{mutation.isPending ? 'Registering...' : 'Register Vendor'}</span>
          </Button>
        </>
      }
    >
      <form
        class={form}
        onSubmit={(e) => {
          e.preventDefault()
          mutation.mutate()
        }}
      >
        <FormField label="Facility Name" required helper="Legal or operational facility name">
          <TextInput
            placeholder="e.g. Berlin Synchron Labs"
            value={name}
            onInput={(e) => handleNameChange((e.target as HTMLInputElement).value)}
            autoFocus
          />
        </FormField>

        <FormField
          label="Identifier Slug"
          required
          helper="Unique entity ID in operational storage"
        >
          <div class={inputWithPrefix}>
            <span class={prefix}>vendor-</span>
            <input
              type="text"
              class={bareInput}
              placeholder="berlin-synchron"
              value={slug}
              onInput={(e) => handleSlugChange((e.target as HTMLInputElement).value)}
            />
          </div>
        </FormField>

        <FormField label="Vendor Specialty" required>
          <SelectInput
            value={specialty}
            onChange={(e) => setSpecialty((e.target as HTMLSelectElement).value)}
            options={[
              { value: 'AUDIO_DUBBING', label: 'Audio Dubbing Studio' },
              { value: 'SUBTITLES', label: 'Subtitling & Captioning Facility' },
              { value: 'QC_LAB', label: 'Quality Control & Inspection Lab' },
              { value: 'METADATA', label: 'Metadata & Packaging House' },
              { value: 'MASTERING', label: 'Master Cut Video Processing' },
            ]}
          />
        </FormField>

        <FormField label="Hourly Rate (USD)" required helper="Standard facility billing rate">
          <TextInput
            type="number"
            min="0"
            step="0.01"
            placeholder="120.00"
            value={hourlyRate}
            onInput={(e) => setHourlyRate((e.target as HTMLInputElement).value)}
          />
        </FormField>

        <FormField label="Turnaround Time (Hours)" required helper="Standard delivery lead time">
          <TextInput
            type="number"
            min="1"
            step="1"
            placeholder="24"
            value={turnaround}
            onInput={(e) => setTurnaround((e.target as HTMLInputElement).value)}
          />
        </FormField>

        <FormField label="Facility Logo / Avatar" optional helper="Max 1MB PNG, JPEG, WebP, or GIF">
          <ImageUpload value={avatarUrl} onChange={setAvatarUrl} disabled={mutation.isPending} />
        </FormField>
      </form>
    </Modal>
  )
}
