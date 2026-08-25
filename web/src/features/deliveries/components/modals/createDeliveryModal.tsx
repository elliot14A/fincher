import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Plus } from 'lucide-preact'
import { useEffect, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Button } from '#/components/ui/button'
import { FormField, SelectInput, TextInput } from '#/components/ui/input'
import { Modal } from '#/components/ui/modal'
import { deliveriesKeys } from '#/features/deliveries/queryKeys'
import { titlesQueryOptions } from '#/features/titles/queryOptions'
import { postDeliveries } from '#/lib/api'
import type { ModelsDeliveryStatus } from '#/lib/api/generated'
import { countryChip, countryChipActive, form, formRow, quickCountryChips } from './modals.css'

const POPULAR_COUNTRIES = [
  { code: 'US', name: 'United States' },
  { code: 'DE', name: 'Germany' },
  { code: 'JP', name: 'Japan' },
  { code: 'FR', name: 'France' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'ES', name: 'Spain' },
  { code: 'AU', name: 'Australia' },
  { code: 'BR', name: 'Brazil' },
  { code: 'CA', name: 'Canada' },
  { code: 'IT', name: 'Italy' },
  { code: 'KR', name: 'South Korea' },
  { code: 'MX', name: 'Mexico' },
]

export type CreateDeliveryModalProps = {
  isOpen: boolean
  onClose: () => void
  initialTitleId?: string
}

function getDefaultDeliveryDate(): string {
  const d = new Date()
  d.setDate(d.getDate() + 7)
  d.setMinutes(d.getMinutes() - d.getTimezoneOffset())
  return d.toISOString().slice(0, 16)
}

export function CreateDeliveryModal({ isOpen, onClose, initialTitleId }: CreateDeliveryModalProps) {
  const queryClient = useQueryClient()

  const { data: titlesData } = useQuery(titlesQueryOptions({ limit: 100 }))
  const titleOptions =
    titlesData?.items?.map((t) => ({
      value: t.id,
      label: `${t.name} (${t.id})`,
    })) ?? []

  const [titleId, setTitleId] = useState(initialTitleId || '')
  const [country, setCountry] = useState('US')
  const [targetDate, setTargetDate] = useState(getDefaultDeliveryDate())
  const [status, setStatus] = useState<ModelsDeliveryStatus>('PENDING')

  useEffect(() => {
    if (initialTitleId) {
      setTitleId(initialTitleId)
    } else if (titleOptions.length > 0 && !titleId) {
      setTitleId(titleOptions[0].value)
    }
  }, [initialTitleId, titleOptions, titleId])

  const calculatedId =
    titleId && country ? `del-${titleId.replace(/^title-/, '')}-${country.toLowerCase()}` : ''

  const resetForm = () => {
    if (titleOptions.length > 0) {
      setTitleId(titleOptions[0].value)
    }
    setCountry('US')
    setTargetDate(getDefaultDeliveryDate())
    setStatus('PENDING')
  }

  const mutation = useMutation({
    mutationFn: async () => {
      if (!titleId) throw new Error('Please select a Title')
      if (!country.trim()) throw new Error('Country code is required')

      const dateObj = new Date(targetDate)
      const isoDate = Number.isNaN(dateObj.getTime())
        ? new Date().toISOString()
        : dateObj.toISOString()

      const { data, error } = await postDeliveries({
        body: {
          id: calculatedId,
          title_id: titleId,
          country: country.trim().toUpperCase(),
          target_date: isoDate,
          status,
        },
      })

      if (error) {
        throw new Error(error.message || 'Failed to create territory delivery')
      }

      return data
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: deliveriesKeys.all })
      toast.success(`Delivery for market ${data?.country ?? country} created`)
      resetForm()
      onClose()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to create delivery')
    },
  })

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Schedule Market Delivery"
      description="Create a localized delivery commitment with a target shipping milestone."
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose} disabled={mutation.isPending}>
            Cancel
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => mutation.mutate()}
            disabled={mutation.isPending || !titleId || !country}
          >
            <Plus size={14} />
            <span>{mutation.isPending ? 'Scheduling...' : 'Schedule Delivery'}</span>
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
        <FormField label="Target Title" required helper="Select the release being delivered">
          <SelectInput
            value={titleId}
            onChange={(e) => setTitleId((e.target as HTMLSelectElement).value)}
            options={
              titleOptions.length > 0
                ? titleOptions
                : [{ value: '', label: 'No titles available (create a title first)' }]
            }
          />
        </FormField>

        <FormField label="Market Country Code" required helper="ISO-3166 2-letter country code">
          <TextInput
            placeholder="e.g. US, DE, JP"
            maxLength={2}
            value={country}
            onInput={(e) => setCountry((e.target as HTMLInputElement).value.toUpperCase())}
          />
          <div class={quickCountryChips}>
            {POPULAR_COUNTRIES.map((c) => (
              <button
                type="button"
                key={c.code}
                class={country === c.code ? `${countryChip} ${countryChipActive}` : countryChip}
                onClick={() => setCountry(c.code)}
              >
                {c.code} {c.name}
              </button>
            ))}
          </div>
        </FormField>

        <div class={formRow}>
          <FormField label="Target Shipping Date" required>
            <TextInput
              type="datetime-local"
              value={targetDate}
              onInput={(e) => setTargetDate((e.target as HTMLInputElement).value)}
            />
          </FormField>

          <FormField label="Initial Status" required>
            <SelectInput
              value={status}
              onChange={(e) =>
                setStatus((e.target as HTMLSelectElement).value as ModelsDeliveryStatus)
              }
              options={[
                { value: 'PENDING', label: 'Pending Processing' },
                { value: 'READY_TO_SHIP', label: 'Ready To Ship' },
                { value: 'HOLD', label: 'Hold / Blocked' },
              ]}
            />
          </FormField>
        </div>

        <FormField label="Derived Delivery ID" helper="System-generated identifier">
          <TextInput value={calculatedId} disabled />
        </FormField>
      </form>
    </Modal>
  )
}
