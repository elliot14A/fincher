import { globalStyle, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const page = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  height: '100%',
  overflowY: 'auto',
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

export const list = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.xs} ${vars.space['2xl']}`,
})

export const row = style({
  position: 'relative',
  display: 'grid',
  gridTemplateColumns: '40px minmax(0, 320px) 1fr 150px 130px 84px',
  alignItems: 'center',
  columnGap: vars.space.lg,
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

export const posterThumb = style({
  width: '40px',
  height: '56px',
  borderRadius: vars.radii.xs,
  objectFit: 'cover',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const nameStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '4px',
  minWidth: 0,
})

export const cardName = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const metaRow = style({
  display: 'flex',
  alignItems: 'baseline',
  gap: vars.space['2xs'],
  overflow: 'hidden',
})

export const metaVersion = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textSecondary,
  letterSpacing: '0.02em',
  flexShrink: 0,
})

export const metaDivider = style({
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.xs,
  flexShrink: 0,
})

export const metaTerritories = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const statusStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '6px',
  minWidth: 0,
})

export const statusBadge = style({
  alignSelf: 'flex-start',
})

export const statusNote = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const scheduleStack = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-end',
  gap: '4px',
})

export const scheduleLabel = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

export const scheduleLabelMuted = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
})

export const countdownValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.lg,
  fontWeight: 500,
  fontVariantNumeric: 'tabular-nums',
  color: vars.color.textPrimary,
  letterSpacing: '0.01em',
})

export const countdownEmpty = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.sm,
  color: vars.color.textTertiary,
})

export const actions = style({
  display: 'flex',
  justifyContent: 'flex-end',
  gap: vars.space.sm,
  opacity: 0,
  transition: 'opacity 0.15s ease',
})

globalStyle(`${row}:hover ${actions}, ${row}:focus-within ${actions}`, {
  opacity: 1,
})

export const actionLink = style({
  appearance: 'none',
  background: 'none',
  border: 'none',
  padding: 0,
  fontFamily: fonts.sans,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  cursor: 'pointer',
  transition: 'color 0.15s ease',
  ':hover': {
    color: vars.color.textPrimary,
  },
  ':focus-visible': {
    color: vars.color.textPrimary,
    outline: `2px solid ${vars.color.primary}`,
    outlineOffset: '2px',
  },
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
})
