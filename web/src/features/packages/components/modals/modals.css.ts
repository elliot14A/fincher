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

export const componentToggleGroup = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(4, 1fr)',
  gap: vars.space.xs,
})

export const componentButton = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: `${vars.space.sm} ${vars.space.xs}`,
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.sm,
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  cursor: 'pointer',
  transition: 'all 0.15s ease',
  gap: vars.space['2xs'],
  ':hover': {
    borderColor: vars.color.primary,
    color: vars.color.textPrimary,
  },
})

export const componentButtonActive = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
  color: vars.color.primary,
  fontWeight: 600,
})
