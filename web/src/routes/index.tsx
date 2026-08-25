import { createFileRoute } from '@tanstack/react-router'
import { AlertTriangle, Clock, Film, Send, Users } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Logo } from '#/components/ui/logo'
import {
  centerArea,
  centerColumn,
  composer,
  composerInput,
  composerPrompt,
  heading,
  page,
  queryIcon,
  queryList,
  queryListLabel,
  queryRow,
  queryTag,
  queryText,
  sendButton,
  subtitle,
  title,
} from '#/styles/routes/index.css'

export const Route = createFileRoute('/')({
  component: OperationsChatPage,
})

const PROMPT_SUGGESTIONS = [
  {
    icon: Film,
    title: 'Audit release readiness for Eclipse Season 1',
    tag: 'titles',
  },
  {
    icon: AlertTriangle,
    title: 'Check audio sync drift on recent dub deliveries',
    tag: 'qc',
  },
  {
    icon: Users,
    title: 'Inspect Vendor A historical on-time and pass rate',
    tag: 'vendors',
  },
  {
    icon: Clock,
    title: 'List master cuts updated within the last 48 hours',
    tag: 'masters',
  },
]

function OperationsChatPage() {
  const [prompt, setPrompt] = useState('')

  const handleSelectQuery = (queryTextValue: string) => {
    setPrompt(queryTextValue)
  }

  const handleSubmit = (event: { preventDefault: () => void }) => {
    event.preventDefault()
    if (!prompt.trim()) return
    // Ready for autonomous agent streaming integration
  }

  return (
    <div class={page}>
      <div class={centerArea}>
        <div class={centerColumn}>
          <div class={heading}>
            <Logo size="lg" />
            <h1 class={title}>How can Fincher assist your release today?</h1>
            <p class={subtitle}>
              Query master cuts, investigate QC drift anomalies, or verify global territory
              packages.
            </p>
          </div>

          <form class={composer} onSubmit={handleSubmit}>
            <span class={composerPrompt}>{'>'}</span>
            <input
              type="text"
              value={prompt}
              onInput={(e) => setPrompt((e.target as HTMLInputElement).value)}
              placeholder="Ask anything about titles, vendor drift, or package runs..."
              class={composerInput}
            />
            <button type="submit" class={sendButton} aria-label="Send">
              <Send size={14} />
            </button>
          </form>

          <div>
            <div class={queryListLabel}>Suggested queries</div>
            <div class={queryList}>
              {PROMPT_SUGGESTIONS.map((item) => {
                const Icon = item.icon
                return (
                  <button
                    key={item.title}
                    type="button"
                    class={queryRow}
                    onClick={() => handleSelectQuery(item.title)}
                  >
                    <Icon class={queryIcon} size={14} />
                    <span class={queryText}>{item.title}</span>
                    <span class={queryTag}>{item.tag}</span>
                  </button>
                )
              })}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
