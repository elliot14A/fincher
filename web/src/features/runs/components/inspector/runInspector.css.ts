import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

const pulseGlow = keyframes({
  '0%, 100%': { opacity: '1', transform: 'scale(1)' },
  '50%': { opacity: '0.4', transform: 'scale(0.85)' },
})

export const inspectorPanel = style({
  width: '540px',
  backgroundColor: vars.color.surface,
  borderLeft: `1px solid ${vars.color.borderSubtle}`,
  display: 'flex',
  flexDirection: 'column',
  flexShrink: 0,
  height: '100%',
  minHeight: 0,
  overflowY: 'auto',
})

export const header = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.md} ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  gap: vars.space.sm,
  flexShrink: 0,
  backgroundColor: vars.color.surfaceElevated,
})

export const headerTopRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.md,
})

export const headerTopRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  flexShrink: 0,
})

export const closeBtn = style({
  background: 'transparent',
  border: 'none',
  padding: vars.space.xs,
  color: vars.color.textTertiary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'color 0.1s ease, background-color 0.1s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const titleStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
})

export const runTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const runIdRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const runIdText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const copyBtn = style({
  background: 'none',
  border: 'none',
  padding: '2px 4px',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  display: 'inline-flex',
  alignItems: 'center',
  gap: '3px',
  fontSize: vars.fontSize['3xs'],
  transition: 'all 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const statBar = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(4, 1fr)',
  gap: vars.space.xs,
  paddingTop: vars.space.xs,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const statCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
})

export const statLabel = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
})

export const statValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const tabsNav = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `0 ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surface,
  flexShrink: 0,
})

export const tabBtn = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  padding: `${vars.space.sm} 0`,
  cursor: 'pointer',
  borderBottom: '2px solid transparent',
  background: 'none',
  border: 'none',
  fontWeight: 500,
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  transition: 'color 0.12s ease',
  ':hover': {
    color: vars.color.textSecondary,
  },
})

export const tabBtnActive = style({
  color: vars.color.textPrimary,
  fontWeight: 600,
  borderBottom: `2px solid ${vars.color.primary}`,
})

export const tabCountBadge = style({
  fontSize: vars.fontSize['3xs'],
  padding: '1px 5px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.surfaceElevated,
  color: vars.color.textSecondary,
  border: `1px solid ${vars.color.border}`,
})

export const scrollArea = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  minHeight: 0,
  overflowY: 'auto',
  padding: vars.space.lg,
  gap: vars.space.md,
})

export const section = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const sectionHeading = style({
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
  color: vars.color.textTertiary,
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  margin: 0,
})

export const rawJsonHeading = style([
  sectionHeading,
  {
    marginTop: vars.space.sm,
  },
])

// Waterfall & Steps
export const waterfallList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const stepCard = style({
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  overflow: 'hidden',
  transition: 'border-color 0.15s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
  },
})

export const stepHeaderBtn = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  width: '100%',
  textAlign: 'left',
  background: 'none',
  border: 'none',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  cursor: 'pointer',
  color: 'inherit',
  font: 'inherit',
})

export const stepHeaderLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  minWidth: 0,
})

export const stepIconBox = style({
  width: '24px',
  height: '24px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceHover,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.textSecondary,
  flexShrink: 0,
})

export const stepIndex = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  width: '14px',
})

export const stepName = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const stepCategoryBadge = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.teal,
  backgroundColor: vars.color.tealMuted,
  padding: '1px 5px',
  borderRadius: vars.radii.xs,
  fontWeight: 500,
})

export const stepHeaderRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  flexShrink: 0,
})

export const latencyBarContainer = style({
  width: '45px',
  height: '4px',
  backgroundColor: vars.color.surfaceActive,
  borderRadius: vars.radii.full,
  overflow: 'hidden',
})

export const latencyBarFill = style({
  height: '100%',
  backgroundColor: vars.color.primary,
  borderRadius: vars.radii.full,
})

export const latencyText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  fontVariantNumeric: 'tabular-nums',
  minWidth: '40px',
  textAlign: 'right',
})

export const statusIconSuccess = style({
  color: vars.color.success,
})

export const statusIconDanger = style({
  color: vars.color.danger,
})

export const statusIconClock = style({
  color: vars.color.textTertiary,
})

export const stepBody = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.xs} ${vars.space.sm} ${vars.space.sm} ${vars.space.sm}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.background,
  gap: vars.space.xs,
})

export const sqlBlock = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.teal,
  backgroundColor: vars.color.surfaceElevated,
  padding: vars.space.xs,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.border}`,
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-all',
  lineHeight: 1.4,
})

export const keyValGrid = style({
  display: 'grid',
  gridTemplateColumns: 'max-content minmax(0, 1fr)',
  gap: `${vars.space['2xs']} ${vars.space.sm}`,
  fontSize: vars.fontSize['2xs'],
  alignItems: 'baseline',
})

export const kvKey = style({
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
  fontWeight: 500,
})

export const kvVal = style({
  color: vars.color.textSecondary,
  wordBreak: 'break-word',
})

export const tagList = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space['2xs'],
})

export const pillTag = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textPrimary,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  padding: '1px 6px',
  borderRadius: vars.radii.xs,
})

// Policy Decisions
export const decisionList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
})

export const decisionCard = style({
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.primaryBorder}`,
  borderRadius: vars.radii.xs,
  padding: vars.space.md,
  gap: vars.space.sm,
})

export const decisionHeader = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
})

export const judgeTitle = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.primary,
  letterSpacing: '0.02em',
})

export const attemptPill = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  fontFamily: fonts.mono,
})

export const outcomeRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const outcomeLabel = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  fontWeight: 500,
})

export const rationaleCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  backgroundColor: vars.color.background,
  borderLeft: `2px solid ${vars.color.primary}`,
  border: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.xs,
  padding: vars.space.sm,
})

export const rationaleText = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: 1.5,
  margin: 0,
})

export const decisionMeta = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

// Context & Telemetry
export const contextCard = style({
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  padding: vars.space.md,
  gap: vars.space.sm,
})

export const rawJsonPre = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  backgroundColor: vars.color.background,
  padding: vars.space.sm,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.borderSubtle}`,
  overflowX: 'auto',
  maxHeight: '260px',
  lineHeight: 1.45,
})

export const pulseDot = style({
  width: '6px',
  height: '6px',
  borderRadius: '50%',
  backgroundColor: vars.color.warning,
  display: 'inline-block',
  marginRight: '6px',
  animation: `${pulseGlow} 1.5s infinite ease-in-out`,
})

export const emptyState = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  textAlign: 'center',
  color: vars.color.textTertiary,
  gap: vars.space.sm,
  minHeight: '380px',
  flex: 1,
})

export const emptyTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 500,
  color: vars.color.textSecondary,
})

export const emptyText = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  maxWidth: '320px',
})

export const footer = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surfaceElevated,
  flexShrink: 0,
})
