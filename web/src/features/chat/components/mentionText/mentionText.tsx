import { AtSign } from 'lucide-preact'
import { mention, paragraph } from './mentionText.css'

export type MentionSegment = { value: string; mention: boolean }

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

export function tokenizeMentions(text: string, identifiers: string[]): MentionSegment[] {
  const unique = Array.from(new Set(identifiers.filter((value) => value.trim().length > 0))).sort(
    (a, b) => b.length - a.length,
  )
  if (unique.length === 0) return [{ value: text, mention: false }]

  const pattern = new RegExp(`@(?:${unique.map(escapeRegExp).join('|')})(?![A-Za-z0-9_-])`, 'g')
  const segments: MentionSegment[] = []
  let lastIndex = 0

  for (const match of text.matchAll(pattern)) {
    const start = match.index ?? 0
    const prevChar = start > 0 ? text[start - 1] : ''
    if (prevChar && /[A-Za-z0-9_]/.test(prevChar)) continue

    if (start > lastIndex) segments.push({ value: text.slice(lastIndex, start), mention: false })
    segments.push({ value: match[0], mention: true })
    lastIndex = start + match[0].length
  }

  if (lastIndex < text.length) segments.push({ value: text.slice(lastIndex), mention: false })
  return segments
}

type MentionTextProps = {
  text: string
  identifiers: string[]
}

export function MentionText({ text, identifiers }: MentionTextProps) {
  const segments = tokenizeMentions(text, identifiers)
  return (
    <p class={paragraph}>
      {segments.map((segment, index) =>
        segment.mention ? (
          <span class={mention} key={`mention-${index}-${segment.value}`}>
            <AtSign size={11} />
            {segment.value.slice(1)}
          </span>
        ) : (
          <span key={`text-${index}-${segment.value.slice(0, 8)}`}>{segment.value}</span>
        ),
      )}
    </p>
  )
}
