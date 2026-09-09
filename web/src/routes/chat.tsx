import { useQuery } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { AlertTriangle, Clock, Film, Plus, Send, Users } from 'lucide-preact'
import { useEffect, useMemo, useRef, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Logo } from '#/components/ui/logo'
import { type ChatMessageView, MentionPicker, MessageList } from '#/features/chat/components'
import { useObjectMentions, useSendMessage } from '#/features/chat/hooks'
import { chatDetailQueryOptions } from '#/features/chat/queryOptions'
import type { ModelsChatCitation, ModelsChatSource, ModelsSearchResult } from '#/lib/api'
import {
  centerArea,
  centerColumn,
  composer,
  composerDock,
  composerDockColumn,
  composerPrompt,
  composerShell,
  composerTextarea,
  conversationArea,
  conversationColumn,
  conversationHeader,
  conversationTitle,
  heading,
  newChatButton,
  page,
  queryIcon,
  queryList,
  queryListLabel,
  queryRow,
  queryTag,
  queryText,
  sendButton,
  sendButtonDisabled,
  subtitle,
  title,
} from '#/styles/routes/chat.css'

type ChatSearch = {
  session?: string
}

export const Route = createFileRoute('/chat')({
  validateSearch: (search: Record<string, unknown>): ChatSearch => ({
    session: typeof search.session === 'string' ? search.session : undefined,
  }),
  component: OperationsChatPage,
})

const PROMPT_SUGGESTIONS = [
  { icon: Film, title: 'Which titles are on HOLD and why?', tag: 'titles' },
  { icon: AlertTriangle, title: 'Which vendors are busy right now?', tag: 'load' },
  { icon: Users, title: 'Who is the best vendor for German audio dubbing?', tag: 'vendors' },
  { icon: Clock, title: 'Which vendor has the highest QC accuracy for audio?', tag: 'qc' },
]

let localMessageId = 0
const nextLocalId = () => {
  localMessageId += 1
  return `local-${localMessageId}`
}

function OperationsChatPage() {
  const { session: sessionParam } = Route.useSearch()
  const navigate = useNavigate()

  const [prompt, setPrompt] = useState('')
  const [localMessages, setLocalMessages] = useState<ChatMessageView[]>([])
  const [tagged, setTagged] = useState<ModelsSearchResult[]>([])
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const sendMessage = useSendMessage()

  const sessionQuery = useQuery({
    ...chatDetailQueryOptions(sessionParam ?? ''),
    enabled: Boolean(sessionParam),
  })

  useEffect(() => {
    if (!sessionParam) {
      setLocalMessages([])
      return
    }
    const session = sessionQuery.data
    if (!session?.messages) return
    setLocalMessages(
      session.messages.map((m) => ({
        id: m.id,
        role: m.role === 'assistant' ? 'assistant' : 'user',
        content: m.content ?? '',
        citations: m.citations ?? [],
        sources: m.sources ?? [],
      })),
    )
  }, [sessionParam, sessionQuery.data])

  const mentions = useObjectMentions({
    value: prompt,
    setValue: setPrompt,
    textareaRef,
    onSelectObject: (obj) => {
      setTagged((prev) =>
        prev.some((t) => t.id === obj.id && t.kind === obj.kind) ? prev : [...prev, obj],
      )
    },
  })

  const mentionIdentifiers = useMemo(
    () => tagged.map((t) => t.identifier ?? t.display_name ?? '').filter(Boolean),
    [tagged],
  )

  const autoGrow = () => {
    const el = textareaRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, 160)}px`
  }

  const handleSelectQuery = (queryTextValue: string) => {
    setPrompt(queryTextValue)
    textareaRef.current?.focus()
  }

  const handleNewChat = () => {
    setLocalMessages([])
    setPrompt('')
    setTagged([])
    navigate({ to: '/chat', search: {} })
    textareaRef.current?.focus()
  }

  const submit = () => {
    const trimmed = prompt.trim()
    if (!trimmed || sendMessage.isPending) return

    const taggedRefs = tagged
      .filter((t) => mentionIdentifiers.includes(t.identifier ?? t.display_name ?? ''))
      .map((t) => ({
        kind: t.kind ?? 'title',
        identifier: t.identifier ?? '',
        display_name: t.display_name ?? '',
      }))

    setLocalMessages((prev) => [
      ...prev,
      { id: nextLocalId(), role: 'user', content: trimmed, citations: [], sources: [] },
    ])
    setPrompt('')
    setTagged([])

    sendMessage.mutate(
      { message: trimmed, sessionId: sessionParam, tagged: taggedRefs },
      {
        onSuccess: (data) => {
          const answer = data.message
          setLocalMessages((prev) => [
            ...prev,
            {
              id: answer?.id ?? nextLocalId(),
              role: 'assistant',
              content: answer?.content ?? 'No answer was returned.',
              citations: (answer?.citations ?? []) as ModelsChatCitation[],
              sources: (answer?.sources ?? []) as ModelsChatSource[],
            },
          ])
          if (data.session_id && data.session_id !== sessionParam) {
            navigate({ to: '/chat', search: { session: data.session_id } })
          }
        },
        onError: (err) => {
          toast.error(err instanceof Error ? err.message : 'The assistant failed to respond')
        },
      },
    )
  }

  const handleSubmit = (event: { preventDefault: () => void }) => {
    event.preventDefault()
    submit()
  }

  const hasConversation = localMessages.length > 0
  const sendDisabled = !prompt.trim() || sendMessage.isPending

  const composerForm = (
    <div class={composerShell}>
      <MentionPicker
        isOpen={mentions.isOpen}
        query={mentions.query}
        results={mentions.results}
        isLoading={mentions.isLoading}
        activeIndex={mentions.activeIndex}
        onActiveIndexChange={mentions.setActiveIndex}
        onSelect={(i) => {
          mentions.setActiveIndex(i)
          mentions.selectActive()
        }}
      />
      <form class={composer} onSubmit={handleSubmit}>
        <span class={composerPrompt}>{'>'}</span>
        <textarea
          ref={textareaRef}
          value={prompt}
          rows={1}
          onInput={(e) => {
            setPrompt((e.target as HTMLTextAreaElement).value)
            autoGrow()
          }}
          onKeyDown={(e) => {
            mentions.handleKeyDown(e as unknown as KeyboardEvent)
            if (e.key === 'Enter' && !e.shiftKey && !mentions.isOpen) {
              e.preventDefault()
              submit()
            }
          }}
          placeholder="Ask anything — type @ to reference a title, vendor, or delivery..."
          class={composerTextarea}
          disabled={sendMessage.isPending}
        />
        <button
          type="submit"
          class={sendDisabled ? `${sendButton} ${sendButtonDisabled}` : sendButton}
          aria-label="Send"
          disabled={sendDisabled}
        >
          <Send size={14} />
        </button>
      </form>
    </div>
  )

  if (hasConversation) {
    return (
      <div class={page}>
        <div class={conversationHeader}>
          <span class={conversationTitle}>Operations assistant</span>
          <button type="button" class={newChatButton} onClick={handleNewChat}>
            <Plus size={12} />
            New chat
          </button>
        </div>

        <div class={conversationArea}>
          <div class={conversationColumn}>
            <MessageList
              messages={localMessages}
              pending={sendMessage.isPending}
              mentionIdentifiers={mentionIdentifiers}
            />
          </div>
        </div>

        <div class={composerDock}>
          <div class={composerDockColumn}>{composerForm}</div>
        </div>
      </div>
    )
  }

  return (
    <div class={page}>
      <div class={centerArea}>
        <div class={centerColumn}>
          <div class={heading}>
            <Logo size="lg" />
            <h1 class={title}>How can Fincher assist your release today?</h1>
            <p class={subtitle}>
              Ask about titles, vendor performance, QC drift, or premiere readiness. Every answer is
              grounded in live operational state and historical evidence.
            </p>
          </div>

          {composerForm}

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
