import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { Check, MessageSquare, Pencil, Trash2, X } from 'lucide-preact'
import { useEffect, useRef, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { useDeleteSession, useRenameSession } from '../../hooks'
import { chatListQueryOptions } from '../../queryOptions'
import {
  emptyState,
  renameInput,
  row,
  rowAction,
  rowActive,
  rowButton,
  rowTitle,
  scrollArea,
  sectionLabel,
  wrapper,
} from './chatHistory.css'

type ChatHistoryProps = {
  activeSessionId?: string
}

export function ChatHistory({ activeSessionId }: ChatHistoryProps) {
  const navigate = useNavigate()
  const { data: sessions } = useQuery(chatListQueryOptions())
  const renameSession = useRenameSession()
  const deleteSession = useDeleteSession()

  const [editingId, setEditingId] = useState<string | null>(null)
  const [draftTitle, setDraftTitle] = useState('')
  const renameInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (editingId) {
      renameInputRef.current?.focus()
      renameInputRef.current?.select()
    }
  }, [editingId])

  const startRename = (id: string, currentTitle: string) => {
    setEditingId(id)
    setDraftTitle(currentTitle)
  }

  const commitRename = (id: string) => {
    const trimmed = draftTitle.trim()
    setEditingId(null)
    if (!trimmed) return
    renameSession.mutate(
      { id, title: trimmed },
      { onError: (err) => toast.error(err instanceof Error ? err.message : 'Rename failed') },
    )
  }

  const handleDelete = (id: string) => {
    deleteSession.mutate(id, {
      onSuccess: () => {
        if (id === activeSessionId) {
          navigate({ to: '/chat', search: {} })
        }
      },
      onError: (err) => toast.error(err instanceof Error ? err.message : 'Delete failed'),
    })
  }

  return (
    <div class={wrapper}>
      <span class={sectionLabel}>Recent chats</span>
      <div class={scrollArea}>
        {!sessions || sessions.length === 0 ? (
          <span class={emptyState}>No conversations yet. Ask a question to start one.</span>
        ) : (
          sessions.map((session) => {
            const isActive = session.id === activeSessionId
            const isEditing = session.id === editingId
            const label = session.title?.trim() || 'Untitled chat'

            if (isEditing) {
              return (
                <div class={isActive ? `${row} ${rowActive}` : row} key={session.id}>
                  <input
                    ref={renameInputRef}
                    type="text"
                    class={renameInput}
                    value={draftTitle}
                    onInput={(e) => setDraftTitle((e.target as HTMLInputElement).value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') commitRename(session.id)
                      if (e.key === 'Escape') setEditingId(null)
                    }}
                  />
                  <button
                    type="button"
                    class={rowAction}
                    aria-label="Save"
                    onClick={() => commitRename(session.id)}
                  >
                    <Check size={13} />
                  </button>
                  <button
                    type="button"
                    class={rowAction}
                    aria-label="Cancel"
                    onClick={() => setEditingId(null)}
                  >
                    <X size={13} />
                  </button>
                </div>
              )
            }

            return (
              <div class={isActive ? `${row} ${rowActive}` : row} key={session.id}>
                <button
                  type="button"
                  class={rowButton}
                  onClick={() => navigate({ to: '/chat', search: { session: session.id } })}
                >
                  <MessageSquare size={13} />
                  <span class={rowTitle}>{label}</span>
                </button>
                <button
                  type="button"
                  class={rowAction}
                  aria-label="Rename chat"
                  onClick={() => startRename(session.id, label)}
                >
                  <Pencil size={12} />
                </button>
                <button
                  type="button"
                  class={rowAction}
                  aria-label="Delete chat"
                  onClick={() => handleDelete(session.id)}
                >
                  <Trash2 size={12} />
                </button>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
