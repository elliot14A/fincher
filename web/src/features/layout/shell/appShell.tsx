import { useRouterState } from '@tanstack/react-router'
import type { ComponentChildren } from 'preact'
import { NavigationSidebar } from '../sidebar'
import { shellContent, shellMainArea, shellRoot, shellRootLanding } from './appShell.css'

export type AppShellProps = {
  children?: ComponentChildren
}

export function AppShell({ children }: AppShellProps) {
  const pathname = useRouterState({ select: (s) => s.location.pathname })
  const isLanding = pathname === '/'

  if (isLanding) {
    return <div class={shellRootLanding}>{children}</div>
  }

  return (
    <div class={shellRoot}>
      <NavigationSidebar />
      <div class={shellMainArea}>
        <main class={shellContent}>{children}</main>
      </div>
    </div>
  )
}
