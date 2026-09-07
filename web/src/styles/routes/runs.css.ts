import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

const pulseGlow = keyframes({
  '0%, 100%': { opacity: '1', transform: 'scale(1)' },
  '50%': { opacity: '0.4', transform: 'scale(0.85)' },
})

export const page = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  height: '100%',
  overflow: 'hidden',
})

export const header = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  padding: `${vars.space.lg} ${vars.space['2xl']}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  flexShrink: 0,
})

export const pageTitle = style({
  fontSize: vars.fontSize.lg,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const pageSubtitle = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textTertiary,
  marginTop: vars.space['3xs'],
  display: 'block',
})

export const toolbar = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.lg,
  padding: `${vars.space.sm} ${vars.space['2xl']}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  flexShrink: 0,
})

export const toolbarGroup = style({
  display: 'flex',
  gap: vars.space.md,
})

export const toolbarTab = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  padding: `${vars.space['2xs']} 0`,
  cursor: 'pointer',
  borderBottom: '2px solid transparent',
  background: 'none',
  border: 'none',
  transition: 'color 0.12s ease',
  ':hover': {
    color: vars.color.textSecondary,
  },
})

export const toolbarTabActive = style({
  color: vars.color.textPrimary,
  fontWeight: 500,
  borderBottom: `2px solid ${vars.color.primary}`,
})

export const contentLayout = style({
  display: 'flex',
  flex: 1,
  minHeight: 0,
  overflow: 'hidden',
})

export const mainListContainer = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  minWidth: 0,
  height: '100%',
  overflow: 'hidden',
})

export const list = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.xs} ${vars.space['2xl']}`,
  flex: 1,
  minHeight: 0,
  overflowY: 'auto',
})

export const row = style({
  position: 'relative',
  display: 'grid',
  gridTemplateColumns: '36px minmax(0, 1fr) 110px minmax(100px, auto)',
  alignItems: 'center',
  columnGap: vars.space.md,
  padding: `${vars.space.md} ${vars.space.sm}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  cursor: 'pointer',
  transition: 'background-color 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
  },
})

export const rowActive = style({
  backgroundColor: vars.color.surfaceElevated,
  '::before': {
    content: '""',
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
    width: '2px',
    backgroundColor: vars.color.primary,
  },
})

export const runIcon = style({
  width: '32px',
  height: '32px',
  borderRadius: vars.radii.xs,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textSecondary,
  flexShrink: 0,
})

export const nameStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '3px',
  minWidth: 0,
})

export const cardName = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const metaRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  overflow: 'hidden',
  flexWrap: 'wrap',
})

export const metaTrigger = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.teal,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
  flexShrink: 0,
})

export const sublabelTag = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  backgroundColor: vars.color.surfaceElevated,
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.border}`,
})

export const metaDivider = style({
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.xs,
  flexShrink: 0,
})

export const metaDate = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  fontVariantNumeric: 'tabular-nums',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const statusStack = style({
  display: 'flex',
  alignItems: 'center',
})

export const timeStack = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-end',
  gap: '2px',
})

export const timeValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  fontVariantNumeric: 'tabular-nums',
})

export const stepsCount = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
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
  maxWidth: '360px',
})

export const loadingState = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.sm,
  minHeight: '380px',
  flex: 1,
})
