import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const container = style({
  position: 'relative',
  width: '100%',
})

export const trigger = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.sm,
  width: '100%',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  borderRadius: vars.radii.sm,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surface,
  color: vars.color.textPrimary,
  fontFamily: fonts.sans,
  fontSize: vars.fontSize.sm,
  cursor: 'pointer',
  textAlign: 'left',
  transition: 'border-color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
  },
})

export const triggerOpen = style({
  borderColor: vars.color.primary,
})

export const triggerLabel = style({
  minWidth: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const placeholder = style({
  color: vars.color.textTertiary,
})

export const chevron = style({
  flexShrink: 0,
  color: vars.color.textTertiary,
})

export const panel = style({
  position: 'absolute',
  top: 'calc(100% + 4px)',
  left: 0,
  right: 0,
  zIndex: 60,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.sm,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
})

export const searchRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  color: vars.color.textTertiary,
})

export const searchInput = style({
  flex: 1,
  minWidth: 0,
  border: 'none',
  outline: 'none',
  background: 'transparent',
  color: vars.color.textPrimary,
  fontFamily: fonts.sans,
  fontSize: vars.fontSize.sm,
  '::placeholder': {
    color: vars.color.textTertiary,
  },
})

export const list = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  padding: '4px',
  maxHeight: '260px',
  overflowY: 'auto',
})

export const option = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.sm,
  width: '100%',
  padding: '6px 8px',
  borderRadius: vars.radii.xs,
  border: 'none',
  background: 'transparent',
  cursor: 'pointer',
  textAlign: 'left',
  color: vars.color.textSecondary,
  transition: 'background-color 0.12s ease, color 0.12s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
  },
})

export const optionActive = style({
  backgroundColor: vars.color.primaryMuted,
  color: vars.color.textPrimary,
})

export const optionMain = style({
  minWidth: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  fontSize: vars.fontSize.sm,
})

export const optionMeta = style({
  flexShrink: 0,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.03em',
})

export const empty = style({
  padding: `${vars.space.sm} ${vars.space.md}`,
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  textAlign: 'center',
})
