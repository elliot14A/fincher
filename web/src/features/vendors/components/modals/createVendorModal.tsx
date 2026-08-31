import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Button } from '#/components/ui/button'
import { FormField, ImageUpload, TextInput } from '#/components/ui/input'
import { Modal } from '#/components/ui/modal'
import { vendorsKeys } from '#/features/vendors/queryKeys'
import { postVendors } from '#/lib/api'
import { slugify } from '#/lib/utils/slugify'
import { bareInput, form, inputWithPrefix, pill, pillActive, pillGroup, prefix } from './modals.css'

export type CreateVendorModalProps = {
  isOpen: boolean
  onClose: () => void
}

const AVAILABLE_COMPONENTS = [
  { id: 'AUDIO', label: 'Audio Dubbing' },
  { id: 'SUBTITLE', label: 'Subtitles & Captions' },
  { id: 'VIDEO', label: 'Video & QC Lab' },
  { id: 'METADATA', label: 'Metadata & Packaging' },
] as const

const AVAILABLE_MARKETS = [
  { id: 'en-US', label: 'en-US (English)' },
  { id: 'de-DE', label: 'de-DE (German)' },
  { id: 'fr-FR', label: 'fr-FR (French)' },
  { id: 'hi-IN', label: 'hi-IN (Hindi)' },
  { id: 'te-IN', label: 'te-IN (Telugu)' },
] as const

export function CreateVendorModal({ isOpen, onClose }: CreateVendorModalProps) {
  const queryClient = useQueryClient()

  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [isSlugManual, setIsSlugManual] = useState(false)
  const [selectedComponents, setSelectedComponents] = useState<string[]>(['AUDIO'])
  const [selectedMarkets, setSelectedMarkets] = useState<string[]>(['en-US'])
  const [hourlyRate, setHourlyRate] = useState('120')
  const [turnaround, setTurnaround] = useState('24')
  const [posterUrl, setPosterUrl] = useState<string | undefined>()

  const toggleComponent = (id: string) => {
    if (selectedComponents.includes(id)) {
      if (selectedComponents.length > 1) {
        setSelectedComponents(selectedComponents.filter((c) => c !== id))
      }
    } else {
      setSelectedComponents([...selectedComponents, id])
    }
  }

  const toggleMarket = (id: string) => {
    if (selectedMarkets.includes(id)) {
      setSelectedMarkets(selectedMarkets.filter((m) => m !== id))
    } else {
      setSelectedMarkets([...selectedMarkets, id])
    }
  }

  const handleNameChange = (val: string) => {
    setName(val)
    if (!isSlugManual) {
      setSlug(slugify(val))
    }
  }

  const handleSlugChange = (val: string) => {
    setIsSlugManual(true)
    setSlug(slugify(val))
  }

  const resetForm = () => {
    setName('')
    setSlug('')
    setIsSlugManual(false)
    setSelectedComponents(['AUDIO'])
    setSelectedMarkets(['en-US'])
    setHourlyRate('120')
    setTurnaround('24')
    setPosterUrl(undefined)
  }

  const mutation = useMutation({
    mutationFn: async () => {
      const finalId = slug ? `vendor-${slug}` : `vendor-${slugify(name)}`
      if (!name.trim()) throw new Error('Vendor name is required')
      if (!finalId.trim()) throw new Error('Vendor identifier is required')
      if (selectedComponents.length === 0) throw new Error('At least one component is required')

      const { data, error } = await postVendors({
        body: {
          id: finalId,
          name: name.trim(),
          components: selectedComponents,
          markets: selectedMarkets,
          hourly_rate_usd: hourlyRate ? parseFloat(hourlyRate) : 0,
          turnaround_hours: turnaround ? parseInt(turnaround, 10) : 24,
          metadata: posterUrl ? { poster_url: posterUrl } : {},
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
      description="Onboard an external dubbing studio, subtitling facility, or QC inspection lab with component and market coverage."
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose} disabled={mutation.isPending}>
            Cancel
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => mutation.mutate()}
            disabled={mutation.isPending || !name.trim() || selectedComponents.length === 0}
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
            placeholder="e.g. Deluxe Media"
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
              placeholder="deluxe-media"
              value={slug}
              onInput={(e) => handleSlugChange((e.target as HTMLInputElement).value)}
            />
          </div>
        </FormField>

        <FormField
          label="Covered Components"
          required
          helper="Select all media types this vendor processes"
        >
          <div class={pillGroup}>
            {AVAILABLE_COMPONENTS.map((comp) => {
              const active = selectedComponents.includes(comp.id)
              return (
                <button
                  key={comp.id}
                  type="button"
                  class={`${pill} ${active ? pillActive : ''}`}
                  onClick={() => toggleComponent(comp.id)}
                >
                  {comp.label}
                </button>
              )
            })}
          </div>
        </FormField>

        <FormField
          label="Market Coverage"
          helper="Select language-market territories this vendor supports (global for video)"
        >
          <div class={pillGroup}>
            {AVAILABLE_MARKETS.map((m) => {
              const active = selectedMarkets.includes(m.id)
              return (
                <button
                  key={m.id}
                  type="button"
                  class={`${pill} ${active ? pillActive : ''}`}
                  onClick={() => toggleMarket(m.id)}
                >
                  {m.label}
                </button>
              )
            })}
          </div>
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

        <FormField label="Facility Logo / Poster" optional helper="Max 1MB PNG, JPEG, WebP, or GIF">
          <ImageUpload value={posterUrl} onChange={setPosterUrl} disabled={mutation.isPending} />
        </FormField>
      </form>
    </Modal>
  )
}
