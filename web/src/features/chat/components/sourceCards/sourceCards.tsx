import { Film } from 'lucide-preact'
import type { ModelsChatSource } from '#/lib/api'
import { body, card, label, list, meta, name, poster, statusPill, wrapper } from './sourceCards.css'

type SourceCardsProps = {
  sources: ModelsChatSource[]
}

export function SourceCards({ sources }: SourceCardsProps) {
  const withImages = sources.filter((s) => s.kind === 'title')
  if (withImages.length === 0) {
    return null
  }

  return (
    <div class={wrapper}>
      <span class={label}>Sources</span>
      <div class={list}>
        {withImages.map((source) => (
          <div class={card} key={source.id}>
            {source.image_url ? (
              <img class={poster} src={source.image_url} alt={source.display_name} loading="lazy" />
            ) : (
              <span class={poster}>
                <Film size={18} />
              </span>
            )}
            <div class={body}>
              <span class={name}>{source.display_name}</span>
              {source.subtitle ? <span class={meta}>{source.subtitle}</span> : null}
              {source.status ? <span class={statusPill}>{source.status}</span> : null}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
