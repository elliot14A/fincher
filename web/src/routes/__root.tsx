import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Toaster } from 'sonner'
import { AppShell } from '#/features/layout/shell'

export const Route = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <AppShell>
      <Outlet />
      <Toaster theme="dark" position="bottom-right" richColors />
    </AppShell>
  )
}
