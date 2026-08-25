import type { UseQueryOptions } from '@tanstack/react-query'
import { useQuery } from '@tanstack/react-query'
import { useState } from 'preact/hooks'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'

export interface TabItem<TId extends string = string> {
  id: TId
  label: string
}

export interface PaginationEnvelope<TItem> {
  items?: TItem[]
  total_items?: number
  page?: number
  limit?: number
  total_pages?: number
  has_next_page?: boolean
  has_prev_page?: boolean
}

// biome-ignore lint/suspicious/noExplicitAny: supports arbitrary TanStack Query queryOptions shapes
export type AnyQueryOptions = UseQueryOptions<any, any, any, any>

export interface UseTabbedQueryListOptions<
  TTabId extends string,
  TOptions extends AnyQueryOptions = AnyQueryOptions,
> {
  tabs: readonly TabItem<TTabId>[]
  initialTab?: TTabId
  initialPage?: number
  limit?: number
  sortOrder?: string
  buildQueryOptions: (args: {
    tab: TTabId
    filter: TTabId extends 'ALL' ? undefined : TTabId
    page: number
    limit: number
    sort_order: string
  }) => TOptions
}

export interface UseTabbedQueryListResult<
  TItem extends { id: string },
  TTabId extends string,
  TData extends PaginationEnvelope<TItem> = PaginationEnvelope<TItem>,
> {
  activeTab: TTabId
  onTabChange: (tab: TTabId) => void
  selectedId: string | null
  setSelectedId: (id: string | null) => void
  currentSelectedId: string | null
  page: number
  setPage: (page: number | ((prev: number) => number)) => void
  onPrevPage: () => void
  onNextPage: () => void
  data: TData | undefined
  items: TItem[]
  totalPages: number
  hasNextPage: boolean
  hasPrevPage: boolean
  isLoading: boolean
  isError: boolean
  error: Error | null
}

export function useTabbedQueryList<
  TItem extends { id: string },
  TTabId extends string,
  TOptions extends AnyQueryOptions = AnyQueryOptions,
  TData extends PaginationEnvelope<TItem> = PaginationEnvelope<TItem>,
>({
  tabs,
  initialTab,
  initialPage = DEFAULT_PAGE,
  limit = DEFAULT_PAGE_LIMIT,
  sortOrder = DEFAULT_SORT_ORDER,
  buildQueryOptions,
}: UseTabbedQueryListOptions<TTabId, TOptions>): UseTabbedQueryListResult<TItem, TTabId, TData> {
  const defaultTab = initialTab ?? tabs[0]?.id ?? ('ALL' as TTabId)
  const [activeTab, setActiveTab] = useState<TTabId>(defaultTab)
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [page, setPage] = useState(initialPage)

  const filter = (activeTab === 'ALL' ? undefined : activeTab) as TTabId extends 'ALL'
    ? undefined
    : TTabId

  const queryOptions = buildQueryOptions({
    tab: activeTab,
    filter,
    page,
    limit,
    sort_order: sortOrder,
  })

  const { data, isLoading, isError, error } = useQuery(queryOptions)

  const items = ((data as PaginationEnvelope<TItem> | undefined)?.items ?? []) as TItem[]
  const totalPages = (data as PaginationEnvelope<TItem> | undefined)?.total_pages ?? 1
  const hasNextPage = (data as PaginationEnvelope<TItem> | undefined)?.has_next_page ?? false
  const hasPrevPage = (data as PaginationEnvelope<TItem> | undefined)?.has_prev_page ?? false

  const isSelectedPresent = items.some((i) => i.id === selectedId)
  const currentSelectedId = isSelectedPresent ? selectedId : (items[0]?.id ?? null)

  const onTabChange = (nextTab: TTabId) => {
    setActiveTab(nextTab)
    setSelectedId(null)
    setPage(DEFAULT_PAGE)
  }

  const onPrevPage = () => {
    setPage((p) => Math.max(1, p - 1))
  }

  const onNextPage = () => {
    setPage((p) => Math.min(totalPages, p + 1))
  }

  return {
    activeTab,
    onTabChange,
    selectedId,
    setSelectedId,
    currentSelectedId,
    page,
    setPage,
    onPrevPage,
    onNextPage,
    data: data as TData | undefined,
    items,
    totalPages,
    hasNextPage,
    hasPrevPage,
    isLoading,
    isError,
    error: (error as Error) ?? null,
  }
}
