import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const sidebarContainer = style({
  width: '244px',
  backgroundColor: vars.color.surface,
  borderRight: `1px solid ${vars.color.borderSubtle}`,
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.md} ${vars.space.sm}`,
  gap: vars.space['3xs'],
  flexShrink: 0,
  height: '100vh',
  overflowY: 'auto',
})

export const brandRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  marginBottom: vars.space.md,
})

export const brandSubtitle = style({
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.xs,
})

export const composeButton = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.xs,
  width: '100%',
  height: '32px',
  backgroundColor: vars.color.primary,
  border: 'none',
  borderRadius: vars.radii.sm,
  color: vars.color.textInverse,
  fontSize: vars.fontSize.sm,
  fontWeight: 500,
  textDecoration: 'none',
  cursor: 'pointer',
  marginBottom: vars.space.sm,
  transition: 'background-color 0.1s ease',
  ':hover': {
    backgroundColor: vars.color.primaryHover,
  },
  ':focus-visible': {
    outline: `2px solid ${vars.color.primary}`,
    outlineOffset: '2px',
  },
})

export const navItem = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  height: '30px',
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.sm,
  borderRadius: vars.radii.sm,
  textDecoration: 'none',
  transition: 'color 0.1s ease, background-color 0.1s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const navItemActive = style({
  backgroundColor: vars.color.surfaceHover,
  color: vars.color.textPrimary,
  fontWeight: 500,
})

export const navItemLabel = style({
  flex: 1,
})
