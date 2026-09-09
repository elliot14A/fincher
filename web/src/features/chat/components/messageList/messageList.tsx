import type { ModelsChatCitation, ModelsChatSource } from '#/lib/api'
import { CitationPills } from '../citationPills'
import { MarkdownAnswer } from '../markdownAnswer'
import { MentionText } from '../mentionText'
import { SourceCards } from '../sourceCards'
import {
  assistantBubble,
  assistantRow,
  list,
  roleLabel,
  thinking,
  thinkingDot,
  userBubble,
  userRow,
} from './messageList.css'

export type ChatMessageView = {
  id: string
  role: 'user' | 'assistant'
  content: string
  citations: ModelsChatCitation[]
  sources: ModelsChatSource[]
}

type MessageListProps = {
  messages: ChatMessageView[]
  pending: boolean
  mentionIdentifiers: string[]
}

export function MessageList({ messages, pending, mentionIdentifiers }: MessageListProps) {
  return (
    <div class={list}>
      {messages.map((message) =>
        message.role === 'user' ? (
          <div class={userRow} key={message.id}>
            <span class={roleLabel}>Operator</span>
            <div class={userBubble}>
              <MentionText text={message.content} identifiers={mentionIdentifiers} />
            </div>
          </div>
        ) : (
          <div class={assistantRow} key={message.id}>
            <span class={roleLabel}>Fincher</span>
            <div class={assistantBubble}>
              <MarkdownAnswer content={message.content} />
              <SourceCards sources={message.sources} />
              <CitationPills citations={message.citations} />
            </div>
          </div>
        ),
      )}

      {pending ? (
        <div class={assistantRow}>
          <span class={roleLabel}>Fincher</span>
          <div class={assistantBubble}>
            <span class={thinking}>
              <span class={thinkingDot} />
              Querying operational state and history...
            </span>
          </div>
        </div>
      ) : null}
    </div>
  )
}
