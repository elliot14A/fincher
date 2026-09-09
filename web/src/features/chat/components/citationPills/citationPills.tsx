import { Database, Terminal } from 'lucide-preact'
import type { ModelsChatCitation } from '#/lib/api'
import { label, list, pill, pillIcon, pillText, wrapper } from './citationPills.css'

type CitationPillsProps = {
  citations: ModelsChatCitation[]
}

export function CitationPills({ citations }: CitationPillsProps) {
  if (citations.length === 0) {
    return null
  }

  return (
    <div class={wrapper}>
      <span class={label}>
        Grounded in {citations.length} {citations.length === 1 ? 'query' : 'queries'}
      </span>
      <div class={list}>
        {citations.map((citation, index) => {
          const isClickhouse = citation.source === 'clickhouse'
          const Icon = isClickhouse ? Database : Terminal
          return (
            <span class={pill} key={`${citation.tool}-${index}`} title={citation.label}>
              <Icon class={pillIcon} size={11} />
              <span class={pillText}>{citation.label}</span>
            </span>
          )
        })}
      </div>
    </div>
  )
}
