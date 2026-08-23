import { pageButton, pageInfo, paginationBar } from './paginationControls.css'

export type PaginationControlsProps = {
  page: number
  totalPages: number
  hasNextPage: boolean
  hasPrevPage: boolean
  onPrevPage: () => void
  onNextPage: () => void
}

export function PaginationControls({
  page,
  totalPages,
  hasNextPage,
  hasPrevPage,
  onPrevPage,
  onNextPage,
}: PaginationControlsProps) {
  return (
    <div class={paginationBar}>
      <button type="button" onClick={onPrevPage} disabled={!hasPrevPage} class={pageButton}>
        Previous
      </button>

      <span class={pageInfo}>
        Page {page} of {totalPages}
      </span>

      <button type="button" onClick={onNextPage} disabled={!hasNextPage} class={pageButton}>
        Next
      </button>
    </div>
  )
}
