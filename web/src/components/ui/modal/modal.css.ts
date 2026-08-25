import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const backdrop = style({
  position: 'fixed',
  inset: 0,
  backgroundColor: 'rgba(5, 5, 8, 0.78)',
  backdropFilter: 'blur(6px)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  zIndex: 1000,
  minHeight: '100vh',
  padding: vars.space.lg,
})

export const modalContainer = style({
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.lg,
  width: '100%',
  maxWidth: '540px',
  boxShadow: '0 24px 48px rgba(0, 0, 0, 0.65), 0 0 0 1px rgba(255, 255, 255, 0.05)',
  display: 'flex',
  flexDirection: 'column',
  maxHeight: 'calc(100vh - 48px)',
  overflow: 'hidden',
})

export const header = style({
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  padding: `${vars.space.lg} ${vars.space.xl}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
})

export const titleGroup = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const title = style({
  fontSize: vars.fontSize.lg,
  fontWeight: 600,
  color: vars.color.textPrimary,
  lineHeight: vars.lineHeight.tight,
  margin: 0,
})

export const description = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: vars.lineHeight.normal,
  margin: 0,
})

export const closeButton = style({
  background: 'transparent',
  border: 'none',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  padding: vars.space['2xs'],
  borderRadius: vars.radii.sm,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'color 0.15s ease, background-color 0.15s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
  ':focus-visible': {
    outline: `2px solid ${vars.color.primary}`,
    outlineOffset: '2px',
  },
})

export const body = style({
  padding: `${vars.space.lg} ${vars.space.xl}`,
  overflowY: 'auto',
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
})

export const footer = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  gap: vars.space.sm,
  padding: `${vars.space.md} ${vars.space.xl}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surface,
})
