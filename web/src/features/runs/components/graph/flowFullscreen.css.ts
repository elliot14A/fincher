import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

const fadeIn = keyframes({
  from: { opacity: '0' },
  to: { opacity: '1' },
})

export const backdrop = style({
  position: 'fixed',
  inset: 0,
  zIndex: 100,
  backgroundColor: vars.color.background,
  display: 'flex',
  flexDirection: 'column',
  animation: `${fadeIn} 0.12s ease`,
})

export const bar = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.md,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surfaceElevated,
  flexShrink: 0,
})

export const barTitleGroup = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  minWidth: 0,
})

export const barTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const barSubtitle = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const barRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const toggleGroup = style({
  display: 'inline-flex',
  padding: '2px',
  borderRadius: vars.radii.sm,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surface,
})

export const toggleBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
  borderRadius: vars.radii.xs,
  border: 'none',
  background: 'none',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  transition: 'color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    color: vars.color.textSecondary,
  },
})

export const toggleBtnActive = style({
  color: vars.color.textPrimary,
  backgroundColor: vars.color.surfaceActive,
})

export const closeBtn = style({
  background: 'transparent',
  border: `1px solid ${vars.color.border}`,
  padding: vars.space.xs,
  color: vars.color.textSecondary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const stage = style({
  flex: 1,
  minHeight: 0,
  position: 'relative',
})

export const loading = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  height: '100%',
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.sm,
})
