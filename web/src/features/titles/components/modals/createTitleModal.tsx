import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Button } from '#/components/ui/button'
import { FormField, ImageUpload, SelectInput, TextInput } from '#/components/ui/input'
import { Modal } from '#/components/ui/modal'
import { titlesKeys } from '#/features/titles/queryKeys'
import { postTitles } from '#/lib/api'
import type { ModelsTitleType } from '#/lib/api/generated'
import { slugify } from '#/lib/utils/slugify'
import { bareInput, form, formRow, inputWithPrefix, prefix } from './modals.css'

export type CreateTitleModalProps = {
  isOpen: boolean
  onClose: () => void
}

function getDefaultPremiereDate(): string {
  const d = new Date()
  d.setDate(d.getDate() + 14)
  d.setMinutes(d.getMinutes() - d.getTimezoneOffset())
  return d.toISOString().slice(0, 16)
}

export function CreateTitleModal({ isOpen, onClose }: CreateTitleModalProps) {
  const queryClient = useQueryClient()

  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [isSlugManual, setIsSlugManual] = useState(false)
  const [type, setType] = useState<ModelsTitleType>('FEATURE')
  const [premiereDate, setPremiereDate] = useState(getDefaultPremiereDate())
  const [territories, setTerritories] = useState('35')
  const [masterVersion, setMasterVersion] = useState('V01')
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
    setType('FEATURE')
    setPremiereDate(getDefaultPremiereDate())
    setTerritories('35')
    setMasterVersion('V01')
    setAvatarUrl(undefined)
  }

  const mutation = useMutation({
    mutationFn: async () => {
      const finalSlug = slug.trim() || slugify(name)
      const finalId = `title-${finalSlug}`
      if (!name.trim()) throw new Error('Title name is required')
      if (!finalId.trim()) throw new Error('Title identifier is required')

      const dateObj = new Date(premiereDate)
      const isoDate = Number.isNaN(dateObj.getTime())
        ? new Date().toISOString()
        : dateObj.toISOString()

      const territoryCount = Number.parseInt(territories, 10) || 1

      const { data, error } = await postTitles({
        body: {
          id: finalId,
          name: name.trim(),
          slug: finalSlug,
          type,
          premiere_date: isoDate,
          territories: territoryCount,
          current_master_version: masterVersion.trim() || 'V01',
          overall_status: 'PROCESSING',
          metadata: avatarUrl ? { avatar_url: avatarUrl } : {},
        },
      })

      if (error) {
        throw new Error(error.message || 'Failed to create title')
      }

      return data
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: titlesKeys.all })
      toast.success(`Title "${data?.name ?? name}" created successfully`)
      resetForm()
      onClose()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to register title')
    },
  })

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Create New Title"
      description="Register a theatrical, episodic, or special media release for production and localization."
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
            <span>{mutation.isPending ? 'Creating...' : 'Create Title'}</span>
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
        <FormField
          label="Title Name"
          required
          helper="The primary public display name of the release"
        >
          <TextInput
            placeholder="e.g. Avatar: Fire & Ash"
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
            <span class={prefix}>title-</span>
            <input
              type="text"
              class={bareInput}
              placeholder="fire-and-ash"
              value={slug}
              onInput={(e) => handleSlugChange((e.target as HTMLInputElement).value)}
            />
          </div>
        </FormField>

        <div class={formRow}>
          <FormField label="Title Type" required>
            <SelectInput
              value={type}
              onChange={(e) => setType((e.target as HTMLSelectElement).value as ModelsTitleType)}
              options={[
                { value: 'FEATURE', label: 'Feature Film' },
                { value: 'SERIES', label: 'Series / Episodic' },
                { value: 'SPECIAL', label: 'Special / Short' },
              ]}
            />
          </FormField>

          <FormField label="Markets Count" required helper="Number of global release markets">
            <TextInput
              type="number"
              min="1"
              max="250"
              value={territories}
              onInput={(e) => setTerritories((e.target as HTMLInputElement).value)}
            />
          </FormField>
        </div>

        <div class={formRow}>
          <FormField label="Premiere Schedule" required>
            <TextInput
              type="datetime-local"
              value={premiereDate}
              onInput={(e) => setPremiereDate((e.target as HTMLInputElement).value)}
            />
          </FormField>

          <FormField label="Initial Master Cut" required helper="e.g. V01">
            <TextInput
              placeholder="V01"
              value={masterVersion}
              onInput={(e) => setMasterVersion((e.target as HTMLInputElement).value)}
            />
          </FormField>
        </div>

        <FormField label="Poster Thumbnail" optional helper="Max 1MB PNG, JPEG, WebP, or GIF">
          <ImageUpload value={avatarUrl} onChange={setAvatarUrl} disabled={mutation.isPending} />
        </FormField>
      </form>
    </Modal>
  )
}
