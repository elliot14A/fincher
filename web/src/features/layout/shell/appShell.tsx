import type { ComponentChildren } from 'preact'
import { NavigationSidebar } from '../sidebar'
import { shellContent, shellMainArea, shellRoot } from './appShell.css'

export type AppShellProps = {
  children?: ComponentChildren
}

export function AppShell({ children }: AppShellProps) {
  return (
    <div class={shellRoot}>
      <NavigationSidebar />
      <div class={shellMainArea}>
        <main class={shellContent}>{children}</main>
      </div>
    </div>
  )
}
