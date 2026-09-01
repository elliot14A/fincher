import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'

export const form = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
})

export const formRow = style({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: vars.space.md,
  '@media': {
    '(max-width: 500px)': {
      gridTemplateColumns: '1fr',
    },
  },
})

export const inputWithPrefix = style({
  display: 'flex',
  alignItems: 'center',
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.sm,
  overflow: 'hidden',
  ':focus-within': {
    borderColor: vars.color.primary,
    boxShadow: `0 0 0 1px ${vars.color.primary}`,
  },
})

export const prefix = style({
  padding: `${vars.space.sm} ${vars.space.md}`,
  backgroundColor: vars.color.surfaceHover,
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.xs,
  fontFamily: 'monospace',
  borderRight: `1px solid ${vars.color.borderSubtle}`,
  userSelect: 'none',
})

export const bareInput = style({
  border: 'none',
  backgroundColor: 'transparent',
  color: vars.color.textPrimary,
  fontSize: vars.fontSize.sm,
  padding: `${vars.space.sm} ${vars.space.md}`,
  width: '100%',
  outline: 'none',
  '::placeholder': {
    color: vars.color.textTertiary,
  },
})

export const pillGroup = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.xs,
})

export const pill = style({
  display: 'inline-flex',
  alignItems: 'center',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  borderRadius: vars.radii.sm,
  border: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surface,
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.xs,
  cursor: 'pointer',
  userSelect: 'none',
  transition: 'all 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
  },
})

export const pillActive = style({
  backgroundColor: vars.color.primaryMuted,
  borderColor: vars.color.primaryBorder,
  color: vars.color.textInverse,
  fontWeight: 500,
})
