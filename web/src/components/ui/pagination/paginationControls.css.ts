import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const paginationBar = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: `${vars.space.md} ${vars.space['2xl']}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.background,
  marginTop: 'auto',
  flexShrink: 0,
})

export const pageButton = style({
  padding: `${vars.space.xs} ${vars.space.md}`,
  borderRadius: vars.radii.sm,
  fontSize: vars.fontSize.xs,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textSecondary,
  backgroundColor: vars.color.surfaceElevated,
  cursor: 'pointer',
  transition: 'color 0.12s ease, border-color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    borderColor: vars.color.borderStrong,
    backgroundColor: vars.color.surfaceHover,
  },
  ':focus-visible': {
    outline: `2px solid ${vars.color.primary}`,
    outlineOffset: '2px',
  },
  ':disabled': {
    opacity: 0.45,
    cursor: 'not-allowed',
    pointerEvents: 'none',
  },
})

export const pageInfo = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  fontVariantNumeric: 'tabular-nums',
})
