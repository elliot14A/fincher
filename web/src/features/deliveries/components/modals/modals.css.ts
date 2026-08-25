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

export const quickCountryChips = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space['2xs'],
  marginTop: vars.space['2xs'],
})

export const countryChip = style({
  background: vars.color.surface,
  border: `1px solid ${vars.color.borderSubtle}`,
  color: vars.color.textSecondary,
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  fontSize: vars.fontSize['2xs'],
  cursor: 'pointer',
  transition: 'all 0.15s ease',
  ':hover': {
    borderColor: vars.color.primary,
    color: vars.color.textPrimary,
  },
})

export const countryChipActive = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
  color: vars.color.primary,
  fontWeight: 600,
})
