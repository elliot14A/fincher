import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const node = style({
  width: '200px',
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  padding: vars.space.sm,
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
  cursor: 'pointer',
  transition: 'border-color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const nodeSelected = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
})

export const accentIncident = style({
  borderTop: `2px solid ${vars.color.danger}`,
})

export const accentAllocation = style({
  borderTop: `2px solid ${vars.color.primary}`,
})

export const accentResolution = style({
  borderTop: `2px solid ${vars.color.success}`,
})

export const accentDefault = style({
  borderTop: `2px solid ${vars.color.borderStrong}`,
})

export const headerRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.xs,
})

export const triggerLabel = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  fontWeight: 700,
  letterSpacing: '0.05em',
  textTransform: 'uppercase',
  color: vars.color.textSecondary,
})

export const statusDotSuccess = style({
  width: '7px',
  height: '7px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.success,
})

export const statusDotDanger = style({
  width: '7px',
  height: '7px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.danger,
})

export const statusDotNeutral = style({
  width: '7px',
  height: '7px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.textTertiary,
})

export const timeText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const metaRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.xs,
})

export const metaText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const expandBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  padding: '2px 6px',
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surface,
  color: vars.color.textSecondary,
  cursor: 'pointer',
  fontSize: vars.fontSize['3xs'],
  ':hover': {
    color: vars.color.textPrimary,
    borderColor: vars.color.borderStrong,
  },
})
