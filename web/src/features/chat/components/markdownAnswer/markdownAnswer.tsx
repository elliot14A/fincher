import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { markdown } from './markdownAnswer.css'

type MarkdownAnswerProps = {
  content: string
}

export function MarkdownAnswer({ content }: MarkdownAnswerProps) {
  return (
    <div class={markdown}>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
    </div>
  )
}
