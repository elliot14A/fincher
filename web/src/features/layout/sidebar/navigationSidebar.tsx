import { Link } from '@tanstack/react-router'
import { FileText, LayoutGrid, MessageSquare, Play, Plus, Search, Users } from 'lucide-preact'
import { Logo } from '#/components/ui/logo'
import {
  brandRow,
  brandSubtitle,
  composeButton,
  kbdHint,
  navItem,
  navItemActive,
  navItemLabel,
  searchLabel,
  searchRow,
  sidebarContainer,
} from './navigationSidebar.css'

const NAV_LINKS = [
  { to: '/chat', label: 'Chat', icon: MessageSquare },
  { to: '/titles', label: 'Titles', icon: FileText },
  { to: '/deliveries', label: 'Deliveries', icon: LayoutGrid },
  { to: '/vendors', label: 'Vendors', icon: Users },
  { to: '/runs', label: 'Runs', icon: Play },
] as const

export function NavigationSidebar() {
  return (
    <aside class={sidebarContainer}>
      <div class={brandRow}>
        <div>
          <Logo size="md" />
          <div class={brandSubtitle}>LUME Studios</div>
        </div>
      </div>

      <Link to="/chat" className={composeButton}>
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
    </aside>
  )
}
