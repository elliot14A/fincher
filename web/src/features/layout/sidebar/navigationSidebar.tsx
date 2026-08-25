import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import {
  ChevronRight,
  FileText,
  LayoutGrid,
  MessageSquare,
  Play,
  Plus,
  Search,
  Users,
} from 'lucide-preact'
import { Logo } from '#/components/ui/logo'
import { titlesQueryOptions } from '#/features/titles'
import type { ModelsTitle } from '#/lib/api'
import {
  activeDot,
  brandRow,
  brandSubtitle,
  composeButton,
  kbdHint,
  navItem,
  navItemActive,
  navItemLabel,
  searchLabel,
  searchRow,
  sectionTitle,
  sidebarContainer,
  threadItem,
  threadItemActive,
  threadItemLabel,
  viewAllRow,
} from './navigationSidebar.css'

const NAV_LINKS = [
  { to: '/', label: 'Chat', icon: MessageSquare },
  { to: '/titles', label: 'Titles', icon: FileText },
  { to: '/deliveries', label: 'Deliveries', icon: LayoutGrid },
  { to: '/vendors', label: 'Vendors', icon: Users },
  { to: '/runs', label: 'Runs', icon: Play },
] as const

type RecentThreadItem = {
  label: string
  active: boolean
  hasHold: boolean
}

export function NavigationSidebar() {
  const { data: titlesResult } = useQuery(titlesQueryOptions({ limit: 4 }))
  const titles = titlesResult?.items ?? []

  const recentItems: RecentThreadItem[] = titles.slice(0, 4).map((t: ModelsTitle, idx: number) => ({
    label: `${t.name} release integrity`,
    active: idx === 0,
    hasHold: t.overall_status === 'HOLD' || t.overall_status === 'AT_RISK',
  }))

  return (
    <aside class={sidebarContainer}>
      <div class={brandRow}>
        <div>
          <Logo size="md" />
          <div class={brandSubtitle}>LUME Studios</div>
        </div>
      </div>

      <Link to="/" className={composeButton}>
        <Plus size={14} />
        <span>New chat</span>
      </Link>

      <div class={searchRow}>
        <Search size={14} />
        <span class={searchLabel}>Search</span>
        <span class={kbdHint}>⌘K</span>
      </div>

      {NAV_LINKS.map(({ to, label, icon: Icon }) => (
        <Link
          key={to}
          to={to}
          className={navItem}
          activeProps={{ className: `${navItem} ${navItemActive}` }}
        >
          <Icon size={15} />
          <span class={navItemLabel}>{label}</span>
        </Link>
      ))}

      {recentItems.length > 0 ? (
        <>
          <div class={sectionTitle}>Recent</div>

          {recentItems.map((thread: RecentThreadItem) => (
            <div
              key={thread.label}
              class={thread.active ? `${threadItem} ${threadItemActive}` : threadItem}
            >
              <span class={threadItemLabel}>{thread.label}</span>
              {(thread.active || thread.hasHold) && <span class={activeDot} />}
            </div>
          ))}

          <div class={viewAllRow}>
            <span>View all</span>
            <ChevronRight size={14} />
          </div>
        </>
      ) : null}
    </aside>
  )
}
